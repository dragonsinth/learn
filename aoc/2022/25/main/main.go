package main

import (
	"fmt"
	"strings"
)

var sample = `
1=-0-2
12111
2=0=
21
2=01
111
20012
112
1=-1=
1-12
12
1=
122
`

func main() {
	run(sample)
}

func run(input string) {
	sum := 0
	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		v := decode(line)
		sum += v
	}
	fmt.Println(sum)
	fmt.Println(encode(sum))
}

var digitVals = map[byte]int{
	'=': -2,
	'-': -1,
	'0': 0,
	'1': 1,
	'2': 2,
}

func decode(line string) int {
	in := []byte(line)
	sum := 0
	mul := 1
	for i := range in {
		c := in[len(in)-i-1]
		if v, ok := digitVals[c]; ok {
			sum += mul * v
		} else {
			panic(c)
		}
		mul *= 5
	}
	return sum
}

func encode(v int) string {
	var out []byte
	for v != 0 {
		digit := v % 5
		out = append(out, "012=-"[digit])
		if digit > 2 {
			digit = digit - 5
		}
		v -= digit
		v = v / 5
	}

	// reverse
	for i := 0; i < len(out)/2; i++ {
		out[i], out[len(out)-i-1] = out[len(out)-i-1], out[i]
	}
	return string(out)
}
