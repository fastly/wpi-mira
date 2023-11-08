package parsing

import (
	"BGPAlert/common"
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func ParseStaticFile(folderDir string, msgChannel chan []common.BGPMessage) {
	fmt.Println("test")
	// Gets all files from directory we want to test
	files, err := os.ReadDir("staticdata/" + folderDir)
	if err != nil {
		log.Fatal(err)
	}

	// Iterate through the files in the directory
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".bz2" {
			filePath := filepath.Join("testdata/"+folderDir, file.Name())
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

func parseBGPFile(filePath string) ([]common.BGPMessage, error) {
	cmd := exec.Command("bgpdump", "-m", "-O", "tempbgpdump.txt", filePath)
	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	file, err := os.Open("tempbgpdump.txt")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var bgpMessages []common.BGPMessage
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		//fmt.Println(line)
		bgpMsg, err := parseBGPMessage(line)
		if err != nil {
			fmt.Printf("Error parsing message: %v\n", err)
			continue // Skip the invalid message
		}
		bgpMessages = append(bgpMessages, bgpMsg)
	}

	os.Remove("tempbgpdump.txt")

	if scanner.Err() != nil {
		return nil, scanner.Err()
	}

	return bgpMessages, nil
}

func parseBGPMessage(data string) (common.BGPMessage, error) {
	fields := strings.Split(data, "|")

	// If message type is withdrawal (W) -> there's only the first 6 fields
	if len(fields) == 6 && fields[2] == "W" {
		return common.BGPMessage{
			Timestamp:      parseTimestamp(fields[1]),
			BGPMessageType: fields[2],
			PeerIP:         fields[3],
			PeerASN:        parseUint32(fields[4]),
			Prefix:         fields[5],
		}, nil
	} else if len(fields) == 15 && fields[2] == "A" {
		// If the message is advertisement (A) -> parse like normal
		return common.BGPMessage{
			Timestamp:      parseTimestamp(fields[1]),
			BGPMessageType: fields[2],
			PeerIP:         fields[3],
			PeerASN:        parseUint32(fields[4]),
			Prefix:         fields[5],
		}, nil
	}

	return common.BGPMessage{}, fmt.Errorf("Invalid BGP message: %s", data)
}

func parseTimestamp(timestampStr string) float64 {
	timestamp, err := strconv.ParseFloat(timestampStr, 64)
	if err != nil {
		fmt.Println("Error Parsing Float from Timestamp: ", err)
		return 0.0
	}
	return timestamp
}

func parseUint32(valueStr string) uint32 {
	value, _ := strconv.ParseUint(valueStr, 10, 32)
	return uint32(value)
}
