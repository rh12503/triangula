package testutils

import (
	"fmt"
	"github.com/RH12503/Triangula/algorithm"
	imageData "github.com/RH12503/Triangula/image"
	"github.com/RH12503/Triangula/render"
	"github.com/RH12503/Triangula/triangulation"
	"github.com/RH12503/draw/draw"
)

var currentAlgorithm algorithm.Algorithm
var repetitions int
var image imageData.Data

func RunWindow(img imageData.Data, algo algorithm.Algorithm, reps, maxSize int) {
	image = img
	currentAlgorithm = algo
	repetitions = reps

	w, h := image.Size()
	if w > h {
		ratio := float64(maxSize) / float64(w)
		w = maxSize
		h = int(float64(h) * ratio)
	} else {
		ratio := float64(maxSize) / float64(h)
		h = maxSize
		w = int(float64(w) * ratio)
	}

	draw.RunWindow("test", w, h, update)
}

func update(window draw.Window) {
	width, height := window.Size()

	window.FillRect(0, 0, width, height, draw.White)

	var timeSum int64
	for i := 0; i < repetitions; i++ {
		currentAlgorithm.Step()
		timeSum += currentAlgorithm.Stats().TimeForGen.Microseconds()
	}

	stats := currentAlgorithm.Stats()
	fmt.Printf("Gen: %v | Fit: %v | Time: %v\n", stats.Generation, stats.BestFitness, float64(timeSum)/(float64(repetitions)*1000.))

	w, h := image.Size()
	triangles, _ := triangulation.Triangulate(currentAlgorithm.Best(), w, h)
	triangleData := render.TrianglesOnImage(triangles, image)

	for _, d := range triangleData {
		drawTriangle(d, window)
	}
}

func drawTriangle(data render.TriangleData, window draw.Window) {
	w, h := window.Size()
	t := data.Triangle
	c := data.Color
	window.FillTriangle(
		scale(t.Points[0].X, w),
		scale(t.Points[0].Y, h),
		scale(t.Points[1].X, w),
		scale(t.Points[1].Y, h),
		scale(t.Points[2].X, w),
		scale(t.Points[2].Y, h),
		draw.RGBA(float32(c.R), float32(c.G), float32(c.B), 1),
	)
}
