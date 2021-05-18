package ACO

// A colony of ants.
type mmasColony struct {
	tauMax float64
	tauMin float64
	colony
}

// Solve ATSP returns optimal cost and solution to the ATSP specified by matrix graph.
// The solution is heavily dependent on specified values for alfa, beta, rho, q, m & tauMin/Max.
func SolveMMAS(graph [][]float64, alfa float64, beta float64, rho float64, q float64, m int, tauMax float64, tauMin float64, iterations int) (float64, []int) {
	parent := colony{graph: graph, alfa: alfa, beta: beta, rho: rho, q: q, m: m}
	c := mmasColony{colony: parent, tauMax: tauMax, tauMin: tauMin}

	c.configurateSolver()

	for gen := 0; gen < iterations; gen++ {
		c.constructAntSolutions()
		c.updateBestPath()
		c.updatePheromones()

	}

	return c.bestCost, c.bestPath
}

// Update pheromones. (Update only for best path)
func (c *mmasColony) updatePheromones() {
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

	c.scalePheromones()
}

// Limits range of pheromone.
func (c *mmasColony) scalePheromones() {
	for i := 0; i < c.size; i++ {
		for j := 0; j < c.size; j++ {
			if c.pheromones[i][j] > c.tauMax {
				c.pheromones[i][j] = c.tauMax
			} else if c.pheromones[i][j] < c.tauMin {
				c.pheromones[i][j] = c.tauMin

			}
		}
	}
}

// Configurate solver.
func (c *mmasColony) configurateSolver() {
	c.size = len(c.graph)
	c.generateMatrices()
	c.makeAnts()

	c.bestPath = make([]int, c.size)
	c.bestCost = 0
}

// Generates pheromone- and weight-matrices.
func (c *mmasColony) generateMatrices() {
	weights := make([][]float64, c.size)
	pheromones := make([][]float64, c.size)

	for i := 0; i < c.size; i++ {
		weights[i] = make([]float64, c.size)
		pheromones[i] = make([]float64, c.size)
	}

	for i := 0; i < c.size; i++ {
		for j := 0; j < c.size; j++ {
			pheromones[i][j] = c.tauMax
			weights[i][j] = c.tauMax / float64(c.graph[i][j])
		}
	}
	c.weights = weights
	c.pheromones = pheromones
}
