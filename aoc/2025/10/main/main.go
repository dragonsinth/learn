package main

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

const sample = `
[.##.] (3) (1,3) (2) (2,3) (0,2) (0,1) {3,5,4,7}
[...#.] (0,2,3,4) (2,3) (0,4) (0,1,2) (1,2,3,4) {7,5,12,7,2}
[.###.#] (0,1,2,3,4) (0,3,4) (0,1,2,4,5) (1,2) {10,11,11,5,10,5}
`

var re = regexp.MustCompile(`^\[([.#]+)]\s+([0-9,()\s]+){([^}]+)}$`)

type puzzle struct {
	len      int
	goal     []bool
	buttons  []button
	joltages []uint16
}

type button []int

func main() {
	run(parse(sample), true)
}

func run(puzzles []puzzle, debug bool) {
	sum := 0
	for _, p := range puzzles {
		result := p.Part1()
		if debug {
			fmt.Println(result)
		}
		sum += result.cost
	}
	fmt.Println(sum)

	sum = 0
	for _, p := range puzzles {
		result := p.Part2()
		if debug {
			fmt.Println(result)
		}
		sum += result
	}
	fmt.Println(sum)
}

func parse(input string) []puzzle {
	var ret []puzzle
	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if !re.MatchString(line) {
			panic(line)
		}

		matches := re.FindStringSubmatch(line)
		if len(matches) != 4 {
			panic(strings.Join(matches, "\n"))
		}

		p := puzzle{
			len: len(matches[1]),
		}

		// .##.
		for _, c := range matches[1] {
			switch c {
			case '.':
				p.goal = append(p.goal, false)
			case '#':
				p.goal = append(p.goal, true)
			default:
				panic(c)
			}
		}

		// (3) (1,3) (2) (2,3) (0,2) (0,1)
		for _, bl := range strings.Fields(matches[2]) {
			// (1,3)
			var b button
			for _, nl := range strings.Split(strings.Trim(bl, "()"), ",") {
				b = append(b, mustInt(nl))
			}
			p.buttons = append(p.buttons, b)
		}

		// {10,11,11,5,10,5}
		for _, jl := range strings.Split(strings.Trim(matches[3], "{}"), ",") {
			p.joltages = append(p.joltages, mustUint16(jl))
		}

		ret = append(ret, p)
	}
	return ret
}

func mustInt(s string) int {
	if v, err := strconv.Atoi(s); err != nil {
		panic(fmt.Sprint(s, err))
	} else {
		return v
	}
}

func mustUint16(s string) uint16 {
	v := mustInt(s)
	if v < 0 || v > math.MaxUint16 {
		panic(v)
	}
	return uint16(v)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func gcd(a, b int) int {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

func lcm(a, b int) int {
	return a * (b / gcd(a, b))
}

// floor(a/b) for integers, with correct behavior for negatives
func divFloor(a, b int) int {
	q := a / b
	r := a % b
	// If there is a remainder and signs differ, round down one more step
	if r != 0 && ((r > 0) != (b > 0)) {
		q--
	}
	return q
}

// ceil(a/b) for integers, with correct behavior for negatives
func divCeil(a, b int) int {
	q := a / b
	r := a % b
	// If there is a remainder and signs are the same, round up
	if r != 0 && ((r > 0) == (b > 0)) {
		q++
	}
	return q
}
