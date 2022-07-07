package main

import (
	"fmt"
)

var samples = []string{
	`109,1,204,-1,1001,100,1,100,1008,100,16,101,1006,101,0,99`,
	`1102,34915192,34915192,7,4,7,99,0`,
	`104,1125899906842624,99`,
}

func main() {
	for _, input := range samples {
		codes := parseIntCodes(input)
		m := NewIntMachine(codes)
		go m.Run()
		for {
			v, ok := <-m.Reader()
			if !ok {
				break
			}
			fmt.Println(v)
		}
		fmt.Println()
	}

	for _, startVal := range []int{1, 2} {
		codes := parseIntCodes(`99`)
		m := NewIntMachine(codes)
		go m.Run()

		m.Writer() <- startVal
		for {
			v, ok := <-m.Reader()
			if !ok {
				break
			}
			fmt.Println(v)
		}
		fmt.Println()
	}
}
