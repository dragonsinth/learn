package main

import (
	"github.com/dragonsinth/learn/aoc/grid"
	"strings"
)

var sample = `
...........
.....###.#.
.###.##..#.
..#.#...#..
....#.#....
.##..S####.
.##..#...#.
.......##..
.##.#.####.
.##..##.##.
...........
`

func main() {
	run(sample, 1, true)
	run(sample, 2, true)
	run(sample, 3, true)
	run(sample, 6, true)

	run2(sample, 6, false)
	run2(sample, 10, false)
	run2(sample, 50, false)
	run2(sample, 100, false)
	run2(sample, 500, false)
	run2(sample, 1000, false)
	run2(sample, 5000, false)
}

func parse(input string) puz {
	p := puz{
		pts: map[pos]byte{},
		sz:  0,
	}
	lines := strings.Split(input, "\n")
	y := 0
	for _, line := range lines {
		if line == "" {
			continue
		}
		if p.sz == 0 {
			p.sz = len(line)
		} else if p.sz != len(line) {
			panic(line)
		}

		for x, c := range []byte(line) {
			pt := pos{x, y}
			if c == '#' || c == '.' {
				p.pts[pt] = c
			} else if c == 'S' {
				p.pts[pt] = '.'
				p.start = pt
			} else {
				panic(c)
			}
		}
		y++
	}
	if y != p.sz {
		panic(y)
	}

	return p
}

type puz struct {
	pts   map[pos]byte
	start pos
	sz    int
}

func (p puz) at(x, y int) byte {
	pt := pos{mod(x, p.sz), mod(y, p.sz)}
	v, ok := p.pts[pt]
	if !ok {
		return '#'
	}
	return v
}

func (p puz) Render() [][]byte {
	buf := grid.Alloc2d(p.sz, p.sz, byte(0))
	for y := 0; y < p.sz; y++ {
		for x := 0; x < p.sz; x++ {
			pt := pos{x, y}
			if pt == p.start {
				buf[y][x] = 'S'
			} else {
				buf[y][x] = p.at(x, y)
			}
		}
	}
	return buf
}

func (p puz) RenderWith(seen map[pos]int, maxDist int) [][]byte {
	buf := p.Render()
	for pt, dist := range seen {
		if dist > maxDist {
			continue
		}
		if dist%2 == maxDist%2 {
			buf[mod(pt.y, p.sz)][mod(pt.x, p.sz)] = 'O'
		}
	}
	return buf
}

func mod(a int, m int) int {
	return (a%m + m) % m
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
