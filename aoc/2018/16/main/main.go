package main

import (
	"fmt"
	"github.com/dragonsinth/learn/aoc/2018/asm"
	"regexp"
	"strconv"
	"strings"
)

const (
	sample = `
Before: [3, 2, 1, 1]
9 2 1 2
After:  [3, 2, 2, 1]
`
	data  = `` // puzzle input
	data2 = `` // puzzle input
)

var (
	re = []*regexp.Regexp{
		regexp.MustCompile(`^Before:\s*\[(\d+), (\d+), (\d+), (\d+)]$`),
		regexp.MustCompile(`^(\d+) (\d+) (\d+) (\d+)$`),
		regexp.MustCompile(`^After:\s*\[(\d+), (\d+), (\d+), (\d+)]$`),
	}
)

func main() {
	// part 1
	part1(sample)
	possible := part1(data)

	// part 2
	codec := solve(possible)
	fmt.Println(codec)
	run(codec, data2)
}

func part1(in string) map[string]map[int]bool {
	var lines []string
	for _, line := range strings.Split(in, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		lines = append(lines, line)
	}
	if len(lines)%3 != 0 {
		panic(len(lines))
	}

	possible := map[string]map[int]bool{}
	for name := range asm.Instrs {
		possible[name] = map[int]bool{}
	}

	sum := 0
	for i := 0; i < len(lines); i += 3 {
		for lines[i] == "" {
			i++
		}
		in := parseQuad(re[0], lines[i])
		code := parseQuad(re[1], lines[i+1])
		out := parseQuad(re[2], lines[i+2])

		matches := 0
		for name, fun := range asm.Instrs {
			opCode := code[0]
			var ops [3]int
			copy(ops[:], code[1:])
			actual := fun(ops, in)
			if actual == out {
				possible[name][opCode] = true
				matches++
			}
		}
		if matches >= 3 {
			sum++
		}
	}
	fmt.Println(sum)
	return possible
}

func solve(possible map[string]map[int]bool) [16]string {
	var ret [16]string
	for len(possible) > 0 {
		// Find the one with only one possibility
		func() {
			for name, options := range possible {
				if len(options) == 1 {
					// match
					for code := range options {
						ret[code] = name
						delete(possible, name)

						// remove this option from every other code
						for _, v := range possible {
							delete(v, code)
						}
						return
					}
				}
			}
			panic("no possibilities")
		}()
	}
	return ret
}

func run(codec [16]string, in string) {
	var reg asm.Registers
	for _, line := range strings.Split(in, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		code := parseQuad(re[1], line)
		opCode := code[0]
		var ops [3]int
		copy(ops[:], code[1:])
		fun := asm.Instrs[codec[opCode]]
		reg = fun(ops, reg)
	}
	fmt.Println(reg)
}

func parseQuad(r *regexp.Regexp, s string) asm.Registers {
	if !r.MatchString(s) {
		panic(s)
	}
	m := r.FindStringSubmatch(s)
	return asm.Registers{
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
