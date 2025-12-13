package astar

import (
	"iter"

	"github.com/dragonsinth/learn/aoc/sliceheap"
)

// RunConsistent runs with a consistent heuristic; panics if the heuristic is inconsistent.
func RunConsistent[T any, K comparable](
	start T,
	key func(st T) K,
	cost func(st T) int,
	heuristic func(st T) int,
	neighbors func(st T) iter.Seq[T],
) T {
	totalCost := func(st T) int {
		return cost(st) + heuristic(st)
	}

	work := sliceheap.New[T](func(a, b T) bool {
		return totalCost(a) < totalCost(b)
	})

	work.Push(start)

	seen := map[K]int{}

	for {
		w := work.Pop()
		k := key(w)
		c := cost(w)
		if prev, ok := seen[k]; ok {
			if prev <= c {
				continue // already found an as-cheap route to this spot
			} else {
				panic("should not replace old cost, inconsistent heuristic")
			}
		}
		seen[k] = c

		if heuristic(w) == 0 {
			return w // found
		}

		// try advancing each next state
		for n := range neighbors(w) {
			work.Push(n)
		}
	}
}
