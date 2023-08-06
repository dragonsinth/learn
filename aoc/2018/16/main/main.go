package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const sample = `
Before: [3, 2, 1, 1]
9 2 1 2
After:  [3, 2, 2, 1]
`

var (
	re = []*regexp.Regexp{
		regexp.MustCompile(`^Before:\s*\[(\d+), (\d+), (\d+), (\d+)]$`),
		regexp.MustCompile(`^(\d+) (\d+) (\d+) (\d+)$`),
		regexp.MustCompile(`^After:\s*\[(\d+), (\d+), (\d+), (\d+)]$`),
	}
)

const prog = ""

func main() {
	var lines []string
	for _, line := range strings.Split(sample, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		lines = append(lines, line)
	}
	if len(lines)%3 != 0 {
		panic(len(lines))
	}

	impossible := [16][16]bool{}

	sum := 0
	for i := 0; i < len(lines); i += 3 {
		for lines[i] == "" {
			i++
		}
		in := parseQuad(re[0], lines[i])
		code := parseQuad(re[1], lines[i+1])
		out := parseQuad(re[2], lines[i+2])

		matches := 0
		for i, inst := range instrs {
			actual := inst.fn(code, in)
			if actual == out {
				matches++
			} else {
				impossible[i][code[0]] = true
			}
		}
		if matches >= 3 {
			sum++
		}
	}
	fmt.Println(sum)

	for i, inst := range instrs {
		if impossible[i][inst.code] {
			panic(i)
		}
	}

	// part 2
	reg := [4]int{}
	for _, line := range strings.Split(prog, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		code := parseQuad(re[1], line)
		reg = instrs[code[0]].fn(code, reg)
	}
	fmt.Println(reg)
}

func parseQuad(r *regexp.Regexp, s string) [4]int {
	if !r.MatchString(s) {
		panic(s)
	}
	m := r.FindStringSubmatch(s)
	return [4]int{
		mustInt(m[1]),
		mustInt(m[2]),
		mustInt(m[3]),
		mustInt(m[4]),
	}
}

func mustInt(s string) int {
	if v, err := strconv.Atoi(s); err != nil {
		panic(fmt.Sprint(s, err))
	} else {
		return v
	}
}

type instr struct {
	code int
	name string
	fn   func(code [4]int, reg [4]int) [4]int
}

var instrs = [16]instr{
	{0, "eqri", eqri},
	{1, "banr", banr},
	{2, "bori", bori},
	{3, "mulr", mulr},
	{4, "seti", seti},
	{5, "bani", bani},
	{6, "muli", muli},
	{7, "gtrr", gtrr},
	{8, "setr", setr},
	{9, "addi", addi},
	{10, "gtir", gtir},
	{11, "borr", borr},
	{12, "addr", addr},
	{13, "eqrr", eqrr},
	{14, "gtri", gtri},
	{15, "eqir", eqir},
}

func addr(code [4]int, reg [4]int) [4]int {
	out := reg
	A := reg[code[1]]
	B := reg[code[2]]
	C := code[3]
	out[C] = A + B
	return out
}

func addi(code [4]int, reg [4]int) [4]int {
	out := reg
	A := reg[code[1]]
	B := code[2]
	C := code[3]
	out[C] = A + B
	return out
}

func mulr(code [4]int, reg [4]int) [4]int {
	out := reg
	A := reg[code[1]]
	B := reg[code[2]]
	C := code[3]
	out[C] = A * B
	return out
}

func muli(code [4]int, reg [4]int) [4]int {
	out := reg
	A := reg[code[1]]
	B := code[2]
	C := code[3]
	out[C] = A * B
	return out
}

func banr(code [4]int, reg [4]int) [4]int {
	out := reg
	A := reg[code[1]]
	B := reg[code[2]]
	C := code[3]
	out[C] = A & B
	return out
}

func bani(code [4]int, reg [4]int) [4]int {
	out := reg
	A := reg[code[1]]
	B := code[2]
	C := code[3]
	out[C] = A & B
	return out
}

func borr(code [4]int, reg [4]int) [4]int {
	out := reg
	A := reg[code[1]]
	B := reg[code[2]]
	C := code[3]
	out[C] = A | B
	return out
}

func bori(code [4]int, reg [4]int) [4]int {
	out := reg
	A := reg[code[1]]
	B := code[2]
	C := code[3]
	out[C] = A | B
	return out
}

func setr(code [4]int, reg [4]int) [4]int {
	out := reg
	A := reg[code[1]]
	C := code[3]
	out[C] = A
	return out
}

func seti(code [4]int, reg [4]int) [4]int {
	out := reg
	A := code[1]
	C := code[3]
	out[C] = A
	return out
}

func gtir(code [4]int, reg [4]int) [4]int {
	out := reg
	A := code[1]
	B := reg[code[2]]
	C := code[3]
	out[C] = toInt(A > B)
	return out
}

func gtri(code [4]int, reg [4]int) [4]int {
	out := reg
	A := reg[code[1]]
	B := code[2]
	C := code[3]
	out[C] = toInt(A > B)
	return out
}

func gtrr(code [4]int, reg [4]int) [4]int {
	out := reg
	A := reg[code[1]]
	B := reg[code[2]]
	C := code[3]
	out[C] = toInt(A > B)
	return out
}

func eqir(code [4]int, reg [4]int) [4]int {
	out := reg
	A := code[1]
	B := reg[code[2]]
	C := code[3]
	out[C] = toInt(A == B)
	return out
}

func eqri(code [4]int, reg [4]int) [4]int {
	out := reg
	A := reg[code[1]]
	B := code[2]
	C := code[3]
	out[C] = toInt(A == B)
	return out
}

func eqrr(code [4]int, reg [4]int) [4]int {
	out := reg
	A := reg[code[1]]
	B := reg[code[2]]
	C := code[3]
	out[C] = toInt(A == B)
	return out
}

func toInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
