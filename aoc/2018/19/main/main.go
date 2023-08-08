package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const sample = `
#ip 0
seti 5 0 1
seti 6 0 2
addi 0 1 0
addr 1 2 3
setr 1 0 0
seti 8 0 4
seti 9 0 5
`

var (
	dirRe  = regexp.MustCompile(`^#ip (\d+)`)
	instRe = regexp.MustCompile(`^(\w+) (\d+) (\d+) (\d+)$`)
)

type prog struct {
	ipReg int
	code  []instr
}

type instr struct {
	name string
	f    fn
	ops  [3]int
}

type regs [6]int

type fn func(ops [3]int, in regs) regs

func main() {
	run(sample, true, 0)
	run(data, true, 0)

	decompile(data)

	// 977, 10551377
	reg := translate(0)
	fmt.Println(reg)
	//reg = translate(1)
	//fmt.Println(reg)

	for _, n := range []int{977, 10551377} {
		pf := primeFactors(n)
		fmt.Println(pf)
		f := combine(pf)
		fmt.Println(f)
		fmt.Println(sum(f))
	}
}

func decompile(input string) {
	p := parse(input)

	names := []string{
		"ax",
		"bx",
		"cx",
		"dx",
		"ex",
		"fx",
	}
	names[p.ipReg] = "pc"
	for i, c := range p.code {
		l := strconv.Itoa(i)
		if len(l) == 1 {
			l = " " + l
		}
		cs := printCode(c, names)
		for len(cs) < 20 {
			cs = cs + " "
		}
		cs = cs + "; "
		fmt.Println(l, ":", cs)
	}
	fmt.Println()
}

func run(input string, debug bool, init int) {
	p := parse(input)
	names := []string{
		"ax",
		"bx",
		"cx",
		"dx",
		"ex",
		"fx",
	}
	names[p.ipReg] = "pc"

	var reg regs
	ip := -1
	reg[0] = init
	for i := 0; true; i++ {
		ip++
		if ip < 0 || ip >= len(p.code) {
			break
		}
		reg[p.ipReg] = ip
		c := p.code[ip]
		if debug && i < 100 {
			fmt.Printf("i=%d ip=%d %+v %s %+v ", i, ip, reg, c.name, c.ops)
		}
		reg = c.f(c.ops, reg)
		if debug && i < 100 {
			cs := printCode(p.code[ip], names)
			fmt.Println(reg, cs)
		}
		ip = reg[p.ipReg]
	}
	fmt.Println(reg)
}

func parse(input string) *prog {
	p := prog{
		ipReg: -1,
	}

	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if p.ipReg < 0 {
			if !dirRe.MatchString(line) {
				panic(line)
			}
			m := dirRe.FindStringSubmatch(line)
			p.ipReg = mustInt(m[1])
		} else {
			if !instRe.MatchString(line) {
				panic(line)
			}
			m := instRe.FindStringSubmatch(line)
			name := m[1]
			f := instrs[name]
			if f == nil {
				panic(name)
			}
			p.code = append(p.code, instr{
				name: name,
				f:    f,
				ops: [3]int{
					mustInt(m[2]),
					mustInt(m[3]),
					mustInt(m[4]),
				},
			})

		}
	}
	return &p
}

func mustInt(s string) int {
	if v, err := strconv.Atoi(s); err != nil {
		panic(fmt.Sprint(s, err))
	} else {
		return v
	}
}
