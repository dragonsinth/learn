package main

import (
	"fmt"
)

const (
	size = 50
)

type puzzle struct {
	v [][][]byte
}

func (p *puzzle) Lit() int64 {
	var ret int64
	for _, plane := range p.v {
		for _, row := range plane {
			for _, v := range row {
				if v > 0 {
					ret++
				}
			}
		}
	}
	return ret
}

func (p *puzzle) Set(v byte, x1, x2, y1, y2, z1, z2 int) {
	x1, x2, ok := orderRangeLimit(x1, x2)
	if !ok {
		return
	}
	y1, y2, ok = orderRangeLimit(y1, y2)
	if !ok {
		return
	}
	z1, z2, ok = orderRangeLimit(z1, z2)
	if !ok {
		return
	}

	for z := z1; z <= z2; z++ {
		for y := y1; y <= y2; y++ {
			for x := x1; x <= x2; x++ {
				// translate -50..50 -> 0...10
				p.v[z+size][y+size][x+size] = v
			}
		}
	}
	fmt.Printf("PUZZLE 1: %d: x=%d..%d,y=%d..%d,z=%d..%d -> %d\n", v, x1, x2, y1, y2, z1, z2, p.Lit())
}

func orderRangeLimit(p1 int, p2 int) (int, int, bool) {
	if p1 > p2 {
		p1, p2 = p2, p1
	}
	if p2 < -size || p1 > size {
		return 0, 0, false
	}
	if p1 < -size {
		p1 = -size
	}
	if p2 > size {
		p2 = size
	}
	return p1, p2, true
}
