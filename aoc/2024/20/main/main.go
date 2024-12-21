package main

import (
	"fmt"
	"strings"
)

const sample = `
###############
#...#...#.....#
#.#.#.#.#.###.#
#S#...#.#.#...#
#######.#.#.###
#######.#.#...#
#######.#.###.#
###..E#...#...#
###.#######.###
#...###...#...#
#.#####.#.###.#
#.#...#.#.#...#
#.#.#.#.#.#.###
#...#...#...###
###############
`

func main() {
	run(sample, true, 2, 2)
	run(sample, true, 20, 50)
}

func run(input string, debug bool, maxWarp int, minIntersting int) {
	p := parse(input)
	rt := p.mapRoute()

	sum := 0
	cheats := make([]int, len(rt.byOrd))

	// compare every point to every other point, for any interesting cheats, range check
	for dStart, pStart := range rt.byOrd {
		for dEnd, pEnd := range rt.byOrd {
			// must be warpable
			realDist := pStart.dist(pEnd)
			if realDist > maxWarp {
				continue
			}

			// must be useful
			logicalDist := dEnd - dStart
			saved := logicalDist - realDist
			if saved >= minIntersting {
				cheats[saved]++
				sum++
			}
		}
	}

	if debug {
		for dist, ch := range cheats {
			if ch > 0 {
				fmt.Printf("There are %d cheats that save %d picoseconds.\n", ch, dist)
			}
		}
	}

	fmt.Println(sum)
}

func parse(input string) *puz {
	p := &puz{
		data: map[pos]byte{},
	}

	lines := strings.Split(input, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if p.w <= 0 {
			p.w = len(line)
		} else if p.w != len(line) {
			panic(line)
		}
		y := p.h
		for x, wat := range line {
			pt := pos{x: x, y: y}
			switch wat {
			case '.', '#':
				p.data[pt] = byte(wat)
			case 'S':
				p.start = pt
				p.data[pt] = '.'
			case 'E':
				p.end = pt
				p.data[pt] = '.'
			default:
				panic(wat)
			}
		}
		p.h++

	}
	return p
}

type puz struct {
	data  map[pos]byte
	w, h  int
	start pos
	end   pos
}

type Route struct {
	byPos map[pos]int
	byOrd []pos
}

func (p puz) mapRoute() Route {
	ret := Route{byPos: map[pos]int{}, byOrd: nil}
	cur := p.start

	for {
		ret.byPos[cur] = len(ret.byOrd)
		ret.byOrd = append(ret.byOrd, cur)

		if cur == p.end {
			if len(ret.byPos) != len(ret.byOrd) {
				panic("mismatch")
			}
			return ret
		}

		// exactly one adjacent point won't have been mapped
		found := false
		for _, n := range cur.adjacent() {
			if _, ok := ret.byPos[n]; p.data[n] == '.' && !ok {
				if found {
					panic("double")
				}
				found = true
				cur = n
			}
		}
		if !found {
			panic("none")
		}
	}

	return ret
}

type pos struct {
	x, y int
}

func (p pos) adjacent() [4]pos {
	return [4]pos{
		{p.x, p.y - 1},
		{p.x + 1, p.y},
		{p.x, p.y + 1},
		{p.x - 1, p.y},
	}
}

func (p pos) warp() [4]pos {
	return [4]pos{
		{p.x, p.y - 2},
		{p.x + 2, p.y},
		{p.x, p.y + 2},
		{p.x - 2, p.y},
	}
}

func (p pos) dist(o pos) int {
	return abs(p.x-o.x) + abs(p.y-o.y)
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}
