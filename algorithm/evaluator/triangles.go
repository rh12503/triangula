package evaluator

import (
	"github.com/RH12503/Triangula/fitness"
	"github.com/RH12503/Triangula/image"
)

// parallel is a fitness evaluator that supports parallel calculations.
// It stores and updates a cache, and contains a fitness.TrianglesImageEvaluator for each member
// to calculate fitnesses.
type parallel struct {
	evaluators []*fitness.TrianglesImageEvaluator

	cache     []fitness.TriFit // The current triangle cache being used by the fitness functions.
	nextCache []fitness.TriFit // The cache for the next generation.
}

func (p parallel) Get(i int) fitness.Function {
	return p.evaluators[i]
}

func (p *parallel) Prepare() {
	p.cache, p.nextCache = p.nextCache, p.cache
}

func (p *parallel) Update(i int) {
	eval := p.evaluators[i]

	// Put triangles that have been calculated from the fitness function into the cache
	for _, d := range eval.TriangleCache {
		p.cache[d.OtherHash] = d
	}

	eval.TriangleCache = p.cache
}

func (p *parallel) SetBase(i, base int) {
	p.evaluators[i].Base = p.evaluators[base].Triangulation
}

func (p *parallel) Swap(i, j int) {
	p.evaluators[i], p.evaluators[j] = p.evaluators[j], p.evaluators[i]
}

// NewParallel creates a new parallel evaluator.
func NewParallel(img image.Data, cachePowerOf2, blockSize, n int) *parallel {
	return &parallel{
		evaluators: fitness.TrianglesImageEvaluators(img, blockSize, n),
		cache:      make([]fitness.TriFit, 1<<cachePowerOf2),
		nextCache:  make([]fitness.TriFit, 1<<cachePowerOf2),
	}
}
