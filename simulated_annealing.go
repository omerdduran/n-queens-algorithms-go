package main

import (
	"fmt"
	"math"
	"math/rand"
)

// SimulatedAnnealingSolver implements simulated annealing search
type SimulatedAnnealingSolver struct {
	n             int
	board         []int
	solution      []int
	solved        bool
	initialTemp   float64
	coolingRate   float64
	minTemp       float64
	maxIterations int
	restarts      int
}

// NewSimulatedAnnealingSolver creates a new simulated annealing solver
func NewSimulatedAnnealingSolver(n int) *SimulatedAnnealingSolver {
	return &SimulatedAnnealingSolver{
		n:             n,
		board:         make([]int, n),
		initialTemp:   float64(n * n), // Scale with problem size
		coolingRate:   0.99,           // Slower cooling for better exploration
		minTemp:       0.01,
		maxIterations: n * 1000, // Scale iterations with problem size
		restarts:      5,        // Multiple restarts for better success rate
	}
}

// Solve attempts to find a solution using simulated annealing with restarts
func (sa *SimulatedAnnealingSolver) Solve() bool {
	for restart := 0; restart < sa.restarts; restart++ {
		if sa.singleRun() {
			return true
		}
	}
	return false
}

// singleRun performs one complete simulated annealing run
func (sa *SimulatedAnnealingSolver) singleRun() bool {
	// Initialize with better starting position
	sa.smartInit()

	temperature := sa.initialTemp
	currentCost := sa.calculateCost()
	bestCost := currentCost
	bestBoard := make([]int, sa.n)
	copy(bestBoard, sa.board)

	for iter := 0; iter < sa.maxIterations && temperature > sa.minTemp; iter++ {
		if currentCost == 0 {
			sa.solution = make([]int, sa.n)
			copy(sa.solution, sa.board)
			sa.solved = true
			return true
		}

		// Generate better neighbor
		neighbor := sa.generateSmartNeighbor()
		neighborCost := sa.calculateCostForBoard(neighbor)

		// Calculate cost difference
		deltaCost := neighborCost - currentCost

		// Accept or reject the neighbor
		if deltaCost <= 0 || sa.acceptanceProbability(deltaCost, temperature) > rand.Float64() {
			copy(sa.board, neighbor)
			currentCost = neighborCost

			// Track best solution found
			if currentCost < bestCost {
				bestCost = currentCost
				copy(bestBoard, sa.board)
			}
		}

		// Adaptive cooling - slow down when making progress
		if iter%100 == 0 && currentCost > bestCost*2 {
			// Restart from best known position if we're doing poorly
			copy(sa.board, bestBoard)
			currentCost = bestCost
			temperature = sa.initialTemp * 0.5 // Restart with lower temperature
		} else {
			temperature *= sa.coolingRate
		}
	}

	// Check if we found a solution
	copy(sa.board, bestBoard)
	if bestCost == 0 {
		sa.solution = make([]int, sa.n)
		copy(sa.solution, bestBoard)
		sa.solved = true
		return true
	}

	return false
}

// smartInit initializes the board with a better starting position
func (sa *SimulatedAnnealingSolver) smartInit() {
	// Always start with permutation (one queen per row)
	perm := rand.Perm(sa.n)
	copy(sa.board, perm)
}

// generateSmartNeighbor creates a neighbor with more intelligent strategies
func (sa *SimulatedAnnealingSolver) generateSmartNeighbor() []int {
	neighbor := make([]int, sa.n)
	copy(neighbor, sa.board)

	strategy := rand.Float64()

	if strategy < 0.6 {
		// Strategy 1: Swap two random queens (most effective for permutations)
		pos1 := rand.Intn(sa.n)
		pos2 := rand.Intn(sa.n)
		for pos1 == pos2 {
			pos2 = rand.Intn(sa.n)
		}
		neighbor[pos1], neighbor[pos2] = neighbor[pos2], neighbor[pos1]
	} else if strategy < 0.8 {
		// Strategy 2: Move a conflicted queen to a better position
		conflictedQueens := sa.findConflictedQueens()
		if len(conflictedQueens) > 0 {
			col := conflictedQueens[rand.Intn(len(conflictedQueens))]
			// Try to find a less conflicted row
			bestRow := rand.Intn(sa.n)
			minConflicts := sa.n * sa.n

			for row := 0; row < sa.n; row++ {
				if row != neighbor[col] {
					tempBoard := make([]int, sa.n)
					copy(tempBoard, neighbor)
					tempBoard[col] = row
					conflicts := sa.calculateConflictsForPosition(tempBoard, col)
					if conflicts < minConflicts {
						minConflicts = conflicts
						bestRow = row
					}
				}
			}
			neighbor[col] = bestRow
		}
	} else {
		// Strategy 3: Local search - try to improve a random position
		col := rand.Intn(sa.n)
		bestRow := neighbor[col]
		minConflicts := sa.calculateConflictsForPosition(neighbor, col)

		for row := 0; row < sa.n; row++ {
			if row != neighbor[col] {
				tempBoard := make([]int, sa.n)
				copy(tempBoard, neighbor)
				tempBoard[col] = row
				conflicts := sa.calculateConflictsForPosition(tempBoard, col)
				if conflicts < minConflicts {
					minConflicts = conflicts
					bestRow = row
				}
			}
		}
		neighbor[col] = bestRow
	}

	return neighbor
}

// findConflictedQueens returns a list of column indices for queens that are in conflict
func (sa *SimulatedAnnealingSolver) findConflictedQueens() []int {
	var conflicted []int
	for i := 0; i < sa.n; i++ {
		if sa.calculateConflictsForPosition(sa.board, i) > 0 {
			conflicted = append(conflicted, i)
		}
	}
	return conflicted
}

// calculateConflictsForPosition calculates conflicts for a queen at a specific position
func (sa *SimulatedAnnealingSolver) calculateConflictsForPosition(board []int, col int) int {
	conflicts := 0
	for j := 0; j < sa.n; j++ {
		if j != col {
			// Check row conflict
			if board[col] == board[j] {
				conflicts++
			}
			// Check diagonal conflict
			if abs(board[col]-board[j]) == abs(col-j) {
				conflicts++
			}
		}
	}
	return conflicts
}

// calculateCost calculates the cost (number of conflicts) for current board
func (sa *SimulatedAnnealingSolver) calculateCost() int {
	return sa.calculateCostForBoard(sa.board)
}

// calculateCostForBoard calculates the cost for a given board configuration
func (sa *SimulatedAnnealingSolver) calculateCostForBoard(board []int) int {
	conflicts := 0
	for i := 0; i < sa.n; i++ {
		for j := i + 1; j < sa.n; j++ {
			// Check row conflict
			if board[i] == board[j] {
				conflicts++
			}
			// Check diagonal conflict
			if abs(board[i]-board[j]) == abs(i-j) {
				conflicts++
			}
		}
	}
	return conflicts
}

// acceptanceProbability calculates the probability of accepting a worse solution
func (sa *SimulatedAnnealingSolver) acceptanceProbability(deltaCost int, temperature float64) float64 {
	return math.Exp(-float64(deltaCost) / temperature)
}

// GetSolution returns the found solution
func (sa *SimulatedAnnealingSolver) GetSolution() []int {
	return sa.solution
}

// PrintSolution prints the solution board
func (sa *SimulatedAnnealingSolver) PrintSolution() {
	if !sa.solved {
		fmt.Println("No solution found")
		return
	}

	fmt.Printf("Simulated Annealing Solution for N=%d:\n", sa.n)
	for i := 0; i < sa.n; i++ {
		for j := 0; j < sa.n; j++ {
			if sa.solution[i] == j {
				fmt.Print("Q ")
			} else {
				fmt.Print(". ")
			}
		}
		fmt.Println()
	}
	fmt.Println()
}
