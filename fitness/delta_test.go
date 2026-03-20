package fitness

import (
	image2 "github.com/RH12503/Triangula/image"
	"github.com/RH12503/Triangula/mutation"
	"github.com/RH12503/Triangula/normgeom"
	"github.com/RH12503/Triangula/random"
	"image"
	"image/color"
	"math"
	"testing"
)

// TestDeltaFitnessMatchesFull verifies that delta-based fitness evaluation produces
// the same result as a full evaluation for the same set of points.
func TestDeltaFitnessMatchesFull(t *testing.T) {
	random.Seed(42)

	const w, h = 80, 80
	const bs = 3

	// Create a random target image
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			img.Set(x, y, color.RGBA{
				R: uint8(random.Intn(math.MaxUint8)),
				G: uint8(random.Intn(math.MaxUint8)),
				B: uint8(random.Intn(math.MaxUint8)),
				A: math.MaxUint8,
			})
		}
	}

	target := image2.ToData(img)

	// Create base points
	basePoints := normgeom.NormPointGroup{
		{0.1, 0.1},
		{0.9, 0.1},
		{0.5, 0.5},
		{0.1, 0.9},
		{0.9, 0.9},
		{0.3, 0.3},
		{0.7, 0.7},
		{0.2, 0.6},
		{0.8, 0.4},
		{0.5, 0.2},
	}

	// Compute full fitness of the base points
	baseEval := NewTrianglesImageFunction(target, bs)
	baseFitness := baseEval.Calculate(PointsData{
		Points:    basePoints,
		Mutations: nil,
	})
	t.Logf("Base fitness: %v", baseFitness)

	// Create mutated points: move point 2 and point 5
	mutatedPoints := basePoints.Copy()
	oldP2 := mutatedPoints[2]
	mutatedPoints[2] = normgeom.NormPoint{X: 0.55, Y: 0.45}
	newP2 := mutatedPoints[2]

	oldP5 := mutatedPoints[5]
	mutatedPoints[5] = normgeom.NormPoint{X: 0.35, Y: 0.25}
	newP5 := mutatedPoints[5]

	mutations := []mutation.Mutation{
		{Index: 2, Old: oldP2, New: newP2},
		{Index: 5, Old: oldP5, New: newP5},
	}

	// Method 1: Full evaluation of mutated points from scratch
	fullEval := NewTrianglesImageFunction(target, bs)
	fullFitness := fullEval.Calculate(PointsData{
		Points:    mutatedPoints.Copy(),
		Mutations: nil,
	})

	// Method 2: Delta evaluation from base
	deltaEval := NewTrianglesImageFunction(target, bs)
	// First, compute the base to populate triangleVariances
	deltaEval.Calculate(PointsData{
		Points:    basePoints.Copy(),
		Mutations: nil,
	})
	// Now set this as the base and compute via delta path
	deltaEval2 := NewTrianglesImageFunction(target, bs)
	deltaEval2.SetBase(deltaEval)
	deltaFitness := deltaEval2.Calculate(PointsData{
		Points:    mutatedPoints.Copy(),
		Mutations: mutations,
	})

	t.Logf("Full fitness:  %v", fullFitness)
	t.Logf("Delta fitness: %v", deltaFitness)

	// They should match within floating-point tolerance
	epsilon := 1e-10
	diff := math.Abs(fullFitness - deltaFitness)
	if diff > epsilon {
		t.Errorf("Delta fitness (%v) does not match full fitness (%v), diff=%v", deltaFitness, fullFitness, diff)
	}
}

// TestDeltaFitnessSingleMutation tests with a single point mutation.
func TestDeltaFitnessSingleMutation(t *testing.T) {
	random.Seed(99)

	const w, h = 60, 60
	const bs = 3

	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			img.Set(x, y, color.RGBA{
				R: uint8(random.Intn(math.MaxUint8)),
				G: uint8(random.Intn(math.MaxUint8)),
				B: uint8(random.Intn(math.MaxUint8)),
				A: math.MaxUint8,
			})
		}
	}

	target := image2.ToData(img)

	basePoints := normgeom.NormPointGroup{
		{0.15, 0.15},
		{0.85, 0.15},
		{0.5, 0.5},
		{0.15, 0.85},
		{0.85, 0.85},
	}

	// Compute base
	baseEval := NewTrianglesImageFunction(target, bs)
	baseEval.Calculate(PointsData{
		Points:    basePoints.Copy(),
		Mutations: nil,
	})

	// Single mutation on point 2
	mutatedPoints := basePoints.Copy()
	oldP := mutatedPoints[2]
	mutatedPoints[2] = normgeom.NormPoint{X: 0.52, Y: 0.48}
	newP := mutatedPoints[2]

	mutations := []mutation.Mutation{
		{Index: 2, Old: oldP, New: newP},
	}

	// Full evaluation
	fullEval := NewTrianglesImageFunction(target, bs)
	fullFitness := fullEval.Calculate(PointsData{
		Points:    mutatedPoints.Copy(),
		Mutations: nil,
	})

	// Delta evaluation
	deltaEval := NewTrianglesImageFunction(target, bs)
	deltaEval.SetBase(baseEval)
	deltaFitness := deltaEval.Calculate(PointsData{
		Points:    mutatedPoints.Copy(),
		Mutations: mutations,
	})

	t.Logf("Full: %v, Delta: %v", fullFitness, deltaFitness)

	epsilon := 1e-10
	diff := math.Abs(fullFitness - deltaFitness)
	if diff > epsilon {
		t.Errorf("Delta fitness (%v) does not match full fitness (%v), diff=%v", deltaFitness, fullFitness, diff)
	}
}
