package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type Configuration struct {
	//cast onto the needed type when processing in algos
	FileInputOption           string `json:"dataOption"`
	StaticFile                string `json:"staticFilePath"`
	OutlierDetectionAlgorithm string `json:"outlierDetectionAlgorithm"`
	MadParameters             string `json:"madParameters"`
	Prefix                    string `json:"prefix"` // can input a list of string with values seperated by a comma
	Asn                       string `json:"asn"`
	PeerIP                    string `json:"peerIP"`
	Connector                 string `json:"connector"`
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

//to parse input of prefix, asn, peer ip, and connector if needed
func parseByComma(data string) []string {
	return strings.Split(data, ",")
}

func ValidDateConfiguration(config *Configuration) {
	//check that the fileInputOption is either live or static
	//convert all strings to lower case to ignore any capitalizations
	fileInputL := strings.ToLower(config.FileInputOption)
	outlierL := strings.ToLower(config.OutlierDetectionAlgorithm)

	if fileInputL == "live" {
		//require prefix and collector
		if len(config.Connector) == 0 || len(config.Prefix) == 0 {
			fmt.Println("Choosing live data input stream requires to input at least one value for the connector and at lease one value for the prefix")
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
