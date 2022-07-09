package main

import (
	"fmt"
	"github.com/dragonsinth/learn/aoc/2019/intcode"
)

var samples = []string{
	`109,1,204,-1,1001,100,1,100,1008,100,16,101,1006,101,0,99`,
	`1102,34915192,34915192,7,4,7,99,0`,
	`104,1125899906842624,99`,
}

func main() {
	for _, input := range samples {
		codes := intcode.Parse(input)
		m := intcode.NewIntMachine(codes, nil, func(v int) {
			fmt.Println(v)
		})
		m.Run()
		fmt.Println()
	}

	for _, startVal := range []int{1, 2} {
		codes := intcode.Parse(`99`)
		m := intcode.NewIntMachine(codes, func() int {
			return startVal
		}, func(v int) {
			fmt.Println(v)
		})
		m.Run()
		fmt.Println()
	}
}
