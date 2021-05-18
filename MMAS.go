package main

// A colony of ants.
type mmasColony struct {
	tauMax float64
	tauMin float64
	tau0   float64 // tau0 = tauMax
	colony
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
