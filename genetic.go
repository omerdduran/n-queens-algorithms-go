package main

import (
	"fmt"
	"math/rand"
	"sort"
)

// Individual represents a chromosome in the genetic algorithm
type Individual struct {
	chromosome []int
	fitness    int
}

// GeneticSolver implements genetic algorithm for N-Queens
type GeneticSolver struct {
	n              int
	populationSize int
	maxGenerations int
	mutationRate   float64
	crossoverRate  float64
	population     []Individual
	solution       []int
	solved         bool
	restarts       int
}

// NewGeneticSolver creates a new genetic algorithm solver
func NewGeneticSolver(n int) *GeneticSolver {
	// Balanced parameters for success rate and speed
	popSize := 80
	if n > 20 {
		popSize = 120
	}
	if n > 40 {
		popSize = 150
	}

	return &GeneticSolver{
		n:              n,
		populationSize: popSize,
		maxGenerations: 200,  // More generations for better results
		mutationRate:   0.15, // Balanced mutation rate
		crossoverRate:  0.85, // Higher crossover rate
		population:     make([]Individual, popSize),
		restarts:       5, // More restarts for much better success
	}
}

// Solve attempts to find a solution using genetic algorithm with restarts
func (ga *GeneticSolver) Solve() bool {
	// Multiple runs for better success rate
	for restart := 0; restart < ga.restarts; restart++ {
		if ga.singleRun() {
			return true
		}
	}
	return false
}

// singleRun performs one complete genetic algorithm run
func (ga *GeneticSolver) singleRun() bool {
	// Initialize population
	ga.initializePopulation()

	generationsWithoutImprovement := 0
	bestFitnessEver := ga.n * ga.n

	for generation := 0; generation < ga.maxGenerations; generation++ {
		// Evaluate fitness for all individuals
		ga.evaluatePopulation()

		// Check if we found a solution
		if ga.population[0].fitness == 0 {
			ga.solution = make([]int, ga.n)
			copy(ga.solution, ga.population[0].chromosome)
			ga.solved = true
			return true
		}

		// Track progress
		currentBest := ga.population[0].fitness
		if currentBest < bestFitnessEver {
			bestFitnessEver = currentBest
			generationsWithoutImprovement = 0
		} else {
			generationsWithoutImprovement++
		}

		// Early termination if stuck too long
		if generationsWithoutImprovement > 50 {
			break
		}

		// Adaptive parameters - more sophisticated approach
		if generationsWithoutImprovement > 20 {
			ga.mutationRate = 0.3 // Higher mutation
			// Add some new random individuals for diversity
			for i := ga.populationSize * 4 / 5; i < ga.populationSize; i++ {
				ga.initializeIndividual(i)
			}
		} else {
			ga.mutationRate = 0.15
		}

		// Create new generation
		newPopulation := ga.createNewGeneration()
		ga.population = newPopulation
	}

	// Check final generation
	ga.evaluatePopulation()
	if ga.population[0].fitness == 0 {
		ga.solution = make([]int, ga.n)
		copy(ga.solution, ga.population[0].chromosome)
		ga.solved = true
		return true
	}

	return false
}

// initializePopulation creates the initial population with better diversity
func (ga *GeneticSolver) initializePopulation() {
	for i := 0; i < ga.populationSize; i++ {
		ga.initializeIndividual(i)
	}
}

// initializeIndividual creates a single individual with permutation-based initialization
func (ga *GeneticSolver) initializeIndividual(index int) {
	chromosome := make([]int, ga.n)

	// Permutation initialization (one queen per row) - most effective for N-Queens
	perm := rand.Perm(ga.n)
	copy(chromosome, perm)

	ga.population[index] = Individual{
		chromosome: chromosome,
		fitness:    0,
	}
}

// evaluatePopulation calculates fitness for all individuals and sorts them
func (ga *GeneticSolver) evaluatePopulation() {
	for i := 0; i < ga.populationSize; i++ {
		ga.population[i].fitness = ga.calculateFitness(ga.population[i].chromosome)
	}

	// Sort population by fitness (ascending - lower is better)
	sort.Slice(ga.population, func(i, j int) bool {
		return ga.population[i].fitness < ga.population[j].fitness
	})
}

// calculateFitness calculates the fitness (number of conflicts) for a chromosome
func (ga *GeneticSolver) calculateFitness(chromosome []int) int {
	conflicts := 0
	for i := 0; i < ga.n; i++ {
		for j := i + 1; j < ga.n; j++ {
			// Check row conflict
			if chromosome[i] == chromosome[j] {
				conflicts++
			}
			// Check diagonal conflict
			if abs(chromosome[i]-chromosome[j]) == abs(i-j) {
				conflicts++
			}
		}
	}
	return conflicts
}

// createNewGeneration creates a new generation through selection, crossover, and mutation
func (ga *GeneticSolver) createNewGeneration() []Individual {
	newPopulation := make([]Individual, ga.populationSize)

	// Elite preservation - balanced approach
	eliteSize := ga.populationSize / 10
	if eliteSize < 2 {
		eliteSize = 2
	}
	if eliteSize > 10 {
		eliteSize = 10
	}
	for i := 0; i < eliteSize; i++ {
		newPopulation[i] = Individual{
			chromosome: make([]int, ga.n),
			fitness:    ga.population[i].fitness,
		}
		copy(newPopulation[i].chromosome, ga.population[i].chromosome)
	}

	// Generate rest of the population
	for i := eliteSize; i < ga.populationSize; i++ {
		if rand.Float64() < ga.crossoverRate {
			// Crossover
			parent1 := ga.tournamentSelection()
			parent2 := ga.tournamentSelection()
			child := ga.smartCrossover(parent1, parent2)

			// Mutation
			if rand.Float64() < ga.mutationRate {
				ga.smartMutation(child)
			}

			newPopulation[i] = Individual{
				chromosome: child,
				fitness:    0,
			}
		} else {
			// Direct selection with possible mutation
			parent := ga.tournamentSelection()
			child := make([]int, ga.n)
			copy(child, parent.chromosome)

			if rand.Float64() < ga.mutationRate {
				ga.smartMutation(child)
			}

			newPopulation[i] = Individual{
				chromosome: child,
				fitness:    0,
			}
		}
	}

	return newPopulation
}

// tournamentSelection selects an individual using tournament selection
func (ga *GeneticSolver) tournamentSelection() Individual {
	tournamentSize := 5 // Balanced tournament size
	if tournamentSize > ga.populationSize {
		tournamentSize = ga.populationSize
	}

	best := ga.population[rand.Intn(ga.populationSize)]
	for i := 1; i < tournamentSize; i++ {
		candidate := ga.population[rand.Intn(ga.populationSize)]
		if candidate.fitness < best.fitness {
			best = candidate
		}
	}

	return best
}

// smartCrossover performs Order Crossover (OX) - more suitable for N-Queens
func (ga *GeneticSolver) smartCrossover(parent1, parent2 Individual) []int {
	return ga.orderCrossover(parent1.chromosome, parent2.chromosome)
}

// orderCrossover implements Order Crossover (OX)
func (ga *GeneticSolver) orderCrossover(parent1, parent2 []int) []int {
	child := make([]int, ga.n)

	// Select a random segment from parent1
	start := rand.Intn(ga.n)
	end := rand.Intn(ga.n)
	if start > end {
		start, end = end, start
	}

	// Copy the segment from parent1
	used := make(map[int]bool)
	for i := start; i <= end; i++ {
		child[i] = parent1[i]
		used[parent1[i]] = true
	}

	// Fill remaining positions with parent2's order
	childIndex := (end + 1) % ga.n
	for i := 0; i < ga.n; i++ {
		parent2Index := (end + 1 + i) % ga.n
		if !used[parent2[parent2Index]] {
			child[childIndex] = parent2[parent2Index]
			childIndex = (childIndex + 1) % ga.n
		}
	}

	return child
}

// smartMutation performs effective mutation
func (ga *GeneticSolver) smartMutation(chromosome []int) {
	strategy := rand.Float64()

	if strategy < 0.5 {
		// Swap mutation (good for permutations)
		pos1 := rand.Intn(ga.n)
		pos2 := rand.Intn(ga.n)
		chromosome[pos1], chromosome[pos2] = chromosome[pos2], chromosome[pos1]
	} else if strategy < 0.8 {
		// Smart mutation - move a conflicted queen
		conflicts := make([]int, 0, ga.n)
		for i := 0; i < ga.n; i++ {
			if ga.calculateConflictsForPosition(chromosome, i) > 0 {
				conflicts = append(conflicts, i)
			}
		}

		if len(conflicts) > 0 {
			col := conflicts[rand.Intn(len(conflicts))]
			// Try a few random positions and pick the best
			bestRow := chromosome[col]
			minConflicts := ga.calculateConflictsForPosition(chromosome, col)

			for attempts := 0; attempts < 3; attempts++ {
				testRow := rand.Intn(ga.n)
				if testRow != chromosome[col] {
					originalRow := chromosome[col]
					chromosome[col] = testRow
					conflicts := ga.calculateConflictsForPosition(chromosome, col)
					if conflicts < minConflicts {
						minConflicts = conflicts
						bestRow = testRow
					}
					chromosome[col] = originalRow
				}
			}
			chromosome[col] = bestRow
		} else {
			// If no conflicts, random mutation
			pos := rand.Intn(ga.n)
			chromosome[pos] = rand.Intn(ga.n)
		}
	} else {
		// Random mutation
		pos := rand.Intn(ga.n)
		chromosome[pos] = rand.Intn(ga.n)
	}
}

// calculateConflictsForPosition calculates conflicts for a queen at a specific position
func (ga *GeneticSolver) calculateConflictsForPosition(chromosome []int, col int) int {
	conflicts := 0
	for j := 0; j < ga.n; j++ {
		if j != col {
			// Check row conflict
			if chromosome[col] == chromosome[j] {
				conflicts++
			}
			// Check diagonal conflict
			if abs(chromosome[col]-chromosome[j]) == abs(col-j) {
				conflicts++
			}
		}
	}
	return conflicts
}

// GetSolution returns the found solution
func (ga *GeneticSolver) GetSolution() []int {
	return ga.solution
}

// PrintSolution prints the solution board
func (ga *GeneticSolver) PrintSolution() {
	if !ga.solved {
		fmt.Println("No solution found")
		return
	}

	fmt.Printf("Genetic Algorithm Solution for N=%d:\n", ga.n)
	for i := 0; i < ga.n; i++ {
		for j := 0; j < ga.n; j++ {
			if ga.solution[i] == j {
				fmt.Print("Q ")
			} else {
				fmt.Print(". ")
			}
		}
		fmt.Println()
	}
	fmt.Println()
}
