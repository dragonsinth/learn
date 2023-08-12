package main

import (
	"fmt"
	"github.com/dragonsinth/learn/aoc/sliceheap"
)

const (
	swapCost = 7
)

type loc struct {
	x, y int
	g    gear
}

func (l loc) adjacent() [4]loc {
	return [4]loc{
		{l.x, l.y - 1, l.g},
		{l.x + 1, l.y, l.g},
		{l.x, l.y + 1, l.g},
		{l.x - 1, l.y, l.g},
	}
}

func (l loc) swapGear() [2]loc {
	return [2]loc{
		{l.x, l.y, (l.g + 1) % gearCount},
		{l.x, l.y, (l.g + 2) % gearCount},
	}
}

func (l loc) String() string {
	return fmt.Sprintf("(%d,%d,%s)", l.x, l.y, l.g)
}

type gear int

const gearCount = 3

const (
	TORCH = gear(iota)
	CLIMB
	NOTHN
)

var gearNames = []string{
	"TORCH",
	"CLIMB",
	"NOTHN",
}

func (g gear) String() string {
	return gearNames[g]
}

type path struct {
	prev *path
	l    loc
	cost int
}

func (p path) Walk() []path {
	if p.prev == nil {
		return []path{p}
	}
	return append(p.prev.Walk(), p)
}

func (p *puz) Astar() path {
	heuristic := func(l loc) int {
		manh := abs(p.dst.x-l.x) + abs(p.dst.y-l.y)
		if l.g != p.dst.g {
			manh += swapCost
		}
		return manh
	}

	cost := func(pt path) int {
		return pt.cost + heuristic(pt.l)
	}

	work := sliceheap.New[path](func(a, b path) bool {
		return cost(a) < cost(b)
	})

	work.Push(path{
		prev: nil,
		l:    p.src,
		cost: 0,
	})

	for {
		w := work.Pop()
		if p.Cost(w.l) <= w.cost {
			continue // already found a cheaper route to this spot
		}
		p.SetCost(w.l, w.cost)

		if w.l == p.dst {
			return w
		}

		// try moving
		for _, l := range w.l.adjacent() {
			if p.passable(l) {
				work.Push(path{
					prev: &w,
					l:    l,
					cost: w.cost + 1,
				})
			}
		}

		// try switching
		for _, l := range w.l.swapGear() {
			if p.passable(l) {
				work.Push(path{
					prev: &w,
					l:    l,
					cost: w.cost + swapCost,
				})
			}
		}
	}
}

func (p *puz) Cost(l loc) int {
	return p.cost[l.g][l.y][l.x]
}

func (p *puz) SetCost(l loc, cost int) {
	p.cost[l.g][l.y][l.x] = cost
}

func (p *puz) passable(l loc) bool {
	if l.x < 0 || l.y < 0 {
		return false // out of bounds
	}
	if l.x >= p.w || l.y >= p.h {
		panic(l.String())
	}
	typ := p.typ(p.Erosion(l.x, l.y))
	switch typ {
	case ROCKY:
		return l.g != NOTHN
	case WATER:
		return l.g != TORCH
	case NARRO:
		return l.g != CLIMB
	default:
		panic(typ)
	}
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
