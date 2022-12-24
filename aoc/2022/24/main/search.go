package main

import (
	"fmt"
	"github.com/dragonsinth/learn/aoc/sliceheap"
)

func findRoute(g *game) *game {
	seen := map[key]struct{}{}
	h := sliceheap.New(func(a *game, b *game) bool {
		return a.cost() < b.cost()
	})

	add := func(g *game) {
		if _, ok := seen[g.key()]; ok {
			return
		}
		h.Push(g)
		seen[g.key()] = struct{}{}
	}
	add(g)

	for {
		g := h.Pop()
		if g.player == g.dest {
			return g
		}
		track(g.turn, g.cost(), len(seen), h.Len())

		// Try each direction + standing still.
		g = g.next()
		for _, pt := range g.player.reachable() {
			wat := g.get(pt)
			if wat == '.' {
				cp := *g
				cp.player = pt
				add(&cp)
			}
		}
	}
}

var (
	stats [4]int
)

func track(stat ...int) {
	for i := 0; i < 4; i++ {
		stats[i] = max(stats[i], stat[i])
	}
}

func printStats() {
	fmt.Println(stats[0], stats[1], stats[2], stats[3])
}
