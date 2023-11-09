package analyzing

import (
	"BGPAlert/blt_mad"
	"BGPAlert/common"
	"BGPAlert/shake_alert"
	"fmt"
	"sort"
)

func AnalyzeBGPMessages(windowChannel chan []common.Window) {
	for windows := range windowChannel {
		for _, w := range windows {
			fmt.Println("Received Window: ")
			bucketMap := w.BucketMap

			// Convert BucketMap to a map of timestamp to length of messages
			lengthMap := make(map[int64]int64)

			for timestamp, messages := range bucketMap {
				lengthMap[timestamp] = int64(len(messages))
			}

			// Turn map into sorted array of frequencies by timestamp
			sortedFrequencies := getSortedFrequencies(lengthMap)

			// Convert to float array for analysis functions
			floatArray := int64ArrayToFloat64Array(sortedFrequencies)
			fmt.Println(floatArray)

			fmt.Println("BLT MAD Outliers: ")
			fmt.Println(blt_mad.BltMad(floatArray, 10))

			fmt.Println("ShakeAlert Outliers: ")
			fmt.Println(shake_alert.FindOutliers(floatArray))
		}
	}

}

func getSortedFrequencies(bucketMap map[int64]int64) []int64 {
	var timestamps []int64

	// Create a slice of timestamps and a corresponding slice of values in the order of timestamps
	for timestamp, _ := range bucketMap {
		timestamps = append(timestamps, timestamp)
	}

	// Sort the timestamps in ascending order
	sort.Slice(timestamps, func(i, j int) bool {
		return timestamps[i] < timestamps[j]
	})

	// Create a new slice of values in the order of sorted timestamps
	sortedValues := make([]int64, len(timestamps))
	for i, timestamp := range timestamps {
		sortedValues[i] = bucketMap[timestamp]
	}

	return sortedValues
}

func int64ArrayToFloat64Array(intArray []int64) []float64 {
	floatArray := make([]float64, len(intArray))
	for i, v := range intArray {
		floatArray[i] = float64(v)
	}
	return floatArray
}
