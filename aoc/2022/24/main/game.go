package main

import "fmt"

type game struct {
	prev *game

	start, end pos

	player, dest pos
	turn         int
	state        *state
	states       *allStates
}

func (g *game) render() [][]byte {
	max := g.states.max
	out := make([][]byte, max.y+2)
	out[0] = []byte(fmt.Sprintf("Turn %d:", g.turn))
	field := out[1:]
	for y := 0; y <= max.y; y++ {
		field[y] = make([]byte, max.x+1)
		for x := 0; x <= max.x; x++ {
			pt := pos{x, y}
			wat := g.get(pt)
			if wat == '.' && pt == g.player {
				wat = 'E'
			}
			field[y][x] = wat
		}
	}
	return out
}

func (g *game) get(pt pos) byte {
	max := g.states.max
	if pt.x < 0 || pt.x > max.x || pt.y < 0 || pt.y > max.y {
		return ' '
	}
	if pt.x == 0 || pt.x == max.x {
		return '#'
	}
	if pt.y == 0 || pt.y == max.y {
		if pt == g.start || pt == g.end {
			return '.'
		} else {
			return '#'
		}
	}

	bz := g.state.blizzardPos[pt]
	switch c := len(bz); c {
	case 0:
		return '.'
	case 1:
		return bz[0].d.byte()
	default:
		if c < 10 {
			return byte(c + '0')
		}
		return '*'
	}
}

func (g *game) next() *game {
	cp := *g
	cp.prev = g
	cp.turn++
	cp.state = cp.states.get(cp.turn)
	return &cp
}

func (g *game) cost() int {
	return g.turn + manhattan(g.player, g.dest)
}

func (g *game) key() key {
	return key{
		player: g.player,
		turn:   g.turn,
	}
}

type key struct {
	player pos
	turn   int
}
