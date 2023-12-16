package main

import (
	"fmt"
	"github.com/dragonsinth/learn/aoc/termbox"
	"os"
	"strings"
)

var sample = `
.|...\....
|.-.\.....
.....|-...
........|.
..........
.........\
..../.\\..
.-.-/..|..
.|....-|.\
..//.|....
`

func main() {
	part1(sample, true)
	part2(sample, true)
}

func part1(input string, debug bool) {
	p := parse(input)
	if debug {
		termbox.RenderPlain(p.data, os.Stdout)
		fmt.Println()
	}
	fmt.Println(p.trace(task{0, 0, E}, debug))
}

func part2(input string, debug bool) {
	p := parse(input)
	if debug {
		termbox.RenderPlain(p.data, os.Stdout)
		fmt.Println()
	}

	mx := 0
	for x := 0; x < p.w; x++ {
		mx = max(mx, p.trace(task{x, 0, S}, debug))
		mx = max(mx, p.trace(task{x, p.h - 1, N}, debug))
	}
	for y := 0; y < p.h; y++ {
		mx = max(mx, p.trace(task{0, y, E}, debug))
		mx = max(mx, p.trace(task{p.w - 1, y, W}, debug))
	}
	fmt.Println(mx)
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

const litMap = ".#"

func renderLit(lit [][]byte) [][]byte {
	out := make([][]byte, len(lit))
	for y, row := range lit {
		out[y] = make([]byte, len(row))
		for x, c := range row {
			out[y][x] = litMap[c]
		}
	}
	return out
}

func (p puz) trace(start task, debug bool) int {
	lit := make([][]byte, p.h)
	for y := range p.data {
		lit[y] = make([]byte, p.w)
	}

	work := []task{start}
	seen := map[task]bool{}
	for len(work) > 0 {
		t := work[0]
		work = work[1:]

		if t.x < 0 || t.x >= p.w {
			continue
		}
		if t.y < 0 || t.y >= p.h {
			continue
		}
		if seen[t] {
			continue
		}
		seen[t] = true
		lit[t.y][t.x] = 1
		switch p.data[t.y][t.x] {
		case '.':
			work = append(work, t.next())
		case '|':
			if t.d == N || t.d == S {
				// pass through
				work = append(work, t.next())
			} else {
				// split
				t.d = N
				work = append(work, t.next())
				t.d = S
				work = append(work, t.next())
			}
		case '-':
			if t.d == E || t.d == W {
				// pass through
				work = append(work, t.next())
			} else {
				// split
				t.d = E
				work = append(work, t.next())
				t.d = W
				work = append(work, t.next())
			}
		case '\\':
			if t.d == N || t.d == S {
				t.d = t.d.left()
				work = append(work, t.next())
			} else {
				t.d = t.d.right()
				work = append(work, t.next())
			}
		case '/':
			if t.d == N || t.d == S {
				t.d = t.d.right()
				work = append(work, t.next())
			} else {
				t.d = t.d.left()
				work = append(work, t.next())
			}
		default:
			panic(p.data[t.y][t.x])
		}
	}

	sum := 0
	for _, row := range lit {
		for _, c := range row {
			if c != 0 {
				sum++
			}
		}
	}

	if debug {
		fmt.Println(start.x, start.y, sum)
		termbox.RenderPlain(renderLit(lit), os.Stdout)
		fmt.Println()
	}

	return sum
}

type task struct {
	x, y int
	d    dir
}

func (t task) next() task {
	switch t.d {
	case N:
		return task{t.x, t.y - 1, t.d}
	case E:
		return task{t.x + 1, t.y, t.d}
	case S:
		return task{t.x, t.y + 1, t.d}
	case W:
		return task{t.x - 1, t.y, t.d}
	default:
		panic(t.d)
	}
}

type dir byte

func (d dir) invert() dir {
	return (d + 2) % 4
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

func max(a int, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}
