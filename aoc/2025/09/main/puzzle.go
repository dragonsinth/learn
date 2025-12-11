package main

import (
	"fmt"
	"slices"
)

type puz struct {
	pts []pos

	hlines []hline
	vlines []vline
}

type hline struct {
	x1, x2 int // exclusive! x1 and x2 are corners; y1 < y2
	y      int
}

type vline struct {
	x      int
	y1, y2 int // exclusive! y1 and y2 are corners; y1 < y2
}

func (p *puz) part1() {
	best := 0
	for i, a := range p.pts {
		for j, b := range p.pts {
			if i < j {
				ar := area(a, b)
				if ar > best {
					best = ar
				}
			}
		}
	}
	fmt.Println(best)
}

func (p *puz) part2() {
	// connect all the edges
	for i := range p.pts {
		if i == 0 {
			p.connect(p.pts[len(p.pts)-1], p.pts[0])
		} else {
			p.connect(p.pts[i-1], p.pts[i])
		}
	}

	slices.SortFunc(p.hlines, func(a, b hline) int {
		return a.y - b.y
	})
	slices.SortFunc(p.vlines, func(a, b vline) int {
		return a.x - b.x
	})

	best := 0
	for i, a := range p.pts {
		for j, b := range p.pts {
			if i < j {
				ar := area(a, b)
				if ar > best && p.isLegal(a, b) {
					best = ar
				}
			}
		}
	}
	fmt.Println(best)
}

func (p *puz) connect(a pos, b pos) {
	if a.x == b.x {
		y1, y2 := sort(a.y, b.y)
		p.vlines = append(p.vlines, vline{a.x, y1, y2})
	} else if a.y == b.y {
		x1, x2 := sort(a.x, b.x)
		p.hlines = append(p.hlines, hline{x1, x2, a.y})
	} else {
		panic(fmt.Sprint(a, b))
	}
}

// Check if the rectangle is legal by checking each vline/hline to be sure it doesn't intersect this rectangle.
func (p *puz) isLegal(a pos, b pos) bool {
	x1, x2 := sort(a.x, b.x)
	y1, y2 := sort(a.y, b.y)

	for _, h := range p.hlines {
		if h.y <= y1 || h.y >= y2 {
			// completely outside of range
			continue
		}
		if h.x1 <= x1 && h.x2 <= x1 {
			// both points are less than left edge
			continue
		}
		if h.x1 >= x2 && h.x2 >= x2 {
			// both points are greater than right edge
			continue
		}
		return false // illegal
	}

	for _, v := range p.vlines {
		if v.x <= x1 || v.x >= x2 {
			// completely outside of range
			continue
		}
		if v.y1 <= y1 && v.y2 <= y1 {
			// both points are less than top edge
			continue
		}
		if v.y1 >= y2 && v.y2 >= y2 {
			// both points are greater than bottom edge
			continue
		}
		return false // illegal
	}

	return true
}
