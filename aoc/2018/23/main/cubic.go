package main

import (
	"fmt"
	"golang.org/x/exp/slices"
)

const (
	XYZ = 0
	XYz = 1
	XyZ = 2
	Xyz = 3
)

type octa struct {
	r   [4]rng
	ids []int
}

func (oc octa) volume() int {
	ret := 1
	for _, r := range oc.r {
		ret *= r.max - r.min
	}
	return ret
}

func (oc octa) inRange(pos [3]int) bool {
	x, y, z := pos[0], pos[1], pos[2]
	test := [4]int{
		XYZ: x + y + z,
		XYz: x + y - z,
		XyZ: x - y + z,
		Xyz: x - y - z,
	}

	for i, v := range test {
		if v < oc.r[i].min || v > oc.r[i].max {
			return false
		}
	}
	return true
}

func (oc octa) intersect(o octa) (octa, bool) {
	// this might be completely wrong
	var ret octa
	for dim := range oc.r {
		r := rng{
			min: max(oc.r[dim].min, o.r[dim].min),
			max: min(oc.r[dim].max, o.r[dim].max),
		}
		//if dim == Xyz {
		//	r.min, r.max = r.max, r.min
		//}
		if r.max < r.min {
			return octa{}, false
		}
		ret.r[dim] = r
	}

	// Are the resulting vertices in range of both of the inputs?
	fmt.Println(oc.r, oc.String(), oc.vertices())
	fmt.Println(o.r, o.String(), o.vertices())
	fmt.Println(ret.r)
	fmt.Println(ret.String())
	fmt.Println(ret.vertices())
	for _, v := range ret.vertices() {
		if !oc.inRange(v) || !o.inRange(v) {
			panic("out of range")
		}
	}

	ret.ids = zipperMerge(oc.ids, o.ids)
	return ret, true
}

func (oc octa) key() [4]rng {
	return oc.r
}

func (oc octa) loc() [3]int {
	// The center point of any pair of vertices should work.
	vertices := oc.vertices()
	loc := [3]int{
		mean(vertices[0][0], vertices[1][0]),
		mean(vertices[0][1], vertices[1][1]),
		mean(vertices[0][2], vertices[1][2]),
	}
	chk1 := [3]int{
		mean(vertices[2][0], vertices[3][0]),
		mean(vertices[2][1], vertices[3][1]),
		mean(vertices[2][2], vertices[3][2]),
	}
	chk2 := [3]int{
		mean(vertices[4][0], vertices[5][0]),
		mean(vertices[4][1], vertices[5][1]),
		mean(vertices[4][2], vertices[5][2]),
	}

	if loc != chk1 || loc != chk2 {
		panic(loc)
	}

	return loc
}

func (oc octa) rads() [4]int {
	var ret [4]int
	for i := range ret {
		ret[i] = oc.r[i].max - oc.r[i].min
	}
	return ret
}

func (oc octa) vertices() [6][3]int {
	// this is not remotely correct

	return [6][3]int{
		intersectPlanes(oc.r[XYZ].min, oc.r[XYz].min, oc.r[XyZ].min, oc.r[Xyz].min), // left
		intersectPlanes(oc.r[XYZ].max, oc.r[XYz].max, oc.r[XyZ].max, oc.r[Xyz].max), // right

		intersectPlanes(oc.r[XYZ].min, oc.r[XYz].min, oc.r[XyZ].max, oc.r[Xyz].max), // bottom
		intersectPlanes(oc.r[XYZ].max, oc.r[XYz].max, oc.r[XyZ].min, oc.r[Xyz].min), // top

		intersectPlanes(oc.r[XYZ].min, oc.r[XYz].max, oc.r[XyZ].min, oc.r[Xyz].max), // near
		intersectPlanes(oc.r[XYZ].max, oc.r[XYz].min, oc.r[XyZ].max, oc.r[Xyz].min), // far
	}
}

func (oc octa) String() string {
	return fmt.Sprintf("loc=%+v, rads=%+v %+v ", oc.loc(), oc.rads(), oc.ids)
}

func collapseRegionsCubic(bots []bot) octa {
	iter := 1
	octas := toOctas(bots)
	orig := octas

	fmt.Printf("iter=%d, len=%d, vol=%d\n", iter, len(octas), volume(octas))
	for len(octas) > 1 {
		printCubic(octas)
		octas = intersectAllOctas(iter, octas)
		octas = dedupOctas(octas)
		iter++
		fmt.Printf("iter=%d, len=%d, vol=%d\n", iter, len(octas), volume(octas))
		printCubic(octas)
	}

	// validate!
	final := octas[0]
	expectInRange := map[int]bool{}
	for _, id := range final.ids {
		expectInRange[id] = true
	}

	// traversing the entire octa seems hard... maybe just spot check the vertices and origin.
	vertices := final.vertices()
	for _, p := range append(vertices[:], final.loc()) {
		for id, r := range orig {
			fmt.Println(id, r, expectInRange[id], r.inRange(p))
			if expectInRange[id] != r.inRange(p) {
				panic(fmt.Sprintf("expectInRange(%t) != r.inRange(%t), id=%d", expectInRange[id], r.inRange(p), id))
			}
		}
	}

	return final
}

func intersectAllOctas(minLen int, oktas []octa) []octa {
	var ret []octa
	for i, dia := range oktas {
		found := false
		for j := i + 1; j < len(oktas); j++ {
			if newDia, ok := dia.intersect(oktas[j]); ok {
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

func dedupOctas(oktas []octa) []octa {
	slices.SortFunc(oktas, func(a, b octa) bool {
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
	last := &oktas[0]
	for i := 1; i < len(oktas); i++ {
		if oktas[i].key() == last.key() {
			last.ids = zipperMerge(last.ids, oktas[i].ids)
		} else {
			oktas[wIdx] = oktas[i]
			last = &oktas[wIdx]
			wIdx++
		}
	}
	return oktas[:wIdx]
}

func toOctas(bots []bot) []octa {
	var ret []octa
	for i, b := range bots {
		x := b.p[0]
		y := b.p[1]
		z := b.p[2]
		r := b.rad

		// quadruple the coordinate system for precise math
		x *= 4
		y *= 4
		z *= 4
		r *= 4

		oct := octa{
			r: [4]rng{
				XYZ: {
					min: x + y + z - r,
					max: x + y + z + r,
				},
				XYz: {
					min: x + y - z - r,
					max: x + y - z + r,
				},
				XyZ: {
					min: x - y + z - r,
					max: x - y + z + r,
				},
				Xyz: {
					min: x - y - z - r,
					max: x - y - z + r,
				}},
			ids: []int{i},
		}
		if oct.loc() != [3]int{x, y, z} {
			panic("wrong loc")
		}
		ret = append(ret, oct)
	}
	return ret
}

func printCubic(octas []octa) {
	for _, oc := range octas {
		fmt.Println(oc.String())
	}
}

func volume(oktas []octa) int {
	sum := 0
	for _, o := range oktas {
		sum += o.volume()
	}
	return sum
}

func intersectPlanes(xpypz, xpynz, xnypz, xnynz int) [3]int {
	x := mean(xpypz, xnynz)
	y := rads(xpypz, xnypz)
	z := rads(xpypz, xpynz)

	if x+y+z != xpypz {
		panic(fmt.Sprint(x+y+z, xpypz))
	}
	if x+y-z != xpynz {
		panic(fmt.Sprint(x+y-z, xpynz))
	}
	if x-y+z != xnypz {
		panic(fmt.Sprint(x-y+z, xnypz))
	}
	if x-y-z != xnynz {
		panic(fmt.Sprint(x-y-z, xnynz))
	}
	return [3]int{x, y, z}
}
