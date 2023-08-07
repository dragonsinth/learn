package main

import (
	"fmt"
	"math"
	"strings"

	"github.com/dragonsinth/learn/aoc/termbox"
)

const sample = `
.#.#...|#.
.....#|##|
.|..|...#.
..|#.....#
#.#|||#|#|
...#.||...
.|....|...
||...#|.#|
|.||||..|.
...#.|..|.
`

func main() {
	run(sample, false)
}

func run(sample string, debug bool) {
	p := parse(sample)
	p.term = termbox.New(false)
	p.Print()

	if debug {
		p.term.Start()
	}
	seen := map[string]int{}
	dest := math.MaxInt
	for p.turn = 1; p.turn <= dest; p.turn++ {
		p.Tick()
		key := p.key()
		if debug {
			p.PrintFrame()
		} else {
			nTrees, nLumber := p.Count()
			fmt.Println(p.turn, nTrees, nLumber, nTrees*nLumber)
		}
		if last := seen[key]; last != 0 {
			mod := p.turn - last
			const goal = 1000000000
			advance := (goal - p.turn) % mod
			dest = p.turn + advance
			if !debug {
				fmt.Println("last seen at turn:", last, "run", advance, "to", dest)
			}
		}
		seen[key] = p.turn
	}
	p.term.Stop()
	p.Print()
	nTrees, nLumber := p.Count()
	fmt.Println(p.turn, nTrees, nLumber, nTrees*nLumber)
}

func parse(input string) puz {
	p := puz{
		data: nil,
		w:    -1,
		h:    -1,
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
		p.data = append(p.data, []byte(line))
	}
	p.h = len(p.data)
	return p
}

type pos struct {
	x, y int
}

func (p pos) adjacent() [8]pos {
	return [8]pos{
		{p.x - 1, p.y - 1},
		{p.x, p.y - 1},
		{p.x + 1, p.y - 1},
		{p.x - 1, p.y},
		{p.x + 1, p.y},
		{p.x - 1, p.y + 1},
		{p.x, p.y + 1},
		{p.x + 1, p.y + 1},
	}
}
