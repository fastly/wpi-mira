package optimization

import "BGPAlert/blt_mad"

func FindTauForMinReqOutput(dataArray []float64, minReqOutData []float64) float64 {
	var tau = 0.0 //need to find the minimum output array that contains all the results 0.0

	for i := 0; i <= 1000; i++ { //check if I need to go above 100 but I am pretty sure that for most cases this
		currentOutput := blt_mad.BltMad(dataArray, tau)
		if len(blt_mad.FindDifferentValues(currentOutput, minReqOutData)) != 0 { //if there are no missing elements; continue
			//return previous tau
			return float64(i - 1)
		} else {
			tau += 1
		}
	}
	return tau
}
