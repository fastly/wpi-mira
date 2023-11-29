package analyze

import (
	"BGPAlert/blt_mad"
	"BGPAlert/common"
	"BGPAlert/shake_alert"
	"fmt"
	"sort"
	"time"
)

var AllResults []common.Result //global so that it can be added onto by parser and seen by dataHandler

// Takes in a Window, parses object into frequency counts, and then calls specified analysis functions
//code to write the frequencies; the outliers; and the minReqs to files
func AnalyzeBGPMessages(window common.Window) common.Result {
	bucketMap := window.BucketMap

	// Convert BucketMap to a map of timestamp to length of messages
	lengthMap := make(map[time.Time]float64)

	for timestamp, messages := range bucketMap {
		lengthMap[timestamp] = float64(len(messages))
	}

	// Turn map into sorted array of frequencies by timestamp
	sortedFrequencies := GetSortedFrequencies(lengthMap)

	//the file names will contain all the timestamps for a given folder that was processed
	bltOutliers, bltOutlierTimes := BltMadWindow(window, 5) //add optimization to here
	fmt.Printf("Sorted Array of Frequencies: \n%+v\n", sortedFrequencies)
	fmt.Printf("BLT MAD Outliers: \n%+v\n", bltOutliers)
	fmt.Printf("ShakeAlert Outliers: \n%+v\n", shake_alert.FindOutliers(sortedFrequencies))

	//put all the results into the Result struct and pass write it out to a json
	r := common.Result{
		Frequencies:          sortedFrequencies,
		MADOutliers:          bltOutliers,
		MADTimestamps:        bltOutlierTimes,
		ShakeAlertOutliers:   shake_alert.FindOutliers(sortedFrequencies),
		ShakeAlertTimestamps: make([]time.Time, 0),
	}
	blt_mad.StoreResultIntoJson(r, "static_data/result.json")
	maxPoints := 100
	//make sure that we do not get more than a threshold number of points
	if len(AllResults)+1 > maxPoints {
		AllResults = append(AllResults[1:], r) //append all the elements except for the first one
	} else {
		AllResults = append(AllResults, r)
	}
	fmt.Println("--------------------------------------AllResult------------------------------------")
	fmt.Println(AllResults)
	fmt.Println("--------------------------------------AllResult------------------------------------")

	//get min reqArray for the 97th percentile
	//minReqArray := blt_mad.GetValuesLargerThanPercentile(sortedFrequencies, 97)
	//minReqOutFileName := fmt.Sprintf("static_data/minReq/minOutFile%s.txt", timeStampsFull)
	return r
}

func BltMadWindow(window common.Window, tau float64) ([]float64, []time.Time) {
	var outliers []float64
	var times []time.Time

	bucketMap := window.BucketMap

	lengthMap := make(map[time.Time]float64)
	for timestamp, messages := range bucketMap {
		lengthMap[timestamp] = float64(len(messages))
	}

	//the frequencies needed to check if something is an outlier
	data := GetSortedFrequencies(lengthMap)

	for timestamp, _ := range lengthMap {
		if blt_mad.IsAnOutlier(data, tau, lengthMap[timestamp]) {
			outliers = append(outliers, lengthMap[timestamp])
			times = append(times, timestamp)
		}
	}
	return outliers, times
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
