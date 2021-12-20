package main

import (
	"fmt"
	"strings"
)

var sample = `
..#.#..#####.#.#.#.###.##.....###.##.#..###.####..#####..#....#..#..##..##
#..######.###...####..#..#####..##..#.#####...##.#.#..#.##..#.#......#.###
.######.###.####...#.##.##..#..#..#####.....#.#....###..#.##......#.....#.
.#..#..##..#...##.######.####.####.#.#...#.......#..#.#.#...####.##.#.....
.#..#...##.#.##..#...##.#.##..###.#......#.#.......#.#.#.####.###.##...#..
...####.#..#..#.##.#....##..#.####....##...##..#...#......#.#.......#.....
..##..####..#...#.#.#...##..#.#..###..#####........#..####......#..#

#..#.
#....
##..#
..#..
..###
`

var (
	input = sample
)

const (
	Ex  = 1
	Ex2 = 2
)


type puzzle struct {
	v             [][]byte
	width, height int
	base          byte
}

func (p puzzle) String() string {
	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("%d, %d (%c)\n", p.width, p.height, toRune(p.base)))
	for _, row := range p.v {
		for _, v := range row {
			buf.WriteRune(toRune(v))
		}
		buf.WriteRune('\n')
	}
	return buf.String()
}

func toRune(v byte) rune {
	if v == 0 {
		return '.'
	} else {
		return '#'
	}
}

func (p puzzle) Lit() int {
	var ret int
	for _, row := range p.v {
		for _, v := range row {
			if v > 0 {
				ret++
			}
		}
	}
	return ret
}

func (p *puzzle) Enhance(key []byte) *puzzle {
	ret := puzzle{
		v:      make([][]byte, p.height),
		width:  p.width,
		height: p.height,
	}

	if p.base == 0 {
		ret.base = key[0]
	} else {
		ret.base = key[511]
	}

	for y := range p.v {
		ret.v[y] = make([]byte, p.width)
		for x := range ret.v[y] {
			val := p.nineSquare(x, y)
			if key[val] > 0 {
				ret.v[y][x] = 1
			}
		}
	}

	return &ret
}

func (p *puzzle) Zoom() {
	baseRow := make([]byte, p.width+Ex2)
	for i := range baseRow {
		baseRow[i] = p.base
	}
	newRows := make([][]byte, 0, p.height+Ex2)
	for i := 0; i < Ex; i++ {
		newRows = append(newRows, baseRow)
	}
	// Embiggen the original rows.
	baseVals := make([]byte, Ex)
	for i := range baseVals {
		baseVals[i] = p.base
	}
	for _, row := range p.v {
		newRow := append(append(baseVals, row...), baseVals...)
		newRows = append(newRows, newRow)
	}
	for i := 0; i < Ex; i++ {
		newRows = append(newRows, baseRow)
	}

	p.width += Ex2
	p.height += Ex2
	p.v = newRows
}

func (p puzzle) nineSquare(px int, py int) int {
	ret := 0
	for y := py - 1; y <= py+1; y++ {
		for x := px - 1; x <= px+1; x++ {
			ret <<= 1
			if p.val(x, y) > 0 {
				ret++
			}
		}
	}

	return ret
}

func (p puzzle) val(x int, y int) byte {
	if x < 0 || x >= p.width {
		return p.base
	}
	if y < 0 || y >= p.height {
		return p.base
	}
	return p.v[y][x]
}

func main() {
	key := make([]byte, 0, 512)
	input := strings.Split(input, "\n")
	for i, line := range input {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		for _, v := range []byte(line) {
			switch v {
			case '.':
				key = append(key, 0)
			case '#':
				key = append(key, 1)
			default:
				panic(line)
			}
		}
		if len(key) == 512 {
			input = input[i+1:]
			break
		}
	}

	var puz puzzle
	for _, line := range input {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		row := make([]byte, 0, 256)
		for _, v := range []byte(line) {
			switch v {
			case '.':
				row = append(row, 0)
			case '#':
				row = append(row, 1)
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
	for i := 0; i < 50; i++ {
		p.Zoom()
		//fmt.Println(p)
		p = p.Enhance(key)
		fmt.Println(p)
		fmt.Println(p.Lit())
	}
}
