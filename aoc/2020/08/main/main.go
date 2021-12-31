package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var sample = `
nop +0
acc +1
jmp +4
acc +3
jmp -3
acc -99
acc +1
jmp -4
acc +6
`

var (
	input = sample
)

type op int

const (
	NOP = op(iota)
	ACC
	JMP
)

func (op op) String() string {
	return []string{
		NOP: "nop",
		ACC: "acc",
		JMP: "jmp",
	}[op]
}

type inst struct {
	code op
	val  int
}

var (
	parse = regexp.MustCompile(`^([a-z]{3}) ([-+]\d+)$`)

	opMap = map[string]op{
		"nop": NOP,
		"acc": ACC,
		"jmp": JMP,
	}
)

func main() {
	var insts []inst

	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if !parse.MatchString(line) {
			panic(line)
		}

		parts := parse.FindStringSubmatch(line)
		if len(parts) != 3 {
			panic(parts)
		}

		code, ok := opMap[parts[1]]
		if !ok {
			panic(parts[1])
		}
		in := inst{
			code: code,
			val:  mustInt(parts[2]),
		}
		insts = append(insts, in)
		fmt.Println(in.code, in.val)
	}

	fmt.Println(runProgram(insts))

	// Try flipping every jmp -> nop, nop -> jump
	for i := range insts {
		in := &insts[i]
		if in.code == NOP {
			in.code = JMP
			v, ok := runProgram(insts)
			fmt.Println(v, ok)
			in.code = NOP
		} else if in.code == JMP {
			in.code = NOP
			v, ok := runProgram(insts)
			fmt.Println(v, ok)
			in.code = JMP
		}
	}
}

func runProgram(insts []inst) (int, bool) {
	seen := make([]bool, len(insts))
	acc := 0
	pc := 0
	for {
		if pc == len(insts) {
			return acc, true
		}
		if seen[pc] {
			return acc, false
		}
		seen[pc] = true

		in := &insts[pc]
		switch in.code {
		case NOP:
		case ACC:
			acc += in.val
		case JMP:
			pc += in.val
			continue
		default:
			panic(in.code)
		}
		pc++
	}
}

func mustInt(s string) int {
	if v, err := strconv.Atoi(s); err != nil {
		panic(fmt.Sprint(s, err))
	} else {
		return v
	}
}
