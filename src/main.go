package main

import (
	"BGPAlert/analyze"
	"BGPAlert/common"
	"BGPAlert/config"
	"BGPAlert/parse"
	"BGPAlert/process"
	"log"
	"sync"
)

func main() {
	configStruct, err := config.LoadConfig("config.json")
	if err != nil {
		log.Fatal("Error loading configuration:", err)
	}
	config.ValidDateConfiguration(configStruct)

	// WaitGroup for waiting on goroutines to finish
	var wg sync.WaitGroup

	// Channel for sending BGP messages between parsing and processing
	msgChannel := make(chan []common.BGPMessage)

	// Channel for sending windows from processing to analyzing
	windowChannel := make(chan []common.Window)

	wg.Add(3)

	// Start the goroutines

	// Can change folder directory to any folder inside of src/staticdata
	go func() {
		parse.ParseStaticFile("bgptest1", msgChannel)
		wg.Done()
	}()

	go func() {
		process.ProcessBGPMessages(msgChannel, windowChannel)
		wg.Done()
	}()

	go func() {
		analyze.AnalyzeBGPMessages(windowChannel)
		wg.Done()
	}()

	// Wait for all goroutines to finish
	wg.Wait()

}
