package main

import (
	"testing"
)

func TestParseBGPFile(t *testing.T) {
	filePath := "testdata/bgptest1/updates.20211204.1845.bz2"
	
	bgpMessages, err := parseBGPFile(filePath)

	if len(bgpMessages) < 1 || err != nil {
		t.Fatalf("ParseBGPFile function failed")
	}

}

func TestBuildFrequencyMap(t *testing.T) {
	filePath := "testdata/bgptest1/updates.20211204.1845.bz2"
	
	bgpMessages, _ := parseBGPFile(filePath)

	freqMap := buildFrequencyMap(bgpMessages)

	totalFrequencies := 0 
	for _, value := range freqMap {
		totalFrequencies += value
	}

	if totalFrequencies != len(bgpMessages) || len(freqMap) < 1 {
		t.Fatalf("BuildFrequencyMap function failed")
	}

}


