package parse

import (
	"BGPAlert/common"
	"encoding/json"
	"fmt"
	"log"
	"net/netip"
	"time"

	"github.com/gorilla/websocket"
)

// RIS Live websocket url
const socketUrl = "ws://ris-live.ripe.net/v1/ws/"

// connects to ris live, parsing messages, and putting messages into msgChannel for processor
func ParseRisLiveData(msgChannel chan []common.BGPMessage) {

	fmt.Println("starting...")

	// create websocket connection to ris live websocket
	conn, _, err := websocket.DefaultDialer.Dial(socketUrl, nil)
	if err != nil {
		log.Fatal("Websocket connection error:", err)
		return
	}
	defer conn.Close()

	fmt.Println("made connection")

	//keep reading in new message from connection
	for {

		fmt.Println("in for loop")

		//take in next msg from connection
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Websocket read error:", err)
			return
		}

		fmt.Println("read message") //NOT REACHING HERE

		//parse message to data structure
		bgpMsgs, err := parseLiveMessage(message)
		if err != nil {
			log.Println("Error parsing BGP message", err)
		} else {
			fmt.Printf("Parsed BGP Message: %+v\n", bgpMsgs) //prints parsed BGP msg
		}

		//put bgp messages into channel
		msgChannel <- bgpMsgs

	}
}

//example ris message
//{
//	"type":"ris_message",
//	"data": {
//		"timestamp":1695269583.730,
//		"peer":"217.29.66.158",
//		"peer_asn":"24482",
//		"id":"217.29.66.158-018ab5f0fb720005",
//		"host":"rrc10.ripe.net",
//		"type":"UPDATE",
//		"path":[24482,6939,38040,23969],
//		"community":[[24482,2],[24482,200],[24482,12000],[24482,12040],[24482,12041],[24482,20100]],
//		"origin":"IGP",
//		"med":0,
//		"announcements":[{"next_hop":"217.29.66.158","prefixes":["1.1.249.0/24"]}],
//		"withdrawals":[]
//	}
//}

type RisAnnouncement struct {
	NextHop  string   `json:next_hop`
	Prefixes []string `json:prefixes`
}

type RisWithdrawal struct {
	Withdrawal string //withdrawn IP prefix
}

type RisLiveMessageData struct {
	Timestamp     float64           `json:timestamp`
	Peer          string            `json:peer`
	PeerAsn       string            `json:"peer_asn"` //,"peer_asn":"396998"
	Id            string            `json:id`
	Host          string            `json:host`
	Type          string            `json:type`
	Path          []int             `json:path`
	Community     [][]int           `json:community`
	Origin        string            `json:origin`
	Med           int               `json:med`
	Announcements []RisAnnouncement `json:announcements`
	Withdrawals   []RisWithdrawal   `json:withdrawals`
}

type RisLiveMessage struct {
	Type string             `json:"type"`
	Data RisLiveMessageData `json:"data"`
}

func parseLiveMessage(data []byte) ([]common.BGPMessage, error) {
	var parsedMsgs []common.BGPMessage
	var parsedMsg common.BGPMessage

	//initial format of ris live message is TYPE and DATA
	var message RisLiveMessage
	err := json.Unmarshal(data, &message) //decode data from JSON to a struct
	if err != nil {
		log.Println("Bad parse:", err)
		log.Println("Original message:", data)
	}

	//check is TYPE is ris message
	if message.Type != "ris_message" {
		log.Println("Received unhandled message:", message)
	} else {
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
					parsedMsg.BGPMessageType = "A"

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
				parsedMsg.BGPMessageType = "W"

				//prefix
				parsedMsg.Prefix, err = netip.ParsePrefix(withdrawal.Withdrawal)
				if err != nil {
					return []common.BGPMessage{}, fmt.Errorf("error parsing prefix: %v", err)
				}

				//add this parsed message to list of parsed messages
				parsedMsgs = append(parsedMsgs, parsedMsg)
			}

		}
	}

	return parsedMsgs, nil
}

// converts float64 to time.Time
func float64ToTime(timestamp float64) (time.Time, error) {
	// Extract seconds and nanoseconds
	seconds := int64(timestamp)
	nanos := int64((timestamp - float64(seconds)) * 1e9)

	// Perform error handling
	if nanos < 0 || nanos >= int64(time.Second) {
		return time.Time{}, fmt.Errorf("invalid timestamp: %f", timestamp)
	}

	// Convert to time.Time
	resultTime := time.Unix(seconds, nanos)
	return resultTime, nil
}
