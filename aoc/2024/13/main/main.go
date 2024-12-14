package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const sample = `
Button A: X+94, Y+34
Button B: X+22, Y+67
Prize: X=8400, Y=5400

Button A: X+26, Y+66
Button B: X+67, Y+21
Prize: X=12748, Y=12176

Button A: X+17, Y+86
Button B: X+84, Y+37
Prize: X=7870, Y=6450

Button A: X+69, Y+23
Button B: X+27, Y+71
Prize: X=18641, Y=10279
`

var (
	reB = regexp.MustCompile(`^Button [AB]: X\+(\d+), Y\+(\d+)$`)
	reP = regexp.MustCompile(`^Prize: X=(\d+), Y=(\d+)`)
)

func main() {
	run(sample, 0, true)
	run(sample, 10000000000000, true)
}

func run(input string, adjust int, debug bool) {
	sum := 0
	lines := strings.Split(input, "\n")
	var A, B, C, M, N, K int
	for _, line := range lines {
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "Button") {
			if !reB.MatchString(line) {
				panic(line)
			}
			matches := reB.FindStringSubmatch(line)
			if len(matches) != 3 {
				panic(matches)
			}

			if line[7] == 'A' {
				A, M = mustInt(matches[1]), mustInt(matches[2])
			} else if line[7] == 'B' {
				B, N = mustInt(matches[1]), mustInt(matches[2])
			} else {
				panic(line)
			}
		} else if strings.HasPrefix(line, "Prize: ") {
			if !reP.MatchString(line) {
				panic(line)
			}
			matches := reP.FindStringSubmatch(line)
			if len(matches) != 3 {
				panic(matches)
			}
			C, K = mustInt(matches[1]), mustInt(matches[2])
			C += adjust
			K += adjust

			a := rat(N*C-B*K, N*A-B*M)
			b := rat(M*C-A*K, M*B-A*N)
			if a.isInt() && b.isInt() {
				cost := a.val()*3 + b.val()
				if debug {
					fmt.Println(a.val(), b.val(), cost)
				}
				sum += cost
			}
		} else {
			panic(line)
		}
	}
	fmt.Println(sum)
}

func rat(n int, d int) rational {
	g := gcd(n, d)
	n = n / g
	d = d / g
	if d < 0 {
		n = -n
		d = -d
	}
	return rational{n: n, d: d}
}

type rational struct {
	n int
	d int
}

func (r rational) isInt() bool {
	return r.d == 1
}

func (r rational) val() int {
	return r.n / r.d
}

func gcd(a, b int) int {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

func mustInt(s string) int {
	if v, err := strconv.Atoi(s); err != nil {
		panic(fmt.Sprint(s, err))
	} else {
		return v
	}
}
