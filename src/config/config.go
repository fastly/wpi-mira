package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type SubscriptionMsg struct {
	Host   string `json:"host,omitempty"` //aka collector
	Peer   string `json:"peer,omitempty"`
	Path   string `json:"path,omitempty"` //aka ASN
	Prefix string `json:"prefix,omitempty"`
}

type Configuration struct {
	//cast onto the needed type when processing in algos
	FileInputOption           string            `json:"dataOption"`
	StaticFile                string            `json:"staticFilePath"`
	URLStaticData             string            `json:"staticFilesLink"`
	OutlierDetectionAlgorithm string            `json:"outlierDetectionAlgorithm"`
	MadParameters             string            `json:"madParameters"`
	Subscriptions             []SubscriptionMsg `json:"subscriptions"`
	WindowSize                string            `json:"windowSize"`
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

func ValidateConfiguration(config *Configuration) {

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
			fmt.Println("Choosing live data input stream requires to input at least one subscription")
		}
	} else if fileInputL == "static" {
		//require valid file path
		_, err := os.Stat(config.StaticFile)
		if os.IsNotExist(err) {
			fmt.Println("Please enter a valid pathway to the static file")
		}
	} else if fileInputL != "live" && fileInputL != "static" {
		fmt.Println("Please enter either live or static as a dataOption in config.json")
	} else if outlierL == "mad" {
		//require mad parameter
	} else if outlierL != "mad" && outlierL != "shakealert" {
		fmt.Println("Please enter either mad or shakeAlert as input for outlierDetectionAlgorithm in config.json")
	}
	fmt.Println("Configuration successful")

}
