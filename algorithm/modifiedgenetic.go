package algorithm

import (
	"github.com/RH12503/Triangula/algorithm/evaluator"
	"github.com/RH12503/Triangula/mutation"
	"github.com/RH12503/Triangula/normgeom"
	"github.com/panjf2000/ants/v2"
	"sort"
	"time"
)

// An modifiedGenetic is a genetic algorithm with a different crossover process.
// It runs in parallel, scaling across multiple cores.
//
// Simply put, each generation a set number of members of the population survive. (defined by cutoff)
// These members are called "bases", and
// all other members are mutations are mutated versions of the bases.
// At the end, for each base, a member is created combining all beneficial mutations for that base.
// For a full explanation on how the algorithm works see: <add>
type modifiedGenetic struct {
	evaluator evaluator.Evaluator // Used to calculate fitnesses. There is an unique evaluator for each member of the population
	mutator   mutation.Method     // Used in newGeneration to mutate members of the population

	population    []normgeom.NormPointGroup // The population of the algorithm, storing all the data of the generator
	newPopulation []normgeom.NormPointGroup // Used in newGeneration

	fitnesses []FitnessData // fitnesses[i] is the fitness of population[i]

	mutations [][]mutation.Mutation // Stores the mutations made in newGeneration so beneficial mutations can later be combined in combineMutations

	beneficialMutations []MutationsData // Stores the beneficial mutations for each base

	best normgeom.NormPointGroup // The member of the population with the highest fitness

	cutoff int // The number of members that survive to the next generation

	stats Stats // Simple statistics relating to the algorithm
}

func (g *modifiedGenetic) Step() {
	t := time.Now() // Measure the time taken for the generation

	// Fill the population with new members
	g.newGeneration()

	// Calculate the fitnesses of the new population and combine all beneficial mutations made

	g.calculateFitnesses()
	g.combineMutations()

	// Update fitnesses in preparation for the next generation
	g.updateFitnesses()

	g.stats.Generation++
	g.stats.TimeForGen = time.Since(t)
}

// newGeneration populates a generation with new members
func (g *modifiedGenetic) newGeneration() {
	i := 0

	// newPopulation is filled with new members, and then is swapped with population (the old population)

	// The bases are guaranteed to survive without any mutations
	for ; i < g.cutoff; i++ {
		g.newPopulation[i].Set(g.population[i])
		g.mutations[i] = g.mutations[i][:0]
	}

	// Fill the rest of the algorithm with mutated versions of the bases
	//
	// The population is only filled to len(g.population)-g.cutoff because there needs to be a member
	// At the end for each base in order to combine beneficial mutations
	for i < len(g.population)-g.cutoff {
		// i = the member, j = the base of the member
		for j := 0; j < g.cutoff && i < len(g.population)-g.cutoff; j++ {
			g.mutations[i] = g.mutations[i][:0] // clear all previous mutations
			g.newPopulation[i].Set(g.population[j])

			// The evaluator need to know the base of each member and any mutations made

			g.evaluator.SetBase(i, j)
			g.mutator.Mutate(g.newPopulation[i], func(mut mutation.Mutation) {
				g.evaluator.Changed(i, mut.Old, mut.New)
				g.mutations[i] = append(g.mutations[i], mut)
			})
			i++
		}
	}

	for i := range g.beneficialMutations {
		g.beneficialMutations[i].Clear() // Clear all previous beneficial mutations
	}

	g.population, g.newPopulation = g.newPopulation, g.population
}

// calculateFitnesses calculates the fitnesses of the current generation
func (g *modifiedGenetic) calculateFitnesses() {
	ch := make(chan FitnessData, len(g.population)) // Buffered channel for performance

	for i := 0; i < len(g.population)-g.cutoff; i++ {
		i := i
		p := g.population[i]
		e := g.evaluator.Get(i)
		// Workers calculate the fitness of each member
		ants.Submit(
			func() {
				fit := e.Calculate(p)
				ch <- FitnessData{
					I:       i,
					Fitness: fit,
				}
			},
		)
		g.fitnesses[i].I = i // Assign an index to each fitness so it can be found after being sorted
	}

	g.evaluator.Prepare()

	done := 0
	for d := range ch {
		g.fitnesses[d.I].Fitness = d.Fitness
		g.evaluator.Update(d.I)

		// If the new fitness of a member is higher than its base, that means its mutations were beneficial
		if d.Fitness > g.fitnesses[g.getBase(d.I)].Fitness {
			g.setBeneficial(d.I)
		}

		done++
		if done == len(g.population)-g.cutoff { // Wait till all the fitnesses are calculated
			close(ch)
		}
	}

}

// setBeneficial adds the mutations of population[index] to beneficialMutations
func (g *modifiedGenetic) setBeneficial(index int) {
	base := g.getBase(index)

	for _, m := range g.mutations[index] {
		// Check if a mutation already exists for a given point (we can't have 2 mutations on one point)
		found := false
		foundIndex := -1
		for i, o := range g.beneficialMutations[base].Mutations {
			if m.Index == o.Index {
				found = true
				foundIndex = i
				break
			}
		}
		if !found {
			// If there are no existing mutation, add it to the list of beneficial mutations
			g.beneficialMutations[base].Mutations = append(g.beneficialMutations[base].Mutations, m)
			g.beneficialMutations[base].Indexes = append(g.beneficialMutations[base].Indexes, index)
		} else {
			// If there is a duplicate mutation, check to see which one is more beneficial, and replace the other
			// with this one if this one is more beneficial
			other := g.beneficialMutations[base].Indexes[foundIndex]
			if g.fitnesses[index].Fitness > g.fitnesses[other].Fitness {
				g.beneficialMutations[base].Mutations[foundIndex] = m
				g.beneficialMutations[base].Indexes[foundIndex] = index
			}
		}
	}
}

// combineMutations combines all the beneficial mutations found together for each base
func (g *modifiedGenetic) combineMutations() {
	// The members with the combined mutations are at the end
	for i := len(g.population) - g.cutoff; i < len(g.population); i++ {
		base := g.getBase(i)
		if g.beneficialMutations[base].Count() > 0 {
			// If there are any beneficial mutations, set the member to its base and perform all the mutations
			g.population[i].Set(g.population[base])
			g.evaluator.SetBase(i, base)

			for _, m := range g.beneficialMutations[base].Mutations {
				g.evaluator.Changed(i, m.Old, m.New)
				g.population[i][m.Index].X = m.New.X
				g.population[i][m.Index].Y = m.New.Y
			}

			g.fitnesses[i].I = i

			// Calculate the fitness of the new member
			e := g.evaluator.Get(i)
			fit := e.Calculate(g.population[i])
			g.fitnesses[i].Fitness = fit
			g.evaluator.Update(i)
		} else {
			// If there aren't any beneficial mutations, leave the member out of the next generation by setting
			// its fitness to 0
			g.fitnesses[i].Fitness = 0
		}
	}
}

// updateFitnesses prepares the calculated fitnesses for the next generation
func (g *modifiedGenetic) updateFitnesses() {
	// Sort the population by fitness with g.population[0] being the best
	sort.Sort(g)

	g.best.Set(g.population[0])
	g.stats.BestFitness = g.fitnesses[0].Fitness
}

// getBase returns the base of a member given the index of that member
func (g modifiedGenetic) getBase(index int) int {
	return index % g.cutoff
}

func (g modifiedGenetic) Best() normgeom.NormPointGroup {
	return g.best
}

func (g modifiedGenetic) Stats() Stats {
	return g.stats
}

// Functions for sorting

func (g modifiedGenetic) Len() int {
	return len(g.fitnesses)
}

func (g modifiedGenetic) Less(i, j int) bool {
	return g.fitnesses[i].Fitness > g.fitnesses[j].Fitness
}

func (g *modifiedGenetic) Swap(i, j int) {
	g.fitnesses[i], g.fitnesses[j] = g.fitnesses[j], g.fitnesses[i]
	g.population[i], g.population[j] = g.population[j], g.population[i]
	g.evaluator.Swap(i, j)
}

// NewModifiedGenetic returns a new modifiedGenetic algorithm
func NewModifiedGenetic(newPointGroup func() normgeom.NormPointGroup, size int, cutoff int,
	newEvaluators func(n int) evaluator.Evaluator, mutator mutation.Method) *modifiedGenetic {

	var algo modifiedGenetic

	// Fill the population with point groups
	for i := 0; i < size; i++ {
		pg := newPointGroup()
		algo.population = append(algo.population, pg)
		algo.newPopulation = append(algo.newPopulation, pg.Copy())
	}

	algo.best = algo.population[0].Copy() // Set a random best member to start off

	algo.evaluator = newEvaluators(size)

	algo.fitnesses = make([]FitnessData, len(algo.population))

	algo.mutations = make([][]mutation.Mutation, len(algo.population))
	algo.beneficialMutations = make([]MutationsData, cutoff)

	algo.mutator = mutator

	algo.cutoff = cutoff

	// Calculate and update fitnesses in preparation for the first generation
	algo.calculateFitnesses()
	algo.updateFitnesses()

	return &algo
}
