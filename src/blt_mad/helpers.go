package blt_mad

import (
	"BGPAlert/common"
	"encoding/csv"
	"fmt"
	"math"
	"os"
	"reflect"
	"sort"
	"strconv"
	"time"
)

func RemoveZeros(data []float64) ([]float64, error) {
	var nonZeros []float64
	for _, value := range data {
		if value != 0.0 {
			nonZeros = append(nonZeros, value)
		}
	}

	if len(nonZeros) == 0 {
		fmt.Println("the slice provided was all zeros; working with an array [0.0] for analysis") //do not raise an error to avoid crashing the program
		return []float64{0.0}, nil
	}

	return nonZeros, nil
}

func WriteToCsv(filename string, messages []common.BGPMessage) error {
	var file *os.File
	var err error

	// Check if file exists
	if _, err = os.Stat(filename); os.IsNotExist(err) {
		// Create the file if it doesn't exist
		file, err = os.Create(filename)
		if err != nil {
			return err
		}
		defer file.Close()

		// Write CSV header
		writer := csv.NewWriter(file)
		defer writer.Flush()

		if err := writer.Write([]string{
			"Timestamp", "BGPMessageType", "PeerIP", "PeerASN", "Prefix",
		}); err != nil {
			return err
		}
	} else {
		// Open the file in append mode
		file, err = os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		defer file.Close()
	}

	// Create a CSV writer
	writer := csv.NewWriter(file)

	for _, message := range messages {
		// Format struct data into []string for CSV
		data := []string{
			message.Timestamp.Format(time.RFC3339),
			message.BGPMessageType,
			message.PeerIP.String(),
			strconv.FormatUint(uint64(message.PeerASN), 10),
			message.Prefix.String(),
		}

		if err := writer.Write(data); err != nil {
			return err
		}
	}

	// Flush and close the writer outside the loop
	writer.Flush()
	if err := writer.Error(); err != nil {
		return err
	}
	return nil
}

func FindMedian(data []float64) float64 {
	var med float64
	sortedData := sortData(data)
	if len(sortedData) == 0 {
		return math.SmallestNonzeroFloat64
	} else if len(sortedData)%2 == 1 {
		//if the length is odd return the number in the middle
		med = sortedData[len(sortedData)/2]
	} else {
		//if the length is even take the findMean of the middle two numbers
		left := sortedData[(len(sortedData)/2)-1]
		right := sortedData[(len(sortedData) / 2)]
		med = (left + right) / 2
	}
	return med
}

//if data needs to be sorted use this function to avoid changes in the original input
func sortData(data []float64) []float64 {
	// Create a copy of the original array to avoid modifying the input slice
	sortedData := make([]float64, len(data))
	copy(sortedData, data)

	// Use the sort package to sort the array
	sort.Float64s(sortedData)
	return sortedData
}

func findMean(data []float64) float64 {
	sum := 0.0
	if len(data) == 0 {
		return math.NaN() //invalid findMean in the context of BGP message counts
	} else {
		for _, num := range data {
			sum += num
		}
	}

	return sum / float64(len(data))
}

func containsValue(arr []float64, target float64) bool {
	for _, value := range arr {
		if value == target {
			return true
		}
	}
	return false
}

//https://medium.com/pragmatic-programmers/testing-floating-point-numbers-in-go-9872fe6de17f
func WithinTolerance(a float64, b float64, e float64) bool {
	if a == b {
		return true
	}
	d := math.Abs(a - b)
	if b == 0 {
		return d < e
	} else {
		return (d / math.Abs(b)) < e
	}
}

func WithinToleranceFloatSlice(a []float64, b []float64, e float64) bool {
	if reflect.DeepEqual(a, b) {
		return true
	} else {
		for i := 0; i < len(a); i++ {
			if !WithinTolerance(a[i], b[i], e) {
				return false
			}
		}
	}
	return true
}
