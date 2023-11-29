package main

import (
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

var allResults []common.Result //global so that it can be added onto by parser and seen by dataHandler

func dataHandler(w http.ResponseWriter, r *http.Request) {
	/*//encode all the values into the json file
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(allResults[0])*/
	d := []common.Result{{
		Frequencies:          []float64{1, 4, 111, 135, 186},
		MADOutliers:          []float64{5, 8},
		MADTimestamps:        []time.Time{},
		ShakeAlertOutliers:   []float64{2, 10},
		ShakeAlertTimestamps: []time.Time{}},
		{
			Frequencies:          []float64{100, 120, 130},
			MADOutliers:          []float64{5, 8},
			MADTimestamps:        []time.Time{},
			ShakeAlertOutliers:   []float64{2, 10},
			ShakeAlertTimestamps: []time.Time{}},
	}
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(d)
	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}
}

func main() {

	startTime := time.Now()

	// define a flag for config file
	configFile := flag.String("config", "default-config.json", "Path to configuration file")

	//parse command line arguments
	flag.Parse()

	//indicate which config is being used
	if *configFile == "default-default-config.json" { //default
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
	msgChannel := make(chan common.Message)

	// Start the goroutines

	// Can change folder directory to any folder inside of src/static_data

	wg.Add(1)

	go func() {
		parse.ParseRisLiveData(msgChannel, configStruct) //this thing returns a list of all the results
		//parse.ParseStaticFile("bgpTest1", msgChannel)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		process.ProcessBGPMessages(msgChannel, configStruct) //error handling done in the processBGPMessage
		//allResults = append(allResults, res)
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
