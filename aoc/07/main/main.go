package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

var sample = `16,1,2,0,4,2,7,1,2,14`

var input = sample

func main() {
	var crabPositions []int
	max := 0
	for _, s := range strings.Split(input, ",") {
		v := mustInt(s)
		if v < 0 {
			panic(v)
		}
		crabPositions = append(crabPositions, v)
		if v > max {
			max = v
		}
	}

	bestPos := 0
	minCost := math.MaxInt
	for pos := 0; pos <= max; pos++ {
		// Compute the cost at this position.
		cost := 0
		for _, crab := range crabPositions {
			// cost := n(n+1)/2
			dist := abs(pos - crab)
			cost += dist * (dist + 1) / 2
 		}
		fmt.Println(pos, cost)
		if minCost > cost {
			bestPos = pos
			minCost = cost
		}
	}
	fmt.Println(bestPos, minCost)
}

func main1() {
	var crabPositions []int
	max := 0
	for _, s := range strings.Split(input, ",") {
		v := mustInt(s)
		if v < 0 {
			panic(v)
		}
		crabPositions = append(crabPositions, v)
		if v > max {
			max = v
		}
	}

	bestPos := 0
	minCost := math.MaxInt
	for pos := 0; pos <= max; pos++ {
		// Compute the cost at this position.
		cost := 0
		for _, crab := range crabPositions {
			cost += abs(pos - crab)
		}
		fmt.Println(pos, cost)
		if minCost > cost {
			bestPos = pos
			minCost = cost
		}
	}
	fmt.Println(bestPos, minCost)
}

func mustInt(s string) int {
	if v, err := strconv.Atoi(s); err != nil {
		panic(fmt.Sprint(s, err))
	} else {
		return v
	}
}

func abs(v int) int {
	if v < 0 {
		return -v
	} else {
		return v
	}
}
