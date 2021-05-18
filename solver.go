// Package ACO solves the Asymmetric Traveling Salesman Problem using Ant Colony Optimization
//
// Ant colony optimization is a probabilistic optimization algorithm useful for getting approximative solutions to computational problems. The implementation in this package solves the Asymmetric Traveling Salesman Problem given by a weighted adjacency matrix.
//
// The ACO algorithm in this package is based on the simplest variant known as Ant System. Other more refined variants include MMAS and Ant Colony System.
// Important guidelines when defining constants are that all constants are required to be strictly positive, alfa < beta, rho < 1.
package ACO

import (
	"bufio"
	"math/rand"
	"os"
	"regexp"
	"strconv"
)

type solver interface {
	updatePheromones()
	constructAntSolutions()
	chooseNode()
}

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

// Visit node.
func (a *ant) visitNode(node int) {
	a.path = append(a.path, node)
	a.visited[node] = true
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
