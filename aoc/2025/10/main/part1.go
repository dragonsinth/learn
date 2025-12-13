package main

import (
	"iter"

	"github.com/dragonsinth/learn/aoc/astar"
)

const sz = 16

type lightState struct {
	curr    [sz]bool
	pressed [sz]byte
	cost    int
}

func (p *puzzle) Part1() lightState {
	return astar.RunConsistent[lightState, [sz]bool](
		lightState{},
		func(st lightState) [sz]bool {
			return st.curr
		},
		func(st lightState) int {
			return st.cost * p.len
		},
		func(st lightState) int {
			sum := 0
			for i := range p.goal {
				if st.curr[i] != p.goal[i] {
					sum++
				}
			}
			return sum
		},
		func(st lightState) iter.Seq[lightState] {
			return func(yield func(lightState) bool) {
				// try pushing each button
				for bi, b := range p.buttons {
					next := st
					for _, which := range b {
						next.curr[which] = !next.curr[which]
					}
					next.pressed[bi]++
					next.cost++
					if !yield(next) {
						return
					}
				}
			}
		},
	)
}
