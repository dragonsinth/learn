package main

import (
	"fmt"
	"strconv"
	"strings"
)

var sample = `
F10
N3
F7
R90
F11
`

var (
	input = sample
)

type dir byte

const (
	EAST = dir(iota)
	SOUTH
	WEST
	NORTH
)

func (d dir) String() string {
	return string(d.toRune())
}

func (d dir) toRune() rune {
	return rune("ESWN"[d])
}

func (d dir) direction() direction {
	return directions[d]
}

type direction struct {
	dx, dy int
}

var (
	directions = []direction{
		{+1, +0},
		{+0, -1},
		{-1, +0},
		{+0, +1},
	}
)

func main() {
	var sx, sy int
	wx, wy := 10, 1
	for _, line := range strings.Split(input, "\n") {
		fmt.Println(sx, sy, wx, wy)

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		val, err := strconv.Atoi(line[1:])
		if err != nil {
			panic(err)
		}

		switch line[0] {
		case 'N':
			wy += val
		case 'S':
			wy -= val
		case 'E':
			wx += val
		case 'W':
			wx -= val
		case 'L':
			if val%90 != 0 {
				panic(val)
			}
			for ; val > 0; val -= 90 {
				wx, wy = -wy, wx
			}
		case 'R':
			if val%90 != 0 {
				panic(val)
			}
			for ; val > 0; val -= 90 {
				wx, wy = wy, -wx
			}
		case 'F':
			sx, sy = sx+val*wx, sy+val*wy
		default:
			panic(line[0])
		}
	}

	fmt.Println(sx, sy, abs(sx)+abs(sy))
}

func main1() {
	var x, y int
	d := EAST
	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		val, err := strconv.Atoi(line[1:])
		if err != nil {
			panic(err)
		}

		switch line[0] {
		case 'N':
			y += val
		case 'S':
			y -= val
		case 'E':
			x += val
		case 'W':
			x -= val
		case 'L':
			if val%90 != 0 {
				panic(val)
			}
			d = d - dir(val/90)
			d = d % 4
		case 'R':
			if val%90 != 0 {
				panic(val)
			}
			d = d + dir(val/90)
			d = d % 4
		case 'F':
			delta := d.direction()
			x, y = x+val*delta.dx, y+val*delta.dy
		default:
			panic(line[0])
		}
	}

	fmt.Println(x, y, abs(x)+abs(y), d.String())
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
