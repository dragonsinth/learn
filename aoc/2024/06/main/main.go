package main

import (
	"fmt"
	"github.com/dragonsinth/learn/aoc/termbox"
	"os"
	"strings"
	"time"
)

const sample = `
....#.....
.........#
..........
..#.......
.......#..
..........
.#..^.....
........#.
#.........
......#...
`

func main() {
	part1(parse(sample), true)
	// part1(parse(data), true)
	part2(parse(sample))
}

func part1(p *puz, debug bool) {
	term := termbox.New(false)

	if debug {
		term.Start()
	}

	r := p.Run(term)
	if r != EXIT {
		panic(r)
	}

	if debug {
		term.Stop()
	}
	term.Render(p.Render(), os.Stdout)
	fmt.Println(len(p.visitPos))
}

func part2(p *puz) {
	sum := 0
	for y := 0; y < p.h; y++ {
		for x := 0; x < p.w; x++ {
			pt := pos{x, y}
			w := p.data[pt]
			if w == '.' && pt != p.guardStart.p {
				p.data[pt] = '#'

				switch p.Run(nil) {
				case CONT:
					panic("here")
				case EXIT:
					// not a hit
				case LOOP:
					fmt.Println("hit", pt)
					sum++
				}

				p.data[pt] = '.'
			}
		}
	}
	fmt.Println(sum)
}

func parse(input string) *puz {
	p := &puz{
		data: map[pos]wat{},
		w:    0,
		h:    0,
	}

	lines := strings.Split(input, "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		if p.w <= 0 {
			p.w = len(line)
		} else if p.w != len(line) {
			panic(line)
		}
		y := p.h
		for x, c := range line {
			pt := pos{x: x, y: y}
			if c == '^' {
				p.guardStart = state{p: pt, d: N}
				p.data[pt] = '.'
			} else {
				p.data[pt] = wat(c)
			}
		}
		p.h++
	}
	return p
}

type wat byte

type pos struct {
	x, y int
}

type dir byte

const (
	N = dir(iota)
	E
	S
	W
)

func (d dir) Next(p pos) pos {
	switch d {
	case N:
		return pos{p.x + 0, p.y - 1}
	case E:
		return pos{p.x + 1, p.y + 0}
	case S:
		return pos{p.x + 0, p.y + 1}
	case W:
		return pos{p.x - 1, p.y + 0}
	default:
		panic(d)
	}
}

type puz struct {
	data       map[pos]wat
	w, h       int
	guardStart state
	guardCur   state
	visitPos   map[pos]bool
	visitState map[state]bool
}

func (p *puz) Render() [][]byte {
	var buf [][]byte
	for y := 0; y < p.h; y++ {
		var line []byte
		for x := 0; x < p.w; x++ {
			pt := pos{x, y}
			if pt == p.guardCur.p {
				line = append(line, "^>V<"[p.guardCur.d])
			} else if p.visitPos[pt] {
				line = append(line, 'X')
			} else {
				line = append(line, byte(p.data[pt]))
			}
		}
		buf = append(buf, line)
	}
	return buf
}

type result int

const (
	CONT = result(iota)
	LOOP
	EXIT
)

type state struct {
	p pos
	d dir
}

func (p *puz) Tick() result {
	p.visitPos[p.guardCur.p] = true
	p.visitState[p.guardCur] = true

	nextPos := p.guardCur.d.Next(p.guardCur.p)
	if nextPos.x < 0 || nextPos.x >= p.w || nextPos.y < 0 || nextPos.y >= p.h {
		return EXIT
	}

	if p.data[nextPos] == '#' {
		p.guardCur.d = (p.guardCur.d + 1) % 4
		return CONT
	}

	nextState := state{p: nextPos, d: p.guardCur.d}
	if p.visitState[nextState] {
		return LOOP
	}

	p.guardCur.p = nextPos
	return CONT
}

func (p *puz) Run(term termbox.Terminal) interface{} {
	p.Reset()
	for {
		if term != nil && term.Enabled() {
			term.Render(p.Render(), os.Stdout)
			time.Sleep(10 * time.Millisecond)
		}
		switch p.Tick() {
		case CONT:
		case EXIT:
			return EXIT
		case LOOP:
			return LOOP
		}
	}
}

func (p *puz) Reset() {
	p.guardCur = p.guardStart
	p.visitPos = map[pos]bool{}
	p.visitState = map[state]bool{}
}
