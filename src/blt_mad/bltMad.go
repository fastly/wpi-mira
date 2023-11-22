package blt_mad

import (
	"math"
)

//Mad
func Mad(data []float64) float64 {
	//implement MAD given the list of the number of messages in each given bucket
	sum := 0.0
	if len(data) == 0 {
		return math.SmallestNonzeroFloat64 //invalid Mad in terms of BGP message counts
	} else {
		for _, value := range data {
			sum += math.Abs(value - findMean(data))
		}
	}
	result := sum / float64(len(data))
	return result
}

//blt Mad

func BltMad(data []float64, tau float64) []float64 {
	var outliers []float64
	noZeroData := removeZeros(data)
	if len(noZeroData) == 0 { //return empty array of outliers if the data array is empty
		print("The are no non-zero message counts")
	} else {
		//calculate blt formula
		//everything is based on the noZeroData since we are looking at spikes rather than lack of messages
		med := FindMedian(noZeroData)
		m := Mad(noZeroData)

		bltScore := math.Abs(med - tau*m)
		for _, value := range noZeroData {
			if value > bltScore {
				outliers = append(outliers, value)
			}
		}
	}

	return outliers
}
