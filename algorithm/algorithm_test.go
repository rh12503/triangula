package algorithm

import (
	"github.com/RH12503/Triangula/algorithm/evaluator"
	"github.com/RH12503/Triangula/fitness"
	"github.com/RH12503/Triangula/generator"
	imageData "github.com/RH12503/Triangula/image"
	"github.com/RH12503/Triangula/mutation"
	"github.com/RH12503/Triangula/normgeom"
	"github.com/RH12503/Triangula/random"
	"image"
	"log"
	"math/rand"
	"os"
	"testing"
	"time"
)

func BenchmarkAlgorithm(b *testing.B) {
	rand.Seed(time.Now().UnixNano())
	random.Seed(time.Now().UnixNano())

	file, err := os.Open("imgs/clown.jpg")

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

		return (generator.RandomGenerator{}).Generate(3000)
	}
	evaluatorFactory := func(n int) evaluator.Evaluator {
		return evaluator.NewParallel(fitness.TrianglesImageFunctions(imgData, 5, n), 22)
	}

	mutator := mutation.NewGaussianMethod(0.001, 0.3)

	algo := NewModifiedGenetic(pointFactory, 500, 5, evaluatorFactory, mutator)

	real := func() {
		for i := 0; i < 10000; i++ {
			algo.Step()
		}
	}
	real()
}
