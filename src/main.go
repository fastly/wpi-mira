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

	// Channel for sending windows from processing to analyzing
	windowChannel := make(chan common.Window)

	wg.Add(2)

	// Start the goroutines

	// Can change folder directory to any folder inside of src/staticdata
	go func() {
		parse.ParseStaticFile("bgptest1", msgChannel)
		//parse.ParseRisLiveData(msgChannel)
		wg.Done()
	}()

	go func() {
		process.ProcessBGPMessagesLive(msgChannel, windowChannel)
		wg.Done()
	}()

	/*
		go func() {
			analyze.AnalyzeBGPMessages(windowChannel)
			wg.Done()
		}()
	*/

	// Wait for all goroutines to finish
	wg.Wait()

	elapsedTime := time.Now().Sub(startTime)
	fmt.Println("Elapsed Time: ", elapsedTime)

}
