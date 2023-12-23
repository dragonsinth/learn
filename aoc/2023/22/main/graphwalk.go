package main

import "golang.org/x/exp/slices"

func graphWalk(deps [][]int, rdeps [][]int) int {
	sum := 0
	for b := range deps {
		sum += blocksKilledWithout(b, deps, rdeps)
	}
	return sum
}

func blocksKilledWithout(b int, deps [][]int, rdeps [][]int) int {
	if len(rdeps[b]) == 0 {
		return 0
	}

	killed := make([]bool, len(deps))
	killed[b] = true

	work := slices.Clone(rdeps[b])
	for i := 0; i < len(work); i++ {
		t := work[i]

		if killed[t] {
			continue // already dead
		}

		// are all my deps dead?
		dead := func() bool {
			for _, d := range deps[t] {
				if !killed[d] {
					return false
				}
			}
			return true // all dead
		}()

		if dead {
			killed[t] = true
			work = append(work, rdeps[t]...)
		}
	}
	sum := 0
	for _, k := range killed {
		if k {
			sum++
		}
	}
	return sum - 1 // don't count myself
}
