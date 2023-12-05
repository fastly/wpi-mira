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
	FileInputOption string            `json:"dataOption"`
	StaticFile      string            `json:"staticFilePath"`  //path to specific static file
	URLStaticData   string            `json:"staticFilesLink"` //link to routeviews bz2 folder
	MadParameters   int               `json:"madParameters"`
	MaxBuckets      int               `json:"maxBuckets"`
	WindowSize      int               `json:"windowSize"`
	Subscriptions   []SubscriptionMsg `json:"subscriptions"`
}

// load config json file
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

// checks if each value in the config file is valid input
func ValidateConfiguration(config *Configuration) error {

	//convert all strings to lower case to ignore any capitalizations
	fileInputL := strings.ToLower(config.FileInputOption)

	//check that the fileInputOption is either live or static
	if fileInputL != "live" && fileInputL != "static" {
		return errors.New("DataOption in config.json must be either 'live' or 'static'")
	}

	//checks for valid input corresponding to choice of live vs static data analysis
	if fileInputL == "live" {
		//require at least 1 subscription
		if len(config.Subscriptions) == 0 {
			return errors.New("choosing live data input stream requires to input at least one subscription")
			return errors.New("choosing live data input stream requires to input at least one subscription")
		}
	} else if fileInputL == "static" {
		//require valid file path
		_, err := os.Stat(config.StaticFile)
		if os.IsNotExist(err) {
			return errors.New("Invalid pathway to static file: " + err.Error())
		}
	}

	//require mad parameter
	//set mad param to default value if no val set in config
	if config.MadParameters <= 0 {
		config.MadParameters = 10
		fmt.Println("Invalid mad parameter. The default mad parameter set to 10")
	}

	//require maxBuckets
	//set maxBuckets to default value if not set in config
	if config.MaxBuckets <= 0 {
		config.MaxBuckets = 20
		fmt.Println("Invalid maxBuckets value. The default maxBuckets set to 20")
	}

	//added to check if no window size was put in
	if config.WindowSize <= 0 {
		config.WindowSize = 360
		fmt.Println("Invalid window size. The default window size was set to 360")
	}

	//check if maxBuckets * windowSize is not greater than 1k - could lead to messy graph
	//if so, give a warning
	if (config.MaxBuckets * config.WindowSize) >= 1000 {
		fmt.Println("WARNING: Plot may crash")
	}

	//no errors found - valid config file
	fmt.Println("Configuration successful")
	return nil

	return nil
}
