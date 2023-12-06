package main

import (
	"fmt"
	"sort"
	"time"
)

func main() {
	countWins := countWinsFast
	fmt.Println(countWins(7, 9))
	fmt.Println(countWins(15, 40))
	fmt.Println(countWins(30, 200))
	fmt.Println(countWins(71530, 940200))
	fmt.Println(countWins(45977295, 305106211101695))
}

func countWinsSlow(t int, toBeat int) int {
	start := time.Now()
	defer func() {
		fmt.Println(time.Since(start))
	}()
	sum := 0
	lo, hi := 0, 0
	for charge := 1; charge < t; charge++ {
		dist := charge * (t - charge)
		if dist > toBeat {
			sum++
			if lo == 0 {
				lo = charge
			}
			hi = charge
		}
	}
	fmt.Println(lo, hi, sum)
	return sum
}

func countWinsFast(t int, toBeat int) int {
	start := time.Now()
	defer func() {
		fmt.Println(time.Since(start))
	}()

	// sort.Search: false, false, false, true, true, true

	// Find the lowest number <= mid that wins
	mid := t / 2
	lo := sort.Search(t, func(n int) bool {
		if n > mid {
			return true
		}
		return n*(t-n) > toBeat
	})

	// Find the lowest number >= mid that loses
	hi := sort.Search(t, func(n int) bool {
		if n < mid {
			return false
		}
		return n*(t-n) <= toBeat
	})

	ret := hi - lo
	fmt.Println(lo, hi, ret)
	return ret
}
