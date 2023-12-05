package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type SubscriptionMsg struct {
	Host   string `json:"host,omitempty"` //aka collector
	Peer   string `json:"peer,omitempty"`
	Asn    int    `json:"asn,omitempty"` //aka ASN
	Prefix string `json:"prefix,omitempty"`
}

type Configuration struct {
	FileInputOption     string            `json:"dataOption"`
	Algorithm           string            `json:"anomalyDetectionAlgo"`
	StaticFile          string            `json:"staticFilePath"`  //path to specific static file
	URLStaticData       string            `json:"staticFilesLink"` //link to routeviews bz2 folder
	MadParameters       int               `json:"madParameters"`
	ShakeAlertParameter int               `json:"shakeAlertParameter"`
	MaxBuckets          int               `json:"maxBuckets"`
	WindowSize          int               `json:"windowSize"`
	Subscriptions       []SubscriptionMsg `json:"subscriptions"`
}

func LoadConfig(filename string) (*Configuration, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Configuration

	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

func ValidateConfiguration(config *Configuration) error {

	//check that the fileInputOption is either live or static
	//convert all strings to lower case to ignore any capitalizations
	fileInputL := strings.ToLower(config.FileInputOption)
	outlierL := strings.ToLower(config.OutlierDetectionAlgorithm)

	//added to check if no window size was put in
	if config.WindowSize == "" {
		config.WindowSize = "360"
		fmt.Println("No window size was passed in. The default window size was set to 360")
	}

	if fileInputL == "live" {
		//require at least 1 subscription
		if len(config.Subscriptions) == 0 {
			return errors.New("choosing live data input stream requires to input at least one subscription")
		}
	} else if fileInputL == "static" {
		//require valid file path
		_, err := os.Stat(config.StaticFile)
		if os.IsNotExist(err) {
			fmt.Println("Please enter a valid pathway to the static file")
		}
	} else if fileInputL != "live" && fileInputL != "static" {
		fmt.Println("Please enter either live or static as a dataOption in default-config.json")
	} else if outlierL == "mad" {
		//require mad parameter
	} else if outlierL != "mad" && outlierL != "shakealert" {
		fmt.Println("Please enter either mad or shakeAlert as input for outlierDetectionAlgorithm in default-config.json")
	}
	fmt.Println("Configuration successful")
	return nil
}
