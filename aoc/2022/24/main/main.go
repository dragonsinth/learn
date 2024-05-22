package main

import (
	"github.com/dragonsinth/learn/aoc/termbox"
	"os"
	"strings"
	"time"
)

const sample = `
#.######
#>>.<^<#
#.<..<<#
#>v.><>#
#<^v^^>#
######.#
`

func main() {
	run(sample, 0)
}

func run(input string, frameDelay time.Duration) {
	max := pos{-1, -1}
	var bliz []blizzard
	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if max.x < 0 {
			max.x = len(line) - 1
		} else if max.x != len(line)-1 {
			panic(line)
		}
		max.y++
		y := max.y
		for x, v := range line {
			pt := pos{x, y}
			switch v {
			case '#', '.':
			case '^':
				bliz = append(bliz, blizzard{pt: pt, d: NORTH})
			case 'v':
				bliz = append(bliz, blizzard{pt: pt, d: SOUTH})
			case '<':
				bliz = append(bliz, blizzard{pt: pt, d: WEST})
			case '>':
				bliz = append(bliz, blizzard{pt: pt, d: EAST})
			default:
				panic(v)
			}
		}
	}

	as := &allStates{
		max:    max,
		states: []state{{blizzards: bliz, blizzardPos: blizzardMap(bliz)}},
	}

	end := pos{max.x - 1, max.y}
	start := pos{1, 0}
	g := &game{
		prev:   nil,
		start:  start,
		end:    end,
		turn:   0,
		state:  as.get(0),
		states: as,
	}

	termbox.RenderPlain(g.render(), os.Stdout)

	g.player = start
	g.dest = end
	g = findRoute(g)
	termbox.RenderPlain(g.render(), os.Stdout)

	g.dest = start
	g = findRoute(g)
	termbox.RenderPlain(g.render(), os.Stdout)

	g.dest = end
	g = findRoute(g)
	termbox.RenderPlain(g.render(), os.Stdout)

	printStats()

	if frameDelay > 0 {
		tb := termbox.New(true)
		defer tb.Stop()
		animate(g, tb, frameDelay)
	}
}

func animate(g *game, tb termbox.Terminal, frameDelay time.Duration) {
	if g.prev != nil {
		animate(g.prev, tb, frameDelay)
	}
	tb.Render(g.render(), os.Stdout)
	time.Sleep(frameDelay)
}

func manhattan(a pos, b pos) int {
	return abs(a.x-b.x) + abs(a.y-b.y)
}

func abs(i int) int {
	if i < 0 {
		return -i
	} else {
		return i
	}
}
