package main

import (
	"fmt"
	"strconv"
)

func main() {
	fmt.Println(power(3, 5, 8))
	fmt.Println(power(122, 79, 57))
	fmt.Println(power(217, 196, 39))
	fmt.Println(power(101, 153, 71))

	fmt.Println(bestOfSize(18, 3))
	fmt.Println(bestOfSize(18, 16))
	fmt.Println(bestOfSize(42, 12))

	fmt.Println(best(18))
	fmt.Println(best(42))
}

func best(serial int) (int, int, int) {
	var best, bestLeft, bestTop, bestSz int
	for sz := 1; sz <= 30; sz++ {
		score, l, t := bestOfSize(serial, sz)
		if score > best {
			best, bestLeft, bestTop, bestSz = score, l, t, sz
		}
	}
	return bestLeft, bestTop, bestSz
}

func bestOfSize(serial int, size int) (int, int, int) {
	max := 300 - size + 1

	var best, bestLeft, bestTop int
	for left := 1; left <= max; left++ {
		for top := 1; top <= max; top++ {
			sum := 0
			for x, c := left, left+size; x < c; x++ {
				for y, c := top, top+size; y < c; y++ {
					sum += power(x, y, serial)
				}
			}
			if sum > best {
				best, bestLeft, bestTop = sum, left, top
			}
		}
	}
	return best, bestLeft, bestTop
}

func power(x int, y int, serial int) int {
	rackId := x + 10
	p := rackId * y
	p += serial
	p *= rackId
	p = (p / 100) % 10
	p = p - 5
	return p
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
func mustInt(s string) int {
	if v, err := strconv.Atoi(s); err != nil {
		panic(fmt.Sprint(s, err))
	} else {
		return v
	}
}
