package main

import (
	"fmt"
	"github.com/dragonsinth/learn/aoc/termbox"
	"os"
	"strconv"
	"strings"
)

const sample = `
5,4
4,2
4,5
3,0
2,1
6,3
2,4
1,5
0,6
3,3
2,6
5,1
1,2
5,5
2,5
6,5
1,4
0,4
6,4
1,1
6,1
1,0
0,5
1,6
2,0
`

func main() {
	part1(sample, 7, 12, true)
	part2(sample, 7, 12, true)
}

func part1(input string, sz int, turns int, debug bool) {
	p := puz{
		data: map[pos]bool{},
		sz:   sz,
	}
	pts := parse(input, sz)
	for _, pt := range pts[:turns] {
		p.data[pt] = true
	}

	if debug {
		termbox.RenderPlain(p.Render(), os.Stdout)
	}

	dist := p.FloodFill()
	fmt.Println(dist[pos{sz - 1, sz - 1}])
	if debug {
		termbox.RenderPlain(p.RenderWithDist(dist), os.Stdout)
	}
}

func part2(input string, sz int, start int, debug bool) {
	pts := parse(input, sz)
	lo := start
	hi := len(pts)
	for lo < hi {
		if lo == hi-1 {
			// found the boundary
			fmt.Println(pts[hi])
			return
		}
		mid := (lo + hi) / 2
		ok := canComplete(pts[:mid+1], sz)
		if debug {
			fmt.Println(lo, mid, hi, ok)
		}
		if ok {
			lo = mid
		} else {
			hi = mid
		}
	}
	panic("here")
}

func canComplete(pts []pos, sz int) bool {
	p := puz{
		data: make(map[pos]bool, len(pts)<<1),
		sz:   sz,
	}
	for _, pt := range pts {
		p.data[pt] = true
	}
	dist := p.FloodFill()
	return dist[pos{sz - 1, sz - 1}] != 0
}

func parse(input string, sz int) []pos {
	var ret []pos
	for _, line := range strings.Split(input, "\n") {
		if line == "" {
			continue
		}
		line = strings.TrimSpace(line)
		parts := strings.Split(line, ",")
		if len(parts) != 2 {
			panic(parts)
		}

		x, y := mustInt(parts[0]), mustInt(parts[1])
		if x < 0 || x >= sz || y < 0 || y >= sz {
			panic(line)
		}
		ret = append(ret, pos{x, y})
	}
	return ret
}

type puz struct {
	data map[pos]bool
	sz   int
}

func (p *puz) Render() [][]byte {
	var buf [][]byte
	for y := 0; y < p.sz; y++ {
		var line []byte
		for x := 0; x < p.sz; x++ {
			pt := pos{x, y}
			if p.data[pt] {
				line = append(line, '#')
			} else {
				line = append(line, '.')
			}
		}
		buf = append(buf, line)
	}
	return buf
}

func (p *puz) RenderWithDist(dist map[pos]int) [][]byte {
	var buf [][]byte
	for y := 0; y < p.sz; y++ {
		var line []byte
		for x := 0; x < p.sz; x++ {
			pt := pos{x, y}
			if d, ok := dist[pt]; ok {
				line = append(line, '0'+byte(d%10))
			} else if p.data[pt] {
				line = append(line, '#')
			} else {
				line = append(line, '.')
			}
		}
		buf = append(buf, line)
	}
	return buf
}

func (p *puz) FloodFill() map[pos]int {
	dist := make(map[pos]int, p.sz*p.sz)
	p.floodFill(dist, pos{0, 0}, 0)
	return dist
}

func (p *puz) floodFill(dist map[pos]int, pt pos, cost int) {
	if pt.x < 0 || pt.x >= p.sz || pt.y < 0 || pt.y >= p.sz {
		return // edge
	}
	if p.data[pt] {
		return // blocked
	}
	if prevDist, ok := dist[pt]; ok && prevDist <= cost {
		return // already found a cheaper route to this spot
	}
	dist[pt] = cost

	for _, n := range pt.adjacent() {
		p.floodFill(dist, n, cost+1)
	}
}

type pos struct {
	x, y int
}

func (p pos) adjacent() [4]pos {
	return [4]pos{
		{p.x, p.y - 1},
		{p.x + 1, p.y},
		{p.x, p.y + 1},
		{p.x - 1, p.y},
	}
}

func mustInt(s string) int {
	if v, err := strconv.Atoi(s); err != nil {
		panic(fmt.Sprint(s, err))
	} else {
		return v
	}
}
