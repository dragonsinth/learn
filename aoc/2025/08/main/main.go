package main

import (
	"fmt"
	"math"
	"slices"
	"strconv"
	"strings"
)

const sample = `
162,817,812
57,618,57
906,360,560
592,479,940
352,342,300
466,668,158
542,29,236
431,825,988
739,650,466
52,470,668
216,146,977
819,987,18
117,168,530
805,96,715
346,949,466
970,615,88
941,993,340
862,61,35
984,92,344
425,690,689
`

type node struct {
	x, y, z int
	circuit *circuit
}

type edge struct {
	a, b *node
	dist int
}

type circuit struct {
	nodes []*node
}

func main() {
	run(parse(sample), 10)
	run(parse(sample), math.MaxInt)
}

func run(nodes []*node, turns int) {
	edges := buildEdges(nodes)
	for i := 0; i < turns; i++ {
		e := edges[i]
		a, b := e.a, e.b
		if a.circuit == b.circuit {
			// nothing happens
			continue
		}
		// arbitrarily merge b into a
		ac, bc := a.circuit, b.circuit
		for _, n := range bc.nodes {
			n.circuit = ac
		}
		ac.nodes = append(ac.nodes, bc.nodes...)
		bc.nodes = nil

		// part 2: termination check
		if len(ac.nodes) == len(nodes) {
			fmt.Println("exhausted", a, b)
			fmt.Println(a.x * b.x)
			return
		}
	}

	// part 1: largest circuits
	seen := map[*circuit]bool{}
	var lens []int
	for _, n := range nodes {
		if seen[n.circuit] {
			continue
		}
		lens = append(lens, len(n.circuit.nodes))
		seen[n.circuit] = true
	}
	slices.Sort(lens)
	slices.Reverse(lens)
	fmt.Println(lens)
	fmt.Println(lens[0] * lens[1] * lens[2])
}

func buildEdges(nodes []*node) []edge {
	var ret []edge
	for i, a := range nodes {
		for j, b := range nodes {
			if i < j {
				ret = append(ret, edge{
					a:    a,
					b:    b,
					dist: calcDist(a, b),
				})
			}
		}
	}
	slices.SortFunc(ret, func(a, b edge) int {
		return a.dist - b.dist
	})
	return ret
}

func calcDist(a *node, b *node) int {
	// don't bother taking the square root, we only need distances to sort
	return sqr(a.x-b.x) + sqr(a.y-b.y) + sqr(a.z-b.z)
}

func sqr(n int) int {
	return n * n
}

func parse(input string) []*node {
	var nodes []*node
	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Split(line, ",")
		if len(parts) != 3 {
			panic(len(parts))
		}
		n := &node{
			x: mustInt(parts[0]),
			y: mustInt(parts[1]),
			z: mustInt(parts[2]),
		}
		n.circuit = &circuit{
			nodes: []*node{n},
		}
		nodes = append(nodes, n)
	}
	return nodes
}

func mustInt(s string) int {
	if v, err := strconv.Atoi(s); err != nil {
		panic(fmt.Sprint(s, err))
	} else {
		return v
	}
}
