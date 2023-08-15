package main

import (
	"bytes"
	"fmt"
	"github.com/dragonsinth/learn/aoc/grid"
	"github.com/dragonsinth/learn/aoc/termbox"
	"golang.org/x/exp/slices"
	"math"
	"os"
)

const (
	XY = 0
	Xy = 1
)

type diamond struct {
	r   [2]rng
	ids []int
}

func (d diamond) area() int {
	ret := 1
	for _, r := range d.r {
		ret *= r.max - r.min
	}
	return ret
}

func (d diamond) inRange(pos [2]int) bool {
	x, y := pos[0], pos[1]
	test := [2]int{
		x + y,
		x - y,
	}

	for i, v := range test {
		if v < d.r[i].min || v > d.r[i].max {
			return false
		}
	}
	return true
}

func (d diamond) intersect(o diamond) (diamond, bool) {
	var ret diamond
	for dim := range d.r {
		r := rng{
			min: max(d.r[dim].min, o.r[dim].min),
			max: min(d.r[dim].max, o.r[dim].max),
		}
		if r.max < r.min {
			return diamond{}, false
		}
		ret.r[dim] = r
	}
	ret.ids = zipperMerge(d.ids, o.ids)
	return ret, true
}

func (d diamond) key() [2]rng {
	return d.r
}

func (d diamond) loc() [2]int {
	// The center point of either pair of vertices should work.
	vertices := d.vertices()
	loc := [2]int{
		mean(vertices[0][0], vertices[1][0]),
		mean(vertices[0][1], vertices[1][1]),
	}
	chk := [2]int{
		mean(vertices[2][0], vertices[3][0]),
		mean(vertices[2][1], vertices[3][1]),
	}

	if loc != chk {
		panic(fmt.Sprint("loc != chk", loc, chk))
	}

	return loc
}

func (d diamond) rads() [2]int {
	var ret [2]int
	for i := range ret {
		ret[i] = rads(d.r[i].max, d.r[i].min)
	}
	return ret
}

func (d diamond) vertices() [4][2]int {
	return [4][2]int{
		intersectLines(d.r[XY].min, d.r[Xy].min), // left
		intersectLines(d.r[XY].max, d.r[Xy].max), // right

		intersectLines(d.r[XY].min, d.r[Xy].max), // bottom
		intersectLines(d.r[XY].max, d.r[Xy].min), // top
	}
}

func collapseRegionsPlanar(dx, dy int, bots []bot) diamond {
	iter := 1
	diamonds := toDiamonds(dx, dy, bots)
	orig := diamonds

	fmt.Printf("iter=%d, len=%d, area=%d\n", iter, len(diamonds), area(diamonds))
	printPlanar(diamonds)
	for len(diamonds) > 1 {
		diamonds = intersectAllDiamonds(iter, diamonds)
		diamonds = dedupDiamonds(diamonds)
		iter++
		fmt.Printf("iter=%d, len=%d, area=%d\n", iter, len(diamonds), area(diamonds))
		printPlanar(diamonds)
	}

	// validate!
	final := diamonds[0]
	expectInRange := map[int]bool{}
	for _, id := range final.ids {
		expectInRange[id] = true
	}

	// find the rect circumscribing the diamond
	xmin, xmax := math.MaxInt, math.MinInt
	ymin, ymax := math.MaxInt, math.MinInt
	for _, v := range final.vertices() {
		xmin = min(xmin, v[0])
		xmax = max(xmax, v[0])
		ymin = min(ymin, v[1])
		ymax = max(ymax, v[1])
	}

	for y := ymin; y <= ymax; y++ {
		for x := xmin; x <= xmax; x++ {
			pt := [2]int{x, y}
			if !final.inRange(pt) {
				continue
			}

			for id, r := range orig {
				if expectInRange[id] != r.inRange(pt) {
					panic(fmt.Sprintf("expectInRange(%v) != r.inRange(%v), id=%d", expectInRange[id], r.inRange(pt), id))
				}
			}
		}
	}

	return final
}

func intersectAllDiamonds(minLen int, dias []diamond) []diamond {
	var ret []diamond
	for i, dia := range dias {
		found := false
		for j := i + 1; j < len(dias); j++ {
			if newDia, ok := dia.intersect(dias[j]); ok {
				if len(newDia.ids) >= minLen {
					ret = append(ret, newDia)
					found = true
				}
			}
		}
		if !found && len(dia.ids) >= minLen {
			ret = append(ret, dia)
		}
	}
	return ret
}

func dedupDiamonds(dias []diamond) []diamond {
	slices.SortFunc(dias, func(a, b diamond) bool {
		for i := range a.r {
			if a.r[i].min != b.r[i].min {
				return a.r[i].min < b.r[i].min
			}
			if a.r[i].max != b.r[i].max {
				return a.r[i].max < b.r[i].max
			}
		}
		return false
	})

	wIdx := 1
	last := &dias[0]
	for i := 1; i < len(dias); i++ {
		if dias[i].key() == last.key() {
			last.ids = zipperMerge(last.ids, dias[i].ids)
		} else {
			dias[wIdx] = dias[i]
			last = &dias[wIdx]
			wIdx++
		}
	}
	return dias[:wIdx]
}

func toDiamonds(dx int, dy int, bots []bot) []diamond {
	var ret []diamond
	for i, b := range bots {
		x := b.p[dx]
		y := b.p[dy]
		r := b.rad

		// quadruple the coordinate system for precise math
		x *= 4
		y *= 4
		r *= 4

		dia := diamond{
			r: [2]rng{
				XY: {
					min: x + y - r,
					max: x + y + r,
				},
				Xy: {
					min: x - y - r,
					max: x - y + r,
				}},
			ids: []int{i},
		}
		if dia.loc() != [2]int{x, y} {
			panic("wrong loc")
		}
		ret = append(ret, dia)
	}
	return ret
}

func printPlanar(dias []diamond) {
	draw := grid.Alloc2d(60, 60, byte(0))
	for _, dia := range dias {
		if dia.area() > 100 {
			continue
		}
		l := byte(len(dia.ids))

		// find the rect circumscribing the diamond
		xmin, xmax := math.MaxInt, math.MinInt
		ymin, ymax := math.MaxInt, math.MinInt
		for _, v := range dia.vertices() {
			xmin = min(xmin, v[0])
			xmax = max(xmax, v[0])
			ymin = min(ymin, v[1])
			ymax = max(ymax, v[1])
		}

		for y := ymin; y <= ymax; y++ {
			for x := xmin; x <= xmax; x++ {
				if !dia.inRange([2]int{x, y}) {
					continue
				}

				if draw[y/2][x/2] < l {
					draw[y/2][x/2] = l
				}
			}
		}
	}

	for y := range draw {
		for x := range draw[0] {
			draw[y][x] = ".:;|oO*"[draw[y][x]]
		}
	}

	termbox.RenderPlain(draw, os.Stdout)

	for _, dia := range dias {
		var buf bytes.Buffer
		buf.WriteString(fmt.Sprintf(" loc=%+v, rads=%+v", dia.loc(), dia.rads()))
		buf.WriteString(fmt.Sprintf(" %+v", dia.ids))
		fmt.Println(buf.String())
	}
}

func area(dias []diamond) int {
	sum := 0
	for _, dia := range dias {
		sum += dia.area()
	}
	return sum
}

func intersectLines(xpy, xny int) [2]int {
	x := mean(xpy, xny)
	y := rads(xpy, xny)
	if x+y != xpy {
		panic(xpy)
	}
	if x-y != xny {
		panic(xny)
	}
	return [2]int{x, y}
}
