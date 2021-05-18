package main

import "math"

type solver interface {
	updatePheromones()
	constructAntSolutions()
	chooseNode()
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
