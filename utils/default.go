package utils

import (
	"github.com/RH12503/Triangula/algorithm"
	"github.com/RH12503/Triangula/algorithm/evaluator"
	"github.com/RH12503/Triangula/fitness"
	"github.com/RH12503/Triangula/generator"
	imageData "github.com/RH12503/Triangula/image"
	"github.com/RH12503/Triangula/mutation"
	"github.com/RH12503/Triangula/normgeom"
	"image"
)

// DefaultAlgorithm returns an algorithm than will be optimal for almost all cases
func DefaultAlgorithm(numPoints int, image image.Image) algorithm.Algorithm{
	img := imageData.ToData(image)

	pointFactory := func() normgeom.NormPointGroup {
		return (generator.RandomGenerator{}).Generate(numPoints)
	}

	evaluatorFactory := func(n int) evaluator.Evaluator {
		return evaluator.NewParallel(fitness.TrianglesImageFunctions(img, 5, n), 22)
	}

	var mutator mutation.Method

	mutator = mutation.DefaultGaussianMethod(numPoints)

	algo := algorithm.NewSimple(pointFactory, 400, 5, evaluatorFactory, mutator)
	return algo
}
