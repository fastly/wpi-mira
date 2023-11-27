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

	// Start the goroutines

	// Can change folder directory to any folder inside of src/static_data
	wg.Add(1)

	go func() {
		parse.ParseStaticFile("bgpTest1", msgChannel)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		process.ProcessBGPMessages(msgChannel, configStruct)
		wg.Done()
	}()

	// Wait for all goroutines to finish
	wg.Wait()

	elapsedTime := time.Since(startTime)
	fmt.Println("Elapsed Time: ", elapsedTime)

}
