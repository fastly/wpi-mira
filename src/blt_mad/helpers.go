package blt_mad

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"reflect"
	"sort"
)

func RemoveZeros(data []int) []float64 {
	var nonZeros []float64
	for _, value := range data {
		if value != 0.0 {
			nonZeros = append(nonZeros, float64(value))
		}
	}
	return nonZeros
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

func calculatePercentile(numbers []float64, percentile float64) float64 {
	sortedList := sortData(numbers)
	index := int(percentile / 100 * float64(len(sortedList)-1)) //index corresponding to the percentile
	lower := sortedList[index]
	upper := sortedList[index+1]
	fractionalPart := percentile/100*float64(len(sortedList)-1) - float64(index)
	value := lower + (upper-lower)*fractionalPart
	return value
}

//used to simplify writing out the results into the output file
func StoreResultIntoJson(data interface{}, filename string) error {
	jsonData, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filename, jsonData, 0644)
	if err != nil {
		return err
	}
	fmt.Println(string(jsonData))
	return nil
}

//unmarshall json here
func WriteCSVFile(data interface{}, filename string) error {
	var file *os.File
	var err error

	if _, statErr := os.Stat(filename); os.IsNotExist(statErr) {
		// If the file does not exist, create it
		file, err = os.Create(filename)
		if err != nil {
			return err
		}
		defer file.Close()

		// Write the initial data to the file
		jsonData, err := json.MarshalIndent([]interface{}{data}, "", "  ")
		if err != nil {
			return err
		}

		if _, err := file.Write(jsonData); err != nil {
			return err
		}
	} else {
		// If the file exists, open it in append mode
		file, err = os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		if err != nil {
			return err
		}
		defer file.Close()

		// Read the existing data from the file
		existingData := []interface{}{}
		decoder := json.NewDecoder(file)
		if err := decoder.Decode(&existingData); err != nil {
			return err
		}

		// Append new data to existing data
		existingData = append(existingData, data)

		// Move file cursor to the beginning
		if _, err := file.Seek(0, 0); err != nil {
			return err
		}

		// Write the updated data to the file
		jsonData, err := json.MarshalIndent(existingData, "", "  ")
		if err != nil {
			return err
		}

		if _, err := file.Write(jsonData); err != nil {
			return err
		}
	}

	return nil
}

//check if one array contains the elements of the other
//possibly simplify this function in the future; i might need to use this specific implementation for comparing outputs of shakeAlert and MAD
func FindDifferentValues(mainArr []float64, subArr []float64) []float64 {
	//use in BGP testing to check if the result contains the minimum outliers of interest given the parameters
	var elementsMissing []float64
	for _, subValue := range subArr {
		if !containsValue(mainArr, subValue) {
			elementsMissing = append(elementsMissing, subValue)
		}
	}
	return elementsMissing
}
