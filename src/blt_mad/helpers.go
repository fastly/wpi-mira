package blt_mad

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"reflect"
	"sort"
	"strconv"
)

func removeZeros(data []float64) []float64 {
	var nonZeros []float64
	for _, value := range data {
		if value != 0.0 {
			nonZeros = append(nonZeros, value)
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

func findMin(data []float64) float64 {
	if len(data) == 0 {
		return math.SmallestNonzeroFloat64
	}
	min := data[0]
	for _, value := range data {
		if value < min {
			min = value
		}
	}
	return min
}

func findMax(data []float64) float64 {
	if len(data) == 0 {
		return math.MaxFloat64
	}
	max := data[0]
	for _, value := range data {
		if value > max {
			max = value
		}
	}
	return max
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

// equalSlices checks if two slices are equal in content.
func equalSlices(a, b []float64) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func containsValue(arr []float64, target float64) bool {
	for _, value := range arr {
		if value == target {
			return true
		}
	}
	return false
}

//check if one array contains the elements of the other
func containAllElements(mainArr []float64, subArr []float64) []float64 {
	//use in BGP testing to check if the result contains the minimum outliers of interest given the parameters
	var elementsMissing []float64
	for _, subValue := range subArr {
		if !containsValue(mainArr, subValue) {
			elementsMissing = append(elementsMissing, subValue)
		}
	}
	return elementsMissing
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

func GetValuesLargerThanPercentile(numbers []float64, percentile float64) []float64 {
	valueAtPercentile := calculatePercentile(numbers, percentile)

	// Get values larger than percentile% of the data
	var largerValues []float64
	for _, num := range numbers {
		if num > valueAtPercentile {
			largerValues = append(largerValues, num)
		}
	}

	return largerValues
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

/*func AppendFloat64ArrayToTxt(fileName string, floatArray []float64) error {
	// Open the file in append mode, create if it doesn't exist
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil
	}
	defer file.Close()

	// Iterate through the float array and write each element as a string to the file
	for _, num := range floatArray {
		// Convert float to string and write to file with newline separator
		if _, err := fmt.Fprintf(file, "%.6f\n", num); err != nil {
			fmt.Println("Error writing to file:", err)
			return nil
		}
	}

	fmt.Println("Float array appended to file successfully.")
	return nil
}*/

func ArrayFromJson(filename string) {

}

func TxtIntoArrayFloat64(inputFile string) ([]float64, error) {
	var floats []float64

	// Open the file
	file, err := os.Open(inputFile)
	defer file.Close()
	if err != nil {
		return nil, err
	}

	//scanner
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		value, err := strconv.ParseFloat(line, 64)
		if err != nil {
			continue //skipping the value and continue in the same line
		}
		floats = append(floats, value)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return floats, nil
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
