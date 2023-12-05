package analyze

import (
	"BGPAlert/blt_mad"
	"BGPAlert/common"
	"BGPAlert/config"
	"BGPAlert/shake_alert"
	"fmt"
	"path/filepath"
	"sort"
	"time"
)

var recentResultsFile, _ = filepath.Abs("static_data/recentFullResult.json")
var outlierMessagesFile, _ = filepath.Abs("outlierMessages.csv")
var AllResults []common.Result //global so that it can be added onto by parser and seen by dataHandler
var maxPoints = 20             //the max number of points to be displayed on the graph;
// make sure to divide this by the number
//of points in each window or specify that maxPoints is the number of buckets that will be processed

// Takes in a Window, parses object into frequency counts, and then calls specified analysis functions
// code to write the frequencies; the outliers; and the minReqs to files
func AnalyzeBGPMessages(window common.Window, config *config.Configuration) common.Result {
	lengthMap := makeLengthMap(window)
	// Turn map into sorted array of frequencies by timestamp
	sortedFrequencies := GetSortedFrequencies(lengthMap) //fix these duplicates with time stamps

	//the file names will contain all the timestamps for a given folder that was processed
	bltOutliers, bltOutlierTimes, bltOutlierMessages := BltMadWindow(sortedFrequencies, window, 5) //add optimization to here
	shakeAlertOutliers, shakeAlertOutlierTime, shakeAlertOutlierMessages := ShakeAlertWindow(sortedFrequencies, window, config)
	fmt.Printf("Sorted Array of Frequencies: \n%+v\n", sortedFrequencies)
	fmt.Printf("BLT MAD Outliers: \n%+v\n", bltOutliers)
	fmt.Printf("ShakeAlert Outliers: \n%+v\n", shakeAlertOutliers)

	//put all the results into the Result struct and pass write it out to a json
	r := common.Result{
		WindowSize: config.WindowSize,

		Frequencies: sortedFrequencies, //make all of these maps and append onto them here; create them as global vars

		MADOutliers:   bltOutliers,
		MADTimestamps: bltOutlierTimes,

		ShakeAlertOutliers:   shakeAlertOutliers,
		ShakeAlertTimestamps: shakeAlertOutlierTime,
	}

	//what is the best way to output this?
	m := common.OutlierMessages{
		MADOutlierMessages: bltOutlierMessages,
		ShakeAlertMessages: shakeAlertOutlierMessages}

	blt_mad.StoreResultIntoJson(r, recentResultsFile) //storing the most recent result
	blt_mad.WriteCSVFile(m, outlierMessagesFile)
	//make sure that we do not get more than a threshold number of points
	if len(AllResults)+1 > maxPoints {
		AllResults = append(AllResults[1:], r) //append all the elements except for the first one
	} else {
		AllResults = append(AllResults, r)
	}
	return r
}

// changed bltMad inputs to get timestamps and the outliers at the same time
// this is also producing duplicates
func BltMadWindow(data []float64, window common.Window, tau float64) ([]float64, []time.Time, [][]common.BGPMessage) {
	var outliers []float64
	var times []time.Time
	var messages [][]common.BGPMessage //array of arrays of messages for a given bucket map
	bucketMap := window.BucketMap

	lengthMap := makeLengthMap(window)

	for timestamp, _ := range lengthMap {
		if blt_mad.IsAnOutlierBLT(data, tau, lengthMap[timestamp]) {
			outliers = append(outliers, lengthMap[timestamp])
			times = append(times, timestamp)
			messages = append(messages, bucketMap[timestamp])
		}
	}
	return outliers, times, messages
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

// changed shakeAlert inputs to get timestamps and the outliers at the same time
func ShakeAlertWindow(data []float64, window common.Window, config *config.Configuration) ([]float64, []time.Time, [][]common.BGPMessage) {
	var outliers []float64
	var times []time.Time
	var messages [][]common.BGPMessage //array of arrays of messages for a given bucket map
	bucketMap := window.BucketMap

	//the frequencies needed to check if something is an outlier
	lengthMap := makeLengthMap(window)

	for timestamp, _ := range lengthMap {
		if shake_alert.IsAnOutlierShakeAlert(data, lengthMap[timestamp], config.ShakeAlertParameter) {
			outliers = append(outliers, lengthMap[timestamp])
			times = append(times, timestamp)
			messages = append(messages, bucketMap[timestamp])
		}
	}
	return outliers, times, messages
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
