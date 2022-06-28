package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
	"time"
)

var (
	fInput  = flag.String("in", `389125467`, "input")
	fRounds = flag.Int("rounds", 100, "rounds")
	fSz     = flag.Int("size", 9, "size")

	w io.Writer
)

func main() {
	flag.Parse()
	start := time.Now()
	defer func() {
		fmt.Println(time.Since(start))
	}()

	if *fSz < 256 {
		//f1, _ := os.Create("1.txt")
		//defer f1.Close()
		//w = f1
		w = io.Writer(os.Stdout)
		run1(*fInput, *fRounds, *fSz)
	}
	//f2, _ := os.Create("2.txt")
	//defer f2.Close()
	//w = f2
	run2(*fInput, *fRounds, *fSz)
	//w = nil
	//run2(*fInput, 10000, 1000000)
	w = nil
	//run2(*fInput, 100000, 1000000)
}

func run2(input string, rounds int, size int) {
	sz := int32(size)
	buf := make([]int32, sz*2)
	state, newState := buf[:sz], buf[sz:]
	for i, c := range input {
		state[i] = c - '0'
	}
	for i := len(input); i < len(state); i++ {
		state[i] = int32(i + 1)
	}

	for move := 0; move < rounds; move++ {
		if rounds-move <= 100 {
			w = os.Stdout
		}
		fmtPrintf("-- move %d --\n", move+1)
		//fmtPrintf("cups: %+v\n", state)

		pickUp := state[1:4]
		fmtPrintf("pick up: %+v\n", pickUp)

	outer:
		for tgt := state[0] - 1; true; tgt = (tgt + sz) % (sz + 1) {
			if contains(pickUp, tgt) || tgt == 0 {
				continue
			}
			for idx, v := range state {
				if v == tgt {
					// found insertion point; insert after
					fmtPrintf("destination: %d at %d\n", v, idx)
					idx++

					// Create a new state which is:
					// 1) Elements 4:idx
					// 2) Elements pickUp
					// 3) Elements idx:
					// 4) Element 0
					concat(newState, state[4:idx], pickUp, state[idx:], state[0:1])

					state, newState = newState, state
					break outer
				}
			}
			panic("should not get here")
		}

		fmtPrintln()
	}

	// Find the 1 cup
	if len(state) <= 100 {
		fmt.Println(state)
	}
	for idx, v := range state {
		if v == 1 {
			a := int32(idx+1) % sz
			b := int32(idx+2) % sz
			fmt.Println(state[a], state[b], int(state[a])*int(state[b]))
			if len(state) <= 100 {
				sort.Slice(state, func(i, j int) bool {
					return state[i] < state[j]
				})
				fmt.Println(state)
			}
			break
		}
	}
}

func contains(up []int32, tgt int32) bool {
	for _, v := range up {
		if v == tgt {
			return true
		}
	}
	return false
}

func concat(dst []int32, in ...[]int32) {
	idx := 0
	for _, e := range in {
		n := copy(dst[idx:], e)
		if n != len(e) {
			panic(n)
		}
		idx += n
	}
	if idx != len(dst) {
		panic(idx)
	}
}

func fmtPrintln(args ...interface{}) {
	if w != nil {
		fmt.Fprintln(w, args...)
	}
}

func fmtPrintf(s string, args ...interface{}) {
	if w != nil {
		fmt.Fprintf(w, s, args...)
	}
}
