package main

import (
	"fmt"
	"slices"
	"strings"
)

const sample1a = `
AAAA
BBCD
BBCC
EEEC
`

const sample1b = `
OOOOO
OXOXO
OOOOO
OXOXO
OOOOO
`

const sample2a = `
EEEEE
EXXXX
EEEEE
EXXXX
EEEEE
`

const sample2b = `
AAAAAA
AAABBA
AAABBA
ABBAAA
ABBAAA
AAAAAA
`

const sampleBig = `
RRRRIICCFF
RRRRIICCCF
VVRRRCCFFF
VVRCCCJFFF
VVVVCJJCFE
VVIVCCJJEE
VVIIICJJEE
MIIIIIJJEE
MIIISIJEEE
MMMISSJEEE
`

func main() {
	run(parse(sample1a))
	run(parse(sample1b))
	run(parse(sample2a))
	run(parse(sample2b))
	run(parse(sampleBig))
}

func run(p *puz) {
	seen := map[pos]*region{}
	var regions []*region

	for x := 0; x < p.w; x++ {
		for y := 0; y < p.h; y++ {
			pt := pos{x, y}
			if seen[pt] == nil {
				r := &region{wat: p.data[pt], elements: map[pos]bool{}}
				regions = append(regions, r)
				p.mapRegion(r, pt)
				for el := range r.elements {
					seen[el] = r
				}
			}
		}
	}

	score := 0
	bulkScore := 0
	for _, r := range regions {
		score += r.score()
		bulkScore += r.bulkScore()
	}
	fmt.Printf("sz=%d, reg=%d, score=%d, bulkScore=%d\n", len(seen), len(regions), score, bulkScore)
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
		N: {p.x, p.y - 1},
		E: {p.x + 1, p.y},
		S: {p.x, p.y + 1},
		W: {p.x - 1, p.y},
	}
}

type puz struct {
	data map[pos]byte
	w, h int
}

func (p puz) mapRegion(r *region, start pos) {
	r.elements[start] = true
	q := []pos{start}
	for len(q) > 0 {
		pt := q[0]
		q = q[1:]
		for _, adj := range pt.adjacent() {
			if p.data[adj] == r.wat && !r.elements[adj] {
				r.elements[adj] = true
				q = append(q, adj)
			}
		}
	}
}

// Part 1: each element side not adjacent to any other element in the same region is a perimeter side.

type region struct {
	wat      byte
	elements map[pos]bool
}

func (r region) score() int {
	area := len(r.elements)

	// each element side not adjacent to any other element is a perimeter side
	perimeter := 0
	for el := range r.elements {
		for _, adj := range el.adjacent() {
			if !r.elements[adj] {
				perimeter++
			}
		}
	}

	return area * perimeter
}

// Part 2: find all perimeter sides, then merge them into fence segments adjacent to a sibling fence segment.

type dir int

const (
	N = dir(iota)
	E
	S
	W
)

type fence struct {
	side dir // which side of the square the fence is on
	p    pos // which square this fence belongs to
}

func (f fence) key() fenceKey {
	return fenceKey{side: f.side, rowCol: f.rowCol()}
}

// either the row or column this fence can merge on
func (f fence) rowCol() int {
	switch f.side {
	case N, S:
		return f.p.y
	case E, W:
		return f.p.x
	default:
		panic(f)
	}
}

// the "index" position of this fence within its rowCol (opposite rowCol)
func (f fence) ordinal() int {
	switch f.side {
	case N, S:
		return f.p.x
	case E, W:
		return f.p.y
	default:
		panic(f)
	}
}

// key type to bucket fences on the same "line" with each other
type fenceKey struct {
	side   dir // which side of the square the fence is on
	rowCol int // either the row or column this fence can merge on
}

func (r region) bulkScore() int {
	area := len(r.elements)

	// collect all fence segments
	var fences []fence
	for el := range r.elements {
		for d, adj := range el.adjacent() {
			if !r.elements[adj] {
				fences = append(fences, fence{dir(d), el})
			}
		}
	}

	// bucket each piece of fence into its own list; each list contains all fence segments the represent the same
	// "direction" plus the same row/col.
	buckets := map[fenceKey][]fence{}
	for _, f := range fences {
		k := f.key()
		buckets[k] = append(buckets[k], f)
	}

	// now count the number of distinct "runs" in each bucket
	sides := 0
	for _, b := range buckets {
		slices.SortFunc(b, func(a, b fence) int {
			return a.ordinal() - b.ordinal()
		})

		last := -2
		for _, f := range b {
			if f.ordinal() != last+1 {
				sides++
			}
			last = f.ordinal()
		}
	}

	return area * sides
}
