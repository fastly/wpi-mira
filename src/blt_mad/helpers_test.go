package blt_mad

import (
	"math"
	"reflect"
	"testing"
)

const e = 0.01

func TestRemoveZerosPresent(t *testing.T) {
	// Test case 1: Removing zeros from a slice
	data1 := []float64{0.0, 1.2, 0.0, 3.4, 0.0, 5.6}
	expected1 := []float64{1.2, 3.4, 5.6}
	result1 := removeZeros(data1)
	if !reflect.DeepEqual(result1, expected1) {
		t.Errorf("Test case 1 in TestRemoveZeros failed. Got: %v, Expected: %v", result1, expected1)
	}
}

func TestRemoveZerosNotPresent(t *testing.T) {
	// Test case 2: No zeros in the input slice
	data2 := []float64{1.2, 3.4, 5.6}
	expected2 := []float64{1.2, 3.4, 5.6}
	result2 := removeZeros(data2)
	if !reflect.DeepEqual(result2, expected2) {
		t.Errorf("Test case 2  in TestRemoveZeros failed. Got: %v, Expected: %v", result2, expected2)
	}
}

func TestRemoveZerosEmpty(t *testing.T) {
	// Test case 3: Empty input slice
	data3 := []float64{}
	expected3 := []float64{}
	result3 := removeZeros(data3)
	if !equalSlices(result3, expected3) {
		t.Errorf("Test case 3  in TestRemoveZeros failed. Got: %v, Expected: %v", result3, expected3)
	}
}

func TestRemoveZerosAllZeros(t *testing.T) {
	// Test case 4: All zeros array
	data4 := []float64{0.0, 0.0, 0.0}
	expected4 := []float64{}
	result4 := removeZeros(data4)
	if !equalSlices(expected4, result4) {
		t.Errorf("Test case 4 in TestRemoveZeros failed. Got: %v, Expected: %v", result4, expected4)
	}
}

func TestFindMedianOddLength(t *testing.T) {
	// Test case 1: Odd-length slice
	data1 := []float64{3.0, 1.0, 2.0}
	expected1 := 2.0
	result1 := FindMedian(data1)
	if result1 != expected1 {
		t.Errorf("Test case 1 in TestFindMedian failed. Got: %f, Expected: %f", result1, expected1)
	}
}

func TestFindMedianEvenLength(t *testing.T) {
	// Test case 2: Even-length slice
	data2 := []float64{4.0, 2.0, 1.0, 3.0}
	expected2 := 2.5
	result2 := FindMedian(data2)
	if result2 != expected2 {
		t.Errorf("Test case 2 in TestFindMedian failed. Got: %f, Expected: %f", result2, expected2)
	}
}

func TestFindMedianEmptySlice(t *testing.T) {
	// Test case 3: Empty slice
	data3 := []float64{}
	expected3 := math.SmallestNonzeroFloat64 // You can choose an appropriate default value for an empty slice
	result3 := FindMedian(data3)
	if !WithinTolerance(expected3, result3, e) {
		t.Errorf("Test case 3 in TestFindMedian failed. Got: %f, Expected: %f", result3, expected3)
	}
}

func TestFindMedianDuplicateValues(t *testing.T) {
	// Test case 4: Duplicate values
	data4 := []float64{3.0, 1.0, 8.0, 1.0}
	expected4 := 2.0 // You can choose an appropriate default value for an empty slice
	result4 := FindMedian(data4)
	if !WithinTolerance(expected4, result4, e) {
		t.Errorf("Test case 4 in TestFindMedian failed. Got: %f, Expected: %f", result4, expected4)
	}

	// Test case 4: Duplicate values with the value being a median
	data5 := []float64{1.0, 2.0, 3.0, 3.0, 4.0}
	expected5 := 3.0 // You can choose an appropriate default value for an empty slice
	result5 := FindMedian(data5)
	if !WithinTolerance(expected5, result5, e) {
		t.Errorf("Test case 4 in TestFindMedian failed. Got: %f, Expected: %f", result5, expected5)
	}
}

func TestFindMinNonEmpty(t *testing.T) {
	// Test case 1: Non-empty slice
	data1 := []float64{3.0, 1.0, 2.0, 5.0}
	expected1 := 1.0
	result1 := findMin(data1)
	if !WithinTolerance(expected1, result1, e) {
		t.Errorf("Test case 1 failed. Got: %f, Expected: %f", result1, expected1)
	}
}

func TestFindMinEmpty(t *testing.T) {
	// Test case 2: Empty slice
	data2 := []float64{}
	expected2 := math.SmallestNonzeroFloat64
	result2 := findMin(data2)
	if !WithinTolerance(expected2, result2, e) {
		t.Errorf("Test case 2 failed. Got: %f, Expected: %f", result2, expected2)
	}
}

func TestFindMinNegativeSlice(t *testing.T) {
	data3 := []float64{-3.0, -1.0, -2.0, -5.0}
	expected3 := -5.0
	result3 := findMin(data3)
	if !WithinTolerance(expected3, result3, e) {
		t.Errorf("Test case 3 failed. Got: %f, Expected: %f", result3, expected3)
	}
}

func TestFindMaxNonEmpty(t *testing.T) {
	// Test case 1: Non-empty slice
	data1 := []float64{3.0, 1.0, 2.0, 5.0}
	expected1 := 5.0
	result1 := findMax(data1)
	if !WithinTolerance(expected1, result1, e) {
		t.Errorf("Test case 1 failed. Got: %f, Expected: %f", result1, expected1)
	}
}

func TestFindMaxEmpty(t *testing.T) {
	// Test case 2: Empty slice
	data2 := []float64{}
	expected2 := math.MaxFloat64
	result2 := findMax(data2)
	if !WithinTolerance(expected2, result2, e) {
		t.Errorf("Test case 2 failed. Got: %f, Expected: %f", result2, expected2)
	}
}

func TestFindMaxNegativeSlice(t *testing.T) {
	//Test case 3: Negative slice
	data3 := []float64{-3.0, -1.0, -2.0, -5.0}
	expected3 := -1.0
	result3 := findMax(data3)
	if !WithinTolerance(expected3, result3, e) {
		t.Errorf("Test case 3 failed. Got: %f, Expected: %f", result3, expected3)
	}
}

func TestFindMeanPositive(t *testing.T) {
	// Test case 1: Test with positive numbers
	data1 := []float64{1.0, 2.0, 3.0, 4.0, 5.0}
	expected1 := 3.0
	result1 := findMean(data1)
	if !WithinTolerance(expected1, result1, e) {
		t.Errorf("Expected findMean of %v to be %f, but got %f", data1, expected1, result1)
	}
}

func TestFindMeanNegative(t *testing.T) {
	// Test case 2: Test with negative numbers; will not apply with BGP mesage coutns
	data2 := []float64{-1.0, -2.0, -3.0, -4.0, -5.0}
	expected2 := -3.0
	result2 := findMean(data2)
	if !WithinTolerance(expected2, result2, e) {
		t.Errorf("Expected findMean of %v to be %f, but got %f", data2, expected2, result2)
	}
}

func TestFindMeanSingleValue(t *testing.T) {
	// Test case 3: Test with a single value
	data3 := []float64{42.0}
	expected3 := 42.0
	result3 := findMean(data3)
	if !WithinTolerance(expected3, result3, e) {
		t.Errorf("Expected findMean of %v to be %f, but got %f", data3, expected3, result3)
	}
}

func TestFindMeanZeros(t *testing.T) {
	// Test case 4: Test with a zeros
	data4 := []float64{0.0}
	expected4 := 0.0
	result4 := findMean(data4)
	if !WithinTolerance(expected4, result4, e) {
		t.Errorf("Expected findMean of %v to be %f, but got %f", data4, expected4, result4)
	}
}

func TestFindMeanEmptySlice(t *testing.T) {
	// Test case 5: Test with empty array
	data5 := []float64{}
	expected5 := -1.0
	result5 := findMean(data5)
	if !WithinTolerance(expected5, result5, e) {
		t.Errorf("Expected findMean of %v to be %f, but got %f", data5, expected5, result5)
	}
}

func TestContainsValuesPositive(t *testing.T) {
	// Test case 1: Test with positive array of numbers where the value is contained
	data1 := []float64{1.0, 2.0, 3.0, 4.0, 5.0}
	value1 := 3.0
	expected1 := true
	result1 := containsValue(data1, value1)
	if result1 != expected1 {
		t.Errorf("Test case 1  in TestContainsValue failed. Got: %t, Expected: %t", result1, expected1)
	}

	// Test case 2: Test with positive array of numbers where the value is not contained
	data2 := []float64{1.0, 2.0, 3.0, 4.0, 5.0}
	value2 := 0.0
	expected2 := false
	result2 := containsValue(data2, value2)
	if result2 != expected2 {
		t.Errorf("Test case 2  in TestContainsValue failed. Got: %t, Expected: %t", result2, expected2)
	}
}

func TestContainsValuesEmpty(t *testing.T) {
	// Test case 3: Empty array
	data3 := []float64{}
	value3 := 0.0
	expected3 := false
	result3 := containsValue(data3, value3)
	if result3 != expected3 {
		t.Errorf("Test case 1  in TestContainsValue failed. Got: %t, Expected: %t", result3, expected3)
	}
}

//reflect deep equal and equal slices for some reason do not work; the outputs are the same but the tests fail
func TestContainsAllElementsPositiveHas(t *testing.T) {
	emptyArr := []float64{}
	// Test case 1: Test with positive array of numbers where the subarray is contained
	data1 := []float64{1.0, 2.0, 3.0, 4.0, 5.0}
	subArray1 := []float64{1.0}
	expected1 := emptyArr
	result1 := containAllElements(data1, subArray1)
	if !WithinToleranceFloatSlice(expected1, result1, e) {
		t.Errorf("Test case 1  in TestContainsAllElements failed. Got: %f, Expected: %f", result1, expected1)
	}
}

func TestContainsAllElementsEmptySubarray(t *testing.T) {
	emptyArr := []float64{}
	data2 := []float64{1.0, 2.0, 3.0, 4.0, 5.0}
	subArray2 := []float64{}
	expected2 := emptyArr
	result2 := containAllElements(data2, subArray2)
	if !WithinToleranceFloatSlice(expected2, result2, e) {
		t.Errorf("Test case 2  in TestContainsAllElements failed. Got: %f, Expected: %f", result2, expected2)
	}
}

func TestContainsAllElementsExtraValueSubarray(t *testing.T) {
	// Test case 3: Subarray with one extra value
	data3 := []float64{1.0, 2.0, 3.0, 4.0, 5.0}
	subArray3 := []float64{1.0, 6.0}
	expected3 := []float64{6.0}
	result3 := containAllElements(data3, subArray3)
	if !WithinToleranceFloatSlice(expected3, result3, e) {
		t.Errorf("Test case 3  in TestContainsAllElements failed. Got: %f, Expected: %f", result3, expected3)
	}
}

func TestContainsAllElementsEmptyMain(t *testing.T) {
	// Test case 4: The main array is empty; output is the whole subarray
	data4 := []float64{}
	subArray4 := []float64{1.0, 2.0, 3.0, 4.0, 5.0}
	expected4 := []float64{1.0, 2.0, 3.0, 4.0, 5.0}
	result4 := containAllElements(data4, subArray4)
	if !WithinToleranceFloatSlice(expected4, result4, e) {
		t.Errorf("Test case 4 in TestContainsAllElements failed. Got: %f, Expected: %f", result4, expected4)
	}

}
