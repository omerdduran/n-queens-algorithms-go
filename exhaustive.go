package main

import "fmt"

// ExhaustiveSearchSolver implements depth-first search with backtracking
type ExhaustiveSearchSolver struct {
	n             int
	board         []int
	solution      []int
	solutionFound bool
}

// NewExhaustiveSearchSolver creates a new exhaustive search solver
func NewExhaustiveSearchSolver(n int) *ExhaustiveSearchSolver {
	return &ExhaustiveSearchSolver{
		n:     n,
		board: make([]int, n),
	}
}

// Solve attempts to find a solution using exhaustive depth-first search
func (e *ExhaustiveSearchSolver) Solve() bool {
	e.solutionFound = false
	e.solveRecursive(0)
	return e.solutionFound
}

// solveRecursive implements the recursive backtracking algorithm
func (e *ExhaustiveSearchSolver) solveRecursive(row int) {
	if e.solutionFound {
		return
	}

	if row == e.n {
		// Found a solution
		e.solution = make([]int, e.n)
		copy(e.solution, e.board)
		e.solutionFound = true
		return
	}

	for col := 0; col < e.n; col++ {
		if e.isSafe(row, col) {
			e.board[row] = col
			e.solveRecursive(row + 1)
			if e.solutionFound {
				return
			}
		}
	}
}

// isSafe checks if placing a queen at (row, col) is safe
func (e *ExhaustiveSearchSolver) isSafe(row, col int) bool {
	for i := 0; i < row; i++ {
		// Check column conflict
		if e.board[i] == col {
			return false
		}
		// Check diagonal conflicts
		if abs(e.board[i]-col) == abs(i-row) {
			return false
		}
	}
	return true
}

// GetSolution returns the found solution
func (e *ExhaustiveSearchSolver) GetSolution() []int {
	return e.solution
}

// PrintSolution prints the solution board
func (e *ExhaustiveSearchSolver) PrintSolution() {
	if !e.solutionFound {
		fmt.Println("No solution found")
		return
	}

	fmt.Printf("Exhaustive Search Solution for N=%d:\n", e.n)
	for i := 0; i < e.n; i++ {
		for j := 0; j < e.n; j++ {
			if e.solution[i] == j {
				fmt.Print("Q ")
			} else {
				fmt.Print(". ")
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
