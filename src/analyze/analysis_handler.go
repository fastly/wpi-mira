package analyze

import (
	"BGPAlert/blt_mad"
	"BGPAlert/common"
	"BGPAlert/config"
	"BGPAlert/shake_alert"
	"fmt"
	"sort"
	"time"
)

//need to create a map of results based on filter -> filtering pr
var AllResults common.Result //global so that it can be added onto by parser and seen by dataHandler
var maxPoints = 5            //i did not want to add this into config yet because it will conflic with jolene's pr

const (
	OPTPARAM   = 5.0
	BLTMAD     = 1
	SHAKEALERT = 2
	BOTHALERTS = 3
)

// Takes in a Window, parses object into frequency counts, and then calls specified analysis functions
/*
TODO: code to write the frequencies; the outliers; and the minReqs to files
*/
func AnalyzeBGPMessages(window common.Window, config *config.Configuration) {
	lengthMap := makeLengthMap(window)
	// Turn map into sorted array of frequencies by timestamp
	sortedFrequencies, sortedTimestamps := GetSortedFrequencies(lengthMap) //fix these duplicates with time stamps
	frequencies := AllResults.AllFreq

	//the file names will contain all the timestamps for a given folder that was processed
	fmt.Printf("Sorted Array of Frequencies: \n%+v\n", sortedFrequencies)
	bltOutliers, bltIndexes := blt_mad.BltMad(sortedFrequencies, OPTPARAM)
	shakeAlertOutliers, shakeAlertIndexes := shake_alert.FindOutliers(sortedFrequencies)
	fmt.Printf("BLT MAD Outliers: \n%+v\n", bltOutliers)
	fmt.Printf("ShakeAlert Outliers: \n%+v\n", shakeAlertOutliers)

	//check the amount of the frequencies in results
	if len(AllResults.AllFreq)+1 > maxPoints { //adding one more frequency would result in more than needed points
		//remove the first item in the freq map; no modifications to the outlier map
		firstResultKey := getSmallestTimestamp(getListOfKeys(frequencies))
		delete(frequencies, *firstResultKey)

		//now update the results
		for timestamp := range lengthMap {
			_, resultVal := frequencies[timestamp]
			if resultVal && frequencies[timestamp] < lengthMap[timestamp] { //if there exists a value at a given time stamp in the final results
				frequencies[timestamp] = lengthMap[timestamp]
			} else { //the timestamp is not in the results so append
				frequencies[timestamp] = lengthMap[timestamp]
			}
		}

		// Update outliers for both BLT MAD and ShakeAlert outliers
		updateOutliers(AllResults, BLTMAD, bltIndexes, sortedTimestamps, sortedFrequencies)
		updateOutliers(AllResults, SHAKEALERT, shakeAlertIndexes, sortedTimestamps, sortedFrequencies)

	} else { //have not reached the max amount of points and can keep adding results
		//modify the frequencies for the final results
		for timestamp := range lengthMap {
			_, resultVal := frequencies[timestamp]
			if resultVal && frequencies[timestamp] < lengthMap[timestamp] { //if there exists a value at a given time stamp in the final results
				frequencies[timestamp] = lengthMap[timestamp]
			} else { //the timestamp is not in the results so append
				frequencies[timestamp] = lengthMap[timestamp]
			}
		}

		// Update outliers for both BLT MAD and ShakeAlert outliers
		updateOutliers(AllResults, BLTMAD, bltIndexes, sortedTimestamps, sortedFrequencies)
		updateOutliers(AllResults, SHAKEALERT, shakeAlertIndexes, sortedTimestamps, sortedFrequencies)
	}

}

func getListOfKeys(dataMap map[time.Time]float64) []time.Time {
	keys := []time.Time{}
	for timestamp := range dataMap {
		keys = append(keys, timestamp)
	}
	return keys
}

//use the olderst timestamp when removing the data points
func getSmallestTimestamp(timestamps []time.Time) *time.Time {
	if len(timestamps) == 0 {
		return nil
	}

	smallest := timestamps[0]
	for _, ts := range timestamps {
		if ts.Before(smallest) {
			smallest = ts
		}
	}
	return &smallest
}

func updateOutliers(currResult common.Result, algorithm int, indexes []int, sortedTimestamps []time.Time, sortedFrequencies []float64) {
	for _, index := range indexes {
		// Get outlier timestamp
		outlierTimestamp := sortedTimestamps[index]
		outlierFrequency := sortedFrequencies[index]

		// Go through current outliers and see if timestamp exists
		currOutlier, exists := currResult.AllOutliers[outlierTimestamp]

		// If the outlier is not already in the map, create a new outlier for it
		if !exists {
			currResult.AllOutliers[outlierTimestamp] = createOutlierStruct(outlierTimestamp, algorithm, outlierFrequency)
		} else { // If outlier not in map, need to possibly update it
			if currOutlier.Count != outlierFrequency {
				currOutlier.Count = outlierFrequency
			}
			// Update algorithm based on conditions
			switch {
			case currOutlier.Algorithm != BOTHALERTS:
				if algorithm == BLTMAD && currOutlier.Algorithm == SHAKEALERT ||
					algorithm == SHAKEALERT && currOutlier.Algorithm == BLTMAD {
					currOutlier.Algorithm = BOTHALERTS
				} else {
					currOutlier.Algorithm = algorithm
				}
			}

			// Store the updated outlier back in the map
			currResult.AllOutliers[outlierTimestamp] = currOutlier
		}
	}
}
func createOutlierStruct(timestamp time.Time, algorithm int, count float64) common.OutlierInfo {
	return common.OutlierInfo{Timestamp: timestamp,
		Algorithm: algorithm,
		Count:     count}
}

// repeated code in three of these functions; moved outside to make the code easier to read
func makeLengthMap(window common.Window) map[time.Time]float64 {
	bucketMap := window.BucketMap
	lengthMap := make(map[time.Time]float64)
	for timestamp, messages := range bucketMap {
		lengthMap[timestamp] = float64(len(messages))
	}
	return lengthMap
}

// Takes in map of time objects to frequencies and puts them into an ordered array of frequencies based on increasing timestamps
func GetSortedFrequencies(bucketMap map[time.Time]float64) ([]float64, []time.Time) {
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

	return sortedValues, timestamps
}
