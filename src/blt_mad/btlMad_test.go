package blt_mad

import (
	"testing"
)

func TestArray0(t *testing.T) {
	tau := 0.0

	//repeating values
	data1 := []float64{3, 1, 8, 1}
	expected1 := []float64{3.0, 8.0}
	result1 := BltMad(data1, tau)
	if !equalSlices(expected1, result1) {
		t.Errorf("Test case 1 failed. Got: %f, Expected: %f", result1, expected1)
	}

	data2 := []float64{3, 1, 8, 1, 5, 7}
	expected2 := []float64{5.0, 7.0, 8.0}
	result2 := BltMad(data2, tau)
	if !equalSlices(expected2, result2) {
		t.Errorf("Test case 2 failed. Got: %f, Expected: %f", result2, expected2)
	}
}

func TestArray3(t *testing.T) {
	tau := 3.0

	//repeating values
	data1 := []float64{3, 1, 8, 1}
	expected1 := []float64{8.0}
	result1 := BltMad(data1, tau)
	if !equalSlices(expected1, result1) {
		t.Errorf("Test case 1 failed. Got: %f, Expected: %f", result1, expected1)
	}

	data2 := []float64{3, 1, 8, 1, 5, 7}
	expected2 := []float64{5.0, 7.0, 8.0}
	result2 := BltMad(data2, tau)
	if !equalSlices(expected2, result2) {
		t.Errorf("Test case 2 failed. Got: %f, Expected: %f", result2, expected2)
	}
}

func TestBGP10(t *testing.T) {
	tau := 10.0

	//Test case 2: bgp data 1
	data2, _ := TxtIntoArrayFloat64("bgpDataCounts1.txt")
	expected2 := []float64{14065.0, 16974.0, 25746.0}
	result2 := BltMad(data2, tau)
	if !equalSlices(expected2, result2) {
		t.Errorf("Test case 2 failed. Got: %f, Expected: %f", result2, expected2)
	}

	//Test case 3: bgp data 2
	data3, _ := TxtIntoArrayFloat64("bgpDataCounts2.txt")
	expected3 := []float64{40552.0, 122393.0, 35956.0, 118534.0}
	result3 := BltMad(data3, tau)
	if !equalSlices(expected3, result3) {
		t.Errorf("Test case 3 failed. Got: %f, Expected: %f", result3, expected3)
	}

	//Test case 4: bgp data 3
	data4, _ := TxtIntoArrayFloat64("bgpDataCounts3.txt")
	expected4 := []float64{2909.0, 2509.0, 4166.0, 3888.0, 2542.0, 2864.0, 2640.0, 2357.0, 2314.0, 2701.0, 2436.0, 2555.0, 2438.0, 2802.0, 2352.0, 2802.0}
	result4 := BltMad(data4, tau)
	if !equalSlices(expected4, result4) {
		t.Errorf("Test case 4 failed. Got: %f, Expected: %f", result4, expected4)
	}

}

func TestBGP100(t *testing.T) {
	tau := 100.0 //the outputs should be empty arrays when the sensitivity parameter is too large
	//Test case 2: bgp data 1
	data2, _ := TxtIntoArrayFloat64("bgpDataCounts1.txt")
	expected2 := []float64{}
	result2 := BltMad(data2, tau)
	if !equalSlices(expected2, result2) {
		t.Errorf("Test case 2 failed. Got: %f, Expected: %f", result2, expected2)
	}
	//Test case 3: bgp data 2
	data3, _ := TxtIntoArrayFloat64("bgpDataCounts2.txt")
	expected3 := []float64{}
	result3 := BltMad(data3, tau)
	if !equalSlices(expected3, result3) {
		t.Errorf("Test case 3 failed. Got: %f, Expected: %f", result3, expected3)
	}

	//Test case 4: bgp data 3
	data4, _ := TxtIntoArrayFloat64("bgpDataCounts3.txt")
	expected4 := []float64{}
	result4 := BltMad(data4, tau)
	if !equalSlices(expected4, result4) {
		t.Errorf("Test case 4 failed. Got: %f, Expected: %f", result4, expected4)
	}

}
