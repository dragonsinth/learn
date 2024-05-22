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
position=< 9,  1> velocity=< 0,  2>
position=< 7,  0> velocity=<-1,  0>
position=< 3, -2> velocity=<-1,  1>
position=< 6, 10> velocity=<-2, -1>
position=< 2, -4> velocity=< 2,  2>
position=<-6, 10> velocity=< 2, -2>
position=< 1,  8> velocity=< 1, -1>
position=< 1,  7> velocity=< 1,  0>
position=<-3, 11> velocity=< 1, -2>
position=< 7,  6> velocity=<-1, -1>
position=<-2,  3> velocity=< 1,  0>
position=<-4,  3> velocity=< 2,  0>
position=<10, -3> velocity=<-1,  1>
position=< 5, 11> velocity=< 1, -2>
position=< 4,  7> velocity=< 0, -1>
position=< 8, -2> velocity=< 0,  1>
position=<15,  0> velocity=<-2,  0>
position=< 1,  6> velocity=< 1,  0>
position=< 8,  9> velocity=< 0, -1>
position=< 3,  3> velocity=<-1,  1>
position=< 0,  5> velocity=< 0, -1>
position=<-2,  2> velocity=< 2,  0>
position=< 5, -2> velocity=< 1,  2>
position=< 1,  4> velocity=< 2,  1>
position=<-2,  7> velocity=< 2, -2>
position=< 3,  6> velocity=<-1, -1>
position=< 5,  0> velocity=< 1,  0>
position=<-6,  0> velocity=< 2,  0>
position=< 5,  9> velocity=< 1, -2>
position=<14,  7> velocity=<-2,  0>
position=<-3,  6> velocity=< 2, -1>
`

var re = regexp.MustCompile(`^position=<\s*(-?\d+),\s*(-?\d+)>\s*velocity=<\s*(-?\d+),\s*(-?\d+)>$`)

type pt struct {
	x, y, dx, dy int
}

func (p *pt) Move() {
	p.x += p.dx
	p.y += p.dy
}

type loc struct {
	x, y int
}

func main() {
	var p puz

	lines := strings.Split(sample, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if !re.MatchString(line) {
			panic(line)
		}
		m := re.FindStringSubmatch(line)
		pt := pt{
			x:  mustInt(m[1]),
			y:  mustInt(m[2]),
			dx: mustInt(m[3]),
			dy: mustInt(m[4]),
		}
		p.pts = append(p.pts, pt)
	}
	p.Compute()

	iter, out := func(debug bool) (int, [][]byte) {
		var term termbox.Terminal
		if debug {
			term = termbox.New(true)
			term.Start()
			defer term.Stop()
		}

		minSize := p.Size()
		var lastRender [][]byte
		i := 0
		for {
			if minSize < 10000 {
				lastRender = p.Render()
				if debug {
					term.Render(lastRender, os.Stdout)
					time.Sleep(100 * time.Millisecond)
				}
			}
			p.Move()
			p.Compute()
			sz := p.Size()
			if sz > minSize {
				return i, lastRender
			}
			minSize = sz
			fmt.Println(i, minSize)
			i++
		}
	}(true)

	fmt.Println(iter)
	termbox.RenderPlain(out, os.Stdout)
}

type puz struct {
	pts                    []pt
	locs                   map[loc]bool
	xmin, xmax, ymin, ymax int
}

func (p *puz) Move() {
	for i := range p.pts {
		pt := &p.pts[i]
		pt.Move()
	}
}

func (p *puz) Compute() {
	p.locs = map[loc]bool{}
	p.xmin, p.xmax, p.ymin, p.ymax = math.MaxInt, math.MinInt, math.MaxInt, math.MinInt
	for _, pt := range p.pts {
		p.locs[loc{pt.x, pt.y}] = true
		p.xmin = min(p.xmin, pt.x)
		p.xmax = max(p.xmax, pt.x)
		p.ymin = min(p.ymin, pt.y)
		p.ymax = max(p.ymax, pt.y)
	}
}

func (p *puz) Render() [][]byte {
	var data [][]byte
	for y := p.ymin; y <= p.ymax; y++ {
		var line []byte
		for x := p.xmin; x <= p.xmax; x++ {
			if p.locs[loc{x, y}] {
				line = append(line, '#')
			} else {
				line = append(line, '.')
			}
		}
		data = append(data, line)
	}
	return data
}

func (p *puz) Size() int {
	return (p.xmax - p.xmin) * (p.ymax - p.ymin)
}

func mustInt(s string) int {
	if v, err := strconv.Atoi(s); err != nil {
		panic(fmt.Sprint(s, err))
	} else {
		return v
	}
}
