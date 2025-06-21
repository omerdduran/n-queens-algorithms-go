# N-Queens Problem Solver

This project implements four different approaches to solve the N-Queens problem in Go:

1. **Exhaustive Depth-First Search** - Complete backtracking algorithm
2. **Greedy Hill Climbing** - Local search optimization
3. **Simulated Annealing** - Probabilistic optimization technique
4. **Genetic Algorithm** - Evolutionary computation approach

## Problem Description

The N-Queens problem is a classic combinatorial optimization problem where the goal is to place N chess queens on an N×N chessboard such that no two queens attack each other. This means:
- No two queens share the same row
- No two queens share the same column  
- No two queens share the same diagonal

## Algorithms Implemented

### 1. Exhaustive Depth-First Search
- **Approach**: Systematic backtracking search
- **Guarantees**: Always finds a solution if one exists
- **Time Complexity**: O(N!) in worst case
- **Best for**: Smaller values of N (≤ 30)

### 2. Greedy Hill Climbing
- **Approach**: Iteratively moves to better neighboring states
- **Guarantees**: May get stuck in local optima
- **Time Complexity**: O(iterations × N²)
- **Best for**: Quick solutions, may need restarts

### 3. Simulated Annealing
- **Approach**: Probabilistic search with cooling schedule
- **Guarantees**: Can escape local optima with decreasing probability
- **Time Complexity**: O(iterations × N²)
- **Best for**: Good balance of solution quality and speed

### 4. Genetic Algorithm
- **Approach**: Evolutionary search with population of solutions
- **Guarantees**: Population-based search, good exploration
- **Time Complexity**: O(generations × population_size × N²)
- **Best for**: Large problem instances, parallel processing potential

## Usage

### Running the Comparison
```bash
go run *.go
```

This will test all four algorithms on N = 10, 30, 50, 100, and 200, measuring:
- Execution time
- Memory usage (both TotalAlloc and HeapAlloc)
- Success rate
- Visual solution boards (ASCII format) for successful solutions

### Example Output
```
N-Queens Problem Solver - Basic Comparison
==========================================

Testing N = 5
--------------------------------------------------
Exhaustive DFS      : Time:     12.5µs, Memory:        0 KB (Heap: 0 KB), Success: true
Exhaustive Search Solution for N=5:
Q . . . .
. . Q . .
. . . . Q
. Q . . .
. . . Q .

Greedy Hill Climbing: Time:   156.25µs, Memory:        1 KB (Heap: 1 KB), Success: true
Greedy Search Solution for N=5:
. Q . . .
. . . Q .
Q . . . .
. . Q . .
. . . . Q

Simulated Annealing : Time:   2.145ms, Memory:      245 KB (Heap: 12 KB), Success: true
Simulated Annealing Solution for N=5:
. . Q . .
Q . . . .
. . . Q .
. Q . . .
. . . . Q

Genetic Algorithm   : Time:   1.789ms, Memory:       89 KB (Heap: 89 KB), Success: true
Genetic Algorithm Solution for N=5:
. . . Q .
Q . . . .
. . Q . .
. . . . Q
. Q . . .

```

## Algorithm Parameters

### Greedy Hill Climbing
- Maximum iterations: 10,000
- Restart on local optima

### Simulated Annealing  
- Initial temperature: N×N (scaled with problem size)
- Cooling rate: 0.99
- Minimum temperature: 0.01
- Maximum iterations: N×1000 (scaled with problem size)
- Number of restarts: 5

### Genetic Algorithm
- Population size: 80 (120 for N>20, 150 for N>40)
- Maximum generations: 200
- Mutation rate: 0.15 (15%)
- Crossover rate: 0.85 (85%)
- Elite preservation: populationSize/10
- Selection: Tournament selection
- Number of restarts: 5

## Performance Analysis

### Expected Performance Characteristics:

1. **N = 10**: All algorithms should solve quickly
2. **N = 30**: Exhaustive search may struggle, heuristics perform well
3. **N = 50**: Only heuristic methods practical
4. **N = 100**: Genetic algorithm and simulated annealing preferred
5. **N = 200**: Genetic algorithm likely performs best

### Memory Usage:
- Exhaustive: O(N) for recursion stack
- Greedy: O(N) for board representation
- Simulated Annealing: O(N) for board representation  
- Genetic: O(population_size × N) for population

## Files Structure

- `main.go` - Main entry point and performance testing
- `exhaustive.go` - Depth-first search implementation
- `greedy.go` - Hill climbing implementation  
- `simulated_annealing.go` - Simulated annealing implementation
- `genetic.go` - Genetic algorithm implementation
- `README.md` - This documentation
- `go.mod` - Go module definition

## Building and Running

```bash
# Build the project
go build -o nqueens *.go

# Run the executable
./nqueens

# Or run directly
go run *.go
```

## Results Analysis

The program outputs timing and memory usage data that can be used for:
- Comparing algorithm efficiency across different problem sizes
- Analyzing scalability characteristics
- Understanding trade-offs between solution quality and computational cost
- Generating data for academic reports and analysis
