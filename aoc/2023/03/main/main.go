package main

import (
	"fmt"
	"strconv"
	"strings"
)

var sample = `
467..114..
...*......
..35..633.
......#...
617*......
.....+.58.
..592.....
......755.
...$.*....
.664.598..
`

type part struct {
	num           int
	xmin, xmax, y int
	active        bool
}

func (p part) adjacent(sym symbol) bool {
	return abs(sym.y-p.y) <= 1 && p.xmin-1 <= sym.x && sym.x <= p.xmax+1
}

type symbol struct {
	c    string
	x, y int
}

func main() {
	y := 0
	var parts []part
	var symbols []symbol
	for _, line := range strings.Split(sample, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		var acc []byte
		flushAcc := func(x int) {
			if len(acc) == 0 {
				return
			}
			num := mustInt(string(acc))

			parts = append(parts, part{
				num:    num,
				xmin:   x - len(acc),
				xmax:   x - 1,
				y:      y,
				active: false,
			})
			acc = nil
		}
		for x, c := range line {
			switch {
			case c == '.':
				flushAcc(x)
			case isDigit(c):
				acc = append(acc, byte(c))
			default:
				// symbol
				flushAcc(x)
				symbols = append(symbols, symbol{string(c), x, y})
			}
		}
		flushAcc(len(line))
		y++
	}

	// Now find which parts we can activate
	for _, sym := range symbols {
		for i := range parts {
			p := &parts[i]
			if p.adjacent(sym) {
				p.active = true
			}
		}
	}

	sum := 0
	for _, p := range parts {
		if p.active {
			sum += p.num
		}
	}
	fmt.Println(sum)

	// Now find any gears and compute the ratio
	ratSum := 0
	for _, sym := range symbols {
		if sym.c != "*" {
			continue
		}
		var adj []part
		for _, p := range parts {
			if p.adjacent(sym) {
				adj = append(adj, p)
			}
		}
		if len(adj) == 2 {
			rat := adj[0].num * adj[1].num
			ratSum += rat
		}
	}
	fmt.Println(ratSum)
}

func isDigit(c rune) bool {
	return c >= '0' && c <= '9'
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func mustInt(s string) int {
	if v, err := strconv.Atoi(s); err != nil {
		panic(fmt.Sprint(s, err))
	} else {
		return v
	}
}
