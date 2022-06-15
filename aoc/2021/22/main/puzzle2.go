package main

import (
	"fmt"
)

type puzzle2 struct {
	cubes []cube
}

type cube struct {
	x1, x2, y1, y2, z1, z2 int
}

func (c cube) volume() int64 {
	return int64(c.x2-c.x1+1) * int64(c.y2-c.y1+1) * int64(c.z2-c.z1+1)
}

func (c cube) isReal() bool {
	return (c.x1 <= c.x2) &&
		(c.y1 <= c.y2) &&
		(c.z1 <= c.z2)
}

func (c cube) contains(o cube) bool {
	return (c.x1 <= o.x1) &&
		(c.y1 <= o.y1) &&
		(c.z1 <= o.z1) &&
		(c.x2 >= o.x2) &&
		(c.y2 >= o.y2) &&
		(c.z2 >= o.z2)
}

func (c cube) breakUp(ic cube) []cube {
	// Break this cube into 27 smaller cubes, throwing away any no-ops.
	type seg struct {
		s, e int
	}

	xSegs := []seg{{c.x1, ic.x1 - 1}, {ic.x1, ic.x2}, {ic.x2 + 1, c.x2}}
	ySegs := []seg{{c.y1, ic.y1 - 1}, {ic.y1, ic.y2}, {ic.y2 + 1, c.y2}}
	zSegs := []seg{{c.z1, ic.z1 - 1}, {ic.z1, ic.z2}, {ic.z2 + 1, c.z2}}

	ret := make([]cube, 0, 27)
	for _, zSeg := range zSegs {
		for _, ySeg := range ySegs {
			for _, xSeg := range xSegs {
				newCube := cube{
					x1: xSeg.s,
					x2: xSeg.e,
					y1: ySeg.s,
					y2: ySeg.e,
					z1: zSeg.s,
					z2: zSeg.e,
				}
				if newCube.isReal() {
					ret = append(ret, newCube)
				}
			}

		}
	}
	return ret
}

func (p *puzzle2) Lit() int64 {
	var sum int64
	for _, c := range p.cubes {
		sum += c.volume()
	}
	return sum
}

func (p *puzzle2) Set(v byte, x1, x2, y1, y2, z1, z2 int) {
	x1, x2 = orderRange(x1, x2)
	y1, y2 = orderRange(y1, y2)
	z1, z2 = orderRange(z1, z2)

	// See if the new cube intersects any existing cubes.
	p.cubes = func(c cube) []cube {
		newCubes := []cube{}

		for _, e := range p.cubes {
			interCube, typ := intersect(c, e)
			//fmt.Println(typ)
			switch typ {
			case NONE:
				newCubes = append(newCubes, e)
			case A_CONTAINS_B:
				// the new cube fully contains an existing cube; either way, remove the existing cube
			case B_CONTAINS_A:
				// the new cube is fully contained by an existing cube
				if v > 0 {
					return p.cubes // this new cube is a no-op, bail on everything
				}
				fallthrough // we're going to have to break the existing cube
			case DIVERGE:
				replacements := e.breakUp(interCube)
				//fmt.Println("broke up: ", e, interCube, len(replacements))
				for _, r := range replacements {
					if r != interCube {
						newCubes = append(newCubes, r)
					}
				}
			}
		}

		if v > 0 {
			newCubes = append(newCubes, c)
		}
		return newCubes
	}(cube{
		x1: x1,
		x2: x2,
		y1: y1,
		y2: y2,
		z1: z1,
		z2: z2,
	})

	fmt.Printf("PUZZLE 2(%d): %d: x=%d..%d,y=%d..%d,z=%d..%d -> %d\n", len(p.cubes), v, x1, x2, y1, y2, z1, z2, p.Lit())

	// Check intersects.
	//for i, c1 := range p.cubes {
	//	for j, c2 := range p.cubes {
	//		if i == j {
	//			continue
	//		}
	//		if _, typ := intersect(c1, c2); typ != NONE {
	//			panic(fmt.Sprint("intersect", c1, c2, typ))
	//		}
	//	}
	//}
}

type intersection int

func (i intersection) String() string {
	switch i {
	case NONE:
		return "NONE"
	case A_CONTAINS_B:
		return "A_CONTAINS_B"
	case B_CONTAINS_A:
		return "B_CONTAINS_A"
	case DIVERGE:
		return "DIVERGE"
	default:
		return "ERROR"
	}
}

const (
	NONE = intersection(iota)
	A_CONTAINS_B
	B_CONTAINS_A
	DIVERGE
)

func intersect(a cube, b cube) (cube, intersection) {
	if a.contains(b) {
		return b, A_CONTAINS_B
	}
	if b.contains(a) {
		return b, B_CONTAINS_A
	}

	ret := cube{
		x1: max(a.x1, b.x1),
		x2: min(a.x2, b.x2),
		y1: max(a.y1, b.y1),
		y2: min(a.y2, b.y2),
		z1: max(a.z1, b.z1),
		z2: min(a.z2, b.z2),
	}

	if ret.isReal() {
		return ret, DIVERGE
	}

	return cube{}, NONE
}

func max(p1 int, p2 int) int {
	if p1 > p2 {
		return p1
	}
	return p2
}

func min(p1 int, p2 int) int {
	if p1 < p2 {
		return p1
	}
	return p2
}

func orderRange(p1 int, p2 int) (int, int) {
	if p1 > p2 {
		return p2, p1
	}
	return p1, p2
}
