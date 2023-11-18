package parse

import (
	"BGPAlert/common"
	"bufio"
	"fmt"
	"log"
	"net/netip"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// Reads in static files from a directory, parses them into BGPMessage struct and puts them into channel for processor to process
func ParseStaticFile(folderDir string, msgChannel chan []common.BGPMessage) {
	const directoryPath = "static_data/"

	// Gets all files from directory we want to test
	files, err := os.ReadDir(directoryPath + folderDir)
	if err != nil {
		log.Fatal(err)
		return
	}

	// Iterate through the files in the directory
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".bz2" {
			filePath := filepath.Join(directoryPath, folderDir, file.Name())
			bgpMessages, err := parseBGPFile(filePath)
			if err != nil {
				fmt.Println("Error parsing file:", err)
				continue
			}

			// Send the parsed BGP messages to the channel
			msgChannel <- bgpMessages
		}
	}

	// Close the channel to signal that all messages have been sent
	close(msgChannel)
}

// Runs BGPDump on individual static bgp files and extracts all BGPMessages from that individual file
func parseBGPFile(filePath string) ([]common.BGPMessage, error) {
	const tempFileName = "tempbgpdump.txt"

	cmd := exec.Command("bgpdump", "-m", "-O", tempFileName, filePath)
	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	file, err := os.Open(tempFileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var bgpMessages []common.BGPMessage
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		bgpMsg, err := parseBGPMessage(line)
		if err != nil {
			fmt.Printf("Error parsing message: %v\n", err)
			continue // Skip the invalid message
		}
		bgpMessages = append(bgpMessages, bgpMsg)
	}

	os.Remove(tempFileName)

	if scanner.Err() != nil {
		return nil, scanner.Err()
	}

	return bgpMessages, nil
}

// Takes a line of BGPDumped data and turns it into the BGPMessage struct
func parseBGPMessage(data string) (common.BGPMessage, error) {
	const announcementFields = 15
	const announcmentType = "A"
	const withdrawalFields = 6
	const withdrawalType = "W"

	fields := strings.Split(data, "|")

	timestamp, err := parseTimestamp(fields[1])
	if err != nil {
		return common.BGPMessage{}, fmt.Errorf("error parsing timestamp: %v", err)
	}

	bgpMessageType := fields[2]

	peerIP, err := netip.ParseAddr(fields[3])
	if err != nil {
		return common.BGPMessage{}, fmt.Errorf("error parsing address: %v", err)
	}

	peerASN, err := parseUint32(fields[4])
	if err != nil {
		return common.BGPMessage{}, fmt.Errorf("error parsing peer asn: %v", err)
	}

	prefix, err := netip.ParsePrefix(fields[5])
	if err != nil {
		return common.BGPMessage{}, fmt.Errorf("error parsing prefix: %v", err)
	}

	// Announcment examples: (Has 15 fields and messageType must be "A")
	// BGP4MP_ET|1638317699.191812|A|206.82.104.185|398465|102.66.116.0/24|398465 5713 37457 37457 37457 37457 37457 37457 37457 328471 328471 328471|IGP|206.82.104.185|0|0|5713:800 65101:1085 65102:1000 65103:276 65104:150|NAG|4200000002 10.102.100.2|
	// BGP4MP_ET|1638317699.744853|A|2001:504:36::6:1481:0:1|398465|2a10:cc42:1bb9::/48|398465 174 1299 20473|IGP|2001:504:36::6:1481:0:1|0|0|174:21000 174:22003|NAG||
	// BGP4MP_ET|1638317700.043880|A|206.82.104.185|398465|102.66.116.0/24|398465 30844 328471 328471 328471|IGP|206.82.104.185|0|0|30844:27 65101:1082 65102:1000 65103:276 65104:150|NAG|4200000002 10.102.100.2|

	// Withdrawal examples: (Has 6 fields and third field messageType must be "W")
	// BGP4MP_ET|1638317694.706880|W|2001:504:36::6:1481:0:1|398465|2a10:cc42:131c::/48
	// BGP4MP_ET|1638317679.511516|W|2001:504:36::6:1481:0:1|398465|2804:2688::/33

	if (len(fields) == announcementFields && bgpMessageType == announcmentType) || (len(fields) == withdrawalFields && bgpMessageType == withdrawalType) {
		return common.BGPMessage{
			Timestamp:      timestamp,
			BGPMessageType: bgpMessageType,
			PeerIP:         peerIP,
			PeerASN:        peerASN,
			Prefix:         prefix,
		}, nil
	}

	return common.BGPMessage{}, fmt.Errorf("invalid bgp message: %s", data)
}

// Converts a string timestamp into a time.Time object
func parseTimestamp(timestampStr string) (time.Time, error) {
	timestamp, err := strconv.ParseFloat(timestampStr, 64)
	if err != nil {
		return time.Time{}, fmt.Errorf("error parsing float from timestamp: %v", err)
	}

	seconds := int64(timestamp)
	nanoseconds := int64((timestamp - float64(seconds)) * 1e9)

	time := time.Unix(seconds, nanoseconds)

	return time, nil
}

// Converts a string uint and converts it into an uint32
func parseUint32(valueStr string) (uint32, error) {
	value, err := strconv.ParseUint(valueStr, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("error parsing uint32: %v", err)
	}

	return uint32(value), nil
}
