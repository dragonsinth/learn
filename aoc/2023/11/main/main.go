package main

import (
	"fmt"
	"github.com/dragonsinth/learn/aoc/termbox"
	"maps"
	"os"
	"strings"
)

var sample = `
...#......
.......#..
#.........
..........
......#...
.#........
.........#
..........
.......#..
#...#.....
`

func main() {
	run(sample, 2, true)
	run(sample, 10, false)
	run(sample, 100, false)
}

func run(input string, fac int, debug bool) {
	p := parse(input)
	if debug {
		fmt.Println("start")
		termbox.RenderPlain(p.Render(), os.Stdout)
	}
	p = p.Embiggen(fac)
	if debug {
		fmt.Println("embiggen")
		termbox.RenderPlain(p.Render(), os.Stdout)
	}

	// compute dists
	sum := 0
	for k1, v1 := range p.pts {
		for k2, v2 := range p.pts {
			if v1 == G && v2 == G {
				sum += k1.dist(k2)
			}
		}
	}
	fmt.Println(sum / 2)
}

func parse(input string) puz {
	p := puz{
		pts: map[pos]wat{},
		w:   -1,
		h:   -1,
	}

	lines := strings.Split(input, "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		if p.w < 0 {
			p.w = len(line)
		} else if p.w != len(line) {
			panic(line)
		}
		p.h++

		for x, c := range []byte(line) {
			w := wat(c)
			if w == G {
				pt := pos{x, p.h}
				p.pts[pt] = w
			} else if w != E {
				panic(c)
			}
		}
	}
	p.h++
	return p
}

type puz struct {
	pts  map[pos]wat
	w, h int
}

func (p puz) at(pt pos) wat {
	v, ok := p.pts[pt]
	if !ok {
		return '.'
	}
	return v
}

func (p puz) Render() [][]byte {
	var buf [][]byte
	for y := 0; y < p.h; y++ {
		var line []byte
		for x := 0; x < p.w; x++ {
			w := p.at(pos{x, y})
			line = append(line, byte(w))
		}
		buf = append(buf, line)
	}
	return buf
}

func (p puz) Embiggen(fac int) puz {
	r := puz{
		pts: maps.Clone(p.pts),
		w:   p.w,
		h:   p.h,
	}

	// map rows
	dx := 0
	rowMap := make([]int, r.w)
	for x := 0; x < r.w; x++ {
		rowMap[x] = dx
		found := false
		for y := 0; y < r.h; y++ {
			if r.at(pos{x, y}) == G {
				found = true
				break
			}
		}
		if found {
			dx++
		} else {
			dx += fac
		}
	}

	// map cols
	dy := 0
	colMap := make([]int, r.h)
	for y := 0; y < r.h; y++ {
		colMap[y] = dy
		found := false
		for x := 0; x < r.w; x++ {
			if r.at(pos{x, y}) == G {
				found = true
				break
			}
		}
		if found {
			dy++
		} else {
			dy += fac
		}
	}

	// map points
	out := make(map[pos]wat, len(r.pts))
	for k, v := range r.pts {
		k.x = rowMap[k.x]
		k.y = colMap[k.y]
		out[k] = v
	}

	r.pts, r.w, r.h = out, dx, dy
	return r
}

type pos struct {
	x, y int
}

func (p pos) dist(o pos) int {
	return abs(p.x-o.x) + abs(p.y-o.y)
}

type wat byte

const (
	G = wat('#')
	E = wat('.')
)

func abs(i int) int {
	if i < 0 {
		return -i
	} else {
		return i
	}
}
