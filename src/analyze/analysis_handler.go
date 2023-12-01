package analyze

import (
	"BGPAlert/blt_mad"
	"BGPAlert/common"
	"BGPAlert/config"
	"BGPAlert/shake_alert"
	"fmt"
	"os"
	"sort"
	"time"
)

var finalResult common.Result //global so that it can be added onto by parser and seen by dataHandler
var maxPoints = 20            //the max number of points to be displayed on the graph; make sure to divide this by the number of points in each window or specify that maxPoints is the number of buckets that will be processed
var optParam = 10.0
// Takes in a Window, parses object into frequency counts, and then calls specified analysis functions
//code to write the frequencies; the outliers; and the minReqs to files
func AnalyzeBGPMessages(window common.Window, config *config.Configuration) common.Result {
	lengthMap := makeLengthMap(window)
	frequencies := finalResult.AllFreq
	// Turn map into sorted array of frequencies by timestamp
	sortedFrequencies := GetSortedFrequencies(lengthMap) //fix these duplicates with time stamps

	//modify the frequencies for the final results
	for timestamp, _ := range lengthMap {
		_, resultVal := frequencies[timestamp]
		if resultVal { //if there exists a value at a given time stamp in the final results
			if frequencies[timestamp] < lengthMap[timestamp] { //check if the existing value is smaller than the incoming
				frequencies[timestamp] = lengthMap[timestamp]
			} else {
				continue//if smaller append onto map else skip
			}

		} else { //the timestamp is not in the results so append
			frequencies[timestamp] = lengthMap[timestamp]
		}
	}

	//modify the outliers for the final result; check the outliers for the incoming window and remove duplicates/update

	//cap the length of frequencies at maxNumWindows *  windowSize

	fmt.Printf("Sorted Array of Frequencies: \n%+v\n", sortedFrequencies)
	fmt.Printf("BLT MAD Outliers: \n%+v\n", bltOutliers)
	fmt.Printf("ShakeAlert Outliers: \n%+v\n", shakeAlertOutliers)

	//put all the results into the Result struct and pass write it out to a json
	//conditionally update all the parameters in the result struct to avoid repeats and account for missed message counts at a given timestamp

	//what is the best way to output this?
	m := common.OutlierMessages{
		MADOutlierMessages: bltOutlierMessages,
		ShakeAlertMessages: shakeAlertOutlierMessages}

	blt_mad.StoreResultIntoJson(r, "static_data/recentFullResult.json") //storing the most recent result
	blt_mad.WriteCSVFile(m, "outlierMessages.csv")
	//make sure that we do not get more than a threshold number of points
	if len(AllResults)+1 > maxPoints {
		AllResults = append(AllResults[1:], r) //append all the elements except for the first one
	} else {
		AllResults = append(AllResults, r)
	}
	return r
}

//for a given window check what the outliers are and record them into a list of structs
func createOutliers (window common.Window) []common.OutlierInfo{
	lengthMap := makeLengthMap(window)
	sortedFrequencies := GetSortedFrequencies(lengthMap)

	//check if each individual entry in the length map is an outlier and create an outlier struct if needed
	for timestamp, _ := range lengthMap {
		//blt outlier
		if blt_mad.IsAnOutlierBLT(sortedFrequencies, optParam, lengthMap[timestamp])
		//madOutlier
	}
}

func createOutlierStruct (timestamp time.Time, algorithm int, count int) common.OutlierInfo{
	o := common.OutlierInfo{
		Timestamp: timestamp,
		Algorithm: algorithm,
		Count: count,
	}
	return o
}

/*//changed bltMad inputs to get timestamps and the outliers at the same time
//this is also producing duplicates
func BltMadWindow(window common.Window, tau float64) ([]float64, []time.Time, [][]common.BGPMessage) {
	var outliers []float64
	var times []time.Time
	var messages [][]common.BGPMessage //array of arrays of messages for a given bucket map
	bucketMap := window.BucketMap

	lengthMap := makeLengthMap(window)
	data := GetSortedFrequencies(lengthMap)

	for timestamp, _ := range lengthMap {
		if blt_mad.IsAnOutlierBLT(data, tau, lengthMap[timestamp]) {
			outliers = append(outliers, lengthMap[timestamp])
			times = append(times, timestamp)
			messages = append(messages, bucketMap[timestamp])
		}
	}
	return outliers, times, messages
}*/

//repeated code in three of these functions; moved outside to make the code easier to read
func makeLengthMap(window common.Window) map[time.Time]int {
	bucketMap := window.BucketMap
	lengthMap := make(map[time.Time]int)
	for timestamp, messages := range bucketMap {
		lengthMap[timestamp] = int(len(messages))
	}
	return lengthMap
}

/*//changed shakeAlert inputs to get timestamps and the outliers at the same time
func ShakeAlertWindow(window common.Window) ([]float64, []time.Time, [][]common.BGPMessage) {
	var outliers []float64
	var times []time.Time
	var messages [][]common.BGPMessage //array of arrays of messages for a given bucket map
	bucketMap := window.BucketMap

	//the frequencies needed to check if something is an outlier
	lengthMap := makeLengthMap(window)
	data := GetSortedFrequencies(lengthMap)

	for timestamp, _ := range lengthMap {
		if shake_alert.IsAnOutlierShakeAlert(data, lengthMap[timestamp]) {
			outliers = append(outliers, lengthMap[timestamp])
			times = append(times, timestamp)
			messages = append(messages, bucketMap[timestamp])
		}
	}
	return outliers, times, messages
}*/

// Takes in map of time objects to frequencies and puts them into an ordered array of frequencies based on increasing timestamps
func GetSortedFrequencies(bucketMap map[time.Time]int) []int {
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
	sortedValues := make([]int, len(timestamps))
	for i, timestamp := range timestamps {
		sortedValues[i] = bucketMap[timestamp]
	}

	return sortedValues
}
