package main

import (
	"fmt"
	"runtime"
	"strings"
	"time"
)

func main() {
	runBasicComparison() // Quick comparison with smaller N values
}

func runBasicComparison() {
	fmt.Println("N-Queens Problem Solver - Basic Comparison")
	fmt.Println("==========================================")

	testSizes := []int{10, 15, 20, 30, 50, 100, 200}
	//testSizes := []int{5, 10, 15, 20, 25}

	for _, n := range testSizes {
		fmt.Printf("\nTesting N = %d\n", n)
		fmt.Println(strings.Repeat("-", 50))

		// Test Exhaustive Search (only for small N)
		if n <= 20 {
			testAlgorithmWithSolution("Exhaustive DFS", n, func() (bool, func()) {
				solver := NewExhaustiveSearchSolver(n)
				success := solver.Solve()
				return success, func() { solver.PrintSolution() }
			})
		} else {
			fmt.Printf("%-20s: Time: %12s, Memory: %8s, Success: %s\n",
				"Exhaustive DFS", "SKIPPED", "N/A", "N/A (too large)")
		}

		// Test Greedy Search
		if n <= 50 {
			testAlgorithmWithSolution("Greedy Hill Climbing", n, func() (bool, func()) {
				solver := NewGreedySolver(n)
				success := solver.Solve()
				return success, func() { solver.PrintSolution() }
			})
		} else {
			fmt.Printf("%-20s: Time: %12s, Memory: %8s, Success: %s\n",
				"Greedy Hill Climbing", "SKIPPED", "N/A", "N/A (too large)")
		}

		// Test Simulated Annealing
		testAlgorithmWithSolution("Simulated Annealing", n, func() (bool, func()) {
			solver := NewSimulatedAnnealingSolver(n)
			success := solver.Solve()
			return success, func() { solver.PrintSolution() }
		})

		// Test Genetic Algorithm
		testAlgorithmWithSolution("Genetic Algorithm", n, func() (bool, func()) {
			solver := NewGeneticSolver(n)
			success := solver.Solve()
			return success, func() { solver.PrintSolution() }
		})
	}
}

func testAlgorithmWithSolution(name string, n int, solveFunc func() (bool, func())) {
	var m1, m2 runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&m1)

	start := time.Now()
	success, printFunc := solveFunc()
	duration := time.Since(start)

	runtime.ReadMemStats(&m2)
	memUsed := m2.TotalAlloc - m1.TotalAlloc
	heapUsed := m2.HeapAlloc - m1.HeapAlloc

	fmt.Printf("%-20s: Time: %12v, Memory: %8d KB (Heap: %d KB), Success: %v\n",
		name, duration, memUsed/1024, heapUsed/1024, success)

	// Show solution for small N values
	if n <= 20 && success {
		printFunc()
	}
}
