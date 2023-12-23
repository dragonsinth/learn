package main

import (
	"fmt"
	"strings"
)

var sample = `
#.#####################
#.......#########...###
#######.#########.#.###
###.....#.>.>.###.#.###
###v#####.#v#.###.#.###
###.>...#.#.#.....#...#
###v###.#.#.#########.#
###...#.#.#.......#...#
#####.#.#.#######.#.###
#.....#.#.#.......#...#
#.#####.#.#.#########v#
#.#...#...#...###...>.#
#.#.#v#######v###.###v#
#...#.>.#...>.>.#.###.#
#####v#.#.###v#.#.###.#
#.....#...#...#.#.#...#
#.#########.###.#.#.###
#...###...#...#...#.###
###.###.#.###v#####v###
#...#...#.#.>.>.#.>.###
#.###.###.#.###.#.#v###
#.....###...###...#...#
#####################.#
`

func main() {
	run(sample, false, true)
	run(sample, true, true)
}

func run(input string, either bool, debug bool) {
	p := parse(input)
	g := p.computeGraph(pos{1, 0}, pos{p.w - 2, p.h - 1}, debug)
	routes := g.findAllRoutes(either, debug)
	if debug {
		for _, r := range routes {
			fmt.Println(r.len)
		}
	}
	fmt.Println(routes[0].len)
}

func parse(input string) puz {
	var p puz
	lines := strings.Split(input, "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		if p.w == 0 {
			p.w = len(line)
		} else if p.w != len(line) {
			panic(line)
		}

		p.data = append(p.data, []byte(line))
		p.h++
	}
	return p
}

type puz struct {
	data [][]byte
	w, h int
}

type pos struct {
	x, y int
}

type dir byte

func (d dir) String() string {
	return "^>v<"[d : d+1]
}

func (d dir) sym() byte {
	return "^>v<"[d]
}

func (d dir) left() dir {
	return (d + 4 - 1) % 4
}

func (d dir) right() dir {
	return (d + 1) % 4
}

const (
	N = dir(iota)
	E
	S
	W
)
