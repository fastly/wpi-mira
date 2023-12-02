package analyze

import (
	"BGPAlert/blt_mad"
	"BGPAlert/common"
	"BGPAlert/config"
	"BGPAlert/shake_alert"
	"sort"
	"time"
)

//var OB = common.NewBuilder()  //builder for outlierInfo
var FinalResult common.Result //global so that it can be added onto by parser and seen by dataHandler
var maxPoints = 20            //the max number of points to be displayed on the graph; make sure to divide this by the number of points in each window or specify that maxPoints is the number of buckets that will be processed
var optParam = 10.0

// Takes in a Window, parses object into frequency counts, and then calls specified analysis functions
//code to write the frequencies; the outliers; and the minReqs to files
func AnalyzeBGPMessages(window common.Window, config *config.Configuration) {
	lengthMap := makeLengthMap(window)
	frequencies := FinalResult.AllFreq
	// Turn map into sorted array of frequencies by timestamp
	//sortedFrequencies := GetSortedFrequencies(lengthMap) //fix these duplicates with time stamps

	//modify the frequencies for the final results
	for timestamp, _ := range lengthMap {
		_, resultVal := frequencies[timestamp]
		if resultVal { //if there exists a value at a given time stamp in the final results
			if frequencies[timestamp] < lengthMap[timestamp] { //check if the existing value is smaller than the incoming
				frequencies[timestamp] = lengthMap[timestamp]
			} else {
				continue //if smaller append onto map else skip
			}
		} else { //the timestamp is not in the results so append
			frequencies[timestamp] = lengthMap[timestamp]
		}
	}

	//modify the outliers for the final result; check the outliers for the incoming window and remove duplicates/update
	windowOutliers := createOutliers(window) //a list of all the outliers in the individual window
	windowOutlierTimes := getListTimes(windowOutliers)
	resultOutlierTimes := getListTimes(FinalResult.AllOutliers)
	for i, val := range windowOutlierTimes {
		if !containsVal(resultOutlierTimes, val) {
			FinalResult.AllOutliers = append(FinalResult.AllOutliers, windowOutliers[i]) //make sure that window outliers and outlier times are the same here
		}
	}

	//cap the length of frequencies at maxNumWindows *  windowSize
}

func containsVal(times []time.Time, specificTime time.Time) bool {
	for _, val := range times {
		if val == specificTime { //check equal here
			return true
		}
	}
	return false
}

//there must be an easier way to do this
func getListTimes(outliers []common.OutlierInfo) []time.Time {
	times := []time.Time{}
	for _, val := range outliers {
		times = append(times, val.Timestamp)
	}
	return times
}

//for a given window check what the outliers are and record them into a list of structs
func createOutliers(window common.Window) []common.OutlierInfo {
	windowOutliers := []common.OutlierInfo{}
	lengthMap := makeLengthMap(window)
	sortedFrequencies := GetSortedFrequencies(lengthMap)

	//check if each individual entry in the length map is an outlier and create an outlier struct if needed
	for timestamp, _ := range lengthMap {
		//blt outlier
		if blt_mad.IsAnOutlierBLT(sortedFrequencies, optParam, lengthMap[timestamp]) {
			oMad := createOutlierStruct(timestamp, 0, lengthMap[timestamp])
			windowOutliers = append(windowOutliers, oMad)
		} else if shake_alert.IsAnOutlierShakeAlert(sortedFrequencies, lengthMap[timestamp]) { //shakeOutlier
			oShake := createOutlierStruct(timestamp, 1, lengthMap[timestamp])
			windowOutliers = append(windowOutliers, oShake)
		} else { //not an outlier
			continue
		}
	}
	return windowOutliers
}

func createOutlierStruct(timestamp time.Time, algorithm int, count float64) common.OutlierInfo {
	o := common.OutlierInfo{
		Timestamp: timestamp,
		Algorithm: algorithm,
		Count:     count,
	}

	return o
}

//repeated code in three of these functions; moved outside to make the code easier to read
func makeLengthMap(window common.Window) map[time.Time]float64 {
	bucketMap := window.BucketMap
	lengthMap := make(map[time.Time]float64)
	for timestamp, messages := range bucketMap {
		lengthMap[timestamp] = float64((len(messages)))
	}
	return lengthMap
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
