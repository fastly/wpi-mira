package analyze

import (
	"BGPAlert/blt_mad"
	"BGPAlert/common"
	"BGPAlert/shake_alert"
	"fmt"
	"sort"
	"time"
)

// Takes in a Window, parses object into frequency counts, and then calls specified analysis functions
func AnalyzeBGPMessages(window common.Window) {
	bucketMap := window.BucketMap

	// Convert BucketMap to a map of timestamp to length of messages
	lengthMap := make(map[time.Time]float64)

	for timestamp, messages := range bucketMap {
		lengthMap[timestamp] = float64(len(messages))
	}

	// Turn map into sorted array of frequencies by timestamp
	sortedFrequencies := GetSortedFrequencies(lengthMap)

	fmt.Printf("Sorted Array of Frequencies: \n%+v\n", sortedFrequencies)
	fmt.Printf("BLT MAD Outliers: \n%+v\n", blt_mad.BltMad(sortedFrequencies, 10))
	fmt.Printf("ShakeAlert Outliers: \n%+v\n", shake_alert.FindOutliers(sortedFrequencies))
}

// Takes in map of time objects to frequencies and puts them into an ordered array of frequencies based on increasing timestamps
func GetSortedFrequencies(bucketMap map[time.Time]float64) []float64 {
	var timestamps []time.Time

	// Create a slice of timestamps and a corresponding slice of values in the order of timestamps
	for timestamp := range bucketMap {
		timestamps = append(timestamps, timestamp)
	}

	// Sort the timestamps in ascending order
	sort.Slice(timestamps, func(i, j int) bool {
		return timestamps[i].Before(timestamps[j])
	})

	// Create a new slice of values in the order of sorted timestamps
	sortedValues := make([]float64, len(timestamps))
	for i, timestamp := range timestamps {
		sortedValues[i] = bucketMap[timestamp]
	}

	return sortedValues
}
