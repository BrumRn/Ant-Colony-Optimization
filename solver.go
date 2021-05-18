// Package ACO solves the Asymmetric Traveling Salesman Problem using Ant Colony Optimization
//
// Ant colony optimization is a probabilistic optimization algorithm useful for getting approximative solutions to computational problems. The implementation in this package solves the Asymmetric Traveling Salesman Problem given by a weighted adjacency matrix.
//
// The ACO algorithm in this package is based on the simplest variant known as Ant System. Other more refined variants include MMAS and Ant Colony System.
// Important guidelines when defining constants are that all constants are required to be strictly positive, alfa < beta, rho < 1.
package main

import (
	"bufio"
	"math"
	"math/rand"
	"os"
	"regexp"
	"strconv"
)

// A colony of ants.
type colony struct {
	graph      [][]float64
	weights    [][]float64
	pheromones [][]float64
	bestPath   []int
	bestCost   float64
	size       int
	alfa       float64
	beta       float64
	rho        float64
	m          int
	q          float64
	ants       []*ant
}

// A Single ant.
type ant struct {
	cost    float64
	path    []int
	visited []bool
	*colony
}

// ReadATSP returns a graph matrix from an .atsp file.
func ReadATSP(fileName string) [][]float64 {
	file, _ := os.Open(fileName)

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	var txt string
	var n int
	var match bool

	scanner.Split(bufio.ScanWords)

	for scanner.Scan() {
		txt = scanner.Text()

		if match {
			n, _ = strconv.Atoi(scanner.Text())
		}
		match, _ = regexp.MatchString(txt, "DIMENSION:")

		if txt == "EDGE_WEIGHT_SECTION" {
			break
		}
	}

	list := make([][]float64, n)

	for i := 0; i < n; i++ {
		list[i] = make([]float64, n)
	}

	var i, j int
	for scanner.Scan() {
		if j == n {
			i++
			j = 0
		}

		if i == n {
			break
		}
		list[i][j], _ = strconv.ParseFloat(scanner.Text(), 64)
		j++
	}
	return list
}

// Solve ATSP returns optimal cost and solution to the ATSP specified by matrix graph.
// The solution is heavily dependent on specified values for alfa, beta, rho, q & m.
func SolveATSP(graph [][]float64, alfa float64, beta float64, rho float64, q float64, m int, iterations int) (float64, []int) {
	c := colony{graph: graph, alfa: alfa, beta: beta, rho: rho, q: q, m: m}

	c.configurateSolver()

	for gen := 0; gen < iterations; gen++ {
		c.constructAntSolutions()
		c.updateBestPath()
		c.updatePheromones()

	}

	return c.bestCost, c.bestPath
}

// Generate solutions for all ants.
func (c *colony) constructAntSolutions() {
	c.resetAnts()
	for _, a := range c.ants {
		a.antSimulation()
	}
}

// Generate solution for a single ant.
func (a *ant) antSimulation() {
	a.visitNode(rand.Intn(a.size))
	for i := 1; i < a.size; i++ {
		a.visitNode(a.chooseNode())
		a.cost += a.graph[a.path[len(a.path)-2]][a.path[len(a.path)-1]]
	}
	a.cost += a.graph[a.path[len(a.path)-1]][a.path[0]]
}

// Visit node.
func (a *ant) visitNode(node int) {
	a.path = append(a.path, node)
	a.visited[node] = true
}

// Choose next node using pheromone trails.
func (a *ant) chooseNode() int {
	probabilities := make([]float64, a.size)
	var sum float64 = 0.0
	var c float64
	for next := 0; next < a.size; next++ {
		previous := a.path[len(a.path)-1]
		if !a.visited[next] {
			c = math.Pow(a.pheromones[previous][next], a.alfa) * math.Pow(a.weights[previous][next], a.beta)
			probabilities[next] = c
			sum += c
		}
	}
	next := simulateChoice(probabilities, sum)
	return next
}

// Simulate choice using probability distribution.
func simulateChoice(probabilities []float64, sum float64) int {
	u := sum * rand.Float64()
	node := 0
	c := probabilities[node]
	for c < u {
		node++
		c += probabilities[node]
	}

	return node
}

// Update pheromones.
func (c *colony) updatePheromones() {
	for i := 0; i < c.size; i++ {
		for j := 0; j < c.size; j++ {
			c.pheromones[i][j] *= (1 - c.rho)
		}
	}
	for _, a := range c.ants {

		weight := c.q / a.cost

		c.pheromones[a.path[c.size-1]][a.path[0]] += weight

		for i := 1; i < c.size; i++ {
			c.pheromones[a.path[i-1]][a.path[i]] += weight
		}
	}
}

// Save best path/cost.
func (c *colony) updateBestPath() {
	for _, a := range c.ants {
		if a.cost < c.bestCost {
			c.bestPath = a.path
			c.bestCost = a.cost

		} else if c.bestCost == 0 {
			c.bestPath = a.path
			c.bestCost = a.cost
		}
	}
}

// Generates pheromone- and weight-matrices.
func (c *colony) generateMatrices() {
	weights := make([][]float64, c.size)
	pheromones := make([][]float64, c.size)

	for i := 0; i < c.size; i++ {
		weights[i] = make([]float64, c.size)
		pheromones[i] = make([]float64, c.size)
	}

	for i := 0; i < c.size; i++ {
		for j := 0; j < c.size; j++ {
			pheromones[i][j] = 1.0
			weights[i][j] = 1.0 / float64(c.graph[i][j])
		}
	}
	c.weights = weights
	c.pheromones = pheromones
}

// Configurate solver.
func (c *colony) configurateSolver() {
	c.size = len(c.graph)
	c.generateMatrices()
	c.makeAnts()

	c.bestPath = make([]int, c.size)
	c.bestCost = 0
}

// Create m ants.
func (c *colony) makeAnts() {
	c.ants = make([]*ant, 0, c.m)

	for i := 0; i < c.m; i++ {
		a := &ant{colony: c}
		a.cost = 0
		a.path = make([]int, 0, a.size)
		a.visited = make([]bool, a.size)
		c.ants = append(c.ants, a)
	}
}

// Reset ants before each generation.
func (c *colony) resetAnts() {
	for _, a := range c.ants {
		a.cost = 0
		a.path = make([]int, 0, a.size)
		a.visited = make([]bool, a.size)
	}
}
