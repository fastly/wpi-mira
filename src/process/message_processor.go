package process

import (
	"BGPAlert/common"
	"fmt"
	"time"
)

func ProcessBGPMessages(msgChannel chan []common.BGPMessage, windowChannel chan []common.Window) {
	var bucketMap = make(map[time.Time][]common.BGPMessage)

	for bgpMessages := range msgChannel {
		for _, msg := range bgpMessages {
			windowSize, _ := time.ParseDuration("60s")

			// Round down the timestamp to the nearest multiple of 60 seconds
			bucketTimestamp := msg.Timestamp.Truncate(windowSize)

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
		fmt.Printf("Timestamp: %s, Frequency: %d\n", timestamp.Format("2006-01-02 15:04:05"), frequency)
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
