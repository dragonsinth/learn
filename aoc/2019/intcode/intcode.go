package intcode

import (
	"fmt"
	"strconv"
	"strings"
)

func NewIntMachine(intCodes []int, input func() int, output func(v int)) *IntMachine {
	p := &IntMachine{
		pc:     0,
		rel:    0,
		mem:    make([]int, allocSize(len(intCodes))),
		input:  input,
		output: output,
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
	pc     int // program counter
	rel    int // relative base
	mem    []int
	input  func() int
	output func(v int)
}

func (m *IntMachine) Read(addr int) int {
	return m.mem[addr]
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
			*c = m.input()
		case 4:
			a := m.load(op / 100)
			m.output(a)
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
			*c = BoolVal(a < b)
		case 8:
			a := m.load(op / 100)
			b := m.load(op / 1000)
			c := m.addr(op / 10000)
			*c = BoolVal(a == b)
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

func BoolVal(v bool) int {
	if v {
		return 1
	} else {
		return 0
	}
}

func Parse(line string) []int {
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
