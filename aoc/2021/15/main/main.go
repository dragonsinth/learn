package main

import (
	"fmt"
	"regexp"
	"strings"
)

var sample = `
1163751742
1381373672
2136511328
3694931569
7463417111
1319128137
1359912421
3125421639
1293138521
2311944581
`

const embiggen = false

var input = sample

type puzzle struct {
	v        [][]byte
	bestRisk [][]int64
	height   int
	width    int
}

func (p *puzzle) Print() {
	for _, row := range p.v {
		fmt.Println(row)
	}
}

func (p *puzzle) PrintRisk() {
	for _, row := range p.bestRisk {
		fmt.Println(row)
	}
}

func (p *puzzle) Embiggen() {
	// Embiggen the original rows.
	for i, row := range p.v {
		for i, c := p.width, 5*p.width; i < c; i++ {
			v := row[i-p.width] + 1
			if v == 10 {
				v = 1
			}
			row = append(row, v)
		}
		p.v[i] = row
		p.bestRisk[i] = make([]int64, len(row))
	}
	p.width *= 5

	// Now add more rows based on the original rows.
	for i, c := p.height, 5*p.height; i < c; i++ {
		copyRow := p.v[i-p.height]
		newRow := make([]byte, p.width)
		for i, v := range copyRow {
			v++
			if v == 10 {
				v = 1
			}
			newRow[i] = v
		}
		p.v = append(p.v, newRow)
		p.bestRisk = append(p.bestRisk, make([]int64, p.width))
	}
	p.height *= 5
}

// Recursively
func (p *puzzle) MarkRiskDfs(r int, c int, risk int64) {
	if r < 0 || r >= p.height {
		return
	}
	if c < 0 || c >= p.width {
		return
	}
	if r == 0 && c == 0 && risk > 0 {
		return
	}

	if r != 0 || c != 0 {
		risk += int64(p.v[r][c])
	}
	if r := p.bestRisk[r][c]; r > 0 && r <= risk {
		return // already marked with a lower risk value
	}
	p.bestRisk[r][c] = risk

	// Try to navigate in each direction.
	p.MarkRiskDfs(r+1, c, risk)
	p.MarkRiskDfs(r, c+1, risk)
	p.MarkRiskDfs(r-1, c, risk)
	p.MarkRiskDfs(r, c-1, risk)
}

// Iteratively.
func (p *puzzle) MarkRiskBfs() {
	type task struct {
		r, c int
		risk int64
	}

	tasks := make([]task, 0, 100000)
	taskIt := 0
	tasks = append(tasks, task{0, 0, 0})

	oldMark := 0
	newMark := 0

	for taskIt < len(tasks) {
		t := tasks[taskIt]
		taskIt++
		r, c, risk := t.r, t.c, t.risk

		if r < 0 || r >= p.height {
			continue
		}
		if c < 0 || c >= p.width {
			continue
		}
		if r == 0 && c == 0 && risk > 0 {
			continue
		}
		if r != 0 || c != 0 {
			risk += int64(p.v[r][c])
		}

		bestRisk := p.bestRisk[r][c]
		if bestRisk == 0 {
			p.bestRisk[r][c] = risk
			newMark++
			if (newMark % 1000) == 0 {
				fmt.Println(newMark, oldMark, len(tasks)-taskIt)
			}
		} else if risk < bestRisk {
			p.bestRisk[r][c] = risk
			oldMark++
		} else {
			continue
		}

		tasks = append(tasks, task{r: r + 1, c: c, risk: risk})
		tasks = append(tasks, task{r: r, c: c + 1, risk: risk})
		tasks = append(tasks, task{r: r - 1, c: c, risk: risk})
		tasks = append(tasks, task{r: r, c: c - 1, risk: risk})

		if taskIt > 100000 {
			tasks = append([]task{}, tasks[taskIt:]...)
			taskIt = 0
		}
	}
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
		// Make a risk row for each value row.
		puz.bestRisk = append(puz.bestRisk, make([]int64, puz.width))
	}
	puz.height = len(puz.v)

	puz.Print()
	if embiggen {
		puz.Embiggen()
	}
	fmt.Println(puz.width, puz.height, puz.width*puz.height)
	puz.MarkRiskBfs()
	puz.PrintRisk()
}
