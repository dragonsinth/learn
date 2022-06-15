package main

import (
	"fmt"
	"strings"
)

var sample = `
v...>>.vv>
.vv>>.vv..
>>.>v>...v
>>v>>.>.v.
v>v.vv.v..
>.>>..v...
.vv..>.>v.
v.v..>>v.v
....v..v.>
`

var (
	input = sample
)

type what byte

const (
	EMPTY = what(iota)
	EAST
	SOUTH
)

type puzzle struct {
	v             [][]what
	width, height int
}

func (w what) toRune() rune {
	switch w {
	case EMPTY:
		return '.'
	case EAST:
		return '>'
	case SOUTH:
		return 'v'
	default:
		panic(w)
	}
}

func (p puzzle) String() string {
	var buf strings.Builder
	for _, row := range p.v {
		for _, v := range row {
			buf.WriteRune(v.toRune())
		}
		buf.WriteRune('\n')
	}
	return buf.String()
}

func (p *puzzle) Next() int {
	ret := 0

	// All EAST cukes try to move east.
	canMove := make([]bool, p.width*p.height)
	for y, row := range p.v {
		for x, wat := range row {
			if wat == EAST {
				if p.v[y][(x+1)%p.width] == EMPTY {
					canMove[y*p.width+x] = true
				}
			}
		}
	}

	// Apply the moves
	for i, v := range canMove {
		if v {
			y := i / p.width
			x := i % p.width
			p.v[y][x] = EMPTY
			p.v[y][(x+1)%p.width] = EAST
			ret++
		}
	}

	// All SOUTH cukes try to move south.
	canMove = make([]bool, p.width*p.height)
	for y, row := range p.v {
		for x, wat := range row {
			if wat == SOUTH {
				if p.v[(y+1)%p.height][x] == EMPTY {
					canMove[y*p.width+x] = true
				}
			}
		}
	}

	// Apply the moves
	for i, v := range canMove {
		if v {
			y := i / p.width
			x := i % p.width
			p.v[y][x] = EMPTY
			p.v[(y+1)%p.height][x] = SOUTH
			ret++
		}
	}

	return ret
}

func main() {
	var puz puzzle
	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		row := make([]what, 0, 256)
		for _, v := range []byte(line) {
			switch v {
			case '.':
				row = append(row, EMPTY)
			case '>':
				row = append(row, EAST)
			case 'v':
				row = append(row, SOUTH)
			default:
				panic(line)
			}
		}
		if puz.width == 0 {
			puz.width = len(row)
		} else if puz.width != len(row) {
			panic(len(row))
		}
		puz.v = append(puz.v, row)
	}
	puz.height = len(puz.v)

	fmt.Println(puz)
	for i := 0; i < 1000; i++ {
		nMoved := puz.Next()
		fmt.Printf("After %d steps:\n", i+1)
		fmt.Println(puz)
		if nMoved == 0 {
			break
		}
	}
}
