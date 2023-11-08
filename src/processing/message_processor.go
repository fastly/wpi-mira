package processing

import (
	"BGPAlert/common"
	"fmt"
)

func ProcessBGPMessages(msgChannel chan []common.BGPMessage) {
	var bucketMap = make(map[int64][]common.BGPMessage)

	for bgpMessages := range msgChannel {
		for _, msg := range bgpMessages {
			timestamp := int64(msg.Timestamp)

			// Round down the timestamp to the nearest multiple of 60 seconds
			bucketTimestamp := (timestamp / 60) * 60

			// Check if a bucket for the rounded timestamp exists, create it if not
			if _, ok := bucketMap[bucketTimestamp]; !ok {
				bucketMap[bucketTimestamp] = make([]common.BGPMessage, 0)
			}

			// Append the message to the corresponding bucket
			bucketMap[bucketTimestamp] = append(bucketMap[bucketTimestamp], msg)

			// Your processing logic here
			fmt.Println("Received BGP Message:", msg)
		}
	}

	for timestamp, messages := range bucketMap {
		frequency := len(messages)
		fmt.Printf("Timestamp: %d, Frequency: %d\n", timestamp, frequency)
	}
	fmt.Println(len(bucketMap))

}
