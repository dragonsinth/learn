package main

import (
	"fmt"
	"strings"
)

const sample = `
..@@.@@@@.
@@@.@.@.@@
@@@@@.@.@@
@.@@@@..@.
@@.@@@@.@@
.@@@@@@@.@
.@.@.@.@@@
@.@@@.@@@@
.@@@@@@@@.
@.@.@@@.@.
`

func main() {
	fmt.Println(parse(sample).removeFree())
	fmt.Println(parse(sample).removeAllFree())
}

func parse(input string) *puz {
	p := puz{
		wat: map[pos]byte{},
	}

	lines := strings.Split(input, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if p.w <= 0 {
			p.w = len(line)
		} else if p.w != len(line) {
			panic(line)
		}
		p.h++

		for x, c := range []byte(line) {
			pt := pos{x, p.h}
			p.wat[pt] = c
		}
	}
	p.h++
	return &p
}

type pos struct {
	x, y int
}

type dir byte

const (
	N = dir(iota)
	NE
	E
	SE
	S
	SW
	W
	NW
)

func (d dir) Next(p pos) pos {
	switch d {
	case N:
		return pos{p.x - 1, p.y + 0}
	case NE:
		return pos{p.x - 1, p.y + 1}
	case E:
		return pos{p.x + 0, p.y + 1}
	case SE:
		return pos{p.x + 1, p.y + 1}
	case S:
		return pos{p.x + 1, p.y + 0}
	case SW:
		return pos{p.x + 1, p.y - 1}
	case W:
		return pos{p.x + 0, p.y - 1}
	case NW:
		return pos{p.x - 1, p.y - 1}
	default:
		panic(d)
	}
}

type puz struct {
	wat  map[pos]byte
	w, h int
}

func (p *puz) removeAllFree() int {
	sum := 0
	for {
		iter := p.removeFree()
		if iter == 0 {
			return sum
		}
		sum += iter
		fmt.Println(iter)
	}
}

func (p *puz) removeFree() int {
	sum := 0
	for y := 0; y < p.h; y++ {
		for x := 0; x < p.w; x++ {
			pt := pos{x, y}
			if p.wat[pt] == '@' {
				// count neighbors
				neighbors := 0
				for d := dir(0); d < dir(8); d++ {
					neighorWat := p.wat[d.Next(pt)]
					if neighorWat == '@' || neighorWat == 'x' {
						neighbors++
					}
				}
				if neighbors < 4 {
					sum++
					p.wat[pt] = 'x'
				}
			}
		}
	}

	// clear the taken
	for k, v := range p.wat {
		if v == 'x' {
			p.wat[k] = '.'
		}
	}
	return sum
}
