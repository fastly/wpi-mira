package blt_mad

import (
	"bufio"
	"fmt"
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

func SaveArrayToFile(fileName string, arr []float64) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, value := range arr {
		_, err := fmt.Fprintf(file, "%f\n", value)
		if err != nil {
			return err
		}
	}

	return nil
}

func TxtIntoArrayFloat64(inputFile string) ([]float64, error) {
	var floats []float64

	// Open the file
	file, err := os.Open(inputFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	//scanner
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		value, err := strconv.ParseFloat(line, 64)
		if err != nil {
			return nil, err
		}
		floats = append(floats, value)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return floats, nil
}

//check if one array contains the elements of the other
func ContainAllElements(mainArr []float64, subArr []float64) []float64 {
	//use in BGP testing to check if the result contains the minimum outliers of interest given the parameters
	var elementsMissing []float64
	for _, subValue := range subArr {
		if !containsValue(mainArr, subValue) {
			elementsMissing = append(elementsMissing, subValue)
		}
	}
	return elementsMissing
}
