package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/dragonsinth/learn/aoc/grid"
	"github.com/dragonsinth/learn/aoc/termbox"
)

var sample = `
.......S.......
...............
.......^.......
...............
......^.^......
...............
.....^.^.^.....
...............
....^.^...^....
...............
...^.^...^.^...
...............
..^...^.....^..
...............
.^.^.^.^.^...^.
...............
`

type pos struct {
	x, y int
}

func main() {
	run(sample, true)
}

func run(input string, debug bool) {
	p := parse(input)
	if debug {
		termbox.RenderPlain(p.Render(), os.Stdout)
		fmt.Println()
	}
	p.run()
	if debug {
		termbox.RenderPlain(p.Render(), os.Stdout)
		fmt.Println()
	}
	// part 1
	fmt.Println(p.split)

	// part 2
	sum := 0
	for x := 0; x < p.w; x++ {
		sum += p.counts[pos{x, p.h - 1}]
	}
	fmt.Println(sum)
}

func parse(input string) puz {
	p := puz{
		data:   map[pos]byte{},
		counts: map[pos]int{},
	}
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

		for x, c := range line {
			pt := pos{x, p.h}
			switch c {
			case '.':
			// nothing
			case 'S':
				p.start = pt
				p.counts[pt] = 1
			case '^':
				p.data[pt] = '^'
			default:
				panic(c)
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
	split int

	counts map[pos]int // part 2
}

func (p *puz) run() {
	for y := 0; y < p.h-1; y++ {
		for x := 0; x < p.w; x++ {
			pt := pos{x, y}
			if pt == p.start || p.data[pt] == '|' {
				// propagate to the next line
				next := pos{x, y + 1}
				switch wat := p.data[next]; wat {
				case 0, '|':
					p.light(pos{x, y + 1}, p.counts[pt])
				case '^':
					p.split++
					p.light(pos{x - 1, y + 1}, p.counts[pt])
					p.light(pos{x + 1, y + 1}, p.counts[pt])
				default:
					panic(string(wat))
				}
			}
		}
	}
}

func (p *puz) light(pt pos, count int) {
	switch p.data[pt] {
	case 0:
		p.data[pt] = '|'
	case '|':
		// already lit
	default:
		panic(string(p.data[pt]))
	}
	p.counts[pt] += count
}

func (p *puz) Render() [][]byte {
	buf := grid.Alloc2d(p.w, p.h, byte('.'))
	for y := 0; y < p.h; y++ {
		for x := 0; x < p.w; x++ {
			pt := pos{x, y}
			if pt == p.start {
				buf[y][x] = 'S'
			} else if p.data[pt] != 0 {
				buf[y][x] = p.data[pt]
			}
		}
	}
	return buf
}
