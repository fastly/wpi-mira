package main

import (
	"BGPAlert/common"
	"BGPAlert/config"
	"BGPAlert/parse"
	"BGPAlert/process"
	"fmt"
	"log"
	"sync"
	"time"
)

func main() {

	/*	mads, _ := blt_mad.TxtIntoArrayFloat64("/home/taya/Fastly-MQP23/src/static_data/madsFound.txt")
		medians, _ := blt_mad.TxtIntoArrayFloat64("/home/taya/Fastly-MQP23/src/static_data/mediansFound.txt")
		xData := blt_mad.ArrayDivision(mads, medians)
		yData, _ := blt_mad.TxtIntoArrayFloat64("/home/taya/Fastly-MQP23/src/static_data/tausFound.txt")
		//normalize the data to have better output paramenters
		yNorm := blt_mad.Normalize(yData)
		xNorm := blt_mad.Normalize(xData)
		len80Percent := float64(len(xNorm)) * 0.8
		xTrain := xNorm[0:int(len80Percent)]
		xTest := xData[int(len80Percent)-1 : len(xNorm)]
		yTrain := yNorm[0:int(len80Percent)]
		yTest := yNorm[int(len80Percent)-1 : len(yNorm)]

		intercept, slope := optimization.LinearRegressionModel(xTrain, yTrain)
		fmt.Println(intercept, slope)
		fmt.Println(len(xTest))
		fmt.Println(len(yTest))
		predictions := optimization.Predict(xTest, slope, intercept) //check how exactly this works
		fmt.Println(predictions)*/

	startTime := time.Now()

	configStruct, err := config.LoadConfig("config.json")
	if err != nil {
		log.Fatal("Error loading configuration:", err)
	}
	config.ValidDateConfiguration(configStruct)

	// WaitGroup for waiting on goroutines to finish
	var wg sync.WaitGroup

	// Channel for sending BGP messages between parsing and processing
	msgChannel := make(chan common.BGPMessage)

	wg.Add(2)

	// Start the goroutines

	// Can change folder directory to any folder inside of src/static_data
	go func() {
		parse.ParseStaticFile("bgptest1", msgChannel)
		//parse.ParseRisLiveData(msgChannel)
		wg.Done()
	}()

	go func() {
		process.ProcessBGPMessages(msgChannel)
		wg.Done()
	}()

	// Wait for all goroutines to finish
	wg.Wait()

	elapsedTime := time.Since(startTime)
	fmt.Println("Elapsed Time: ", elapsedTime)
}
