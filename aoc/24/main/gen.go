package main

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	inputRegex = regexp.MustCompile(`^inp ([w-z])$`)
	otherRegex = regexp.MustCompile(`^(mul|add|div|mod|eql) ([w-z]) ([w-z]|-?\d+)$`)
)

func generate() {
	var code []inst

	for _, line := range strings.Split(data, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "inp") {
			subs := inputRegex.FindStringSubmatch(line)
			if len(subs) != 2 {
				panic(subs)
			}
			code = append(code, inst{
				op:     inp,
				outReg: register(subs[1][0] - 'w'),
			})
			continue
		}

		if !otherRegex.MatchString(line) {
			panic(line)
		}

		subs := otherRegex.FindStringSubmatch(line)
		if len(subs) != 4 {
			panic(subs)
		}

		code = append(code, inst{
			op:     opCodes[subs[1]],
			outReg: register(subs[2][0] - 'w'),
			input:  subs[3],
		})
	}

	for _, c := range code {
		fmt.Println(c)
	}
}

type opCode int

const (
	inp = opCode(iota)
	add
	mul
	div
	mod
	eql
)

func (c opCode) Symbol() string {
	switch c {
	case add:
		return "+"
	case mul:
		return "*"
	case div:
		return "/"
	case mod:
		return "%"
	default:
		panic(c)
	}
}

var opCodes = map[string]opCode{
	"inp": inp,
	"add": add,
	"mul": mul,
	"div": div,
	"mod": mod,
	"eql": eql,
}

type register int

const (
	W = register(iota)
	X
	Y
	Z
)

func (r register) String() string {
	return string(rune(r + 'w'))
}

type inst struct {
	op     opCode
	outReg register
	input  string
}

func (i inst) String() string {
	switch i.op {
	case inp:
		return fmt.Sprintf("%s = in.next()", i.outReg)
	case add, mul, div, mod:
		return fmt.Sprintf("%s = %s %s %s", i.outReg, i.outReg, i.op.Symbol(), i.input)
	case eql:
		return fmt.Sprintf("%s = eq(%s, %s)", i.outReg, i.outReg, i.input)
	default:
		panic(i.op)
	}
}
