package blt_mad

import (
	"math"
	"sort"
	"strconv"
	"strings"
)

//convert a string of numbers into an array []float64 -> it actually works; all the functions can be called
func convertToFloat64Array(str string) []float64 {
	var numbers []float64
	numbersStr := strings.Fields(str)
	//fmt.Println(reflect.TypeOf(numbersStr[0]))
	//f, _ := strconv.ParseFloat(numbersStr[0], 64)
	//fmt.Println(reflect.TypeOf(f))

	for i, _ := range numbersStr {
		f, _ := strconv.ParseFloat(numbersStr[i], 64)
		numbers = append(numbers, f)
	}
	return numbers
}

func removeZeros(data []float64) []float64 {
	var nonZeros []float64
	for _, value := range data {
		if value != 0.0 {
			nonZeros = append(nonZeros, value)
		}
	}
	return nonZeros
}

func findMedian(data []float64) float64 {
	var med float64
	sort.Float64s(data)
	if len(data) == 0 {
		return math.SmallestNonzeroFloat64
	} else if len(data)%2 == 1 {
		//if the length is odd return the number in the middle
		med = data[len(data)/2]
	} else {
		//if the lenght is even take the findMean of the middle two numbers
		left := data[(len(data)/2)-1]
		right := data[(len(data) / 2)]
		med = (left + right) / 2
	}

	return med
}

func findMin(data []float64) float64 {
	var min float64
	if len(data) == 0 {
		return 0.0
	} else {
		sort.Float64s(data)
		min = data[0]
	}
	return min
}

func findMax(data []float64) float64 {
	var max float64
	if len(data) == 0 {
		return 0.0 //invalid findMean in the context of BGP message counts
	} else {
		sort.Float64s(data)
		max = data[len(data)-1]
	}

	return max
}

func findMean(data []float64) float64 {
	sum := 0.0
	if len(data) == 0 {
		return -1.0 //invalid findMean in the context of BGP message counts
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
