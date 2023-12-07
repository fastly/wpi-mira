package main

import (
	"BGPAlert/analyze"
	"BGPAlert/common"
	"BGPAlert/config"
	"BGPAlert/parse"
	"BGPAlert/process"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

func dataHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(analyze.AllResults)
	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}
	//fmt.Println(analyze.AllResults)
}

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

	err = configStruct.ValidateConfiguration()
	if err != nil {
		log.Fatalf("Failed to validate configuration: %v", err)
	}

	// WaitGroup for waiting on goroutines to finish
	var wg sync.WaitGroup

	// Channel for sending BGP messages between parsing and processing
	msgChannel := make(chan common.Message)
	freqInit := make(map[time.Time]float64)
	outlierInit := make(map[time.Time]common.OutlierInfo)
	analyze.AllResults.AllOutliers = outlierInit //init map of outliers
	analyze.AllResults.AllFreq = freqInit        //initialize map frequencies

	// Start the goroutines

	// Can change folder directory to any folder inside of src/static_data

	wg.Add(1)

	if configStruct.FileInputOption == "live" {
		//starts a go routine to parse live data if live option selected
		go func() {
			parse.ParseRisLiveData(msgChannel, configStruct)
			wg.Done()
		}()
	} else {
		//starts a go routine to parse static data if static option selected
		go func() {
			parse.ParseStaticFile(configStruct.URLStaticData, msgChannel)
			wg.Done()
		}()
	}

	wg.Add(1)
	go func() {
		process.ProcessBGPMessages(msgChannel, configStruct) //error handling done in the processBGPMessage
		wg.Done()
	}()

	//http start
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)
	http.HandleFunc("/data", dataHandler) //data handler to write data onto the local server

	log.Println("Server started on port 8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Server error:", err)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	elapsedTime := time.Since(startTime)
	fmt.Println("Elapsed Time: ", elapsedTime)

}
