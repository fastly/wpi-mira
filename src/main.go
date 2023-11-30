package main

import (
	"BGPAlert/common"
	"BGPAlert/config"
	"BGPAlert/parse"
	"BGPAlert/process"
	"flag"
	"fmt"
	"log"
	"sync"
	"time"
)

func main() {

	startTime := time.Now()

	// define a flag for config file
	configFile := flag.String("config", "default-config.json", "Path to configuration file")

	//parse command line arguments
	flag.Parse()

	//indicate which config is being used
	if *configFile == "default-config.json" { //default
		fmt.Printf("No config file specified. Using default config file: %s\n", *configFile)
	} else { //user input config file
		fmt.Printf("Using config file: %s\n", *configFile)
	}

	configStruct, err := config.LoadConfig(*configFile)
	if err != nil {
		log.Fatal("Error loading configuration:", err)
	}
	config.ValidateConfiguration(configStruct)

	// WaitGroup for waiting on goroutines to finish
	var wg sync.WaitGroup

	// Channel for sending BGP messages between parsing and processing
	msgChannel := make(chan common.Message)

	// Start the goroutines

	// Can change folder directory to any folder inside of src/static_data
	wg.Add(1)

	if configStruct.FileInputOption == "live" {
		go func() {
			parse.ParseRisLiveData(msgChannel, configStruct)
			wg.Done()
		}()
	} else {
		go func() {
			parse.ParseStaticFile(configStruct.URLStaticData, msgChannel)
			wg.Done()
		}()
	}

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
