package blt_mad

import "math"

func mad(data []float64) float64 {
	//implement MAD given the list of the number of messages in each given bucket
	sum := 0.0
	if len(data) == 0 {
		return math.SmallestNonzeroFloat64 //invalid mad in terms of BGP message counts
	} else {
		for _, value := range data {
			sum += math.Abs(value - findMean(data))
		}
	}
	result := sum / float64(len(data))
	return result
}
