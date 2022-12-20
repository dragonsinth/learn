package main

import (
	"fmt"
)

var (
	turns int
	done  int
)

type res int

const (
	ORE = res(iota)
	CLAY
	OBSIDIAN
	GEODE
	MAX
)

func (r res) String() string {
	switch r {
	case ORE:
		return "ore"
	case CLAY:
		return "clay"
	case OBSIDIAN:
		return "obsidian"
	case GEODE:
		return "geode"
	default:
		panic(r)
	}
}

// [outBot][inRes]
type costsheet [4][4]int

type game struct {
	prev    *game
	cs      *costsheet
	turn    int
	maxTurn int

	resources [4]int
	bots      [4]int
	newBots   [4]int
}

func (g *game) run() *game {
	var best *game
	for nextBot := ORE; nextBot < MAX; nextBot++ {
		if !g.canEventuallyBuild(nextBot) {
			continue
		}

		if nextBot != GEODE {
			bots := g.bots[nextBot]
			stock := g.resources[nextBot]
			left := g.maxTurn - g.turn
			maxStock := bots*left + stock

			maxCost := 0
			for anyBot := ORE; anyBot < MAX; anyBot++ {
				maxCost = max(maxCost, g.cs[anyBot][nextBot])
			}

			if maxStock >= maxCost*left {
				continue // never need more
			}
		}

		g2 := *g
		g2.prev = g
		result := g2.runUntil(nextBot)
		if result.betterThan(best) {
			best = result
		}
	}
	return best
}

func (g *game) runUntil(wantBot res) *game {
	for g.turn < g.maxTurn {
		g.turn++
		turns++

		// Spend
		canBuy := g.canBuy(wantBot)
		if canBuy {
			for i := ORE; i < MAX; i++ {
				g.resources[i] -= g.cs[wantBot][i]
			}
			g.newBots[wantBot]++
		}

		// Collect
		for i := ORE; i < MAX; i++ {
			if g.bots[i] > 0 {
				g.resources[i] += g.bots[i]
			}
		}

		// New bots
		for i := ORE; i < MAX; i++ {
			if g.newBots[i] > 0 {
				g.bots[i] += g.newBots[i]
				g.newBots[i] = 0
			}
		}

		if canBuy {
			return g.run()
		}
	}

	done++
	return g
}

func (g *game) canBuy(wantBot res) bool {
	for i := ORE; i < MAX; i++ {
		if g.resources[i] < g.cs[wantBot][i] {
			return false
		}
	}
	return true
}

func (g *game) canEventuallyBuild(wantBot res) bool {
	for i := ORE; i < MAX; i++ {
		if g.cs[wantBot][i] > 0 {
			if g.bots[i]+g.newBots[i] == 0 {
				return false // we can never build this input
			}
		}
	}
	return true
}

func (g *game) betterThan(o *game) bool {
	if o == nil {
		return true
	}
	for i := GEODE; i >= ORE; i-- {
		if g.resources[i] != o.resources[i] {
			return g.resources[i] > o.resources[i]
		}
	}
	return false
}

func (g *game) summary() {
	if g == nil {
		return
	}
	g.prev.summary()
	fmt.Printf("%2d res=%+v, bot=%+v\n", g.turn, g.resources, g.bots)
}
