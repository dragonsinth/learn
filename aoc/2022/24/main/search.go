package main

import (
	"fmt"
	"github.com/dragonsinth/learn/aoc/sliceheap"
	"sync/atomic"
	"time"
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
	stats  [4]int
	pStats atomic.Pointer[[]int]
)

func track(stat ...int) {
	for i := 0; i < 4; i++ {
		stats[i] = max(stats[i], stat[i])
	}
	cp := append([]int{}, stats[:]...)
	pStats.Store(&cp)
}

func printStats() {
	for {
		if p := pStats.Load(); p != nil {
			st := *p
			fmt.Println(st[0], st[1], st[2], st[3])
		}
		time.Sleep(time.Second)
	}
}
