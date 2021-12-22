package main

import "fmt"

func main() {
	//p1pos, p2pos := 0, 9
	ret := computeWins(key{
		pos:   [2]int{3, 7},
		score: [2]int{0, 0},
		turn:  0,
	})

	fmt.Println(ret)
}

type key struct {
	pos   [2]int
	score [2]int
	turn  int
}

type wins [2]int

var memoize = map[key]wins{}

func computeWins(state key) wins {
	if r, ok := memoize[state]; ok {
		return r
	}

	turn := state.turn
	curPos := state.pos[turn]
	curScore := state.score[turn]

	var ret wins
	for i := 1; i <= 3; i++ {
		for j := 1; j <= 3; j++ {
			for k := 1; k <= 3; k++ {
				roll := i + j + k
				nextPos := curPos + roll
				nextPos %= 10
				nextScore := curScore + nextPos + 1
				if nextScore >= 21 {
					// this player won
					ret[turn]++
				} else {
					// other player's turn
					nextState := state
					nextState.pos[turn] = nextPos
					nextState.score[turn] = nextScore
					nextState.turn = (turn + 1) % 2

					shardWins := computeWins(nextState)
					ret[0] += shardWins[0]
					ret[1] += shardWins[1]
				}
			}
		}
	}

	memoize[state] = ret
	return ret
}
