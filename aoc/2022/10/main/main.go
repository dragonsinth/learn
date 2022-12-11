package main

import (
	"fmt"
	"strconv"
	"strings"
)

func main() {
	var instr []int
	for _, line := range strings.Split(sample, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		switch {
		case line == "noop":
			instr = append(instr, 0) // 1-cycle
		case strings.HasPrefix(line, "addx "):
			v := mustInt(strings.TrimPrefix(line, "addx "))
			instr = append(instr, 0, v) // 2-cycle
		default:
			panic(line)
		}
	}

	sum := 0
	x := 1
	for i, v := range instr {
		cycle := i + 1
		if cycle%40 == 20 {
			strength := cycle * x
			sum += strength
			fmt.Println(cycle, x, strength)
		}
		x += v
	}
	fmt.Println(sum)

	x = 1
	for i, v := range instr {
		pos := i % 40
		if pos == 0 {
			fmt.Println()
		}
		if abs(x-pos) < 2 {
			fmt.Print("#")
		} else {
			fmt.Print(".")
		}
		x += v
	}
	fmt.Println()
}

func mustInt(s string) int {
	if v, err := strconv.Atoi(s); err != nil {
		panic(fmt.Sprint(s, err))
	} else {
		return v
	}
}

func abs(i int) int {
	if i < 0 {
		return -i
	} else {
		return i
	}
}
