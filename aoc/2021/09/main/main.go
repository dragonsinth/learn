package main

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

var sample = `
2199943210
3987894921
9856789892
8767896789
9899965678
`

var input = sample

type puzzle struct {
	v     [][]byte
	basin [][]byte
	width int
}

func (p *puzzle) Print() {
	for _, row := range p.v {
		fmt.Println(row)
	}
}

func (p *puzzle) PrintBasins() {
	for _, row := range p.basin {
		fmt.Println(row)
	}
}
func (p *puzzle) IsMin(row, col int) (byte, bool) {
	myScore := byte(0)
	minScore := byte(127)
	for r := row - 1; r <= row+1; r++ {
		for c := col - 1; c <= col+1; c++ {
			score := p.Score(r, c)
			if r == row && c == col {
				myScore = score
			} else {
				if score < minScore {
					minScore = score
				}
			}
		}
	}
	return myScore, myScore < minScore
}

func (p *puzzle) Score(r int, c int) byte {
	if r < 0 || r >= len(p.v) {
		return 128
	}
	if c < 0 || c >= p.width {
		return 128
	}
	return p.v[r][c]
}

// Recursively
func (p *puzzle) MarkBasin(r int, c int, basin byte) bool {
	if r < 0 || r >= len(p.v) {
		return false
	}
	if c < 0 || c >= p.width {
		return false
	}

	if p.v[r][c] == 9 {
		return false // not a basin
	}
	if p.basin[r][c] > 0 {
		return false // already marked
	}
	p.basin[r][c] = basin

	// Recursively flood-fill neighbors.
	p.MarkBasin(r-1, c, basin)
	p.MarkBasin(r+1, c, basin)
	p.MarkBasin(r, c-1, basin)
	p.MarkBasin(r, c+1, basin)
	return true
}

var parseLine = regexp.MustCompile(`^\d+$`)

func main() {
	var puz puzzle
	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if !parseLine.MatchString(line) {
			panic(line)
		}

		buf := []byte(line)
		if puz.width == 0 {
			puz.width = len(buf)
		} else if puz.width != len(buf) {
			panic(line)
		}

		row := make([]byte, puz.width)
		for i, c := range buf {
			row[i] = c - '0'
		}
		puz.v = append(puz.v, row)
		// Make a basin row for each value row.
		puz.basin = append(puz.basin, make([]byte, puz.width))
	}

	puz.Print()

	nextBasin := byte(1)
	for r, row := range puz.v {
		for c := range row {
			if puz.MarkBasin(r, c, nextBasin) {
				nextBasin++
			}
		}
	}

	puz.PrintBasins()

	// Count the number of nodes in each basin.
	counts := make([]int, nextBasin)
	for _, row := range puz.basin {
		for _, basin := range row {
			if basin > 0 {
				counts[basin]++
			}
		}
	}

	// Sort asc
	sort.Ints(counts)
	fmt.Println(counts)
	topThree := counts[len(counts)-3:]
	fmt.Println(topThree)
	fmt.Println(topThree[0] * topThree[1] * topThree[2])
}

func main1() {
	var puz puzzle
	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if !parseLine.MatchString(line) {
			panic(line)
		}

		buf := []byte(line)
		if puz.width == 0 {
			puz.width = len(buf)
		} else if puz.width != len(buf) {
			panic(line)
		}

		row := make([]byte, puz.width)
		for i, c := range buf {
			row[i] = c - '0'
		}
		puz.v = append(puz.v, row)
	}

	puz.Print()

	sum := 0
	for r, row := range puz.v {
		for c := range row {
			score, isMin := puz.IsMin(r, c)
			if isMin {
				fmt.Println(r, c, "->", score)
				sum += int(score) + 1
			}
		}
	}

	fmt.Println(sum)
}
