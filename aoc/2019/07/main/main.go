package main

import (
	"fmt"
	"strconv"
	"strings"
)

var samples1 = []string{
	`3,15,3,16,1002,16,10,16,1,16,15,15,4,15,99,0,0`,
	`3,23,3,24,1002,24,10,24,1002,23,-1,23,101,5,23,23,1,24,23,23,4,23,99,0,0`,
	`3,31,3,32,1002,32,10,32,1001,31,-2,31,1007,31,0,33,1002,33,7,33,1,33,31,31,1,32,31,31,4,31,99,0,0,0`,
}

var samples2 = []string{
	`3,26,1001,26,-4,26,3,27,1002,27,2,27,1,27,26,27,4,27,1001,28,-1,28,1005,28,6,99,0,0,5`,
	`3,52,1001,52,-5,52,3,53,1,52,56,54,1007,54,5,55,1005,55,26,1001,54,-5,54,1105,1,12,1,53,54,53,1008,54,0,55,1001,55,1,55,2,53,55,53,4,53,1001,56,-1,56,1005,56,6,99,0,0,0,0,10`,
}

func main() {
	main1()
	main2()
}

func main2() {
	for _, in := range samples2 {
		codes := parse(in)

		bestVal := 0
		var best []int
		Perm([]int{5, 6, 7, 8, 9}, func(inputs []int) {
			in := make(chan int, 1)
			last := in
			var progs []*prog
			for _, input := range inputs {
				last <- input
				next := make(chan int, 1)
				p := &prog{
					pc:     0,
					codes:  append([]int{}, codes...),
					input:  last,
					output: next,
				}
				last = next
				progs = append(progs, p)
			}

			for _, p := range progs {
				go p.run()
			}

			answer := 0
			for {
				in <- answer
				v, ok := <-last
				if !ok {
					break
				}
				answer = v
			}
			if answer > bestVal {
				bestVal, best = answer, append([]int{}, inputs...)
			}
		})

		fmt.Println(bestVal, best)
	}
}

func main1() {
	for _, in := range samples1 {
		codes := parse(in)

		bestVal := 0
		var best []int
		Perm([]int{0, 1, 2, 3, 4}, func(inputs []int) {
			in := make(chan int, 1)
			last := in
			var progs []*prog
			for _, input := range inputs {
				last <- input
				next := make(chan int, 1)
				p := &prog{
					pc:     0,
					codes:  append([]int{}, codes...),
					input:  last,
					output: next,
				}
				last = next
				progs = append(progs, p)
			}

			for _, p := range progs {
				go p.run()
			}

			in <- 0
			v := <-last
			if v > bestVal {
				bestVal, best = v, append([]int{}, inputs...)
			}
		})

		fmt.Println(bestVal, best)
	}
}

func parse(line string) []int {
	vals := strings.Split(line, ",")
	var codes []int
	for _, v := range vals {
		codes = append(codes, mustInt(v))
	}
	return codes
}

type prog struct {
	pc     int
	codes  []int
	input  <-chan int
	output chan<- int
}

func (p *prog) run() {
	defer close(p.output)
	for {
		op := p.codes[p.pc]
		p.pc++
		if op == 99 {
			return
		}
		switch op % 100 {
		case 1:
			a := p.load(op / 100)
			b := p.load(op / 1000)
			c := p.addr(op / 10000)
			*c = a + b
		case 2:
			a := p.load(op / 100)
			b := p.load(op / 1000)
			c := p.addr(op / 10000)
			*c = a * b
		case 3:
			c := p.addr(op / 100)
			*c = <-p.input
		case 4:
			a := p.load(op / 100)
			p.output <- a
		case 5:
			a := p.load(op / 100)
			b := p.load(op / 1000)
			if a != 0 {
				p.pc = b
			}
		case 6:
			a := p.load(op / 100)
			b := p.load(op / 1000)
			if a == 0 {
				p.pc = b
			}
		case 7:
			a := p.load(op / 100)
			b := p.load(op / 1000)
			c := p.addr(op / 10000)
			*c = boolVal(a < b)
		case 8:
			a := p.load(op / 100)
			b := p.load(op / 1000)
			c := p.addr(op / 10000)
			*c = boolVal(a == b)
		default:
			panic(op)
		}
	}
}

func (p *prog) load(op int) int {
	v := p.codes[p.pc]
	p.pc++
	switch op % 10 {
	case 0:
		return p.codes[v]
	case 1:
		return v
	default:
		panic(op)
	}
}

func (p *prog) addr(op int) *int {
	a := p.codes[p.pc]
	p.pc++
	switch op % 10 {
	case 0:
		return &p.codes[a]
	default:
		panic(op)
	}
}

func boolVal(v bool) int {
	if v {
		return 1
	} else {
		return 0
	}
}

func mustInt(s string) int {
	if v, err := strconv.Atoi(s); err != nil {
		panic(fmt.Sprint(s, err))
	} else {
		return v
	}
}
