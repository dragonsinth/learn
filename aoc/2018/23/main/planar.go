package main

import (
	"bytes"
	"fmt"
	"github.com/dragonsinth/learn/aoc/grid"
	"github.com/dragonsinth/learn/aoc/termbox"
	"os"
)

func collapseRegionsPlanar(dx, dy int, bots []bot) region {
	iter := 1
	regions := toRegions(bots)

	for i := range regions {
		r := &regions[i]
		for d := range (pos{}) {
			if d != dx && d != dy {
				r.min[d] = 0
				r.max[d] = 0
			}
		}
	}
	orig := regions

	fmt.Printf("iter=%d, len=%d, vol=%d\n", iter, len(regions), volume(regions))
	printPlanar(dx, dy, regions)
	for len(regions) > 1 {
		regions = intersectAll(regions)
		regions = dedup(regions)
		iter++
		fmt.Printf("iter=%d, len=%d, vol=%d\n", iter, len(regions), volume(regions))
		printPlanar(dx, dy, regions)
	}

	// validate!
	final := regions[0]
	expectInRange := map[int]bool{}
	for _, id := range final.ids {
		expectInRange[id] = true
	}

	for y := final.min[dy]; y <= final.max[dy]; y++ {
		for x := final.min[dx]; x <= final.max[dx]; x++ {
			for id, r := range orig {
				p := r.min
				p[dx] = x
				p[dy] = y
				if expectInRange[id] != r.inRange(p) {
					panic(fmt.Sprintf("expectInRange(%v) != r.inRange(%v), id=%d", expectInRange[id], r.inRange(p), id))
				}
			}
		}
	}

	return final
}

func printPlanar(dx, dy int, regions []region) {
	draw := grid.Alloc2d(30, 30, byte('.'))
	for _, r := range regions {
		if r.volume() > 100 {
			continue
		}

		for y := max(0, r.min[dy]); y <= min(29, r.max[dy]); y++ {
			for x := max(0, r.min[dx]); x <= min(29, r.max[dx]); x++ {
				draw[y][x] = '*'
			}
		}
	}
	termbox.RenderPlain(draw, os.Stdout)

	for _, r := range regions {
		var buf bytes.Buffer
		buf.WriteString(fmt.Sprintf(" [%d,%d->%d,%d]", r.min[dx], r.min[dy], r.max[dx], r.max[dx]))
		buf.WriteString(fmt.Sprintf(" %+v", r.ids))
		fmt.Println(buf.String())
	}
}
