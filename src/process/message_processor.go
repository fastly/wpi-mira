package process

import (
	"BGPAlert/analyze"
	"BGPAlert/common"
	"errors"
	"fmt"
	"strconv"
	"time"
)

var windowSize = 43

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

func ProcessBGPMessagesLive(msgChannel chan common.BGPMessage, windowChannel chan common.Window) error {
	asn := 54113
	//var bucketMap = make(map[time.Time][]common.BGPMessage)

	window := common.Window{
		Filter:    strconv.Itoa(asn),
		BucketMap: make(map[time.Time][]common.BGPMessage),
	}

	// Constantly read messages from channel
	for msg := range msgChannel {
		// Each bucket is a 60s or 1 minute long bin
		parseDuration, err := time.ParseDuration("60s")
		if err != nil {
			return errors.New("error creating parsing bucket window, " + err.Error())
		}

		// Round down the timestamp to the nearest multiple of 60 seconds to place in correct bin
		bucketTimestamp := msg.Timestamp.Truncate(parseDuration)

		// Check if a bucket for the rounded timestamp exists, create it if not
		if _, ok := window.BucketMap[bucketTimestamp]; !ok {

			// If we're making a new timestamp a minute went by -> check if window is big enough for analysis
			if len(window.BucketMap) >= 1 {
				// Get how many minutes are currently stored in map, 0 counts aren't stored in map so len(map) won't always work
				minTime, maxTime := getMapMinAndMax(window.BucketMap)
				duration := int((maxTime.Sub(minTime)).Minutes()) + 1

				// Is our window big enough for analysis?
				if duration >= windowSize {
					// If our window is bigger we should clean it up to only be windSize first
					if duration > windowSize {
						// Need to find what minTime should be for duration == windSize
						diff := duration - windowSize
						//newMinTime := minTime.Add(time.Duration(diff) * time.Minute)

						// Delete minute buckets starting from minTime iterating by difference in minutes
						for i := 0; i < diff; i++ {
							delete(window.BucketMap, minTime)
							minTime = minTime.Add(time.Minute)
						}

						// Also need to make sure that newMinTime has an entry in hashmap so it can be the minimum next time
						if _, ok := window.BucketMap[minTime]; !ok {
							window.BucketMap[minTime] = make([]common.BGPMessage, 0)
						}
					}

					// Now we know our map is exactly == windSize -> send it to channel for analysis
					fmt.Println("BucketMap length before channel: ", len(window.BucketMap))

					analyze.AnalyzeBGPMessages(window)
				} else { // If our window size is not big enough for analysis, just create the new bin
					window.BucketMap[bucketTimestamp] = make([]common.BGPMessage, 0)
				}
			} else { // If our window size is less than one, just create a new bin
				window.BucketMap[bucketTimestamp] = make([]common.BGPMessage, 0)
			}
		}

		// Append the message to the corresponding bucket
		window.BucketMap[bucketTimestamp] = append(window.BucketMap[bucketTimestamp], msg)
	}

	close(windowChannel)

	return nil
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
	//fmt.Println(minTime)
	//fmt.Println(maxTime)

	return minTime, maxTime
}
