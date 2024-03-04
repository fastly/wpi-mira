package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"mira/analyze"
	"mira/common"
	"mira/config"
	"mira/parse"
	"mira/process"
	"net/http"
	"strings"
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

	err = configStruct.ValidateConfiguration()
	if err != nil {
		log.Fatalf("Failed to validate configuration: %v", err)
	}

	// WaitGroup for waiting on goroutines to finish
	var wg sync.WaitGroup

	// Channel for sending BGP messages between parsing and processing
	msgChannel := make(chan common.Message)
	/*freqInit := make(map[time.Time]float64)
	outlierInit := make(map[time.Time]common.OutlierInfo)*/
	analyze.ResultMap = make(map[string]*common.Result)
	/*analyze.AllResults.AllOutliers = outlierInit //init map of outliers
	analyze.AllResults.AllFreq = freqInit        //initialize map frequencies
	*/
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
	http.HandleFunc("/data", getData) //data handler to write data onto the local server
	http.HandleFunc("/subscriptions", getSubscriptions)
	http.HandleFunc("/frequencies", getFrequenciesFromSubscription)
	http.HandleFunc("/outliers", getOutliersFromSubscription)

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

func getData(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("Data request")
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(analyze.ResultMap)
	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}
}

func getSubscriptions(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Subscriptions request")
	w.Header().Set("Content-Type", "application/json")

	var subscriptions []string
	for sub, _ := range analyze.ResultMap {
		subscriptions = append(subscriptions, sub)
	}

	err := json.NewEncoder(w).Encode(subscriptions)
	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}
}

func getFrequenciesFromSubscription(w http.ResponseWriter, r *http.Request) {
	subscription := r.URL.Query().Get("subscription")

	// Check if the required parameter is missing
	if subscription == "" {
		http.Error(w, "Missing required parameter 'subscription'", http.StatusBadRequest)
		return
	}

	// Removes backslashes from subscription string
	cleanedSubscription := strings.Replace(subscription, "\\", "", -1)
	fmt.Println("Frequency request for ", cleanedSubscription)

	result, exists := analyze.ResultMap[cleanedSubscription]

	var frequencies map[time.Time]float64

	if exists {
		frequencies = result.AllFreq
	} else {
		// Subscription not found, use an empty map
		frequencies = map[time.Time]float64{}
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(frequencies)
	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}
}

func getOutliersFromSubscription(w http.ResponseWriter, r *http.Request) {
	subscription := r.URL.Query().Get("subscription")

	// Check if the required parameter is missing
	if subscription == "" {
		http.Error(w, "Missing required parameter 'subscription'", http.StatusBadRequest)
		return
	}

	// Removes backslashes from subscription string
	cleanedSubscription := strings.Replace(subscription, "\\", "", -1)
	fmt.Println("Outlier request for ", cleanedSubscription)

	result, exists := analyze.ResultMap[cleanedSubscription]

	var outliers map[time.Time]common.OutlierInfo

	if exists {
		outliers = result.AllOutliers
	} else {
		// Subscription not found, use an empty map
		outliers = map[time.Time]common.OutlierInfo{}
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(outliers)
	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}
}
