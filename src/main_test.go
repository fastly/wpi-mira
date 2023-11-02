package main

import (
	"testing"
)

func TestOutlierDetection(t *testing.T) {

	actualOutliers := findOutliers([]int{1, 1, 1, 1, 10, 1, 1, 1, 1})
	expectedOutlier := []int{10}

	// Check if the arrays are equal
	isEqual := AreArraysEqual(actualOutliers, expectedOutlier)

	if !isEqual {
		t.Errorf("Got %+v want %+v\n", actualOutliers, expectedOutlier)
	}

}

// AreArraysEqual checks if two arrays of integers are equal
func AreArraysEqual(arr1, arr2 []int) bool {
	if len(arr1) != len(arr2) {
		return false
	}

	for i := 0; i < len(arr1); i++ {
		if arr1[i] != arr2[i] {
			return false
		}
	}

	return true
}
