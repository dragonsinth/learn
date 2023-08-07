package main

import (
	"fmt"
	"math"
	"strings"
)

var sample = `
Sabqponm
abcryxxl
accszExk
acctuvwj
abdefghi
`

func main() {
	var p puzzle
	var start, end pos

	for _, line := range strings.Split(sample, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		y := len(p.m)
		if y == 0 {
			p.w = len(line)
		} else if p.w != len(line) {
			panic(line)
		}

		var row []byte
		for x, c := range line {
			switch c {
			case 'S':
				row = append(row, 0)
				start = pos{x, y}
			case 'E':
				row = append(row, 25)
				end = pos{x, y}
			default:
				if c < 'a' || c > 'z' {
					panic(x)
				}
				row = append(row, byte(c-'a'))
			}
		}

		p.m = append(p.m, row)
	}
	p.h = len(p.m)
	fmt.Println(p.walk(start, end))
}

type pos struct {
	x, y int
}

type puzzle struct {
	m    [][]byte
	w, h int
}

func (p puzzle) walk(start pos, end pos) (int, int) {
	dist := map[pos]int{}
	p.descend(end, 0, dist)

	best := math.MaxInt32
	var bestPos pos
	for y := 0; y < p.h; y++ {
		for x := 0; x < p.w; x++ {
			loc := pos{x, y}
			d, ok := dist[loc]
			if !ok {
				fmt.Print(".")
				continue
			}

			d = d % 52
			if d < 26 {
				fmt.Print(string(rune(d) + 'a'))
			} else {
				fmt.Print(string(rune(d-26) + 'A'))
			}

			if p.height(loc) == 0 && dist[loc] < best {
				best, bestPos = dist[loc], loc
			}
		}
		fmt.Println()
	}

	fmt.Println(best, bestPos)

	ret, ok := dist[start]
	if !ok {
		panic("unreachable")
	}
	return ret, best
}

func (p puzzle) descend(loc pos, cur int, dist map[pos]int) {
	existing, ok := dist[loc]
	if ok && existing <= cur {
		return
	}
	dist[loc] = cur

	for _, dst := range []pos{
		{loc.x, loc.y - 1},
		{loc.x, loc.y + 1},
		{loc.x + 1, loc.y},
		{loc.x - 1, loc.y},
	} {
		if dst.x < 0 || dst.x >= p.w || dst.y < 0 || dst.y >= p.h {
			continue
		}
		if p.height(loc)-p.height(dst) > 1 {
			continue
		}
		p.descend(dst, cur+1, dist)
	}
}

func (p puzzle) height(loc pos) int {
	return int(p.m[loc.y][loc.x])
}
