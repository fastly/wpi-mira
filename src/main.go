package main

import (
	"math"
	"sort"
)

// constants
const (
	windowSize = 360 //360 min
	k          = 5   //# of neighbors
)

// calculate median of a slice of integers
func getMedian(slice []int) int {
	sort.Ints(slice)
	mid := len(slice) / 2
	return slice[mid]
}

// detect outliers using the density-based algorithm
// input: array of ints, where each index represents a time bin and value is the count for that time bin
// output: values of outliers aka counts
func findOutliers(counts []int) []int {
	var outliers []int

	for i := windowSize; i < len(counts); i++ {
		currentBin := counts[i]
		recentBins := counts[i-windowSize : i]

		// calculate the density-based radius R
		var counts []int
		for _, bin := range recentBins {
			counts = append(counts, bin)
		}
		R := int(math.Abs(float64(getMedian(counts) - currentBin)))

		// count the number of neighbors within radius R
		neighborCount := 0
		for _, bin := range recentBins {
			if int(math.Abs(float64(bin-currentBin))) < R {
				neighborCount++
			}
		}

		// if there are fewer than k neighbors, the current bin is an outlier
		if neighborCount < k {
			outliers = append(outliers, currentBin)
		}
	}

	return outliers
}
