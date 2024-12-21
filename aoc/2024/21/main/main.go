package main

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
)

const sample = `
029A
980A
179A
456A
379A
`

func main() {
	{
		robot1Pad := computeShortestPaths(directionalPad, humanCosts{})
		robot2Pad := computeShortestPaths(directionalPad, robotCosts{robot1Pad})
		finalPad := computeShortestPaths(codePad, robotCosts{robot2Pad})
		run(sample, finalPad)
		// run(data, finalPad)
	}

	{
		upper := UpperPad(humanCosts{})
		for i := 0; i < 25; i++ {
			pad := computeShortestPaths(directionalPad, upper)
			upper = robotCosts{pad}
		}
		finalPad := computeShortestPaths(codePad, upper)
		run(sample, finalPad)
		// run(data, finalPad)
	}
}

func run(input string, pad map[pathKey]path) {
	sum := 0
	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		p := compute(pad, line)
		if p.path != "" {
			fmt.Println(len(p.path), p.path)
		}
		val := mustInt(line[:3])
		score := val * p.size
		sum += score
	}
	fmt.Println(sum)
}

func compute(pad map[pathKey]path, input string) path {
	var ret path
	last := byte('A')
	for _, c := range input {
		next := byte(c)
		p := pad[pathKey{last, next}]
		ret = ret.add(p)
		last = next
	}
	return ret
}

type path struct {
	path string
	size int
}

func (p path) add(other path) path {
	size := p.size + other.size
	var pathStr string
	if size < 100 {
		pathStr = p.path + other.path
	}
	return path{pathStr, size}
}

type pathKey struct {
	start, end byte
}

type pos struct {
	x, y int
}

/*
+---+---+---+
| 7 | 8 | 9 |
+---+---+---+
| 4 | 5 | 6 |
+---+---+---+
| 1 | 2 | 3 |
+---+---+---+
|   | 0 | A |
+---+---+---+
*/
var codePad = map[byte]pos{
	'7': {0, 0},
	'8': {1, 0},
	'9': {2, 0},
	'4': {0, 1},
	'5': {1, 1},
	'6': {2, 1},
	'1': {0, 2},
	'2': {1, 2},
	'3': {2, 2},
	'0': {1, 3},
	'A': {2, 3},
}

/*
+---+---+---+
|   | ^ | A |
+---+---+---+
| < | v | > |
+---+---+---+
*/
var directionalPad = map[byte]pos{
	'^': {1, 0},
	'A': {2, 0},
	'<': {0, 1},
	'v': {1, 1},
	'>': {2, 1},
}

type UpperPad interface {
	Press(start, end byte) path
}

type humanCosts struct{}

func (humanCosts) Press(start, end byte) path {
	return path{
		path: string(rune(end)),
		size: 1,
	}
}

type robotCosts struct {
	paths map[pathKey]path
}

func (r robotCosts) Press(start, end byte) path {
	return r.paths[pathKey{start, end}]
}

func computeShortestPaths(pad map[byte]pos, upper UpperPad) map[pathKey]path {
	// Compute the shortest path from every button to every other button
	posMap := map[pos]byte{}
	for k, pt := range pad {
		posMap[pt] = k
	}

	// precompute a route from every key to every other key
	ret := map[pathKey]path{}
	for k1, pt1 := range pad {
		for k2, pt2 := range pad {
			pk := pathKey{k1, k2}
			ret[pk] = findShortestPath(pt1, pt2, posMap, upper, 'A')
		}
	}

	return ret

}

func findShortestPath(pt1 pos, pt2 pos, posMap map[pos]byte, upper UpperPad, upperPos byte) path {
	// base case
	if pt1 == pt2 {
		return upper.Press(upperPos, 'A')
	}

	var options []path
	if pt1.x < pt2.x {
		n := pt1
		n.x++
		if posMap[n] != 0 {
			nextUpperPos := byte('>')
			p := upper.Press(upperPos, nextUpperPos)
			p = p.add(findShortestPath(n, pt2, posMap, upper, nextUpperPos))
			options = append(options, p)
		}
	}
	if pt1.x > pt2.x {
		n := pt1
		n.x--
		if posMap[n] != 0 {
			nextUpperPos := byte('<')
			p := upper.Press(upperPos, nextUpperPos)
			p = p.add(findShortestPath(n, pt2, posMap, upper, nextUpperPos))
			options = append(options, p)
		}
	}
	if pt1.y < pt2.y {
		n := pt1
		n.y++
		if posMap[n] != 0 {
			nextUpperPos := byte('v')
			p := upper.Press(upperPos, nextUpperPos)
			p = p.add(findShortestPath(n, pt2, posMap, upper, nextUpperPos))
			options = append(options, p)
		}
	}
	if pt1.y > pt2.y {
		n := pt1
		n.y--
		if posMap[n] != 0 {
			nextUpperPos := byte('^')
			p := upper.Press(upperPos, nextUpperPos)
			p = p.add(findShortestPath(n, pt2, posMap, upper, nextUpperPos))
			options = append(options, p)
		}
	}

	slices.SortFunc(options, func(a, b path) int {
		return a.size - b.size
	})
	return options[0]
}

func mustInt(s string) int {
	if v, err := strconv.Atoi(s); err != nil {
		panic(fmt.Sprint(s, err))
	} else {
		return v
	}
}
