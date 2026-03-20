package algorithm

import (
	"github.com/RH12503/Triangula/algorithm/evaluator"
	"github.com/RH12503/Triangula/fitness"
	"github.com/RH12503/Triangula/generator"
	imageData "github.com/RH12503/Triangula/image"
	"github.com/RH12503/Triangula/mutation"
	"github.com/RH12503/Triangula/normgeom"
	"github.com/RH12503/Triangula/random"
	"github.com/stretchr/testify/assert"
	"image"
	"image/color"
	_ "image/jpeg"
	"log"
	"math/rand"
	"os"
	"testing"
	"time"
)

// makeSmallAlgo creates a fast, small algorithm instance for unit testing.
// It uses a 20x20 solid-color image and a small population to keep tests quick.
func makeSmallAlgo() *modifiedGenetic {
	random.Seed(42)
	img := image.NewRGBA(image.Rect(0, 0, 20, 20))
	for x := 0; x < 20; x++ {
		for y := 0; y < 20; y++ {
			img.Set(x, y, color.RGBA{R: 120, G: 80, B: 40, A: 255})
		}
	}
	imgData := imageData.ToData(img)

	pointFactory := func() normgeom.NormPointGroup {
		return (generator.RandomGenerator{}).Generate(10)
	}
	evalFactory := func(n int) evaluator.Evaluator {
		return evaluator.NewParallel(fitness.TrianglesImageFunctions(imgData, 3, n), 4)
	}
	mutator := mutation.NewGaussianMethod(0.5, 0.3)
	return NewModifiedGenetic(pointFactory, 10, 2, evalFactory, mutator)
}

func makeSmallSimpleAlgo() *simple {
	random.Seed(42)
	img := image.NewRGBA(image.Rect(0, 0, 20, 20))
	for x := 0; x < 20; x++ {
		for y := 0; y < 20; y++ {
			img.Set(x, y, color.RGBA{R: 120, G: 80, B: 40, A: 255})
		}
	}
	imgData := imageData.ToData(img)

	pointFactory := func() normgeom.NormPointGroup {
		return (generator.RandomGenerator{}).Generate(10)
	}
	evalFactory := func(n int) evaluator.Evaluator {
		return evaluator.NewParallel(fitness.TrianglesImageFunctions(imgData, 3, n), 4)
	}
	mutator := mutation.NewGaussianMethod(0.5, 0.3)
	return NewSimple(pointFactory, 10, 2, evalFactory, mutator)
}

func TestModifiedGenetic_InitialGenerationIsZero(t *testing.T) {
	algo := makeSmallAlgo()
	assert.Equal(t, algo.Stats().Generation, 0)
}

func TestModifiedGenetic_StatsGenerationIncrementsEachStep(t *testing.T) {
	algo := makeSmallAlgo()
	for i := 1; i <= 5; i++ {
		algo.Step()
		assert.Equal(t, algo.Stats().Generation, i)
	}
}

func TestModifiedGenetic_BestFitnessIsInValidRange(t *testing.T) {
	algo := makeSmallAlgo()
	algo.Step()
	f := algo.Stats().BestFitness
	assert.True(t, f >= 0.0 && f <= 1.0, "fitness %v not in [0,1]", f)
}

func TestModifiedGenetic_BestPointGroupHasCorrectSize(t *testing.T) {
	algo := makeSmallAlgo()
	algo.Step()
	assert.Equal(t, len(algo.Best()), 10)
}

func TestSimple_StatsGenerationIncrementsEachStep(t *testing.T) {
	algo := makeSmallSimpleAlgo()
	for i := 1; i <= 5; i++ {
		algo.Step()
		assert.Equal(t, algo.Stats().Generation, i)
	}
}

func TestSimple_BestFitnessIsInValidRange(t *testing.T) {
	algo := makeSmallSimpleAlgo()
	algo.Step()
	f := algo.Stats().BestFitness
	assert.True(t, f >= 0.0 && f <= 1.0, "fitness %v not in [0,1]", f)
}

func BenchmarkAlgorithm(b *testing.B) {
	rand.Seed(time.Now().UnixNano())
	random.Seed(time.Now().UnixNano())

	file, err := os.Open("../imgs/clown.jpg")

	if err != nil {
		panic(err)
	}

	imageFile, _, err := image.Decode(file)

	file.Close()

	if err != nil {
		log.Fatal(err)
	}

	imgData := imageData.ToData(imageFile)

	if err != nil {
		log.Fatal("Arg #2 not an integer")
	}

	pointFactory := func() normgeom.NormPointGroup {

		return (generator.RandomGenerator{}).Generate(1000)
	}
	evaluatorFactory := func(n int) evaluator.Evaluator {
		return evaluator.NewParallel(fitness.TrianglesImageFunctions(imgData, 5, n), 22)
	}

	mutator := mutation.NewGaussianMethod(2/1000, 0.3)

	algo := NewModifiedGenetic(pointFactory, 400, 5, evaluatorFactory, mutator)

	real := func() {
		for i := 0; i < 3000; i++ {
			algo.Step()
		}
	}
	real()
}
