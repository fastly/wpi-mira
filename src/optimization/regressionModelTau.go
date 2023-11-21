package optimization

/*import (
	"fmt"
	"log"
	"math"
	"os"

	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/stat"
	"gonum.org/v1/gonum/stat/regression"
)

func main() {
	// Load and split the data into testing and training
	taus := loadFile("taus.txt")
	mads := loadFile("mads.txt")
	medians := loadFile("medians.txt")
	ratio := make([]float64, len(mads))
	normalizedRatio := make([]float64, len(ratio))
	normalizedTaus := make([]float64, len(taus))

	for i := range ratio {
		ratio[i] = mads[i] / medians[i]
		normalizedRatio[i] = (ratio[i] - floats.Min(ratio)) / (floats.Max(ratio) - floats.Min(ratio))
		normalizedTaus[i] = (taus[i] - floats.Min(taus)) / (floats.Max(taus) - floats.Min(taus))
	}

	x := mat.NewDense(len(normalizedRatio), 1, normalizedRatio)
	y := mat.NewDense(len(normalizedTaus), 1, normalizedTaus)

	X_train, X_test, y_train, y_test := splitData(x, y, 0.2)

	// Fit the model (linear regression)
	model := new(regression.Linear)
	model.Fit(X_train, y_train)

	// Predict on the test set
	predictions := mat.NewDense(0, 0, nil)
	model.Predict(predictions, X_test)

	// Check the results and calculate percentage
	arrayDif := mat.NewDense(0, 0, nil)
	arrayDif.Sub(y_test, predictions)
	arrayDif.Apply(func(i, j int, v float64) float64 {
		return math.Abs(v)
	}, arrayDif)
	_, col := arrayDif.Dims()
	diffValues := make([]float64, col)
	for i := 0; i < col; i++ {
		diffValues[i] = arrayDif.At(0, i)
	}
	percentageGreaterThan1 := stat.Mean(diffValues, nil) * 100
	fmt.Println(percentageGreaterThan1)
}

func loadFile(fileName string) []float64 {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var data []float64
	for {
		var value float64
		_, err := fmt.Fscanf(file, "%f\n", &value)
		if err != nil {
			break
		}
		data = append(data, value)
	}
	return data
}

func splitData(x, y *mat.Dense, testSize float64) (X_train, X_test, y_train, y_test *mat.Dense) {
	rows, _ := x.Dims()
	testRows := int(float64(rows) * testSize)
	trainRows := rows - testRows

	X_train = mat.NewDense(trainRows, 1, nil)
	X_test = mat.NewDense(testRows, 1, nil)
	y_train = mat.NewDense(trainRows, 1, nil)
	y_test = mat.NewDense(testRows, 1, nil)

	X_train.Slice(0, trainRows, 0, 1).(*mat.Dense).Copy(x.Slice(0, trainRows, 0, 1))
	X_test.Slice(0, testRows, 0, 1).(*mat.Dense).Copy(x.Slice(trainRows, rows, 0, 1))
	y_train.Slice(0, trainRows, 0, 1).(*mat.Dense).Copy(y.Slice(0, trainRows, 0, 1))
	y_test.Slice(0, testRows, 0, 1).(*mat.Dense).Copy(y.Slice(trainRows, rows, 0, 1))

	return X_train, X_test, y_train, y_test
}*/

