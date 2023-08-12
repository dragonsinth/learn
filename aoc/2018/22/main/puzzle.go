package main

import "github.com/dragonsinth/learn/aoc/grid"

type wat int

const (
	ROCKY = wat(iota)
	WATER
	NARRO
)

const watCount = 3

func (w wat) risk() int {
	return int(w)
}

func (w wat) char() byte {
	return ".=|"[w]
}

type puz struct {
	depth int
	w, h  int
	src   loc
	dst   loc

	erosion [][]int
	cost    [gearCount][][]int
}

func (p *puz) Render() [][]byte {
	data := grid.Alloc2d(p.dst.x+1, p.dst.y+1, byte(0))
	for y := 0; y <= p.dst.y; y++ {
		for x := 0; x <= p.dst.x; x++ {
			data[y][x] = p.display(x, y)
		}
	}
	return data
}

func (p *puz) Risk() (ret int) {
	for y := 0; y <= p.dst.y; y++ {
		for x := 0; x <= p.dst.x; x++ {
			e := p.Erosion(x, y)
			ret += p.typ(e).risk()
		}
	}
	return
}

func (p *puz) GeologicIndex(x, y int) int {
	if x == 0 && y == 0 {
		return 0
	}
	if x == p.dst.x && y == p.dst.y {
		return 0
	}
	if y == 0 {
		return x * 16807
	}
	if x == 0 {
		return y * 48271
	}

	return p.Erosion(x-1, y) * p.Erosion(x, y-1)
}

func (p *puz) Erosion(x, y int) int {
	if e := p.erosion[y][x]; e != 0 {
		return e
	}

	e := (p.GeologicIndex(x, y) + p.depth) % 20183
	p.erosion[y][x] = e
	return e
}

func (p *puz) display(x, y int) byte {
	e := p.Erosion(x, y)
	if x == p.src.x && y == p.src.y {
		return 'M'
	}
	if x == p.dst.x && y == p.dst.y {
		return 'T'
	}
	return p.typ(e).char()
}

func (p *puz) typ(erosion int) wat {
	return wat(erosion % watCount)
}
