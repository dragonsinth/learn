package main

import (
	"fmt"
	"strconv"
	"strings"
)

var sample = `
???.### 1,1,3
.??..??...?##. 1,1,3
?#?#?#?#?#?#?#? 1,3,1,6
????.#...#... 4,1,1
????.######..#####. 1,6,5
?###???????? 3,2,1
`

func main() {
	run(sample, false)
	run(sample, true)
}

func run(input string, big bool) {
	sum := 0
	lines := strings.Split(input, "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}

		parts := strings.Split(line, " ")
		if len(parts) != 2 {
			panic(line)
		}

		in := parts[0]
		expect := parseNums(parts[1])

		if big {
			in, expect = embiggen(in, expect)
		}

		ans := newState(in, expect).Count(map[state]int{})
		sum += ans
	}
	fmt.Println(sum)
}

func embiggen(in string, expect []int) (string, []int) {
	var bigData string
	var bigExpect []int
	for i := 0; i < 5; i++ {
		if i > 0 {
			bigData += "?"
		}
		bigData += in
		bigExpect = append(bigExpect, expect...)
	}
	return bigData, bigExpect
}

func parseNums(s string) []int {
	var ret []int
	for _, p := range strings.Split(s, ",") {
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
