package evaluator

import (
	"github.com/RH12503/Triangula/fitness"
)

// parallel is a fitness evaluator that supports parallel calculations.
// It stores and updates a cache, and contains a fitness.TrianglesImageFunction for each member
// to calculate fitnesses.
type parallel struct {
	evaluators []fitness.CacheFunction

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
	for _, d := range eval.Cache() {
		p.cache[d.OtherHash] = d
	}

	eval.SetCache(p.cache)
}

func (p *parallel) SetBase(i, base int) {
	p.evaluators[i].SetBase(p.evaluators[base])
}

func (p *parallel) Swap(i, j int) {
	p.evaluators[i], p.evaluators[j] = p.evaluators[j], p.evaluators[i]
}

// NewParallel creates a new parallel evaluator.
func NewParallel(fitnessFuncs []fitness.CacheFunction, cachePowerOf2 int) *parallel {

	return &parallel{
		evaluators: fitnessFuncs,
		cache:      make([]fitness.TriFit, 1<<cachePowerOf2),
		nextCache:  make([]fitness.TriFit, 1<<cachePowerOf2),
	}
}
