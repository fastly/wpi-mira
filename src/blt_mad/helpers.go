package blt_mad

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"reflect"
	"sort"
	"strconv"
	//"https://github.com/blend/go-sdk.git"
)

//use the mathutil to see if i can simplify all the functions

//convert a string of numbers into an array []float64 -> it actually works; all the functions can be called
//this does not actually work??????????????????
/*func convertToFloat64Array(str string, err error) []float64 {
	var numbers []float64
	numbersStr := strings.Fields(str)
	for i, _ := range numbersStr {
		f, _ := strconv.ParseFloat(numbersStr[i], 64)
		numbers = append(numbers, f)
	}
	return numbers
}*/

func removeZeros(data []float64) []float64 {
	var nonZeros []float64
	for _, value := range data {
		if value != 0.0 {
			nonZeros = append(nonZeros, value)
		}
	}
	return nonZeros
}

/*sort the array in ascending order and all the processing that will be done with an array will be done with it in ascending order*/
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
func withinTolerance(a float64, b float64, e float64) bool {
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

func withinToleranceFloatSlice(a []float64, b []float64, e float64) bool {
	if reflect.DeepEqual(a, b) {
		return true
	} else {
		for i := 0; i < len(a); i++ {
			if !withinTolerance(a[i], b[i], e) {
				return false
			}
		}
	}
	return true
}

func readTxtToString(filePath string) (string, error) {
	// Read the entire file into a byte slice
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	// Convert the byte slice to a string
	return string(content), nil
}

//get some percentile to get the minimum required output
func GetNumbersInPercentile(numbers []float64, percentile float64) []float64 {
	// Sort the array in ascending order
	sort.Float64s(numbers)

	// Calculate the index corresponding to the percentile
	index := int(percentile / 100 * float64(len(numbers)-1))

	// Interpolate the value at the calculated index
	lower := numbers[index]
	upper := numbers[index+1]

	// Interpolation formula: lower + (upper - lower) * fractional part
	fractionalPart := percentile/100*float64(len(numbers)-1) - float64(index)
	value := lower + (upper-lower)*fractionalPart

	// Get all numbers that fall into the 90th percentile
	var percentileNumbers []float64
	for _, num := range numbers {
		if num <= value {
			percentileNumbers = append(percentileNumbers, num)
		}
	}

	return percentileNumbers
}

func calculatePercentile(numbers []float64, percentile float64) float64 {
	// Sort the array in ascending order
	sort.Float64s(numbers)

	// Calculate the index corresponding to the percentile
	index := int(percentile / 100 * float64(len(numbers)-1))

	// Interpolate the value at the calculated index
	lower := numbers[index]
	upper := numbers[index+1]

	// Interpolation formula: lower + (upper - lower) * fractional part
	fractionalPart := percentile/100*float64(len(numbers)-1) - float64(index)
	value := lower + (upper-lower)*fractionalPart

	return value
}

func GetValuesLargerThanPercentile(numbers []float64, percentile float64) []float64 {
	valueAtPercentile := calculatePercentile(numbers, percentile)

	// Get values larger than 90% of the data
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

//running the go routines through one folder and saving the output into the outFile

func TxtIntoArrayFloat64(inputFile string) ([]float64, error) {
	var floats []float64

	// Open the file
	file, err := os.Open(inputFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Parse each line as a float64
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

func ArrayDivision(arr1, arr2 []float64) []float64 {
	result := make([]float64, len(arr1))
	if len(arr1) != len(arr2) {
		return nil
	} else {
		for i := 0; i < len(arr1); i++ {
			result[i] = arr1[i] / arr2[i]
		}
	}
	return result
}

//both of the functions are from https://github.com/blend/go-sdk/blob/v1.20220411.3/mathutil/round.go
//add testing here
/*func Normalize(values []float64) []float64 {
	var total float64
	for _, v := range values {
		total += v
	}
	output := make([]float64, len(values))
	for x, v := range values {
		output[x] = RoundDown(v/total, 0.0001)
	}
	return output
}

//add testing here
func RoundDown(value, roundTo float64) float64 {
	d1 := math.Floor(value / roundTo)
	return d1 * roundTo
}*/

func Normalize(arr []float64) []float64 {
	maxVal := findMax(arr)
	minVal := findMin(arr)

	// Find the minimum and maximum values in the array
	for _, v := range arr {
		if v < minVal {
			minVal = v
		}
		if v > maxVal {
			maxVal = v
		}
	}

	// Normalize the array to values between 0 and 1
	normalized := make([]float64, len(arr))
	for i, v := range arr {
		normalized[i] = (v - minVal) / (maxVal - minVal)
	}

	return normalized
}
