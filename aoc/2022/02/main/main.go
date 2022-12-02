package main

import (
	"fmt"
	"strconv"
	"strings"
)

var sample = `
A Y
B X
C Z
`

type choice int

const (
	rock = choice(iota)
	paper
	scissors
)

func (c choice) score() int {
	return int(c) + 1
}

func (c choice) outcome(other choice) outcome {
	if c == other {
		return draw
	}
	if (c+1)%3 == other {
		// e.g. me rock vs. them paper
		return lose
	}
	if c == (other+1)%3 {
		// e.g. me paper vs. them rock
		return win
	}
	panic("should not")
}

type outcome int

const (
	lose = outcome(iota)
	draw
	win
)

func (oc outcome) score() int {
	switch oc {
	case lose:
		return 0
	case draw:
		return 3
	case win:
		return 6
	default:
		panic(oc)
	}
}

func (oc outcome) versus(opponent choice) choice {
	switch oc {
	case lose:
		// lose to paper with rock
		return (opponent + 2) % 3
	case draw:
		return opponent
	case win:
		// beat to paper with scissors
		return (opponent + 1) % 3
	default:
		panic(oc)
	}
}

func main() {
	main1()

	score := 0

	for _, line := range strings.Split(sample, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		opponent := choice(line[0] - 'A')
		if opponent < rock || opponent > scissors {
			panic(opponent)
		}

		wantOutcome := outcome(line[2] - 'X')
		if wantOutcome < lose || wantOutcome > win {
			panic(wantOutcome)
		}

		myself := wantOutcome.versus(opponent)
		fmt.Println(myself.score(), wantOutcome.score())
		score += myself.score() + wantOutcome.score()
	}

	fmt.Println(score)
}

func main1() {
	score := 0

	for _, line := range strings.Split(sample, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		opponent := choice(line[0] - 'A')
		if opponent < rock || opponent > scissors {
			panic(opponent)
		}

		myself := choice(line[2] - 'X')
		if myself < rock || myself > scissors {
			panic(myself)
		}

		oc := myself.outcome(opponent)

		fmt.Println(myself.score(), oc.score())
		score += myself.score() + oc.score()
	}

	fmt.Println(score)
}

func mustInt(s string) int {
	if v, err := strconv.Atoi(s); err != nil {
		panic(fmt.Sprint(s, err))
	} else {
		return v
	}
}
