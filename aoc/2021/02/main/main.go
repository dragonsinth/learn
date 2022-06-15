package main

import (
	"fmt"
	"strconv"
	"strings"
)

var sample = `
forward 5
down 5
forward 8
up 3
down 8
forward 2
`

func main1() {
	x, y := 0, 0
	for _, line := range strings.Split(sample, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.Split(line, " ")
		if len(parts) != 2 {
			panic(line)
		}

		dist, err := strconv.Atoi(parts[1])
		if err != nil {
			panic(fmt.Sprint(parts[1], err))
		}

		switch parts[0] {
		case "forward":
			x += dist
		case "down":
			y += dist
		case "up":
			y -= dist
			if y < 0 {
				panic(y)
			}
		}
	}

	fmt.Println(x, y, x*y)
}

func main() {
	x, y, aim := 0, 0, 0
	for _, line := range strings.Split(sample, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.Split(line, " ")
		if len(parts) != 2 {
			panic(line)
		}

		dist, err := strconv.Atoi(parts[1])
		if err != nil {
			panic(fmt.Sprint(parts[1], err))
		}

		switch parts[0] {
		case "forward":
			x += dist
			y += dist * aim
			if y < 0 {
				panic(y)
			}
		case "down":
			aim += dist
		case "up":
			aim -= dist
		}
	}

	fmt.Println(x, y, x*y)
}
