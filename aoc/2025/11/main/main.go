package main

import (
	"fmt"
	"strings"
)

const sample1 = `
aaa: you hhh
you: bbb ccc
bbb: ddd eee
ccc: ddd eee fff
ddd: ggg
eee: out
fff: out
ggg: out
hhh: ccc fff iii
iii: out
`

const sample2 = `
svr: aaa bbb
aaa: fft
fft: ccc
bbb: tty
tty: ccc
ccc: ddd eee
ddd: hub
hub: fff
eee: dac
dac: fff
fff: ggg hhh
ggg: out
hhh: out
`

type node struct {
	name    string
	forward []*node
	reverse []*node
}

type puzzle struct {
	nodes map[string]*node
}

func main() {
	part1(parse(sample1), true)
	part2(parse(sample2), true)
}

func parse(input string) *puzzle {
	ret := &puzzle{
		nodes: map[string]*node{
			"you": {name: "you"},
			"out": {name: "out"},
		},
	}

	ensureNode := func(name string) *node {
		if _, ok := ret.nodes[name]; !ok {
			ret.nodes[name] = &node{name: name}
		}
		return ret.nodes[name]
	}

	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.Split(line, ":")
		if len(parts) != 2 {
			panic(parts)
		}

		n := ensureNode(parts[0])
		edges := strings.Split(parts[1], " ")
		for _, edge := range edges {
			dest := ensureNode(edge)
			n.forward = append(n.forward, dest)
			dest.reverse = append(dest.reverse, n)
		}
	}
	return ret
}

func part1(puz *puzzle, debug bool) {
	// simple recursive with map?
	fmt.Println(puz.countRoutes(map[string]int{
		"out": 1,
	}, puz.nodes["you"]))
}

func part2(puz *puzzle, debug bool) {
	// Count routes from dac -> fft, and fft -> dac
	dac2fft := puz.countRoutes(map[string]int{
		"dac": 1,
	}, puz.nodes["fft"])

	fft2dac := puz.countRoutes(map[string]int{
		"fft": 1,
	}, puz.nodes["dac"])

	first, second := "dac", "fft"
	if dac2fft > fft2dac {
		first, second = "fft", "dac"
	}
	middle := max(dac2fft, fft2dac)

	head := puz.countRoutes(map[string]int{
		first: 1,
	}, puz.nodes["svr"])

	tail := puz.countRoutes(map[string]int{
		"out": 1,
	}, puz.nodes[second])

	fmt.Println(dac2fft, fft2dac, head, tail, head*middle*tail)
}

func (p *puzzle) countRoutes(cache map[string]int, n *node) int {
	if v, ok := cache[n.name]; ok {
		return v
	}
	sum := 0
	for _, e := range n.forward {
		sum += p.countRoutes(cache, e)
	}
	cache[n.name] = sum
	return sum
}
