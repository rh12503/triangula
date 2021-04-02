package algorithm

import (
	"Triangula/algorithm/evaluator"
	"Triangula/mutation"
	"Triangula/normgeom"
	"github.com/panjf2000/ants/v2"
	"sort"
	"time"
)

// An simple is a simplified genetic algorithm.
// It runs in parallel, scaling across multiple cores.
//
// Firstly, a new generation is chosen based on fitness.
// Secondly, the fitnesses of the new generation are calculated.
// Lastly, the fitnesses are sorted in preparation for the next generation.
type simple struct {
	evaluator evaluator.Evaluator // Used to calculate fitnesses. There is an unique evaluator for each member of the population
	mutator   mutation.Method     // Used in newGeneration to mutate members of the population

	population    []normgeom.NormPointGroup // The population of the algorithm, storing all the data of the generator
	newPopulation []normgeom.NormPointGroup // Used in newGeneration

	fitnesses []FitnessData // fitnesses[i] is the fitness of population[i]

	mutations [][]mutation.Mutation // Stores the mutations made in newGeneration so beneficial mutations can later be combined in combineMutations

	best normgeom.NormPointGroup // The member of the population with the highest fitness

	cutoff int // The number of members that survive to the next generation

	stats Stats // Simple statistics relating to the algorithm
}

func (s *simple) Step() {
	t := time.Now() // Measure the time taken for the generation

	// Fill the population with new members
	s.newGeneration()

	// Calculate and update the fitnesses of the new population

	s.calculateFitnesses()
	s.updateFitnesses()

	s.stats.Generation++
	s.stats.TimeForGen = time.Since(t)
}

// calculateFitnesses calculates the fitnesses of the current generation
func (s *simple) calculateFitnesses() {
	ch := make(chan FitnessData, len(s.population))

	for i, p := range s.population {
		i := i
		p := p
		e := s.evaluator.Get(i)
		ants.Submit(
			func() {
				fit := e.Calculate(p)
				ch <- FitnessData{
					I:       i,
					Fitness: fit,
				}
			},
		)
		s.fitnesses[i].I = i
	}

	s.evaluator.Prepare()

	done := 0
	for d := range ch {
		s.fitnesses[d.I].Fitness = d.Fitness
		s.evaluator.Update(d.I)

		done++
		if done == len(s.population) {
			close(ch)
		}
	}
}

// updateFitnesses prepares the calculated fitnesses for the next generation
func (s *simple) updateFitnesses() {
	sort.Sort(s)

	s.best.Set(s.population[0])
	s.stats.BestFitness = s.fitnesses[0].Fitness
}

// newGeneration populates a generation with new members
func (s *simple) newGeneration() {
	i := 0

	for ; i < s.cutoff; i++ {
		s.newPopulation[i].Set(s.population[i])
	}

	for i < len(s.population) {
		for j := 0; j < s.cutoff && i < len(s.population); j++ {
			s.newPopulation[i].Set(s.population[j])

			s.evaluator.SetBase(i, j)
			s.mutator.Mutate(s.newPopulation[i], func(mutation mutation.Mutation) {
				s.evaluator.Changed(i, mutation.Old, mutation.New)
			})
			i++
		}
	}

	s.population, s.newPopulation = s.newPopulation, s.population
}

func (s simple) Best() normgeom.NormPointGroup {
	return s.best
}

func (s simple) Stats() Stats {
	return s.stats
}

func (s simple) Len() int {
	return len(s.fitnesses)
}

func (s simple) Less(i, j int) bool {
	return s.fitnesses[i].Fitness > s.fitnesses[j].Fitness
}

func (s *simple) Swap(i, j int) {
	s.fitnesses[i], s.fitnesses[j] = s.fitnesses[j], s.fitnesses[i]
	s.population[i], s.population[j] = s.population[j], s.population[i]
	s.evaluator.Swap(i, j)
}

// NewSimple returns a new Simple algorithm
func NewSimple(newPointGroup func() normgeom.NormPointGroup, size int, cutoff int,
	newEvaluators func(n int) evaluator.Evaluator, mutator mutation.Method) *simple {
	var algo simple

	for i := 0; i < size; i++ {
		pg := newPointGroup()
		algo.population = append(algo.population, pg)
		algo.newPopulation = append(algo.newPopulation, pg.Copy())
	}

	algo.best = algo.population[0].Copy()

	algo.evaluator = newEvaluators(size)

	algo.fitnesses = make([]FitnessData, len(algo.population))

	algo.mutator = mutator

	algo.cutoff = cutoff

	algo.calculateFitnesses()

	algo.updateFitnesses()

	return &algo
}

