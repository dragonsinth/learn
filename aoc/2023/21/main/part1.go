package main

import (
	"fmt"
	"github.com/dragonsinth/learn/aoc/termbox"
	"os"
)

func run(input string, maxDist int, debug bool) {
	p := parse(input)
	if debug {
		termbox.RenderPlain(p.Render(), os.Stdout)
		fmt.Println()
	}

	seen := p.walk(maxDist)
	if debug {
		termbox.RenderPlain(p.RenderWith(seen, maxDist), os.Stdout)
		fmt.Println()
	}
	fmt.Println(countWalkable(seen, maxDist))
}

func countWalkable(seen map[pos]int, maxDist int) int {
	sum := 0
	for _, dist := range seen {
		if dist > maxDist {
			continue
		}
		if dist%2 == maxDist%2 {
			sum++
		}
	}
	return sum
}

func (p puz) walk(maxDist int) map[pos]int {
	type task struct {
		pos
		dist int
	}
	work := []task{{p.start, 0}}
	seen := map[pos]int{p.start: 0}

	maybeWalk := func(x, y, dist int) {
		pt := pos{x, y}
		if v, ok := seen[pt]; ok && v <= dist {
			return
		}
		if p.at(x, y) == '.' {
			seen[pt] = dist
			work = append(work, task{pt, dist})
		}
	}

	for len(work) > 0 {
		t := work[0]
		work = work[1:]
		x, y, dist := t.x, t.y, t.dist+1
		if dist > maxDist {
			continue
		}
		maybeWalk(x, y-1, dist)
		maybeWalk(x+1, y, dist)
		maybeWalk(x, y+1, dist)
		maybeWalk(x-1, y, dist)
	}

	return seen
}

type pos struct {
	x, y int
}
