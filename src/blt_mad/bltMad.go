package blt_mad

import (
	"math"
)

//blt mad

func BltMad(data []float64, tau float64) []float64 {
	var outliers []float64
	noZeroData := removeZeros(data) //add somethingt to do here if the length is 0 aka we dont have any non zero values
	if len(noZeroData) == 0 {
		print("The are no non-zero message counts")
	} else {
		//calculate blt formula
		//everything is based on the noZeroData since we are looking at spikes rather than lack of messages
		med := findMedian(noZeroData)
		m := mad(noZeroData)

		bltScore := math.Abs(med - tau*m)
		for _, value := range noZeroData {
			if value > bltScore {
				outliers = append(outliers, value)
			}
		}
	}

	return outliers
}
