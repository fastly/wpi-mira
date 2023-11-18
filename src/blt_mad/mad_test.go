package blt_mad

import (
	"math"
	"reflect"
	"testing"
)

func TestMad(t *testing.T) {
	// Test case 1: Empty data
	data1 := []float64{}
	expected1 := math.SmallestNonzeroFloat64
	result1 := Mad(data1)
	if !reflect.DeepEqual(result1, expected1) {
		t.Errorf("Test case 1 in TestMad failed. Got: %v, Expected: %v", result1, expected1)
	}

	// Test case 2: Negative data
	data2 := []float64{-1.0, -2.0, 0.0, -1.0}
	expected2 := 0.5
	result2 := Mad(data2)
	if !reflect.DeepEqual(result2, expected2) {
		t.Errorf("Test case 2 in TestMad failed. Got: %v, Expected: %v", result2, expected2)
	}

	// Test case 3: All zero data
	data3 := []float64{1, 2, 3, 4, 5}
	expected3 := 1.2
	result3 := Mad(data3)
	if !reflect.DeepEqual(result3, expected3) {
		t.Errorf("Test case 3 in TestMad failed. Got: %v, Expected: %v", result3, expected3)
	}

	// Test case 4: Outlier test case
	data4 := []float64{1, 1, 1, 1, 10}
	expected4 := 2.88
	result4 := Mad(data4)
	if !reflect.DeepEqual(result4, expected4) {
		t.Errorf("Test case 4 in TestMad failed. Got: %v, Expected: %v", result4, expected4)
	}
}
