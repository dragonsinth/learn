package main

import (
	"bytes"
	"fmt"
	"slices"
	"strconv"
	"strings"
)

var (
	sample1 = []byte{0, 1, 5, 4, 3, 0}
	sample2 = []byte{0, 3, 5, 4, 3, 0}
)

func main() {
	part1(sample1, 729)
	fmt.Println()

	dump(sample2)
	fmt.Println()

	part2(sample2)
	fmt.Println()
}

func part1(ops []byte, a int) {
	p := &Program{
		Ops: ops,
		A:   a,
	}
	fmt.Println(stringBytes(p.Run()))
}

func part2(ops []byte) {
	best = ""
	fmt.Println(findSolutionBits(ops, 0, 3))
}

func dump(ops []byte) {
	opString := []string{
		0: "adv",
		1: "bxl",
		2: "bst",
		3: "jnz",
		4: "bxc",
		5: "out",
		6: "bdv",
		7: "cdv",
	}

	for ip := 0; ip < len(ops); ip += 2 {
		code, op := ops[ip], ops[ip+1]
		fmt.Println(ip, ":", opString[code], opDesc(code, op))
	}
}

func opDesc(code byte, op byte) string {
	switch code {
	case 0, 6, 7:
		if op < 4 {
			return strconv.Itoa(1 << int(op))
		}
		reg := op - 4
		return "2^" + "ABC"[reg:reg+1]
	case 1, 3:
		return strconv.Itoa(int(op))
	case 2, 5:
		if op < 4 {
			return strconv.Itoa(int(op))
		}
		reg := op - 4
		return "ABC"[reg:reg+1] + " % 8"
	case 4:
		return "-"
	default:
		panic(code)
	}
}

var best string

func findSolution(ops []byte, guess int) int {
	if guess > 1<<60 {
		return 0 // give up
	}

	p := &Program{
		Ops: ops,
		A:   guess,
	}
	sol := p.Run()

	if !bytes.HasSuffix(ops, sol) {
		return 0
	}
	if bytes.Equal(ops, sol) {
		return guess
	}

	dbg := stringBytesValue(guess)
	if len(dbg) > len(best) {
		best = dbg
		fmt.Println(dbg)
	}

	if ans := findSolutionBits(ops, guess, 3); ans != 0 {
		return ans
	}
	if ans := findSolutionBits(ops, guess, 6); ans != 0 {
		return ans
	}
	if ans := findSolutionBits(ops, guess, 9); ans != 0 {
		return ans
	}
	if ans := findSolutionBits(ops, guess, 12); ans != 0 {
		return ans
	}

	// give up
	return 0
}

func findSolutionBits(ops []byte, guess int, bits int) int {
	newGuess := guess << bits
	for i := 0; i < 1<<bits; i++ {
		next := newGuess | i
		if next == 0 {
			continue
		}
		if ans := findSolution(ops, next); ans != 0 {
			return ans
		}
	}
	return 0
}

func stringBytes(in []byte) string {
	var out strings.Builder
	for i, e := range in {
		if i > 0 {
			out.WriteRune(',')
		}
		out.WriteRune(rune('0' + e))
	}
	return out.String()
}

func stringBytesValue(v int) string {
	var b []byte
	for v > 0 {
		b = append(b, byte(v%8))
		v /= 8
	}
	slices.Reverse(b)
	return stringBytes(b)
}
