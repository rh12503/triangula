package mutation

import (
	"Triangula/normgeom"
	"Triangula/random"
)

// gaussianMethod uses gaussian random number while calculating the magnitude of a mutation.
// It typically provides better results that a randomMethod
type gaussianMethod struct {
	rate   float32 // The probability of a point being mutated
	amount float64 // The amount a point's coordinates are changed
}

func (g gaussianMethod) Mutate(points normgeom.NormPointGroup, mutated func(mutation Mutation)) {
	for i := range points {
		if random.Float32() < g.rate {
			old := points[i]

			points[i].X += (random.NormFloat64() - 0.5) * g.amount
			points[i].Y += (random.NormFloat64() - 0.5) * g.amount

			points[i].Constrain()

			mutated(Mutation{
				Old:      old,
				New:      points[i],
				Index: i,
			})
		}
	}
}

// NewGaussianMethod returns a gaussianMethod with specified a mutation rate and amount
func NewGaussianMethod(rate float64, amount float64) gaussianMethod {
	return gaussianMethod{rate: float32(rate), amount: amount}
}