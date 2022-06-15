package main

import (
	"fmt"
	"strconv"
	"strings"
)

var sample = `
7,4,9,5,11,17,23,2,0,14,21,24,10,16,13,6,15,25,12,22,18,20,8,19,3,26,1

22 13 17 11  0
 8  2 23  4 24
21  9 14 16  7
 6 10  3 18  5
 1 12 20 15 19

 3 15  0  2 22
 9 18 13 17  5
19  8  7 25 23
20 11 10 24  4
14 21 16 12  6

14 21 17 24  4
10 16 15  9 19
18  8 23 26 20
22 11 13  6  5
 2  0 12  3  7
`

const (
	width = 5
)

type puzzle struct {
	numbers [width * width]int
	marks   [width * width]bool
}

func (p puzzle) Solved() bool {
	for r := 0; r < width; r++ {
		if p.RowSolved(r) {
			return true
		}
	}
	for c := 0; c < width; c++ {
		if p.ColSolved(c) {
			return true
		}
	}
	return false
}

func (p puzzle) RowSolved(r int) bool {
	for c := 0; c < width; c++ {
		if !p.marks[r*width+c] {
			return false
		}
	}
	return true
}

func (p puzzle) ColSolved(c int) bool {
	for r := 0; r < width; r++ {
		if !p.marks[r*width+c] {
			return false
		}
	}
	return true
}

func main1() {
	readState := 0
	var puzzles []*puzzle
	var sequence []int
	for _, line := range strings.Split(sample, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if readState == 0 {
			// sequence data
			for _, s := range strings.Split(line, ",") {
				val, err := strconv.Atoi(s)
				if err != nil {
					panic(fmt.Sprint(s, err))
				}
				sequence = append(sequence, val)
			}
			if len(sequence) < 3 {
				panic(line)
			}
		} else {
			// puzzle data
			puzzleNum := (readState - 1) / width
			if len(puzzles) <= puzzleNum {
				puzzles = append(puzzles, &puzzle{})
			}
			puz := puzzles[puzzleNum]
			row := (readState - 1) % width

			split := strings.Split(line, " ")
			i := 0
			for _, s := range split {
				if s == "" {
					continue
				}
				val, err := strconv.Atoi(s)
				if err != nil {
					panic(fmt.Sprint(s, err))
				}
				puz.numbers[row*width+i] = val
				i++
			}
			if i != width {
				panic(line)
			}
		}

		readState++
	}

	if readState != 1+width*len(puzzles) {
		panic(readState)
	}

	fmt.Println(sequence)
	fmt.Println(puzzles)

	// Run the game.
	for _, val := range sequence {
		for _, puz := range puzzles {
			for i, v := range puz.numbers {
				if val == v {
					puz.marks[i] = true
					break
				}
			}
		}

		for p, puz := range puzzles {
			if puz.Solved() {
				fmt.Println("solved", p, puz)
				// Sum it up
				sum := 0
				for i, v := range puz.numbers {
					if !puz.marks[i] {
						sum += v
					}
				}

				fmt.Println(val, sum, val*sum)
				return
			}
		}
	}
	panic("no puzzles solved")
}

func main() {
	readState := 0
	var puzzles []*puzzle
	var sequence []int
	for _, line := range strings.Split(sample, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if readState == 0 {
			// sequence data
			for _, s := range strings.Split(line, ",") {
				val, err := strconv.Atoi(s)
				if err != nil {
					panic(fmt.Sprint(s, err))
				}
				sequence = append(sequence, val)
			}
			if len(sequence) < 3 {
				panic(line)
			}
		} else {
			// puzzle data
			puzzleNum := (readState - 1) / width
			if len(puzzles) <= puzzleNum {
				puzzles = append(puzzles, &puzzle{})
			}
			puz := puzzles[puzzleNum]
			row := (readState - 1) % width

			split := strings.Split(line, " ")
			i := 0
			for _, s := range split {
				if s == "" {
					continue
				}
				val, err := strconv.Atoi(s)
				if err != nil {
					panic(fmt.Sprint(s, err))
				}
				puz.numbers[row*width+i] = val
				i++
			}
			if i != width {
				panic(line)
			}
		}

		readState++
	}

	if readState != 1+width*len(puzzles) {
		panic(readState)
	}

	fmt.Println(sequence)
	fmt.Println(puzzles)

	// Run the game.
	solved := make([]bool, len(puzzles))
	solvedCount := 0

	for _, val := range sequence {
		for _, puz := range puzzles {
			for i, v := range puz.numbers {
				if val == v {
					puz.marks[i] = true
					break
				}
			}
		}

		for p, puz := range puzzles {
			if !solved[p] && puz.Solved() {
				// A new puzzle was solved.
				fmt.Println("solved", p, puz)
				solved[p] = true
				solvedCount++
				if solvedCount == len(solved) {
					// Sum it up
					sum := 0
					for i, v := range puz.numbers {
						if !puz.marks[i] {
							sum += v
						}
					}

					fmt.Println(val, sum, val*sum)
					return
				}
			}
		}
	}
	panic("not all puzzles solved")
}
