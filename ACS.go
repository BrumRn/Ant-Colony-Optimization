package ACO

import (
	"math"
	"math/rand"
)

type antColonySystem struct {
	tau0 float64
	phi  float64
	q0   float64
	colony
}

// Solve ATSP returns optimal cost and solution to the ATSP specified by matrix graph.
// The solution is heavily dependent on specified values for alfa, beta, rho, q & m.
func SolveACS(graph [][]float64, alfa float64, beta float64, rho float64, q float64, m int, tau0 float64, phi float64, q0 float64, iterations int) (float64, []int) {
	parent := colony{graph: graph, alfa: alfa, beta: beta, rho: rho, q: q, m: m}
	c := antColonySystem{colony: parent, tau0: tau0, phi: phi, q0: q0}

	c.configurateSolver()

	for gen := 0; gen < iterations; gen++ {
		c.constructAntSolutions()
		c.updateBestPath()
		c.updatePheromones()

	}

	return c.bestCost, c.bestPath
}

// Update pheromones. (Update only for best path)
func (c *antColonySystem) updatePheromones() {
	for i := 0; i < c.size; i++ {
		for j := 0; j < c.size; j++ {
			c.pheromones[i][j] *= (1 - c.rho)
		}
	}

	weight := c.q / c.bestCost

	c.pheromones[c.bestPath[c.size-1]][c.bestPath[0]] += weight

	for i := 1; i < c.size; i++ {
		c.pheromones[c.bestPath[i-1]][c.bestPath[i]] += weight
	}
}

// Update pheromones locally
func (c *antColonySystem) localUpdatePheromone(a *ant) {
	c.pheromones[a.path[c.size-1]][a.path[0]] *= (1 - c.phi)
	c.pheromones[a.path[c.size-1]][a.path[0]] += +c.phi * c.tau0

	for i := 1; i < c.size; i++ {
		c.pheromones[a.path[i-1]][a.path[i]] *= (1 - c.phi)
		c.pheromones[a.path[i-1]][a.path[i]] += +c.phi * c.tau0
	}
}

// Generate solutions for all ants.
func (c *antColonySystem) constructAntSolutions() {
	c.resetAnts()
	for _, a := range c.ants {
		c.antSimulation(a)
		c.localUpdatePheromone(a)
	}
}

// Generate solution for a single ant.
func (c *antColonySystem) antSimulation(a *ant) {
	a.visitNode(rand.Intn(a.size))
	for i := 1; i < a.size; i++ {
		a.visitNode(c.chooseNode(a))
		a.cost += a.graph[a.path[len(a.path)-2]][a.path[len(a.path)-1]]
	}
	a.cost += a.graph[a.path[len(a.path)-1]][a.path[0]]
}

// Choose next node using pheromone trails.
func (c *antColonySystem) chooseNode(a *ant) int {
	probabilities := make([]float64, c.size)
	var sum float64 = 0.0
	var p float64
	var pMax float64

	for next := 0; next < a.size; next++ {
		previous := a.path[len(a.path)-1]

		if !a.visited[next] {
			p = math.Pow(a.pheromones[previous][next], c.alfa) * math.Pow(a.weights[previous][next], c.beta)
			probabilities[next] = p
			sum += p
			if p > pMax {
				pMax = p
			}
		}
	}

	probabilities, sum = c.exploitPheromones(a, probabilities, pMax, sum)

	next := simulateChoice(probabilities, sum)
	return next
}

// Modify pheromones for exploitation.
func (c *antColonySystem) exploitPheromones(a *ant, probabilities []float64, pMax float64, sum float64) ([]float64, float64) {
	var newSum float64
	var q float64

	for i := range probabilities {
		q = rand.Float64()
		probabilities[i] /= sum
		if q < c.q0 {
			newSum = pMax - probabilities[i]
			probabilities[i] = pMax
		}
	}
	return probabilities, newSum
}

// Configurate solver.
func (c *antColonySystem) configurateSolver() {
	c.size = len(c.graph)
	c.generateMatrices()
	c.makeAnts()

	c.bestPath = make([]int, c.size)
	c.bestCost = 0
}

// Generates pheromone- and weight-matrices.
func (c *antColonySystem) generateMatrices() {
	weights := make([][]float64, c.size)
	pheromones := make([][]float64, c.size)

	for i := 0; i < c.size; i++ {
		weights[i] = make([]float64, c.size)
		pheromones[i] = make([]float64, c.size)
	}

	for i := 0; i < c.size; i++ {
		for j := 0; j < c.size; j++ {
			pheromones[i][j] = c.tau0
			weights[i][j] = c.tau0 / float64(c.graph[i][j])
		}
	}
	c.weights = weights
	c.pheromones = pheromones
}
