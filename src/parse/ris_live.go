package parse

import (
	"BGPAlert/common"
	"BGPAlert/config"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/netip"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

const (
	socketUrl  = "ws://ris-live.ripe.net/v1/ws/" // RIS Live websocket url
	risMsgType = "ris_message"
)

var done chan interface{}
var interrupt chan os.Signal

type RisMessageData struct {
	Host   string `json:"host,omitempty"` //aka collector
	Peer   string `json:"peer,omitempty"`
	Path   string `json:"path,omitempty"` //aka ASN
	Prefix string `json:"prefix,omitempty"`
}

type RisMessage struct {
	Type string          `json:"type"`
	Data *RisMessageData `json:"data"`
}

/*example ris message
{
	"type":"ris_message",
	"data": {
		"timestamp":1695269583.730,
		"peer":"217.29.66.158",
		"peer_asn":"24482",
		"id":"217.29.66.158-018ab5f0fb720005",
		"host":"rrc10.ripe.net",
		"type":"UPDATE",
		"path":[24482,6939,38040,23969],
		"community":[[24482,2],[24482,200],[24482,12000],[24482,12040],[24482,12041],[24482,20100]],
		"origin":"IGP",
		"med":0,
		"announcements":[{"next_hop":"217.29.66.158","prefixes":["1.1.249.0/24"]}],
		"withdrawals":[]
	}
}
*/

type RisAnnouncement struct {
	NextHop  string   `json:next_hop`
	Prefixes []string `json:prefixes`
}

type RisLiveMessageData struct {
	Timestamp     float64           `json:timestamp`
	Peer          string            `json:peer`
	PeerAsn       string            `json:"peer_asn"` //"peer_asn":"396998"
	Id            string            `json:id`
	Host          string            `json:host`
	Type          string            `json:type`
	Path          []int             `json:path`
	Community     [][]int           `json:community`
	Origin        string            `json:origin`
	Med           int               `json:med`
	Announcements []RisAnnouncement `json:announcements`
	Withdrawals   []string          `json:withdrawals` //string of prefixes being withdrawn
}

type RisLiveMessage struct {
	Type string             `json:"type"`
	Data RisLiveMessageData `json:"data"`
}

// connects to ris live, starts go routine receiverHandler, manages connection and subscription
func ParseRisLiveData(msgChannel chan common.Message, config *config.Configuration) {

	fmt.Println("starting...")

	//for each subscription
	for _, subscription := range config.Subscriptions {
		//go routine handleSubscription
		go handleSubscription(msgChannel, subscription)
	}

	// alternatives:
	// this would listen to one of Fastly's blocks of address space, from all collectors:
	//subscription1 := RisMessage{"ris_subscribe", &RisMessageData{"", "151.101.0.0/16"}}
	// this would listen to all of the IPv4 address space, but from only one collector:
	//subscription1 := RisMessage{"ris_subscribe", &RisMessageData{"rrc21", "0.0.0.0/0"}}

}

// handles the connection for each subscription
func handleSubscription(msgChannel chan common.Message, subscription config.SubscriptionMsg) error {

	// create websocket connection to ris live websocket
	conn, _, err := websocket.DefaultDialer.Dial(socketUrl, nil)
	if err != nil {
		return errors.New("websocket connection err, " + err.Error())
	}
	defer conn.Close()

	//call receive handler
	go receiveHandler(msgChannel, conn, subscriptionToString(subscription))

	fmt.Println("made connection")

	//subscribe
	sub := RisMessage{"ris_subscribe", &RisMessageData{subscription.Host, subscription.Peer, subscription.Path, subscription.Prefix}}
	out, err := json.Marshal(sub)
	if err != nil {
		return errors.New("Error marshalling subscription message, " + err.Error())
	}
	log.Println("Subscribing to: ", subscriptionToString(subscription))
	conn.WriteMessage(websocket.TextMessage, out)

	//manage connection
	/* Ping message (re-send this every minute or so */
	ping := RisMessage{"ping", nil}
	pingstr, err := json.Marshal(ping)
	if err != nil {
		return errors.New("Error marshalling ping message (!), " + err.Error())
	}

	for {
		select {
		case <-time.After(time.Duration(60) * time.Millisecond * 1000):
			// Send an echo packet 60 seconds
			err := conn.WriteMessage(websocket.TextMessage, pingstr)
			if err != nil {
				return errors.New("Error during writing to websocket " + err.Error())
			}

		case <-interrupt:
			// We received a SIGINT; clean up
			log.Println("Received SIGINT interrupt signal. Closing all pending connections")
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				return errors.New("Error during closing websocket: " + err.Error())
			}

			select {
			case <-done:
				log.Println("Receiver channel closed, exiting")
			case <-time.After(time.Duration(1) * time.Second):
				log.Println("Timeout in closing receiving channel; exiting")
			}
			return nil
		}
	}

}

// keep reading in new message from connection, send msg to parser, put parsed messages into channel
func receiveHandler(msgChannel chan common.Message, conn *websocket.Conn, subscription string) {
	var labeledMsg common.Message
	labeledMsg.Filter = subscription

	for {
		//take in next msg from connection
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Websocket read error:", err)
			return
		}

		//parse message to data structure
		bgpMsgs, err := parseLiveMessage(message)
		if err != nil {
			log.Println("Error parsing BGP message", err)
		} else {
			fmt.Printf("Parsed BGP Message: %+v\n", bgpMsgs) //prints parsed BGP msg
		}

		//put bgp messages into channel
		for _, msg := range bgpMsgs {
			labeledMsg.BGPMessage = msg
			msgChannel <- labeledMsg
		}

	}
}

// parses json message into common.BGPmessage struct
// returns array of common.BGPMessage because it separates by prefix and type of update (A or W)
func parseLiveMessage(data []byte) ([]common.BGPMessage, error) {
	var parsedMsgs []common.BGPMessage
	var parsedMsg common.BGPMessage

	//initial format of ris live message is TYPE and DATA
	var message RisLiveMessage
	err := json.Unmarshal(data, &message) //decode data from JSON to a struct
	if err != nil {
		log.Println("Bad parse:", err)
		return nil, err
		//log.Println("Original message:", data)
	}

	//check is TYPE is ris message
	if message.Type == risMsgType {
		payload := message.Data

		//check if message data is of type UPDATE - meaning announcement or withdrawal (or both)
		if payload.Type == "UPDATE" {

			//timestamp
			parsedMsg.Timestamp, err = float64ToTime(payload.Timestamp)
			if err != nil {
				return []common.BGPMessage{}, fmt.Errorf("error parsing timestamp: %v", err)
			}

			//peerIP
			parsedMsg.PeerIP, err = netip.ParseAddr(payload.Peer)
			if err != nil {
				return []common.BGPMessage{}, fmt.Errorf("error parsing address: %v", err)
			}

			//peerASN
			parsedMsg.PeerASN, err = parseUint32(payload.PeerAsn)
			if err != nil {
				return []common.BGPMessage{}, fmt.Errorf("error parsing peer asn: %v", err)
			}

			//PREFIXES

			//for each announcement in the JSON message - make a new parsed message
			for _, announcement := range payload.Announcements {
				//for each prefix of each announcement
				for _, prefix := range announcement.Prefixes {

					//bgpmessagetype
					parsedMsg.BGPMessageType = common.AnnouncementType

					//prefix
					parsedMsg.Prefix, err = netip.ParsePrefix(prefix)
					if err != nil {
						return []common.BGPMessage{}, fmt.Errorf("error parsing prefix: %v", err)
					}

					//add this parsed message to list of parsed messages
					parsedMsgs = append(parsedMsgs, parsedMsg)
				}
			}

			//for each withdrawal in the JSON message - make a new parsed message
			for _, withdrawal := range payload.Withdrawals {
				//bgpmessagetype
				parsedMsg.BGPMessageType = common.WithdrawalType

				//prefix
				parsedMsg.Prefix, err = netip.ParsePrefix(withdrawal)
				if err != nil {
					return []common.BGPMessage{}, fmt.Errorf("error parsing prefix: %v", err)
				}

				//add this parsed message to list of parsed messages
				parsedMsgs = append(parsedMsgs, parsedMsg)
			}

		}
	} else { //if not ris message type
		log.Println("Received unhandled message:", message)
	}

	return parsedMsgs, nil
}

// converts float64 to time.Time
func float64ToTime(timestamp float64) (time.Time, error) {
	// Extract seconds and nanoseconds
	seconds := int64(timestamp)
	nanos := int64((timestamp - float64(seconds)) * 1e9)

	// Perform error handling
	if nanos < 0 || nanos >= 1e9 {
		return time.Time{}, fmt.Errorf("invalid timestamp: %f", timestamp)
	}

	// Convert to time.Time
	resultTime := time.Unix(seconds, nanos)
	return resultTime, nil
}

// toString for subscription struct
// used for labeling each subscription filter
func subscriptionToString(sub config.SubscriptionMsg) string {
	result := "Subscription{"

	if sub.Host != "" {
		result += fmt.Sprintf("Host: %q, ", sub.Host)
	}

	if sub.Peer != "" {
		result += fmt.Sprintf("Peer: %q, ", sub.Peer)
	}

	if sub.Path != "" {
		result += fmt.Sprintf("Path: %q, ", sub.Path)
	}

	if sub.Prefix != "" {
		result += fmt.Sprintf("Prefix: %q, ", sub.Prefix)
	}

	//remove the trailing comma and space if there is at least one field
	if result != "MyStruct{" {
		result = result[:len(result)-2]
	}

	result += "}"

	return result
}
