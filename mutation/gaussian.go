package mutation

import (
	"github.com/RH12503/Triangula/normgeom"
	"github.com/RH12503/Triangula/random"
)

// gaussianMethod uses a gaussian random number while calculating the magnitude of a mutation.
// It typically provides better results than a randomMethod.
type gaussianMethod struct {
	rate   float32 // The probability of a point being mutated.
	amount float64 // The amount a point's coordinates are changed.
}

func (g gaussianMethod) Mutate(points normgeom.NormPointGroup, mutated func(mutation Mutation)) {
	for i := range points {
		if random.Float32() < g.rate {
			old := points[i]

			points[i].X += random.NormFloat64() * g.amount * 0.5
			points[i].Y += random.NormFloat64() * g.amount * 0.5

			points[i].Constrain()

			mutated(Mutation{
				Old:   old,
				New:   points[i],
				Index: i,
			})
		}
	}
}

// NewGaussianMethod returns a gaussianMethod with specified a mutation rate and amount.
func NewGaussianMethod(rate float64, amount float64) gaussianMethod {
	return gaussianMethod{rate: float32(rate), amount: amount}
}

// DefaultGaussianMethod returns a gaussianMethod which is optimal for almost all cases.
func DefaultGaussianMethod(numPoints int) gaussianMethod {
	return gaussianMethod{rate: 2 / float32(numPoints), amount: 0.3}
}
