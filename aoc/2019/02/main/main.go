package main

import (
	"fmt"
	"github.com/dragonsinth/learn/aoc/2019/intcode"
	"strings"
)

var sample = `
1,9,10,3,2,3,11,0,99,30,40,50
1,0,0,0,99
2,3,0,3,99
2,4,4,5,99,0
1,1,1,4,99,5,6,0,99
`

var input = `1,0,0,3,1,1,2,3,1,3,4,3,1,5,0,3,2,1,9,19,1,19,5,23,2,23,13,27,1,10,27,31,2,31,6,35,1,5,35,39,1,39,10,43,2,9,43,47,1,47,5,51,2,51,9,55,1,13,55,59,1,13,59,63,1,6,63,67,2,13,67,71,1,10,71,75,2,13,75,79,1,5,79,83,2,83,9,87,2,87,13,91,1,91,5,95,2,9,95,99,1,99,5,103,1,2,103,107,1,10,107,0,99,2,14,0,0`

func main() {
	lines := strings.Split(sample, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		codes := intcode.Parse(line)
		m := intcode.NewIntMachine(codes, nil, nil)
		m.Run()
		fmt.Println(m.Read(0))
	}

	for i := 0; i < 100; i++ {
		for j := 0; j < 100; j++ {
			codes := intcode.Parse(input)
			codes[1] = i
			codes[2] = j
			m := intcode.NewIntMachine(codes, nil, nil)
			m.Run()
			if v := m.Read(0); v == 19690720 {
				fmt.Println(i, j)
				return
			}
		}
	}
}
