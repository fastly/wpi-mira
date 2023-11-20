package process

import (
	"BGPAlert/analyze"
	"BGPAlert/common"
	"errors"
	"fmt"
	"time"
)

// Constantly read messages from channel, build up windows with frequency maps, calling analysis on windows when full
func ProcessBGPMessagesLive(msgChannel chan common.BGPMessage) error {
	maximumBuckets := 43
	// Each bucket is 60s
	parseDuration, err := time.ParseDuration("60s")
	if err != nil {
		return errors.New("error creating parsing bucket window, " + err.Error())
	}
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
				fmt.Println("len(window.BucketMap):", len(window.BucketMap), " maximumBuckets:", maximumBuckets)

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

// Constantly reads channel of messages and stores them in Window objects to send through windowChannel for analysis
func ProcessBGPMessages(msgChannel chan common.BGPMessage, windowChannel chan common.Window) error {
	var bucketMap = make(map[time.Time][]common.BGPMessage)

	for msg := range msgChannel {
		parseDuration, err := time.ParseDuration("60s")
		if err != nil {
			return errors.New("error creating parsing window, " + err.Error())
		}

		// Round down the timestamp to the nearest multiple of 60 seconds
		bucketTimestamp := msg.Timestamp.Truncate(parseDuration)

		// Check if a bucket for the rounded timestamp exists, create it if not
		if _, ok := bucketMap[bucketTimestamp]; !ok {
			bucketMap[bucketTimestamp] = make([]common.BGPMessage, 0)
		}

		// Append the message to the corresponding bucket
		bucketMap[bucketTimestamp] = append(bucketMap[bucketTimestamp], msg)
	}

	for timestamp, messages := range bucketMap {
		frequency := len(messages)
		fmt.Printf("Timestamp: %s, Frequency: %d\n", timestamp.Format("2006-01-02 15:04:05"), frequency)
	}
	fmt.Println("Minutes Analyzed: ", len(bucketMap))

	minTime, maxTime := getMapMinAndMax(bucketMap)
	duration := uint32((maxTime.Sub(minTime)).Minutes()) + 1
	fmt.Println(duration)

	window := common.Window{
		Filter:    "none",
		BucketMap: bucketMap,
	}

	windowChannel <- window

	close(windowChannel)

	return nil
}

// Returns the minimum key value of a bucketMap
func getBucketMapMin(bucketMap map[time.Time][]common.BGPMessage) time.Time {
	var minTimestamp time.Time
	firstIteration := true
	for timestamp := range bucketMap {
		if firstIteration {
			minTimestamp = timestamp
			firstIteration = false
		} else {
			if timestamp.Before(minTimestamp) {
				minTimestamp = timestamp
			}
		}
	}

	return minTimestamp
}

// Returns the minimium and maximum keys of map
func getMapMinAndMax(timestampMap map[time.Time][]common.BGPMessage) (time.Time, time.Time) {
	var minTime, maxTime time.Time
	firstIteration := true

	for timestamp := range timestampMap {
		if firstIteration {
			minTime, maxTime = timestamp, timestamp
			firstIteration = false
		} else {
			if timestamp.Before(minTime) {
				minTime = timestamp
			}
			if timestamp.After(maxTime) {
				maxTime = timestamp
			}
		}
	}

	return minTime, maxTime
}
