package optimization

import (
	"gonum.org/v1/gonum/stat"
)

//use the code below to get the data for training the  model and getting the outputs
/*mads, _ := blt_mad.TxtIntoArrayFloat64("/home/taya/Fastly-MQP23/src/static_data/madsFound.txt")
medians, _ := blt_mad.TxtIntoArrayFloat64("/home/taya/Fastly-MQP23/src/static_data/mediansFound.txt")
xData := blt_mad.ArrayDivision(mads, medians)
yData, _ := blt_mad.TxtIntoArrayFloat64("/home/taya/Fastly-MQP23/src/static_data/tausFound.txt")
//normalize the data to have better output paramenters
yNorm := blt_mad.Normalize(yData)
xNorm := blt_mad.Normalize(xData)*/

//use the code below if the data needs to be split up into testing and training data
/*len80Percent := float64(len(xData)) * 0.8
xTrain := xData[0:int(len80Percent)]
xTest := xData[int(len80Percent)-1 : len(xData)]
yTrain := yData[0:int(len80Percent)]
yTest := yData[int(len80Percent)-1 : len(yData)]*/

func LinearRegressionModel(x []float64, y []float64) (float64, float64) {
	b, a := stat.LinearRegression(x, y, nil, false)
	return a, b

}

func Predict(x []float64, slope float64, intercept float64) []float64 {
	predictions := []float64{}
	for i := 0; i < len(x); i++ {
		predictions[i] = slope*x[i] + intercept
	}
	return predictions
}
