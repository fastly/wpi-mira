package shake_alert

import (
	"math"
	"sort"
)

// constants
const (
	windowSize = 391
	k          = 5 //# of neighbors
	//r          = 15000
)

// finds p-th percentile of data
func findPercentile(data []float64, p float64) float64 {
	if len(data) == 0 {
		return 0.0
	}

	// sort data in ascending order
	sort.Float64s(data)

	// calc rank
	n := float64(len(data))
	rank := (p / 100) * (n + 1)

	// Interpolate to find the requested percentile.
	lowerIdx := int(rank)
	upperIdx := lowerIdx + 1

	if lowerIdx == 0 {
		return data[0]
	}

	if upperIdx >= len(data) {
		return data[len(data)-1]
	}

	fraction := rank - float64(lowerIdx)
	lowerValue := data[lowerIdx-1]
	upperValue := data[upperIdx-1]

	return lowerValue + fraction*(upperValue-lowerValue)
}

// finds R radius value
// R = distance between 5th and 95th percentile
func findR(data []float64) float64 {

	percentile_5 := findPercentile(data, 5.0)
	percentile_95 := findPercentile(data, 95.0)

	R := percentile_95 - percentile_5

	return R
}

// detect outliers using the density-based algorithm
// input: array of ints, where each index represents a time bin and value is the count for that time bin
// output: values of outliers aka counts
func FindOutliers(data []float64) []float64 {
	//result of outliers
	var outliers []float64

	//iterate thru each time bin
	for i := 0; i < len(data); i++ {
		currentBin := data[i]
		//recentBins := bins[i-windowSize : i]

		//fmt.Println(currentBin)
		//fmt.Println(recentBins)

		// calculate the density-based radius R
		/*
			var counts []float64
			for _, bin := range recentBins {
				counts = append(counts, bin)
				//fmt.Println(counts)
			}
		*/
		//radius := int(math.Abs(float64(median(counts) - currentBin)))

		//fmt.Println(radius)

		// count the number of neighbors within radius R
		R := findR(data)
		neighborCount := 0
		for _, bin := range data {
			if math.Abs(float64(bin-currentBin)) < R {
				neighborCount++
			}
		}

		//fmt.Println(neighborCount)

		// if there are fewer than k neighbors, the current bin is an outlier
		if neighborCount < k {
			outliers = append(outliers, currentBin)
		}
	}

	return outliers
}
