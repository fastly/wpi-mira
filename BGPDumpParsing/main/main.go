package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"time"
)

func main() {
	startTime := time.Now()

	const bgpFilePath = "bgpfiles/updates.20160629.0600.bz2"
	count := 0

	defer os.Remove("tempbgpdump.txt")

	cmd := exec.Command("bgpdump", "-m", "-O", "tempbgpdump.txt", bgpFilePath)

	err := cmd.Run()
	if err != nil {
		fmt.Println("Error running bgpdump:", err)
	}

	file, err := os.Open("tempbgpdump.txt")
	if err != nil {
		fmt.Println("Error opening tempbgpdump.txt:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)
		count += 1
	}

	if scanner.Err() != nil {
		fmt.Println("Error reading tempbgpdump.txt:", scanner.Err())
	}

	elapsedTime := time.Since(startTime)

	fmt.Println("BGPDump Method: ")
	fmt.Println("Message Count: ", count)
	fmt.Println("Total execution time: ", elapsedTime)

}
