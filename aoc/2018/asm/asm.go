package asm

import (
	"fmt"
	"regexp"
)

var (
	dirRe  = regexp.MustCompile(`^#ip (\d+)`)
	instRe = regexp.MustCompile(`^(\w+) (\d+) (\d+) (\d+)$`)
)

type Prog struct {
	IpReg int
	Code  []Instr
	Names [6]string
}

type Instr struct {
	Name string
	Fun  OpFun
	Ops  [3]int
}

type Registers [6]int

type OpFun func(ops [3]int, in Registers) Registers

func (p *Prog) Decompile() {
	for i, c := range p.Code {
		fmt.Printf("%2d: %20s//\n", i, c.ToString(p.Names))
	}
	fmt.Println()
}

func (p *Prog) Run(debug bool, init int) {
	var reg Registers
	ip := -1
	reg[0] = init
	for i := 0; true; i++ {
		ip++
		if ip < 0 || ip >= len(p.Code) {
			break
		}
		reg[p.IpReg] = ip
		c := p.Code[ip]
		shouldDebug := debug && i < 100
		if shouldDebug {
			fmt.Printf("i=%d ip=%d %+v %s %+v ", i, ip, reg, c.Name, c.Ops)
		}
		reg = c.Fun(c.Ops, reg)
		if shouldDebug {
			fmt.Println(reg, p.Code[ip].ToString(p.Names))
		}
		ip = reg[p.IpReg]
	}
	fmt.Println(reg)
}
