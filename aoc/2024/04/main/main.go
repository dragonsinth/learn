package main

import (
	"fmt"
	"strings"
)

const sample = `
MMMSXXMASM
MSAMXMSMSA
AMXSXMAAMM
MSAMASMSMX
XMASAMXAMM
XXAMMXXAMA
SMSMSASXSS
SAXAMASAAA
MAMMMXMMMM
MXMXAXMASX
`

func main() {
	// part 1
	fmt.Println(parse(sample).countXmas())
	// part 2
	fmt.Println(parse(sample).countX_mas())
}

func parse(input string) puz {
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
	return p
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

func (p puz) countXmas() int {
	sum := 0
	for y := 0; y < p.h; y++ {
		for x := 0; x < p.w; x++ {
			pt := pos{x, y}
			if p.wat[pt] == 'X' { // optimization
				for d := dir(0); d < dir(8); d++ {
					if p.isXmas(pt, d) {
						sum++
					}
				}
			}
		}
	}
	return sum
}

func (p puz) isXmas(pt1 pos, d dir) bool {
	pt2 := d.Next(pt1)
	pt3 := d.Next(pt2)
	pt4 := d.Next(pt3)
	return p.wat[pt1] == 'X' && p.wat[pt2] == 'M' && p.wat[pt3] == 'A' && p.wat[pt4] == 'S'
}

func (p puz) countX_mas() int {
	sum := 0
	for y := 0; y < p.h; y++ {
		for x := 0; x < p.w; x++ {
			pt := pos{x, y}
			if p.isX_mas(pt) {
				sum++
			}
		}
	}
	return sum
}

func (p puz) isX_mas(pt pos) bool {
	if p.wat[pt] != 'A' {
		return false
	}

	if !p.isMS(NW.Next(pt), SE.Next(pt)) {
		return false
	}

	if !p.isMS(NE.Next(pt), SW.Next(pt)) {
		return false
	}

	return true
}

func (p puz) isMS(a pos, b pos) bool {
	mCount, sCount := 0, 0
	switch p.wat[a] {
	case 'M':
		mCount++
	case 'S':
		sCount++
	}
	switch p.wat[b] {
	case 'M':
		mCount++
	case 'S':
		sCount++
	}
	return mCount == 1 && sCount == 1
}
