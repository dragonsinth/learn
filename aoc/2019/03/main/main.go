package main

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
)

var samples = []string{`R8,U5,L5,D3
U7,R6,D4,L4`,
	`R75,D30,R83,U83,L12,D49,R71,U7,L72
U62,R66,U55,R34,D71,R55,D58,R83`,
	`R98,U47,R26,D63,R33,U87,L62,D20,R33,U53,R51
U98,R91,D20,R16,D67,R40,U7,R15,U6,R7`,
}

func main() {
	run(samples[0])
	run(samples[1])
	run(samples[2])
}

type wire struct {
	h []hseg
	v []vseg

	hhead *hseg
	vhead *vseg
}

func (w *wire) crawl(x int, y int) int {
	sum := 0
	hhead, vhead := w.hhead, w.vhead
	for {
		if hhead != nil {
			if hhead.has(x, y) {
				sum += abs(x - hhead.x1)
				return sum
			}
			sum += abs(hhead.x1 - hhead.x2)
			vhead = hhead.next
			hhead = nil
		} else if vhead != nil {
			if vhead.has(x, y) {
				sum += abs(y - vhead.y1)
				return sum
			}
			sum += abs(vhead.y1 - vhead.y2)
			hhead = vhead.next
			vhead = nil
		} else {
			panic("fell off")
		}
	}
}

type hseg struct {
	x1, x2 int
	y      int

	next *vseg
}

func (h hseg) has(x int, y int) bool {
	if h.y != y {
		return false
	}
	if h.x1 < h.x2 {
		return h.x1 <= x && x <= h.x2
	} else {
		return h.x2 <= x && x <= h.x1
	}
}

type vseg struct {
	x      int
	y1, y2 int

	next *hseg
}

func (v vseg) has(x int, y int) bool {
	if v.x != x {
		return false
	}
	if v.y1 < v.y2 {
		return v.y1 <= y && y <= v.y2
	} else {
		return v.y2 <= y && y <= v.y1
	}
}

func run(input string) {
	lines := strings.Split(input, "\n")
	wires := make([]*wire, 2)
	wIdx := 0
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		wires[wIdx] = parse(line)
		wIdx++
	}

	// Try to intersect each wire.
	int1 := intersect(wires[0].h, wires[1].v)
	int2 := intersect(wires[1].h, wires[0].v)
	ints := append(int1, int2...)
	sort.Slice(ints, func(i, j int) bool {
		return ints[i].d < ints[j].d
	})

	fmt.Println(ints[0].x, ints[0].y, ints[0].d)

	best := math.MaxInt32
	for _, i := range ints {
		s1 := wires[0].crawl(i.x, i.y)
		s2 := wires[1].crawl(i.x, i.y)
		if s1+s2 < best {
			best = s1 + s2
			fmt.Println(i.x, i.y, s1, s2, best)
		}
	}

	fmt.Println()
}

type intersection struct {
	x, y, d int
}

func intersect(h []hseg, v []vseg) []intersection {
	var ret []intersection
	for _, hs := range h {
		for _, vs := range v {
			if hs.has(vs.x, hs.y) && vs.has(vs.x, hs.y) {
				if vs.x != 0 || hs.y != 0 {
					ret = append(ret, intersection{
						x: vs.x,
						y: hs.y,
						d: abs(vs.x) + abs(hs.y),
					})
				}
			}
		}
	}
	return ret
}

func abs(n int) int {
	if n >= 0 {
		return n
	} else {
		return -n
	}
}

func parse(line string) *wire {
	var ret wire
	x, y := 0, 0
	insegs := strings.Split(line, ",")
	lastV := &vseg{}
	lastH := &hseg{}
	for i, inseg := range insegs {
		v := mustInt(inseg[1:])
		switch inseg[0] {
		case 'U':
			ret.v = append(ret.v, vseg{
				x:  x,
				y1: y,
				y2: y + v,
			})
			y = y + v
			lastV = &ret.v[len(ret.v)-1]
			lastH.next = lastV
			lastH = nil
			if i == 0 {
				ret.vhead = lastV
			}
		case 'L':
			ret.h = append(ret.h, hseg{
				x1: x,
				x2: x - v,
				y:  y,
			})
			x = x - v
			lastH = &ret.h[len(ret.h)-1]
			lastV.next = lastH
			lastV = nil
			if i == 0 {
				ret.hhead = lastH
			}
		case 'D':
			ret.v = append(ret.v, vseg{
				x:  x,
				y1: y,
				y2: y - v,
			})
			y = y - v
			lastV = &ret.v[len(ret.v)-1]
			lastH.next = lastV
			lastH = nil
			if i == 0 {
				ret.vhead = lastV
			}
		case 'R':
			ret.h = append(ret.h, hseg{
				x1: x,
				x2: x + v,
				y:  y,
			})
			x = x + v
			lastH = &ret.h[len(ret.h)-1]
			lastV.next = lastH
			lastV = nil
			if i == 0 {
				ret.hhead = lastH
			}
		}
	}
	return &ret
}

func mustInt(s string) int {
	if v, err := strconv.Atoi(s); err != nil {
		panic(fmt.Sprint(s, err))
	} else {
		return v
	}
}
