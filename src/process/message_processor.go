package process

import (
	"BGPAlert/analyze"
	"BGPAlert/common"
	"BGPAlert/config"
	"fmt"
	"strconv"
	"time"
)

// Constantly read messages from channel, build up windows with frequency maps, calling analysis on windows when full
func ProcessBGPMessages(msgChannel chan common.BGPMessage, config *config.Configuration) error {
	maximumBuckets, _ := strconv.Atoi(config.WindowSize)
	// Each bucket is 60s
	parseDuration := 60 * time.Second
	maximumTimespan := time.Duration(maximumBuckets*60) * time.Second

	window := common.Window{
		Filter:    "none",
		BucketMap: make(map[time.Time][]common.BGPMessage),
	}

	// Read messages from channel
	for msg := range msgChannel {
		// Round timestamp to size of parseDuration to place in correct bucket
		messageBucket := msg.Timestamp.Truncate(parseDuration)

		// If a bucket with this timestamp doesn't exist, need to create one
		if _, ok := window.BucketMap[messageBucket]; !ok {
			fmt.Println("----------------------------------------------------------------------------------")

			// If there's at least one timestamp, may have to fill in missing zeroes in map
			if len(window.BucketMap) > 1 {
				// Get minimum timestamp from map
				minTimestamp := getBucketMapMin(window.BucketMap)

				// Walk through map and see if any zeroes are missing
				for tempTimestamp := minTimestamp; tempTimestamp.Before(messageBucket); tempTimestamp = tempTimestamp.Add(time.Minute) {
					if _, ok := window.BucketMap[tempTimestamp]; !ok {
						fmt.Println("Appended a 0 to map at ", tempTimestamp)
						window.BucketMap[tempTimestamp] = make([]common.BGPMessage, 0)
					}
				}
			}

			// If we at least maximumBuckets in bucketMap we can perform analysis
			if len(window.BucketMap) >= maximumBuckets {
				fmt.Println(fmt.Sprintf("len(window.BucketMap): %d maximumBuckets: %d", len(window.BucketMap), maximumBuckets))

				// First want to remove timestamps out of scope so len(bucketMap) == maximumBuckets
				minimumTimestamp := messageBucket.Add(-maximumTimespan)
				for timestamp := range window.BucketMap {
					if timestamp.Before(minimumTimestamp) {
						fmt.Println("Expired bucket: ", timestamp)
						delete(window.BucketMap, timestamp)
					}
				}

				// Now window is ready for analysis
				analyze.AnalyzeBGPMessages(window)
			}

			// Create new bucket for new timestamp
			fmt.Println("Creating bucket: ", messageBucket)
			window.BucketMap[messageBucket] = make([]common.BGPMessage, 0)
		}

		// Append the message to the corresponding bucket
		window.BucketMap[messageBucket] = append(window.BucketMap[messageBucket], msg)
	}

	return nil
}

// Returns the minimum key value of a bucketMap
func getBucketMapMin(bucketMap map[time.Time][]common.BGPMessage) time.Time {
	var minTimestamp time.Time

	for timestamp := range bucketMap {
		minTimestamp = timestamp
		break
	}

	for timestamp := range bucketMap {
		if timestamp.Before(minTimestamp) {
			minTimestamp = timestamp
		}
	}

	return minTimestamp
}
