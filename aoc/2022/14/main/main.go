package main

import (
	"fmt"
	"github.com/dragonsinth/learn/aoc/termbox"
	"os"
	"strconv"
	"strings"
)

var sample = `
498,4 -> 498,6 -> 496,6
503,4 -> 502,4 -> 502,9 -> 494,9
`

func main() {
	p := puzzle{
		min:  pos{x: 499, y: 0},
		max:  pos{x: 501, y: 1},
		drop: pos{500, 0},
	}

	for _, line := range strings.Split(sample, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		var lastPt pos
		for i, pair := range strings.Split(line, " -> ") {
			parts := strings.Split(pair, ",")
			if len(parts) != 2 {
				panic(pair)
			}
			pt := pos{
				x: mustInt(parts[0]),
				y: mustInt(parts[1]),
			}
			p.min.x = min(p.min.x, pt.x-1)
			p.max.x = max(p.max.x, pt.x+2)
			p.min.y = min(p.min.y, pt.y-1)
			p.max.y = max(p.max.y, pt.y+2)

			if i > 0 {
				switch {
				case lastPt.x == pt.x:
					p.vlines = append(p.vlines, vline{
						x:  pt.x,
						y1: min(lastPt.y, pt.y),
						y2: max(lastPt.y, pt.y),
					})
				case lastPt.y == pt.y:
					p.hlines = append(p.hlines, hline{
						x1: min(lastPt.x, pt.x),
						x2: max(lastPt.x, pt.x),
						y:  pt.y,
					})
				default:
					panic(pt)
				}
			}
			lastPt = pt
		}
	}

	fmt.Println(p.min, p.max)
	p.drawLines()
	termbox.RenderPlain(p.renderView(), os.Stdout)
	fmt.Println()

	sum := 0
	for p.addSandToFloor() {
		sum++
	}
	termbox.RenderPlain(p.renderView(), os.Stdout)
	fmt.Println(sum)

	for p.addSandToFill() {
		sum++
	}
	termbox.RenderPlain(p.renderView(), os.Stdout)
	fmt.Println(sum)
}

type pos struct {
	x, y int
}

func (p pos) String() string {
	return fmt.Sprintf("(%d,%d)", p.x, p.y)
}

type hline struct {
	x1, x2, y int
}

type vline struct {
	x, y1, y2 int
}

type puzzle struct {
	hlines []hline
	vlines []vline

	min, max pos
	drop     pos

	data [][]byte
}

func (p *puzzle) addPoint(pt pos) {
	p.min.x = min(p.min.x, pt.x-1)
	p.max.x = max(p.max.x, pt.x+2)
	p.min.y = min(p.min.y, pt.y)
	p.max.y = max(p.max.y, pt.y+2)
}

func (p *puzzle) drawLines() {
	// widen the simulation to allow a complete fill
	p.min.x = min(p.min.x, p.drop.x-p.max.y)
	p.max.x = max(p.min.x, p.drop.x+p.max.y+1)

	for y := 0; y <= p.max.y; y++ {
		buf := make([]byte, p.max.x+1)
		for i := range buf {
			buf[i] = '.'
		}
		p.data = append(p.data, buf)
	}
	p.data[p.drop.y][p.drop.x] = '+'

	for _, hl := range p.hlines {
		for x := hl.x1; x <= hl.x2; x++ {
			p.data[hl.y][x] = '#'
		}
	}
	for _, vl := range p.vlines {
		for y := vl.y1; y <= vl.y2; y++ {
			p.data[y][vl.x] = '#'
		}
	}

	// floor
	for x := p.min.x; x <= p.max.x; x++ {
		p.data[p.max.y][x] = '#'
	}
}

func (p *puzzle) renderView() [][]byte {
	var ret [][]byte
	for y := p.min.y; y <= p.max.y; y++ {
		buf := make([]byte, p.max.x-p.min.x)
		for i := range buf {
			buf[i] = p.data[y][i+p.min.x]
		}
		ret = append(ret, buf)
	}
	return ret
}

func (p *puzzle) addSandToFloor() bool {
	pt := p.drop
	for pt.y+1 < p.max.y {
		var moved bool
		for _, dst := range []pos{
			{pt.x, pt.y + 1},
			{pt.x - 1, pt.y + 1},
			{pt.x + 1, pt.y + 1},
		} {
			if p.data[dst.y][dst.x] == '.' {
				pt = dst
				moved = true
				break
			}
		}
		if !moved {
			p.data[pt.y][pt.x] = 'o'
			return true
		}
	}
	return false
}

func (p *puzzle) addSandToFill() bool {
	pt := p.drop
	if p.data[pt.y][pt.x] != '+' {
		return false
	}
	for {
		var moved bool
		for _, dst := range []pos{
			{pt.x, pt.y + 1},
			{pt.x - 1, pt.y + 1},
			{pt.x + 1, pt.y + 1},
		} {
			if p.data[dst.y][dst.x] == '.' {
				pt = dst
				moved = true
				break
			}
		}
		if !moved {
			p.data[pt.y][pt.x] = 'o'
			return true
		}
	}
}

func mustInt(s string) int {
	if v, err := strconv.Atoi(s); err != nil {
		panic(fmt.Sprint(s, err))
	} else {
		return v
	}
}

func min(a int, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

func max(a int, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}
