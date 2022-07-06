package main

import (
	"fmt"
	"strconv"
	"strings"
)

var samples = []string{
	`109,1,204,-1,1001,100,1,100,1008,100,16,101,1006,101,0,99`,
	`1102,34915192,34915192,7,4,7,99,0`,
	`104,1125899906842624,99`,
}

func main() {
	for _, input := range samples {
		codes := parse(input)
		mem := make([]int, 1024)
		copy(mem, codes)
		in, out := make(chan int), make(chan int)
		p := &prog{
			pc:     0,
			rel:    0,
			mem:    mem,
			input:  in,
			output: out,
		}
		go p.run()
		for {
			v, ok := <-out
			if !ok {
				break
			}
			fmt.Println(v)
		}
		fmt.Println()
	}

	for _, startVal := range []int{1, 2} {
		codes := parse("99")
		mem := make([]int, 16*1024)
		copy(mem, codes)
		in, out := make(chan int), make(chan int)
		p := &prog{
			pc:     0,
			rel:    0,
			mem:    mem,
			input:  in,
			output: out,
		}
		go p.run()
		in <- startVal
		for {
			v, ok := <-out
			if !ok {
				break
			}
			fmt.Println(v)
		}
		fmt.Println()
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
	pc     int // program counter
	rel    int // relative base
	mem    []int
	input  <-chan int
	output chan<- int
}

func (p *prog) run() {
	defer close(p.output)
	for {
		op := p.mem[p.pc]
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
		case 9:
			a := p.load(op / 100)
			p.rel += a
		default:
			panic(op)
		}
	}
}

func (p *prog) load(op int) int {
	v := p.mem[p.pc]
	p.pc++
	switch op % 10 {
	case 0:
		return p.mem[v]
	case 1:
		return v
	case 2:
		return p.mem[v+p.rel]
	default:
		panic(op)
	}
}

func (p *prog) addr(op int) *int {
	a := p.mem[p.pc]
	p.pc++
	switch op % 10 {
	case 0:
		return &p.mem[a]
	case 2:
		return &p.mem[a+p.rel]
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
