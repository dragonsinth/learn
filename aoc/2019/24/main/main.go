package main

import (
	"fmt"
	"strings"
)

var samples = []string{
	`
....#
#..#.
#.?##
..#..
#....`,
	`
.###.
##...
..?##
.#.#.
#.#.#`,
}

const (
	w   = 5
	h   = 5
	mid = 2
	sz  = w * h
)

func main() {
	for _, sample := range samples {
		run1(sample)
	}

	run2(samples[0], 10)
	run2(samples[1], 200)
}

func parse(input string) []byte {
	var data []byte
	for _, line := range strings.Split(input, "\n") {
		if line == "" {
			continue
		}
		data = append(data, []byte(line)...)
	}
	if len(data) != sz {
		panic(len(data))
	}
	return data
}

func run2(input string, iters int) {
	p2 := puz2{
		data: [401][]byte{},
		dnz:  -1,
		dxz:  1,
	}
	p2.data[200] = parse(input)
	for i := 0; i < iters; i++ {
		p2 = p2.next()
	}
	fmt.Println(p2.Score())
	for z := p2.minZ(); z <= p2.maxZ(); z++ {
		fmt.Printf("Depth %d:\n", z-200)
		for y := 0; y < h; y++ {
			fmt.Println(string(p2.data[z][y*w : y*w+w]))
		}
		fmt.Println()
	}
}

type puz2 struct {
	data     [401][]byte
	dnz, dxz int
}

func (p *puz2) minZ() int {
	return p.dnz/2 + 200
}

func (p *puz2) maxZ() int {
	return p.dxz/2 + 200
}

func (p *puz2) val(x, y, z int) int {
	if x == mid && y == mid {
		panic("here")
	}
	if z < p.minZ() || z > p.maxZ() {
		return 0
	}
	if p.data[z][y*w+x] == '#' {
		return 1
	}
	return 0
}

func (p *puz2) next() puz2 {
	if p.data[p.minZ()] == nil || p.data[p.maxZ()] == nil {
		panic("here")
	}
	ret := puz2{
		data: [401][]byte{},
		dnz:  p.dnz - 1,
		dxz:  p.dxz + 1,
	}
	for z := ret.minZ(); z <= ret.maxZ(); z++ {
		layer := make([]byte, sz)
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				if x == mid && y == mid {
					layer[y*w+x] = '?'
					continue
				}
				v := byte('.')
				sum := p.north(x, y, z) + p.south(x, y, z) + p.east(x, y, z) + p.west(x, y, z)
				switch p.val(x, y, z) {
				case 0:
					if sum == 1 || sum == 2 {
						v = '#'
					}
				case 1:
					if sum == 1 {
						v = '#'
					}
				}
				layer[y*w+x] = v
			}
		}
		ret.data[z] = layer
	}
	return ret
}

func (p *puz2) String() string {
	return string(p.data[200])
}

func (p *puz2) Score() int {
	sum := 0
	for _, layer := range p.data {
		for _, v := range layer {
			if v == '#' {
				sum++
			}
		}
	}
	return sum
}

func (p *puz2) north(x int, y int, z int) int {
	if y == 0 {
		// space north of center
		return p.val(mid, mid-1, z-1)
	}
	if x == mid && y-1 == mid {
		// return the bottom row
		return p.row(h-1, z+1)
	}
	return p.val(x, y-1, z)
}

func (p *puz2) east(x int, y int, z int) int {
	if x+1 == w {
		// space east of center
		return p.val(mid+1, mid, z-1)
	}
	if x+1 == mid && y == mid {
		// return the left col
		return p.col(0, z+1)
	}
	return p.val(x+1, y, z)
}

func (p *puz2) south(x int, y int, z int) int {
	if y+1 == h {
		// space south of center
		return p.val(mid, mid+1, z-1)
	}
	if x == mid && y+1 == mid {
		// return the top row
		return p.row(0, z+1)
	}
	return p.val(x, y+1, z)
}

func (p *puz2) west(x int, y int, z int) int {
	if x == 0 {
		// space west of center
		return p.val(mid-1, mid, z-1)
	}
	if x-1 == mid && y == mid {
		// return the right col
		return p.col(w-1, z+1)
	}
	return p.val(x-1, y, z)
}

func (p *puz2) row(y int, z int) int {
	sum := 0
	for x := 0; x < w; x++ {
		sum += p.val(x, y, z)
	}
	return sum
}

func (p *puz2) col(x int, z int) int {
	sum := 0
	for y := 0; y < w; y++ {
		sum += p.val(x, y, z)
	}
	return sum
}
