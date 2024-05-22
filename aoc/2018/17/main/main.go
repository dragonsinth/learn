package main

import (
	"fmt"
	"github.com/dragonsinth/learn/aoc/termbox"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const sample = `
x=495, y=2..7
y=7, x=495..501
x=501, y=3..7
x=498, y=2..4
x=506, y=1..2
x=498, y=10..13
x=504, y=10..13
y=13, x=498..504
`

var (
	re   = regexp.MustCompile(`^([xy])=(\d+), ([xy])=(\d+)\.\.(\d+)`)
	drop = pos{500, 0}
)

func main() {
	run(sample, true)
	//run(data, false)
}

func run(in string, debug bool) {
	p := puzzle{
		wat: map[pos]byte{},
		min: pos{x: math.MaxInt, y: math.MaxInt},
		max: pos{x: math.MinInt, y: math.MinInt},
	}

	for _, line := range strings.Split(in, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if !re.MatchString(line) {
			panic(line)
		}

		m := re.FindStringSubmatch(line)
		switch m[1] {
		case "x":
			x, ymin, ymax := mustInt(m[2]), mustInt(m[4]), mustInt(m[5])
			for y := ymin; y <= ymax; y++ {
				p.wat[pos{x, y}] = '#'
			}
			p.min.x = min(p.min.x, x)
			p.max.x = max(p.max.x, x)
			p.min.y = min(p.min.y, ymin)
			p.max.y = max(p.max.y, ymax)
		case "y":
			y, xmin, xmax := mustInt(m[2]), mustInt(m[4]), mustInt(m[5])
			for x := xmin; x <= xmax; x++ {
				p.wat[pos{x, y}] = '#'
			}
			p.min.x = min(p.min.x, xmin)
			p.max.x = max(p.max.x, xmax)
			p.min.y = min(p.min.y, y)
			p.max.y = max(p.max.y, y)
		}
	}

	p.min.x--
	p.max.x++
	termbox.RenderPlain(p.Render(p.min.up(), p.max), os.Stdout)
	fmt.Println(p.min, p.max)

	term := termbox.New(debug)
	p.wat[drop] = '+'
	p.Run(term, drop)
	term.Stop()
	termbox.RenderPlain(p.Render(p.min.up(), p.max), os.Stdout)
	fmt.Println(p.Wet(), p.Soaked())
}

type pos struct {
	x, y int
}

func (p pos) left() pos {
	return pos{p.x - 1, p.y}
}

func (p pos) right() pos {
	return pos{p.x + 1, p.y}
}

func (p pos) up() pos {
	return pos{p.x, p.y - 1}
}

func (p pos) down() pos {
	return pos{p.x, p.y + 1}
}

func (p pos) String() string {
	return fmt.Sprintf("(%d,%d)", p.x, p.y)
}

type puzzle struct {
	wat      map[pos]byte
	min, max pos
	inv      int
}

func (p *puzzle) Run(term termbox.Terminal, pt pos) {
	p.inv++
	if p.inv > 10000 {
		term.Stop()
		termbox.RenderPlain(p.Render(p.min, p.max), os.Stdout)
		panic(pt)
	}

	// fall until we hit a surface
	for {
		if pt.y == p.max.y {
			return // fell off map
		}
		next := pos{pt.x, pt.y + 1}
		if !p.empty(next) {
			break
		}
		pt = next
		p.wat[pt] = '|'
	}

	for {
		if term.Enabled() {
			vmin, vmax := p.min, p.max
			vmin.y = 0
			if vmax.y-vmin.y > 80 {
				vmin.y = pt.y - 40
				vmax.y = pt.y + 40
			}
			term.Render(p.Render(vmin, vmax), os.Stdout)
			time.Sleep(100 * time.Millisecond)
		}

		// scan left and right to see what happens
		left, leftWall := p.findWall(pt, pos.left)
		right, rightWall := p.findWall(pt, pos.right)

		if leftWall && rightWall {
			// fill the layer with still water
			for x := left.x; x <= right.x; x++ {
				p.wat[pos{x, pt.y}] = '~'
			}
			pt.y-- // back up and do it again
			continue
		}

		// fill the layer with flowing water
		for x := left.x + 1; x <= right.x-1; x++ {
			p.wat[pos{x, pt.y}] = '|'
		}

		// Now recursively drop water from either open side.
		// The order of operations is very important here to prevent crazy recursion.
		if leftWall {
			p.wat[left] = '|'
		}
		if rightWall {
			p.wat[right] = '|'
		}

		doLeft := p.wat[left] == 0
		if doLeft {
			p.wat[left] = '+'
		}
		doRight := p.wat[right] == 0
		if doRight {
			p.wat[right] = '+'
		}
		if doLeft {
			p.Run(term, left)
		}
		if doRight {
			p.Run(term, right)
		}
		return
	}
}

func (p *puzzle) findWall(pt pos, dir func(pos) pos) (pos, bool) {
	for {
		// Does the floor drop?
		if p.empty(pos{pt.x, pt.y + 1}) {
			return pt, false
		}
		// Is the next tile a wall?
		next := dir(pt)
		if !p.empty(next) {
			return pt, true
		}
		pt = next
	}
}

func (p *puzzle) Render(vmin, vmax pos) [][]byte {
	var ret [][]byte
	for y := vmin.y; y <= vmax.y; y++ {
		buf := make([]byte, 0, vmax.x-vmin.x+1)
		for x := vmin.x; x <= vmax.x; x++ {
			if v := p.wat[pos{x, y}]; v != 0 {
				buf = append(buf, v)
			} else {
				buf = append(buf, '.')
			}
		}
		ret = append(ret, buf)
	}
	return ret
}

func (p *puzzle) Wet() int {
	sum := 0
	for y := p.min.y; y <= p.max.y; y++ {
		for x := p.min.x; x <= p.max.x; x++ {
			wat := p.wat[pos{x, y}]
			if wat == '~' || wat == '|' || wat == '+' {
				sum++
			}
		}
	}
	return sum
}

func (p *puzzle) Soaked() int {
	sum := 0
	for y := p.min.y; y <= p.max.y; y++ {
		for x := p.min.x; x <= p.max.x; x++ {
			wat := p.wat[pos{x, y}]
			if wat == '~' {
				sum++
			}
		}
	}
	return sum
}

func (p *puzzle) empty(pt pos) bool {
	wat := p.wat[pt]
	return wat == 0 || wat == '|' || wat == '+'
}

func mustInt(s string) int {
	if v, err := strconv.Atoi(s); err != nil {
		panic(fmt.Sprint(s, err))
	} else {
		return v
	}
}
