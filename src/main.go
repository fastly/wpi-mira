package main

import (
	"BGPAlert/analyzing"
	"BGPAlert/common"
	"BGPAlert/config"
	"BGPAlert/parsing"
	"BGPAlert/processing"
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
		parsing.ParseStaticFile("bgptest1", msgChannel)
		wg.Done()
	}()

	go func() {
		processing.ProcessBGPMessages(msgChannel, windowChannel)
		wg.Done()
	}()

	go func() {
		analyzing.AnalyzeBGPMessages(windowChannel)
		wg.Done()
	}()

	// Wait for all goroutines to finish
	wg.Wait()
}
