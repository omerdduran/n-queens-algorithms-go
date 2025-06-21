package main

import (
	"fmt"
	"math/rand"
)

// GreedySolver implements hill climbing greedy search
type GreedySolver struct {
	n             int
	board         []int
	solution      []int
	solved        bool
	maxIterations int
}

// NewGreedySolver creates a new greedy solver
func NewGreedySolver(n int) *GreedySolver {
	return &GreedySolver{
		n:             n,
		board:         make([]int, n),
		maxIterations: 10000, // Prevent infinite loops
	}
}

// Solve attempts to find a solution using hill climbing
func (g *GreedySolver) Solve() bool {
	// Initialize with random positions
	g.randomInit()

	for iter := 0; iter < g.maxIterations; iter++ {
		conflicts := g.countConflicts()
		if conflicts == 0 {
			g.solution = make([]int, g.n)
			copy(g.solution, g.board)
			g.solved = true
			return true
		}

		// Find the best neighbor
		bestBoard := make([]int, g.n)
		copy(bestBoard, g.board)
		bestConflicts := conflicts

		for col := 0; col < g.n; col++ {
			for newRow := 0; newRow < g.n; newRow++ {
				if g.board[col] == newRow {
					continue
				}

				// Try moving queen in column col to row newRow
				originalRow := g.board[col]
				g.board[col] = newRow

				newConflicts := g.countConflicts()
				if newConflicts < bestConflicts {
					bestConflicts = newConflicts
					copy(bestBoard, g.board)
				}

				// Restore original position
				g.board[col] = originalRow
			}
		}

		// If no improvement found, restart with random configuration
		if bestConflicts >= conflicts {
			g.randomInit()
		} else {
			copy(g.board, bestBoard)
		}
	}

	return false
}

// randomInit initializes the board with random queen positions
func (g *GreedySolver) randomInit() {
	for i := 0; i < g.n; i++ {
		g.board[i] = rand.Intn(g.n)
	}
}

// countConflicts counts the total number of conflicts on the board
func (g *GreedySolver) countConflicts() int {
	conflicts := 0
	for i := 0; i < g.n; i++ {
		for j := i + 1; j < g.n; j++ {
			// Check row conflict
			if g.board[i] == g.board[j] {
				conflicts++
			}
			// Check diagonal conflict
			if abs(g.board[i]-g.board[j]) == abs(i-j) {
				conflicts++
			}
		}
	}
	return conflicts
}

// GetSolution returns the found solution
func (g *GreedySolver) GetSolution() []int {
	return g.solution
}

// PrintSolution prints the solution board
func (g *GreedySolver) PrintSolution() {
	if !g.solved {
		fmt.Println("No solution found")
		return
	}

	fmt.Printf("Greedy Search Solution for N=%d:\n", g.n)
	for i := 0; i < g.n; i++ {
		for j := 0; j < g.n; j++ {
			if g.solution[i] == j {
				fmt.Print("Q ")
			} else {
				fmt.Print(". ")
			}
		}
		fmt.Println()
	}
	fmt.Println()
}
