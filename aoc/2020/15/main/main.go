package main

import "fmt"

var (
	inputs = [][]int{
		{0, 3, 6},
		{1, 3, 2},
		{2, 1, 3},
		{1, 2, 3},
		{2, 3, 1},
		{3, 2, 1},
		{3, 1, 2},
		{0, 1, 5, 10, 3, 12, 19},
	}
)

func main() {
	for i, input := range inputs {
		r := result(input, 30000000)
		fmt.Println(i, r)
	}
}

func result(input []int, nTurns int) int {
	lastTurnNum := make([]int, nTurns)
	lastNum := 0
	lastNumPrevious := 0

	speak := func(turn, what int) {
		lastNum = what
		lastNumPrevious = lastTurnNum[what]
		lastTurnNum[what] = turn
	}

	for turn := 1; turn <= nTurns; turn++ {
		if turn <= len(input) {
			speak(turn, input[turn-1])
			continue
		}
		if lastNumPrevious == 0 {
			speak(turn, 0)
		} else {
			speak(turn, turn-lastNumPrevious-1)
		}
	}

	return lastNum
}
