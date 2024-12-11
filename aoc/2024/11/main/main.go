package main

import (
	"fmt"
	"strconv"
	"strings"
)

func main() {
	run("125 17", 25)
	run("125 17", 75)
}

func run(input string, steps int) {
	sum := 0
	p := puz{cache: map[entry]int{}}
	for _, n := range parseNums(input) {
		sum += p.countBlocks(n, steps)
	}
	fmt.Println(sum)
}

type entry struct {
	n     int
	steps int
}

type puz struct {
	cache map[entry]int
}

func (p *puz) countBlocks(n int, steps int) int {
	if steps == 0 {
		return 1
	}
	e := entry{n: n, steps: steps}
	if v, ok := p.cache[e]; ok {
		return v
	}

	// must calc
	ret := func() int {
		steps--
		if n == 0 {
			return p.countBlocks(1, steps)
		}
		str := strconv.Itoa(n)
		if len(str)%2 == 0 {
			sz := len(str) / 2
			left := mustInt(str[:sz])
			right := mustInt(str[sz:])
			return p.countBlocks(left, steps) + p.countBlocks(right, steps)
		}
		return p.countBlocks(n*2024, steps)
	}()

	p.cache[e] = ret
	return ret
}

func parseNums(s string) []int {
	var ret []int
	for _, p := range strings.Fields(s) {
		ret = append(ret, mustInt(p))
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
