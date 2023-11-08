package blt_mad

import (
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

func generatePlot(data []float64, fileName string) {
	// Create a new plot
	p := plot.New()
	/*if err != nil {
		panic(err)
	}*/

	// Generate data points for the bell curve
	points := make(plotter.XYs, len(data))
	for i, value := range data {
		points[i].X = float64(i) // X values as indices
		points[i].Y = value      // Y values from the data
	}

	// Create a scatter plot using the generated data
	s, _ := plotter.NewScatter(points)

	// Add the scatter plot to the plot
	p.Add(s)

	// Create a horizontal line at the mean and color it red
	hLine, _ := plotter.NewLine(plotter.XYs{{0, findMean(data)}, {float64(len(data) - 1), findMean(data)}})
	hLine.LineStyle.Width = 3
	hLine.LineStyle.Color = plotutil.Color(0) // Red color

	// Add the horizontal line to the plot
	p.Add(hLine)

	// Create a horizontal line at the median and color it red
	hLine1, _ := plotter.NewLine(plotter.XYs{{0, findMedian(data)}, {float64(len(data) - 1), findMean(data)}})
	hLine1.LineStyle.Width = 3
	hLine1.LineStyle.Color = plotutil.Color(1) // Green color

	// Add the horizontal line to the plot
	p.Add(hLine1)

	//add a label distinguishing the lines -> need to do

	// Set the title and labels
	p.Title.Text = "The count of messages in each bucket"
	p.X.Label.Text = "Bucket number"
	p.Y.Label.Text = "The number of messages"

	// Save the plot to an image file (e.g., PNG)
	if err := p.Save(5*vg.Inch, 5*vg.Inch, fileName); err != nil {
		panic(err)
	}
}
