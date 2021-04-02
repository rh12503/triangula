// Package evaluator provides an interface for a fitness evaluator and
// different implementations for common use cases.

package evaluator

import (
	"github.com/RH12503/Triangula/fitness"
	"github.com/RH12503/Triangula/normgeom"
)

// An Evaluator contains an arbitrary number of fitness.Function's and is used by an algorithm.Algorithm
// to evaluate the fitness of its members
type Evaluator interface {
	// Get returns a fitness.Function given the index of a member
	Get(i int) fitness.Function

	// Update does any updates necessary for a specified member
	Update(i int)

	// Prepare does any preparations necessary for the evaluator
	Prepare()

	// SetBase indicates the base of member i
	SetBase(i, base int)

	// Changed should be called if a point has been changed for member i
	Changed(i int, old, new normgeom.NormPoint)

	// Swap should be called if members i and j are swapped
	Swap(i, j int)
}

// one is a generic evaluator for a fitness function that can be shared between many members
type one struct {
	evaluator fitness.Function
}

func (o one) SetBase(i, base int) {
}

func (o one) Changed(i int, old, new normgeom.NormPoint) {
}

func (o one) Swap(i, j int) {
}

func (o one) Prepare() {
}

func (o one) Update(i int) {
}

func (o one) Get(i int) fitness.Function {
	// Return the same fitness.Function every time
	return o.evaluator
}

// NewOne returns a new one fitness evaluator
func NewOne(evaluator fitness.Function) one {
	return one{evaluator: evaluator}
}

// many is a generic evaluator for storing many fitness functions
type many struct {
	evaluators []fitness.Function
}

func (m many) SetBase(i, base int) {
}

func (m many) Changed(i int, old, new normgeom.NormPoint) {
}

func (m many) Swap(i, j int) {
}

func (m many) Prepare() {
}

func (m many) Update(i int) {
}

func (m many) Get(i int) fitness.Function {
	return m.evaluators[i]
}

// NewOne returns a new many fitness evaluator
func NewMany(newEvaluator func() fitness.Function, n int) many {
	var evaluators []fitness.Function

	for i := 0; i < n; i++ {
		evaluators = append(evaluators, newEvaluator())
	}

	return many{evaluators: evaluators}
}
