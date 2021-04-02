package algorithm

import "Triangula/mutation"

// FitnessData stores the fitness of a point group
type FitnessData struct {
	Fitness float64
	I       int // The index of the point group
}

// MutationsData stores a set of mutations and the indexes of the member which each mutation came from
type MutationsData struct {
	Mutations []mutation.Mutation
	Indexes   []int // The indexes of the point group where the mutation came from
}

func (m *MutationsData) Clear() {
	m.Mutations = m.Mutations[:0]
	m.Indexes = m.Indexes[:0]
}

// Count returns the number of mutations stored
func (m MutationsData) Count() int {
	return len(m.Mutations)
}
