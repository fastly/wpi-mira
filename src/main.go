package main

import (
	"BGPAlert/common"
	"BGPAlert/parsing"
	"BGPAlert/processing"
	"fmt"
	"sync"
)

func main() {
	fmt.Println("In main")

	// WaitGroup for waiting on goroutines to finish
	var wg sync.WaitGroup

	// Channel for sending BGP messages between parsing and processing
	msgChannel := make(chan []common.BGPMessage)

	wg.Add(2)

	// Start the goroutines
	go func() {
		parsing.ParseStaticFile("bgptest1", msgChannel)
		wg.Done()
	}()

	go func() {
		processing.ProcessBGPMessages(msgChannel)
		wg.Done()
	}()

	// Wait for all goroutines to finish
	wg.Wait()
}
