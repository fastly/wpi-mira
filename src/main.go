package main

import (
	"BGPAlert/analyze"
	"BGPAlert/common"
	"BGPAlert/config"
	"BGPAlert/parse"
	"BGPAlert/process"
	"flag"
	"fmt"
	"log"
	"sync"
)

func main() {

	// define a flag for config file
	configFile := flag.String("config", "config.json", "Path to configuration file")

	//parse command line arguments
	flag.Parse()

	//indicate which config is being used
	if *configFile == "config.json" { //default
		fmt.Printf("No config file specified. Using default config file: %s\n", *configFile)
	} else { //user input config file
		fmt.Printf("Using config file: %s\n", *configFile)
	}

	configStruct, err := config.LoadConfig(*configFile)
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

	wg.Add(3)

	// Start the goroutines

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
