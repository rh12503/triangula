// Package evaluator provides an interface for a fitness evaluator and
// some implementations for common use cases.
package evaluator

import (
	"github.com/RH12503/Triangula/fitness"
)

// An Evaluator is used by an algorithm to evaluate the fitness of its members.
type Evaluator interface {
	// Get returns a fitness function given the index of a member.
	Get(i int) fitness.Function

	// Update carries out any updates necessary for a specified member.
	Update(i int)

	// Prepare does any preparations necessary for the evaluator.
	Prepare()

	// SetBase should be called to indicate the base of member i.
	SetBase(i, base int)

	// Swap should be called if members i and j are swapped.
	Swap(i, j int)
}

// one is a generic evaluator which has one fitness function for each member.
type one struct {
	evaluator fitness.Function
}

func (o one) SetBase(i, base int) {
}

func (o one) Swap(i, j int) {
}

func (o one) Prepare() {
}

func (o one) Update(i int) {
}

func (o one) Get(i int) fitness.Function {
	// Return the same fitness function every time
	return o.evaluator
}

// NewOne returns a new "one" fitness evaluator.
func NewOne(evaluator fitness.Function) one {
	return one{evaluator: evaluator}
}

// many is a generic evaluator where there is one fitness function for each member.
type many struct {
	evaluators []fitness.Function
}

func (m many) SetBase(i, base int) {
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

// NewMany returns a new "many" fitness evaluator.
func NewMany(newEvaluator func() fitness.Function, n int) many {
	var evaluators []fitness.Function

	for i := 0; i < n; i++ {
		evaluators = append(evaluators, newEvaluator())
	}

	return many{evaluators: evaluators}
}
