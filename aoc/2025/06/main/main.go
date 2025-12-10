package main

import (
	"fmt"
	"strconv"
	"strings"
)

const sample = `
123 328  51 64 
 45 64  387 23 
  6 98  215 314
*   +   *   +  
`

type problem struct {
	vals []int
	op   byte
}

func (p problem) solve() int {
	switch p.op {
	case '+':
		sum := 0
		for _, i := range p.vals {
			sum += i
		}
		return sum
	case '*':
		prod := 1
		for _, i := range p.vals {
			prod *= i
		}
		return prod
	default:
		panic("invalid op")
	}
}

func main() {
	fmt.Println(solve(parse1(sample)))
	fmt.Println(solve(parse2(sample)))
}

func solve(problems []problem) int {
	sum := 0
	for _, pro := range problems {
		sum += pro.solve()
	}
	return sum
}

func parse1(input string) []problem {
	var problems []problem
	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if problems == nil {
			problems = make([]problem, len(parts))
		} else if len(parts) != len(problems) {
			panic(len(parts))
		}
		if parts[0] == "*" || parts[0] == "+" {
			for i := range problems {
				problems[i].op = parts[i][0]
			}
		} else {
			for i := range problems {
				problems[i].vals = append(problems[i].vals, mustInt(parts[i]))
			}
		}
	}
	return problems
}

func parse2(input string) []problem {
	var problems []problem

	lines := strings.Split(strings.Trim(input, "\n"), "\n")
	sz := len(lines[0])
	var p *problem
	for col := sz - 1; col >= 0; col-- {
		empty := true
		acc := 0
		for i := range lines {
			switch c := lines[i][col]; c {
			case ' ':
				// continue
			case '+', '*':
				p.op = c
				empty = false
			default:
				if c < '0' || c > '9' {
					panic(c)
				}
				if p == nil {
					p = &problem{}
				}
				acc *= 10
				acc += int(c - '0')
				empty = false
			}
		}
		if !empty {
			p.vals = append(p.vals, acc)
			if p.op != 0 {
				problems = append(problems, *p)
				p = nil
			}
		}
	}
	return problems
}

func mustInt(s string) int {
	if v, err := strconv.Atoi(s); err != nil {
		panic(fmt.Sprint(s, err))
	} else {
		return v
	}
}
