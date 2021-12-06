package main

import (
	"fmt"
	"strconv"
	"strings"
)

var sample = `3,4,3,1,2`

var input = sample

func main() {
	var fishByState [9]int64
	for _, s := range strings.Split(input, ",") {
		v := mustInt(s)
		if v < 0 || v > 6 {
			panic(v)
		}
		fishByState[v]++
	}

	fmt.Println(fishByState)
	for day := 1; day <= 256; day++ {
		// Save off the ready fish (state 0)
		readyFish := fishByState[0]

		// Move fish counts down a state.
		copy(fishByState[0:], fishByState[1:])
		fishByState[8] = 0

		// Now copy the readyFish inot both state 6 and 8.
		fishByState[6] += readyFish
		fishByState[8] = readyFish

		fmt.Println(day, fishByState)
	}
	sum := int64(0)
	for _, v := range fishByState {
		sum += v
	}
	fmt.Println(sum)
}

func mustInt(s string) int {
	if v, err := strconv.Atoi(s); err != nil {
		panic(fmt.Sprint(s, err))
	} else {
		return v
	}
}
