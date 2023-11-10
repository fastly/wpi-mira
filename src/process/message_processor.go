package process

import (
	"BGPAlert/common"
	"fmt"
)

func ProcessBGPMessages(msgChannel chan []common.BGPMessage, windowChannel chan []common.Window) {
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

			//fmt.Println("Received BGP Message:", msg)
		}
	}

	for timestamp, messages := range bucketMap {
		frequency := len(messages)
		fmt.Printf("Timestamp: %d, Frequency: %d\n", timestamp, frequency)
	}
	fmt.Println("Minutes Analyzed: ", len(bucketMap))

	var windows []common.Window

	window := common.Window{
		Filter:    "none",
		BucketMap: bucketMap,
	}

	windows = append(windows, window)

	windowChannel <- windows

	close(windowChannel)

}
