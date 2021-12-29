package main

import (
	"fmt"
	"strings"
)

var sample = `
..##.......
#...#...#..
.#....#..#.
..#.#...#.#
.#...##..#.
..#.##.....
.#.#.#....#
.#........#
#.##...#...
#...##....#
.#..#...#.#
`

var (
	input = sample
)

type puzzle struct {
	v             [][]byte
	width, height int
}

func (p puzzle) String() string {
	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("%d, %d\n", p.width, p.height))
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

func main() {
	var puz puzzle
	for _, line := range strings.Split(input, "\n") {
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

	tcs := []struct {
		dx, dy int
	}{
		{1, 1},
		{3, 1},
		{5, 1},
		{7, 1},
		{1, 2},
	}

	prod := 1
	for _, tc := range tcs {
		x, y := 0, 0
		trees := 0
		for y < puz.height {
			if puz.v[y][x%puz.width] != 0 {
				trees++
			}
			x += tc.dx
			y += tc.dy
		}
		fmt.Println(trees)
		prod *= trees
	}

	fmt.Println(prod)
}
