package main

import (
	"fmt"
	"github.com/dragonsinth/learn/aoc/grid"
	"github.com/dragonsinth/learn/aoc/termbox"
	"math"
	"os"
)

var (
	sample = puz{
		depth: 510,
		src:   loc{0, 0, TORCH},
		dst:   loc{10, 10, TORCH},
		w:     50,
		h:     50,
	}
)

func main() {
	run(&sample, true)
}

func run(p *puz, debug bool) {
	p.erosion = grid.Alloc2d(p.w, p.h, 0)
	fmt.Println(p.Risk())
	if debug {
		termbox.RenderPlain(p.Render(), os.Stdout)
	}

	for i := 0; i < 3; i++ {
		p.cost[i] = grid.Alloc2d(p.w, p.h, math.MaxInt)
	}

	pth := p.Astar()
	if debug {
		for i, pt := range pth.Walk() {
			if debug {
				fmt.Println(i, pt.cost, pt.l)
			}
		}
	} else {
		fmt.Println(pth.cost, pth.l)
	}
}
