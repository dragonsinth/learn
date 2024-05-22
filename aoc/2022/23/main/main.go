package main

import (
	"fmt"
	"github.com/dragonsinth/learn/aoc/termbox"
	"os"
	"strings"
	"time"
)

const sample1 = `
##
#.
..
##
`

const sample2 = `
....#..
..###.#
#...#.#
.#...##
#.###..
##.#.##
.#..#..
`

func main() {
	run(sample1, 0, 10)
	run(sample2, 0, 10)
	run(sample1, 0, 1000)
	run(sample2, 0, 1000)
}

func run(input string, frameDelay time.Duration, count int) {
	p := puzzle{
		elves: map[pos]bool{},
		max:   pos{-1, -1},
		turn:  0,
	}
	for _, line := range strings.Split(input, "\n") {
		if line == "" {
			continue
		}

		if p.max.x < 0 {
			p.max.x = len(line) - 1
		} else if p.max.x != len(line)-1 {
			panic(line)
		}
		p.max.y++
		y := p.max.y
		for x, v := range line {
			switch v {
			case '.':
			case '#':
				p.elves[pos{x, y}] = true
			default:
				panic(v)
			}
		}
	}
	p.nElves = len(p.elves)

	tb := termbox.New(frameDelay != 0)
	tb.Render(p.render(), os.Stdout)

	for p.turn < count {
		if !p.runTurn() {
			break
		}
		if tb.Enabled() {
			tb.Render(p.render(), os.Stdout)
			time.Sleep(frameDelay)
		}
	}
	tb.Stop()
	tb.Render(p.render(), os.Stdout)
	if count == 10 {
		w := p.max.x - p.min.x + 1
		h := p.max.y - p.min.y + 1
		fmt.Println(w, h, w*h-len(p.elves))
	}
}

type puzzle struct {
	elves    map[pos]bool
	nElves   int
	min, max pos
	turn     int
}

func (p *puzzle) render() [][]byte {
	var out [][]byte
	for y := p.min.y; y <= p.max.y; y++ {
		var line []byte
		for x := p.min.x; x <= p.max.x; x++ {
			line = append(line, p.get(pos{x, y}))
		}
		out = append(out, line)
	}
	statusLine := fmt.Sprintf("%d: %s -> %s", p.turn, p.min, p.max)
	out = append(out, []byte(statusLine))
	return out
}

func (p *puzzle) runTurn() bool {
	anyMoved := false
	dests := map[pos]int{}
	destByElf := map[pos]pos{}
	newElves := map[pos]bool{}

	setCount := 0
	setDst := func(elf pos) {
		if setCount != len(newElves) {
			panic("here")
		}
		if newElves[elf] {
			panic("here")
		}
		newElves[elf] = true
		setCount++
		if setCount != len(newElves) {
			panic("here")
		}
	}

	// Consider
	for elf := range p.elves {
		moved := false

		anyAdjacent := false
		for _, pt := range elf.allAdjacent() {
			anyAdjacent = anyAdjacent || p.get(pt) == '#'
		}
		if anyAdjacent {
			for i := 0; i < 4; i++ {
				d := dir(p.turn+i) % 4
				blocked := false
				for _, pt := range elf.adjacent(d) {
					blocked = blocked || p.get(pt) == '#'
				}
				if !blocked {
					dest := elf.advance(d)
					if p.get(dest) == '#' {
						panic("here")
					}

					destByElf[elf] = dest
					dests[dest]++
					moved = true
					break
				}
			}
		}

		if !moved {
			setDst(elf)
		}
	}

	// Move
	for elf, dest := range destByElf {
		switch dests[dest] {
		case 0:
			panic("here")
		case 1:
			if p.get(dest) == '#' {
				panic("here")
			}

			// move
			setDst(dest)
			p.min.x = min(p.min.x, dest.x)
			p.max.x = max(p.max.x, dest.x)
			p.min.y = min(p.min.y, dest.y)
			p.max.y = max(p.max.y, dest.y)
			anyMoved = true
		default:
			setDst(elf)
		}
	}

	if p.nElves != len(newElves) {
		panic("here")
	}

	p.elves = newElves
	p.turn++
	return anyMoved
}

func (p *puzzle) get(pt pos) byte {
	if p.elves[pt] {
		return '#'
	} else {
		return '.'
	}
}
