package mutation

import (
	"github.com/RH12503/Triangula/normgeom"
	"math/rand"
)

// randomMethod is a simple implementation of a Method.
type randomMethod struct {
	rate   float64 // The probability of a point being mutated.
	amount float64 // The amount a point's coordinates are changed.
}

func (r randomMethod) Mutate(points normgeom.NormPointGroup, mutated func(mutation Mutation)) {
	for i := range points {
		if rand.Float64() < r.rate {
			old := points[i]
			points[i].X += (rand.Float64() - 0.5) * r.amount
			points[i].Y += (rand.Float64() - 0.5) * r.amount
			points[i].Constrain()
			mutated(Mutation{
				Old:   old,
				New:   points[i],
				Index: i,
			})
		}
	}
}

// NewRandomMethod returns a randomMethod with specified a mutation rate and amount.
func NewRandomMethod(rate float64, amount float64) randomMethod {
	return randomMethod{rate: rate, amount: amount}
}
