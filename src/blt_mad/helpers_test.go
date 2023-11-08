package blt_mad

import (
	"math"
	"reflect"
	"testing"
)

//compaerring empty slices -> correct results but they do not register as equal
func TestRemoveZeros(t *testing.T) {
	// Test case 1: Removing zeros from a slice
	data1 := []float64{0.0, 1.2, 0.0, 3.4, 0.0, 5.6}
	expected1 := []float64{1.2, 3.4, 5.6}
	result1 := removeZeros(data1)
	if !reflect.DeepEqual(result1, expected1) {
		t.Errorf("Test case 1 in TestRemoveZeros failed. Got: %v, Expected: %v", result1, expected1)
	}

	// Test case 2: No zeros in the input slice
	data2 := []float64{1.2, 3.4, 5.6}
	expected2 := []float64{1.2, 3.4, 5.6}
	result2 := removeZeros(data2)
	if !reflect.DeepEqual(result2, expected2) {
		t.Errorf("Test case 2  in TestRemoveZeros failed. Got: %v, Expected: %v", result2, expected2)
	}

	// Test case 3: Empty input slice
	data3 := []float64{}
	expected3 := []float64{}
	result3 := removeZeros(data3)
	if !equalSlices(result3, expected3) {
		t.Errorf("Test case 3  in TestRemoveZeros failed. Got: %v, Expected: %v", result3, expected3)
	}

	// Test case 4: All zeros array
	data4 := []float64{0.0, 0.0, 0.0}
	expected4 := []float64{}
	result4 := removeZeros(data4)
	if !equalSlices(expected4, result4) {
		t.Errorf("Test case 4 in TestRemoveZeros failed. Got: %v, Expected: %v", result3, expected3)
	}
}

func TestFindMedian(t *testing.T) {
	// Test case 1: Odd-length slice
	data1 := []float64{3.0, 1.0, 2.0}
	expected1 := 2.0
	result1 := findMedian(data1)
	if result1 != expected1 {
		t.Errorf("Test case 1 in TestFindMedian failed. Got: %f, Expected: %f", result1, expected1)
	}

	// Test case 2: Even-length slice
	data2 := []float64{4.0, 2.0, 1.0, 3.0}
	expected2 := 2.5
	result2 := findMedian(data2)
	if result2 != expected2 {
		t.Errorf("Test case 2 in TestFindMedian failed. Got: %f, Expected: %f", result2, expected2)
	}

	// Test case 3: Empty slice
	data3 := []float64{}
	expected3 := math.SmallestNonzeroFloat64 // You can choose an appropriate default value for an empty slice
	result3 := findMedian(data3)
	if result3 != expected3 {
		t.Errorf("Test case 3 in TestFindMedian failed. Got: %f, Expected: %f", result3, expected3)
	}
}

func TestFindMin(t *testing.T) {
	// Test case 1: Non-empty slice
	data1 := []float64{3.0, 1.0, 2.0, 5.0}
	expected1 := 1.0
	result1 := findMin(data1)
	if result1 != expected1 {
		t.Errorf("Test case 1 failed. Got: %f, Expected: %f", result1, expected1)
	}

	// Test case 2: Empty slice
	data2 := []float64{}
	expected2 := 0.0
	result2 := findMin(data2)
	if result2 != expected2 {
		t.Errorf("Test case 2 failed. Got: %f, Expected: %f", result2, expected2)
	}

	//Test case 3: Negative slice
	data3 := []float64{-3.0, -1.0, -2.0, -5.0}
	expected3 := -5.0
	result3 := findMin(data3)
	if result3 != expected3 {
		t.Errorf("Test case 3 failed. Got: %f, Expected: %f", result3, expected3)
	}
}

func TestFindMean(t *testing.T) {
	// Test case 1: Test with positive numbers
	data1 := []float64{1.0, 2.0, 3.0, 4.0, 5.0}
	expected1 := 3.0
	result1 := findMean(data1)
	if result1 != expected1 {
		t.Errorf("Expected findMean of %v to be %f, but got %f", data1, expected1, result1)
	}

	// Test case 2: Test with negative numbers; will not apply with BGP mesage coutns
	data2 := []float64{-1.0, -2.0, -3.0, -4.0, -5.0}
	expected2 := -3.0
	result2 := findMean(data2)
	if result2 != expected2 {
		t.Errorf("Expected findMean of %v to be %f, but got %f", data2, expected2, result2)
	}

	// Test case 3: Test with a single value
	data3 := []float64{42.0}
	expected3 := 42.0
	result3 := findMean(data3)
	if result3 != expected3 {
		t.Errorf("Expected findMean of %v to be %f, but got %f", data3, expected3, result3)
	}

	// Test case 4: Test with a zeros
	data4 := []float64{0.0}
	expected4 := 0.0
	result4 := findMean(data4)
	if result4 != expected4 {
		t.Errorf("Expected findMean of %v to be %f, but got %f", data3, expected3, result3)
	}

	// Test case 5: Test with empty array
	data5 := []float64{}
	expected5 := -1.0
	result5 := findMean(data5)
	if result5 != expected5 {
		t.Errorf("Expected findMean of %v to be %f, but got %f", data3, expected3, result3)
	}
}
