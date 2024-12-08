package main

import (
	"fmt"
	"strings"
)

const sample = `
............
........0...
.....0......
.......0....
....0.......
......A.....
............
............
........A...
.........A..
............
............
`

func main() {
	p := parse(sample)
	part1(p)
	part2(p)
}

func part1(p *puz) {
	found := map[pos]bool{}
	for _, v := range p.data {
		v.findAntiNodes(found, p.w, p.h)
	}
	fmt.Println(len(found))
}

func part2(p *puz) {
	found := map[pos]bool{}
	for _, v := range p.data {
		v.findAntiNodesWithHarmonics(found, p.w, p.h)
	}
	fmt.Println(len(found))
}

func parse(input string) *puz {
	p := &puz{
		data: map[byte]*pointSet{},
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
			if wat != '.' {
				ps := p.getPointSet(byte(wat))
				ps.pts = append(ps.pts, pt)
			}
		}
		p.h++
	}
	return p
}

type pos struct {
	x, y int
}

type pointSet struct {
	wat byte
	pts []pos
}

func (ps pointSet) findAntiNodes(found map[pos]bool, w int, h int) {
	maybeAdd := func(pt pos) {
		if pt.x < 0 || pt.x >= w || pt.y < 0 || pt.y >= h {
			return
		}
		found[pt] = true
	}

	for i, a := range ps.pts {
		for j, b := range ps.pts {
			if i == j {
				continue
			}
			dx := b.x - a.x
			dy := b.y - a.y
			maybeAdd(pos{a.x - dx, a.y - dy})
			maybeAdd(pos{b.x + dx, b.y + dy})
		}
	}
}

func (ps pointSet) findAntiNodesWithHarmonics(found map[pos]bool, w int, h int) {
	maybeAdd := func(pt pos) {
		if pt.x < 0 || pt.x >= w || pt.y < 0 || pt.y >= h {
			return
		}
		found[pt] = true
	}

	for i, a := range ps.pts {
		for j, b := range ps.pts {
			if i == j {
				continue
			}
			dx := b.x - a.x
			dy := b.y - a.y

			// use gcd to reduce e.g. 3, 9 -> 1, 3
			g := gcd(dx, dy)
			dx /= g
			dy /= g

			// hack, just iterate a bunch of linear points, don't bother precomputing bounds
			for i := -50; i <= 50; i++ {
				maybeAdd(pos{a.x + dx*i, a.y + dy*i})
			}
		}
	}
}

type puz struct {
	data map[byte]*pointSet
	w, h int
}

func (p *puz) getPointSet(wat byte) *pointSet {
	if p.data[wat] == nil {
		p.data[wat] = &pointSet{wat: wat}
	}
	return p.data[wat]
}

func gcd(a, b int) int {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}
