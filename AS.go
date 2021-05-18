package main

import (
	"math"
	"math/rand"
)

// Solve ATSP returns optimal cost and solution to the ATSP specified by matrix graph.
// The solution is heavily dependent on specified values for alfa, beta, rho, q & m.
func SolveAS(graph [][]float64, alfa float64, beta float64, rho float64, q float64, m int, iterations int) (float64, []int) {
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
		c.antSimulation(a)
	}
}

// Generate solution for a single ant.
func (c *colony) antSimulation(a *ant) {
	a.visitNode(rand.Intn(a.size))
	for i := 1; i < a.size; i++ {
		a.visitNode(c.chooseNode(a))
		a.cost += a.graph[a.path[len(a.path)-2]][a.path[len(a.path)-1]]
	}
	a.cost += a.graph[a.path[len(a.path)-1]][a.path[0]]
}

// Choose next node using pheromone trails.
func (c *colony) chooseNode(a *ant) int {
	probabilities := make([]float64, c.size)
	var sum float64 = 0.0
	var p float64
	for next := 0; next < a.size; next++ {
		previous := a.path[len(a.path)-1]
		if !a.visited[next] {
			p = math.Pow(a.pheromones[previous][next], a.alfa) * math.Pow(a.weights[previous][next], a.beta)
			probabilities[next] = p
			sum += p
		}
	}
	next := simulateChoice(probabilities, sum)
	return next
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
