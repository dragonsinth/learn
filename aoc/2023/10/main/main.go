package main

import (
	"fmt"
	"github.com/dragonsinth/learn/aoc/termbox"
	"os"
	"strings"
)

var samples = []string{`
.....
.S-7.
.|.|.
.L-J.
.....
`, `
-L|F7
7S-7|
L|7||
-L-J|
L|-JF
`, `
..F7.
.FJ|.
SJ.L7
|F--J
LJ...
`, `
7-F7-
.FJ|7
SJLL7
|F--J
LJ.LJ
`, `
...........
.S-------7.
.|F-----7|.
.||.....||.
.||.....||.
.|L-7.F-J|.
.|..|.|..|.
.L--J.L--J.
...........
`, `
.F----7F7F7F7F-7....
.|F--7||||||||FJ....
.||.FJ||||||||L7....
FJL7L7LJLJ||LJ.L-7..
L--J.L7...LJS7F-7L7.
....F-J..F7FJ|L7L7L7
....L7.F7||L7|.L7L7|
.....|FJLJ|FJ|F7|.LJ
....FJL-7.||.||||...
....L---J.LJ.LJLJ...
`, `
FF7FSF7F7F7F7F7F---7
L|LJ||||||||||||F--J
FL-7LJLJ||||||LJL-77
F--JF--7||LJLJ7F7FJ-
L---JF-JLJ.||-FJLJJ7
|F|F-JF---7F7-L7L|7|
|FFJF7L7F-JF7|JL---7
7-L-JL7||F7|L7F-7F7|
L.L7LFJ|||||FJL7||LJ
L7JLJL-JLJLJL--JLJ.L
`}

func main() {
	for _, s := range samples {
		run(s, true)
	}
}

func run(sample string, debug bool) {
	p, start := parse(sample)
	p.resolveStart(start)
	if debug {
		fmt.Println("start")
		termbox.RenderPlain(p.Render(), os.Stdout)
	}

	// do the traversal
	seen := map[pos]int{start: 0}
	work := []pos{start}
	for len(work) > 0 {
		cur := work[0]
		dist := seen[cur]
		for i, pt := range cur.adjacent() {
			if _, ok := seen[pt]; !ok {
				if p.canTraverse(cur, dir(i)) {
					work = append(work, pt)
					seen[pt] = dist + 1
				}
			}
		}
		work = work[1:]
	}

	best := 0
	var bestPos pos
	for k, v := range seen {
		if v > best {
			bestPos, best = k, v
		}
	}
	fmt.Println(bestPos, best)

	// clear out any spaces not part of the loop
	for k := range p.pts {
		if _, ok := seen[k]; !ok {
			delete(p.pts, k)
		}
	}
	if debug {
		fmt.Println("clear")
		termbox.RenderPlain(p.Render(), os.Stdout)
	}

	pBig := p.Embiggen()
	if debug {
		fmt.Println("embiggen")
		termbox.RenderPlain(pBig.Render(), os.Stdout)
	}

	pBig.floodFill()
	if debug {
		fmt.Println("flood")
		termbox.RenderPlain(pBig.Render(), os.Stdout)
	}

	pSmall := pBig.Shrink()
	if debug {
		fmt.Println("shrink")
		termbox.RenderPlain(pSmall.Render(), os.Stdout)
	}

	fmt.Println("sum: ", pSmall.countEmpty())
}

var valid = map[byte]bool{
	'S': true,
	'F': true,
	'7': true,
	'L': true,
	'J': true,
	'-': true,
	'|': true,
	'.': true,
}

func parse(input string) (puz, pos) {
	p := puz{
		pts: map[pos]wat{},
		w:   -1,
		h:   -1,
	}
	var start pos

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
			pt := pos{x, p.h}
			if !valid[c] {
				panic(string(c))
			}
			p.pts[pt] = wat(c)
			if c == 'S' {
				start = pt
			}
		}
	}
	p.h++
	return p, start
}

type puz struct {
	pts  map[pos]wat
	w, h int
}

func (p puz) canTraverse(pt pos, d dir) bool {
	w := p.at(pt)
	switch w {
	case 'F':
		return d == E || d == S
	case '7':
		return d == W || d == S
	case 'L':
		return d == N || d == E
	case 'J':
		return d == N || d == W
	case '-':
		return d == E || d == W
	case '|':
		return d == N || d == S
	case '.':
		return false
	default:
		panic(w)
	}
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

func (p puz) resolveStart(start pos) {
	// figure out the orientation of the start pos
	val := 0
	for i, pt := range start.adjacent() {
		val <<= 1
		d := dir(i).invert()
		if p.canTraverse(pt, d) {
			val |= 1
		}
	}

	// NESW
	var w wat
	switch val {
	case 0b0110:
		w = 'F'
	case 0b0011:
		w = '7'
	case 0b1100:
		w = 'L'
	case 0b1001:
		w = 'J'
	case 0b0101:
		w = '-'
	case 0b1010:
		w = '|'
	default:
		panic(val)
	}
	p.pts[start] = w
}

func (p puz) Embiggen() puz {
	ret := puz{
		pts: map[pos]wat{},
		w:   p.w * 2,
		h:   p.h * 2,
	}

	for k, v := range p.pts {
		k.x *= 2
		k.y *= 2
		ret.pts[k] = v

		for i, pt := range k.adjacent() {
			d := dir(i)
			if ret.canTraverse(k, d) {
				if d == N || d == S {
					ret.pts[pt] = '|'
				} else {
					ret.pts[pt] = '-'
				}
			}
		}
	}
	return ret
}

func (p puz) Shrink() puz {
	ret := puz{
		pts: map[pos]wat{},
		w:   p.w / 2,
		h:   p.h / 2,
	}

	for k, v := range p.pts {
		if k.x%2 == 0 && k.y%2 == 0 {
			k.x /= 2
			k.y /= 2
			ret.pts[k] = v
		}
	}
	return ret
}

func (p puz) floodFill() {
	// flood fill from the edges
	var work []pos
	maybeAdd := func(x, y int) {
		pt := pos{x, y}
		if p.at(pt) == '.' {
			p.pts[pt] = 'X'
			work = append(work, pt)
		}
	}

	for x := 0; x < p.w; x++ {
		maybeAdd(x, 0)
		maybeAdd(x, p.h-1)
	}
	for y := 0; y < p.h; y++ {
		maybeAdd(0, y)
		maybeAdd(p.w-1, y)
	}

	for len(work) > 0 {
		cur := work[0]
		work = work[1:]
		for _, pt := range cur.adjacent() {
			x, y := pt.x, pt.y
			if x > 0 && y > 0 && x < p.w && y < p.h {
				maybeAdd(x, y)
			}
		}
	}
}

func (p puz) countEmpty() int {
	sum := 0
	for y := 0; y < p.h; y++ {
		for x := 0; x < p.w; x++ {
			if p.at(pos{x, y}) == '.' {
				sum++
			}
		}
	}
	return sum
}

type pos struct {
	x, y int
}

func (p pos) adjacent() [4]pos {
	return [4]pos{
		N: {p.x, p.y - 1},
		E: {p.x + 1, p.y},
		S: {p.x, p.y + 1},
		W: {p.x - 1, p.y},
	}
}

type dir byte

func (d dir) invert() dir {
	return (d + 2) % 4
}

const (
	N = dir(iota)
	E
	S
	W
)

type wat byte
