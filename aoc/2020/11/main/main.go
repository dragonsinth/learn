package main

import (
	"fmt"
	"strings"
)

var sample = `
L.LL.LL.LL
LLLLLLL.LL
L.L.L..L..
LLLL.LL.LL
L.LL.LL.LL
L.LLLLL.LL
..L.L.....
LLLLLLLLLL
L.LLLLLL.L
L.LLLLL.LL
`

var (
	input = sample
)

type wat byte

const (
	FLOOR = wat(iota)
	EMPTY
	FULL
)

func (w wat) String() string {
	return string(w.toRune())
}

func (w wat) toRune() rune {
	switch w {
	case FLOOR:
		return '.'
	case EMPTY:
		return 'L'
	case FULL:
		return '#'
	default:
		panic(w)
	}
}

type puzzle struct {
	v             [][]wat
	width, height int
}

func (p *puzzle) String() string {
	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("%d, %d\n", p.width, p.height))
	for _, row := range p.v {
		for _, v := range row {
			buf.WriteRune(v.toRune())
		}
		buf.WriteRune('\n')
	}
	return buf.String()
}

func (p *puzzle) Occupied() int {
	var ret int
	for _, row := range p.v {
		for _, v := range row {
			if v == FULL {
				ret++
			}
		}
	}
	return ret
}

func (p *puzzle) Next(withSight bool) (*puzzle, bool) {
	ret := puzzle{
		v:      make([][]wat, p.height),
		width:  p.width,
		height: p.height,
	}

	changed := false
	for y := range p.v {
		ret.v[y] = make([]wat, p.width)
		for x := range ret.v[y] {
			if withSight {
				ret.v[y][x] = p.nextWithSight(x, y)
			} else {
				ret.v[y][x] = p.next(x, y)
			}
			if ret.v[y][x] != p.v[y][x] {
				changed = true
			}
		}
	}

	return &ret, changed
}

func (p *puzzle) next(px int, py int) wat {
	this := p.v[py][px]
	if this == FLOOR {
		return FLOOR
	}

	count := 0
	for y := py - 1; y <= py+1; y++ {
		for x := px - 1; x <= px+1; x++ {
			if x == px && y == py {
				continue
			}
			if p.val(x, y) == FULL {
				count++
			}
		}
	}

	if this == EMPTY && count == 0 {
		return FULL
	}

	if this == FULL && count >= 4 {
		return EMPTY
	}

	return this
}

type direction struct {
	dx, dy int
}

var (
	directions = []direction{
		{-1, -1},
		{-1, +0},
		{-1, +1},
		{+0, -1},
		{+0, 1},
		{+1, -1},
		{+1, +0},
		{+1, +1},
	}
)

// Count the number of visible squares in each direction.
func (p *puzzle) nextWithSight(px int, py int) wat {
	this := p.v[py][px]
	if this == FLOOR {
		return FLOOR
	}

	count := 0
	for _, dir := range directions {
		if seen := p.seen(px, py, dir); seen == FULL {
			count++
		}
	}

	if this == EMPTY && count == 0 {
		return FULL
	}

	if this == FULL && count >= 5 {
		return EMPTY
	}

	return this
}

func (p *puzzle) val(x int, y int) wat {
	if x < 0 || x >= p.width {
		return FLOOR
	}
	if y < 0 || y >= p.height {
		return FLOOR
	}
	return p.v[y][x]
}

func (p *puzzle) seen(x int, y int, dir direction) wat {
	for {
		x += dir.dx
		y += dir.dy
		if x < 0 || x >= p.width {
			return FLOOR
		}
		if y < 0 || y >= p.height {
			return FLOOR
		}
		this := p.v[y][x]
		if this != FLOOR {
			return this
		}
	}
}

func main() {
	var puz puzzle
	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		row := make([]wat, 0, 256)
		for _, v := range []byte(line) {
			switch v {
			case '.':
				row = append(row, FLOOR)
			case 'L':
				row = append(row, FULL)
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
	p := &puz
	for {
		var changed bool
		p, changed = p.Next(false)
		if !changed {
			break
		}
		//fmt.Println(p)
		fmt.Println(p.Occupied())
	}

	fmt.Println()
	p = &puz
	for {
		var changed bool
		p, changed = p.Next(true)
		if !changed {
			break
		}
		//fmt.Println(p)
		fmt.Println(p.Occupied())
	}
}
