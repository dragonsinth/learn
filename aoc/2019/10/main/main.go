package main

import (
	"fmt"
	"math"
	"sort"
	"strings"
)

var samples = []string{
	`
......#.#.
#..#.#....
..#######.
.#.#.###..
.#..#.....
..#....#.#
#..#....#.
.##.#..###
##...#..#.
.#....####`,
	`
#.#...#.#.
.###....#.
.#....#...
##.#.#.#.#
....#.#.#.
.##..###.#
..#...##..
..##....##
......#...
.####.###.`,
	`
.#..#..###
####.###.#
....###.#.
..###.##.#
##.##.#.#.
....###..#
..#.#..#.#
#..#.#.###
.##...##.#
.....#.#..`,
	`
.#..##.###...#######
##.############..##.
.#.######.########.#
.###.#######.####.#.
#####.##.#.##.###.##
..#####..#.#########
####################
#.####....###.#.#.##
##.#################
#####.##.###..####..
..######..##.#######
####.##.####...##..#
.#####..#.######.###
##...#.##########...
#.##########.#######
.####.#.###.###.#.##
....##.##.###..#####
.#.#.###########.###
#.#.#.#####.####.###
###.##.####.##.#..##`,
}

func main() {
	// test theta
	for _, v := range []vec{
		{00, -1},
		{+1, -1},
		{+1, 00},
		{+1, +1},
		{00, +1},
		{-1, +1},
		{-1, 00},
		{-1, -1},
	} {
		fmt.Println(v.x, v.y, v.theta())
	}

	for _, input := range samples[3:] {
		p := parse(input)
		p.run()
	}
}

func parse(input string) puzzle {
	var data [][]byte
	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		data = append(data, []byte(line))
	}
	return puzzle{d: data}
}

type puzzle struct {
	d [][]byte
}

func (p puzzle) run() {
	bx, by, bc := 0, 0, 0
	for y, line := range p.d {
		for x, v := range line {
			if v != '#' {
				continue
			}

			// consider this location
			c := p.visible(x, y)
			if c > bc {
				bx, by, bc = x, y, c
			}

		}
	}
	fmt.Println(bx, by, bc)

	// Run the zapper program.
	lvs := p.lengthVectors(bx, by)
	var last vec
	i := 0
	z := 0
	for len(lvs) > 0 {
		if i == len(lvs) {
			i = 0
			last = vec{}
		}

		v := lvs[i]
		if v.norm() == last.norm() {
			i++
			continue
		}

		// Zap!
		last = v
		lvs = append(lvs[:i], lvs[i+1:]...)
		z++
		if z <= 10 || z == 20 || z == 50 || abs(200-z) < 2 || len(lvs) == 0 {
			fmt.Printf("Zap %d: %d,%d\n", z, v.x+bx, v.y+by)
		}
		// don't increment i
	}
}

type vec struct {
	x, y int
}

func (v vec) norm() vec {
	// normalize the vector
	g := abs(gcd(v.x, v.y))
	if g != 0 {
		return vec{
			x: v.x / g,
			y: v.y / g,
		}
	}
	return v
}

func (v vec) len() int {
	// hacky
	return abs(v.x) + abs(v.y)
}

func (v vec) theta() float64 {
	// Map atan2 such that -1, 0 => 0 and values go clockwise
	return -math.Atan2(float64(v.x), float64(v.y)) + math.Pi
}

func (p puzzle) visible(srcx int, srcy int) int {
	all := map[vec]bool{}
	for y, line := range p.d {
		for x, v := range line {
			if x == srcx && y == srcy {
				continue
			}
			if v != '#' {
				continue
			}

			// consider this vector
			dx := x - srcx
			dy := y - srcy
			v := vec{
				x: dx,
				y: dy,
			}

			// normalize the vector
			g := abs(gcd(dx, dy))
			if g != 0 {
				dx, dy = dx/g, dy/g
			}
			all[v.norm()] = true

		}
	}
	return len(all)
}

func (p puzzle) lengthVectors(srcx int, srcy int) []vec {
	var ret []vec
	for y, line := range p.d {
		for x, v := range line {
			if x == srcx && y == srcy {
				continue
			}
			if v != '#' {
				continue
			}

			dx := x - srcx
			dy := y - srcy
			ret = append(ret, vec{
				x: dx,
				y: dy,
			})
		}
	}

	sort.Slice(ret, func(i, j int) bool {
		return less(ret[i], ret[j])
	})
	return ret
}

func less(a vec, b vec) bool {
	if a == b {
		return false
	}

	if a.norm() == b.norm() {
		return a.len() < b.len()
	}

	return a.theta() < b.theta()
}

func gcd(a, b int) int {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
