package algorithm

import "github.com/RH12503/Triangula/mutation"

// FitnessData stores the fitness of a point group.
type FitnessData struct {
	Fitness float64
	I       int // The index of the point group in an Algorithm.
}

// MutationsData stores a set of mutations and the index of the member which each mutation came from.
type MutationsData struct {
	Mutations []mutation.Mutation
	Indexes   []int // The index of the member where the mutation came from. Indexes[i] has the mutation Mutations[i].
}

// Clear clears all data from the struct.
func (m *MutationsData) Clear() {
	m.Mutations = m.Mutations[:0]
	m.Indexes = m.Indexes[:0]
}

// Count returns the number of mutations stored.
func (m MutationsData) Count() int {
	return len(m.Mutations)
}
