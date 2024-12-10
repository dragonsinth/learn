package main

import (
	"fmt"
	"strings"
)

const sample = `
89010123
78121874
87430965
96549874
45678903
32019012
01329801
10456732
`

func main() {
	p := parse(sample)
	part1(p)
	part2(p)
}

func part1(p *puz) {
	sum := 0
	for x := 0; x < p.w; x++ {
		for y := 0; y < p.h; y++ {
			sum += p.score(pos{x, y})
		}
	}
	fmt.Println(sum)
}

func part2(p *puz) {
	sum := 0
	for x := 0; x < p.w; x++ {
		for y := 0; y < p.h; y++ {
			sum += p.rating(pos{x, y})
		}
	}
	fmt.Println(sum)
}

func parse(input string) *puz {
	p := &puz{
		data: map[pos]byte{},
	}

	lines := strings.Split(input, "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		if p.w <= 0 {
			p.w = len(line)
		} else if p.w != len(line) {
			panic(line)
		}
		y := p.h
		for x, wat := range line {
			pt := pos{x: x, y: y}
			p.data[pt] = byte(wat)
		}
		p.h++
	}
	return p
}

type pos struct {
	x, y int
}

func (p pos) adjacent() [4]pos {
	return [4]pos{
		{p.x, p.y - 1},
		{p.x + 1, p.y},
		{p.x, p.y + 1},
		{p.x - 1, p.y},
	}
}

type puz struct {
	data map[pos]byte
	w, h int
}

func (p puz) score(pt pos) int {
	seen := map[pos]bool{}
	return p.scoreRecurse(seen, pt, '0')
}

func (p puz) scoreRecurse(seen map[pos]bool, pt pos, test byte) int {
	if seen[pt] {
		return 0
	}
	if p.data[pt] != test {
		return 0
	}
	seen[pt] = true
	if test == '9' {
		return 1
	}
	sum := 0
	for _, next := range pt.adjacent() {
		sum += p.scoreRecurse(seen, next, test+1)
	}
	return sum
}

func (p puz) rating(pt pos) int {
	return p.ratingRecurse(pt, '0')
}

func (p puz) ratingRecurse(pt pos, test byte) int {
	if p.data[pt] != test {
		return 0
	}
	if test == '9' {
		return 1
	}
	sum := 0
	for _, next := range pt.adjacent() {
		sum += p.ratingRecurse(next, test+1)
	}
	return sum
}
