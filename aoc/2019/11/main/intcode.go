package main

import (
	"fmt"
	"strconv"
	"strings"
)

type input func() int
type output func(v int)

func NewIntMachine(intCodes []int, in input, out output) *IntMachine {
	p := &IntMachine{
		pc:  0,
		rel: 0,
		mem: make([]int, allocSize(len(intCodes))),
		in:  in,
		out: out,
	}
	copy(p.mem, intCodes)
	return p
}

func allocSize(min int) int {
	ret := 1024
	for ret < min {
		ret <<= 1
	}
	return ret
}

type IntMachine struct {
	pc  int // program counter
	rel int // relative base
	mem []int
	in  input
	out output
}

func (m *IntMachine) Run() {
	for {
		op := m.mem[m.pc]
		m.pc++
		if op == 99 {
			return
		}
		switch op % 100 {
		case 1:
			a := m.load(op / 100)
			b := m.load(op / 1000)
			c := m.addr(op / 10000)
			*c = a + b
		case 2:
			a := m.load(op / 100)
			b := m.load(op / 1000)
			c := m.addr(op / 10000)
			*c = a * b
		case 3:
			c := m.addr(op / 100)
			*c = m.in()
		case 4:
			a := m.load(op / 100)
			m.out(a)
		case 5:
			a := m.load(op / 100)
			b := m.load(op / 1000)
			if a != 0 {
				m.pc = b
			}
		case 6:
			a := m.load(op / 100)
			b := m.load(op / 1000)
			if a == 0 {
				m.pc = b
			}
		case 7:
			a := m.load(op / 100)
			b := m.load(op / 1000)
			c := m.addr(op / 10000)
			*c = boolVal(a < b)
		case 8:
			a := m.load(op / 100)
			b := m.load(op / 1000)
			c := m.addr(op / 10000)
			*c = boolVal(a == b)
		case 9:
			a := m.load(op / 100)
			m.rel += a
		default:
			panic(op)
		}
	}
}

func (m *IntMachine) load(op int) int {
	v := m.mem[m.pc]
	m.pc++
	switch op % 10 {
	case 0:
		return *m.ensureMem(v)
	case 1:
		return v
	case 2:
		return *m.ensureMem(v + m.rel)
	default:
		panic(op)
	}
}

func (m *IntMachine) addr(op int) *int {
	a := m.mem[m.pc]
	m.pc++
	switch op % 10 {
	case 0:
		return m.ensureMem(a)
	case 2:
		return m.ensureMem(a + m.rel)
	default:
		panic(op)
	}
}

func (m *IntMachine) ensureMem(addr int) *int {
	if addr >= len(m.mem) {
		// realloc
		newBuf := make([]int, allocSize(addr))
		copy(newBuf, m.mem)
		m.mem = newBuf
	}
	return &m.mem[addr]
}

func boolVal(v bool) int {
	if v {
		return 1
	} else {
		return 0
	}
}

func parseIntCodes(line string) []int {
	vals := strings.Split(line, ",")
	var codes []int
	for _, v := range vals {
		codes = append(codes, mustInt(v))
	}
	return codes
}

func mustInt(s string) int {
	if v, err := strconv.Atoi(s); err != nil {
		panic(fmt.Sprint(s, err))
	} else {
		return v
	}
}
