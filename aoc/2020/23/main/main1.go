package main

import (
	"fmt"
	"sort"
)

func run1(input string, rounds int, size int) {
	state := make([]byte, size)
	for i, c := range input {
		state[i] = byte(c - '0')
	}
	for i := len(input); i < len(state); i++ {
		state[i] = byte(i + 1)
	}
	sz := byte(size)

	for move := 0; move < rounds; move++ {
		fmtPrintf("-- move %d --\n", move+1)
		fmtPrintf("cups: %+v\n", state)

		pickUp := state[1:4]
		fmtPrintf("pick up: %+v\n", pickUp)
		newState := concat2(state[0:1], state[4:])

	outer:
		for tgt := state[0] - 1; true; tgt = (tgt + sz) % (sz + 1) {
			for idx, v := range newState {
				if v == tgt {
					// found insertion point; insert after
					idx++
					fmtPrintf("destination: %d\n", v)
					newState = concat2(newState[:idx], pickUp, newState[idx:])
					// rotate
					state = append(newState[1:], newState[0])
					break outer
				}
			}
		}

		fmtPrintln()
	}

	// Find the 1 cup
	fmt.Println(state)
	for idx, v := range state {
		if v == 1 {
			a := byte(idx+1) % sz
			b := byte(idx+2) % sz
			fmt.Println(state[a], state[b], int(state[a])*int(state[b]))
			sort.Slice(state, func(i, j int) bool {
				return state[i] < state[j]
			})
			fmt.Println(state)
			break
		}
	}
}

func concat2(in ...[]byte) []byte {
	ret := make([]byte, 0, 10)
	for _, e := range in {
		ret = append(ret, e...)
	}
	return ret
}
