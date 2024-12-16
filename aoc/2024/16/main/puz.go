package main

import (
	"fmt"
	"github.com/dragonsinth/learn/aoc/sliceheap"
	"strings"
)

func run(input string) {
	p := parse(input)
	routes := p.Best()

	seats := map[pos]bool{}
	for _, r := range routes {
		for c := r; c != nil; c = c.prev {
			seats[c.loc] = true
		}
	}
	fmt.Println(routes[0].cost, len(seats))
}

func parse(input string) *puz {
	p := &puz{
		data: map[pos]byte{},
	}

	lines := strings.Split(input, "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		line = strings.TrimSpace(line)
		if p.w <= 0 {
			p.w = len(line)
		} else if p.w != len(line) {
			panic(line)
		}
		y := p.h
		for x, wat := range line {
			pt := pos{x: x, y: y}
			switch wat {
			case '.', '#':
				p.data[pt] = byte(wat)
			case 'S':
				p.start = pt
				p.data[pt] = '.'
			case 'E':
				p.end = pt
				p.data[pt] = '.'
			default:
				panic(wat)
			}
		}
		p.h++

	}
	return p
}

type puz struct {
	data  map[pos]byte
	w, h  int
	start pos
	end   pos
}

func (p *puz) Best() []*Route {
	var ret []*Route
	bestCost := map[State]int{}

	heuristic := func(l pos) int {
		return abs(p.end.x-l.x) + abs(p.end.y-l.y)
	}

	cost := func(r *Route) int {
		return r.cost + heuristic(r.loc)
	}

	work := sliceheap.New[*Route](func(a, b *Route) bool {
		return cost(a) < cost(b)
	})

	startState := State{
		loc:  p.start,
		face: E,
	}
	work.Push(&Route{
		prev:  nil,
		cost:  0,
		State: startState,
	})
	bestCost[startState] = 0

	for {
		w := work.Pop()

		if len(ret) > 0 && ret[0].cost < w.cost {
			return ret // all paths found
		}

		if w.loc == p.end {
			ret = append(ret, w)
		}

		maybeAdd := func(st State, cost int) {
			if bc, ok := bestCost[st]; ok && bc < cost {
				return // already found a cheaper route to this spot
			}
			bestCost[st] = cost
			work.Push(&Route{
				prev:  w,
				cost:  cost,
				State: st,
			})
		}

		// try moving
		if next := w.loc.next(w.face); p.data[next] == '.' {
			maybeAdd(State{loc: next, face: w.face}, w.cost+1)
		}

		// try turning
		maybeAdd(State{loc: w.loc, face: w.face.left()}, w.cost+1000)
		maybeAdd(State{loc: w.loc, face: w.face.right()}, w.cost+1000)
	}
}

type Route struct {
	prev *Route
	cost int
	State
}

type State struct {
	loc  pos
	face dir
}
