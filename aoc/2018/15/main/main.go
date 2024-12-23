package main

import (
	"fmt"
	"strings"

	"github.com/dragonsinth/learn/aoc/termbox"
)

var samples = []string{`
#######
#.G...#
#...EG#
#.#.#G#
#..G#E#
#.....#
#######
`, `
#######
#G..#E#
#E#E.E#
#G.##.#
#...#E#
#...E.#
#######
`, `
#######
#E..EG#
#.#G.E#
#E.##E#
#G..#.#
#..E#.#
#######
`, `
#######
#E.G#.#
#.#G..#
#G.#.G#
#G..#.#
#...E.#
#######
`, `
#######
#.E...#
#.#..G#
#.###.#
#E#G#G#
#...#G#
#######
`, `
#########
#G......#
#.E.#...#
#..##..G#
#...##..#
#...#...#
#.G...G.#
#.....G.#
#########
`}

var expect = []int{
	27730,
	36334,
	39514,
	27755,
	28944,
	18740,
}

var expect2 = []int{
	4988,
	29064,
	31284,
	3478,
	6474,
	1140,
}

func main() {
	for i, s := range samples {
		p := run(s, false)
		sum := p.hpSum()
		actual := p.turn * sum
		fmt.Println(p.turn, p.hpSum(), actual)
		if actual != expect[i] {
			panic(actual)
		}
	}

	for i, s := range samples {
		p := run2(s, false)
		sum := p.hpSum()
		actual := p.turn * sum
		fmt.Println(p.elfAtk, p.turn, p.hpSum(), actual)
		if actual != expect2[i] {
			panic(actual)
		}
	}
}

func run(sample string, debug bool) *puz {
	p := parse(sample)
	p.elfAtk = 3
	p.debug = debug
	p.term = termbox.New(false)
	//p.Print()
	for {
		live := p.Tick()
		if live {
			p.PrintFrame()
		} else {
			p.term.Stop()
			p.Print()
			return &p
		}
	}
}

func run2(sample string, debug bool) *puz {
	for atk := 4; true; atk++ {
		p := parse(sample)
		p.elfAtk = atk
		p.debug = debug
		p.term = termbox.New(false)
		//p.Print()
		nElf := p.nElf
		for {
			live := p.Tick()
			if live {
				p.PrintFrame()
			} else {
				p.term.Stop()
				p.Print()
				if nElf == p.nElf {
					return &p
				} else {
					break
				}
			}
		}
	}
	return nil
}

func parse(input string) puz {
	p := puz{
		units: map[pos]*unit{},
		walls: map[pos]bool{},
		w:     -1,
		h:     -1,
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
			pt := pos{x, p.h}
			switch c {
			case '#':
				p.walls[pt] = true
			case '.':
			case 'E':
				p.units[pt] = &unit{p: pt, hp: 200, typ: ELF}
				p.nElf++
			case 'G':
				p.units[pt] = &unit{p: pt, hp: 200, typ: GOBLIN}
				p.nGob++
			default:
				panic(string(c))
			}
		}
	}
	p.h++
	return p
}

type pos struct {
	x, y int
}

func (p pos) adjacent() [4]pos {
	return [4]pos{
		{p.x, p.y - 1},
		{p.x - 1, p.y},
		{p.x + 1, p.y},
		{p.x, p.y + 1},
	}
}

func readOrderPosLess(a, b pos) int {
	if a.y != b.y {
		return a.y - b.y
	}
	return a.x - b.x
}

type wat byte

const (
	ELF    = wat('E')
	GOBLIN = wat('G')
)

type unit struct {
	p   pos
	hp  int
	typ wat
}

func mapKeys[K comparable, V any](in map[K]V) []K {
	var ret []K
	for k := range in {
		ret = append(ret, k)
	}
	return ret
}
