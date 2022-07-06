package main

import (
	"fmt"
	"strconv"
	"strings"
)

var samples = []string{
	`3,9,8,9,10,9,4,9,99,-1,8`,
	`3,9,7,9,10,9,4,9,99,-1,8`,
	`3,3,1108,-1,8,3,4,3,99`,
	`3,3,1107,-1,8,3,4,3,99`,
	`3,12,6,12,15,1,13,14,13,4,13,99,-1,0,1,9`,
	`3,3,1105,-1,9,1101,0,0,12,4,12,99,1`,
	`3,21,1008,21,8,20,1005,20,22,107,8,21,20,1006,20,31,1106,0,36,98,0,0,1002,21,125,20,4,20,1105,1,46,104,999,1105,1,46,1101,1000,1,20,4,20,1105,1,46,98,99`,
}

func main() {
	for _, in := range samples {
		codes := parse(in)
		p := &prog{
			pc:    0,
			codes: codes,
			input: []int{8},
		}
		p.run()
		fmt.Println(p.output)
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
	input  []int
	output []int
}

func (p *prog) run() {
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
			*c = p.input[0]
			p.input = p.input[1:]
		case 4:
			a := p.load(op / 100)
			p.output = append(p.output, a)
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
