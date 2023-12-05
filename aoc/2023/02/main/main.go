package main

import (
	"fmt"
	"strconv"
	"strings"
)

var sample = `
Game 1: 3 blue, 4 red; 1 red, 2 green, 6 blue; 2 green
Game 2: 1 blue, 2 green; 3 green, 4 blue, 1 red; 1 green, 1 blue
Game 3: 8 green, 6 blue, 20 red; 5 blue, 4 red, 13 green; 5 green, 1 red
Game 4: 1 green, 3 red, 6 blue; 3 green, 6 red; 3 green, 15 blue, 14 red
Game 5: 6 red, 1 blue, 3 green; 2 blue, 1 red, 2 green
`

func main() {
	sum, powSum := 0, 0
	for _, line := range strings.Split(sample, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.Split(line, ": ")
		if len(parts) != 2 {
			panic(parts)
		}
		id := mustInt(strings.Split(parts[0], " ")[1])

		limit := []int{12, 13, 14}
		inLimit := true
		mins := []int{0, 0, 0}

		draws := strings.Split(parts[1], "; ")
		for _, draw := range draws {
			shards := strings.Split(draw, ", ")
			for _, shard := range shards {
				parts := strings.Split(shard, " ")
				count := mustInt(parts[0])
				color := mustColor(parts[1])
				mins[color] = max(mins[color], count)
				if count > limit[color] {
					inLimit = false
				}
			}
		}

		pow := mins[0] * mins[1] * mins[2]
		fmt.Println("id", id, "=", inLimit, pow)
		if inLimit {
			sum += id
		}
		powSum += pow

	}
	fmt.Println(sum, powSum)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func mustInt(s string) int {
	if v, err := strconv.Atoi(s); err != nil {
		panic(fmt.Sprint(s, err))
	} else {
		return v
	}
}

type color int

const (
	RED = color(iota)
	GREEN
	BLUE
)

var colors = []string{"red", "green", "blue"}

func (c color) String() string {
	return colors[c]
}

func mustColor(s string) color {
	for i, c := range colors {
		if s == c {
			return color(i)
		}
	}
	panic(s)
}
