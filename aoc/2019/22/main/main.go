package main

import (
	"fmt"
	"strconv"
	"strings"
)

var samples = []string{
	`
deal with increment 7
deal into new stack
deal into new stack`,
	`
cut 6
deal with increment 7
deal into new stack`,
	`
deal with increment 7
deal with increment 9
cut -2`,
	`
deal into new stack
cut -2
deal with increment 7
cut 8
cut -4
deal with increment 7
cut 3
deal with increment 9
deal with increment 3
cut -1`,
}

const (
	count = 119315717514047
	times = 101741582076661
)

func main() {
	for _, sample := range samples {
		codes := parse(sample)
		d := newDeck(10).run(codes)
		fmt.Println(d)
	}

	codes := parse(samples[3])

	d := newDeck(10007).run(codes)
	var toReverse int64
	for i := 0; i < len(d); i++ {
		if d[i] == 2019 {
			fmt.Println(i)
			toReverse = int64(i)
			break
		}
	}
	rsmall := newReverseDeck(codes, 10007)
	fmt.Println(rsmall.reverse(toReverse))
	fmt.Println(rsmall.Solve(toReverse, 0))
	fmt.Println(rsmall.Solve(toReverse, 1))

	rbig := newReverseDeck(codes, count)
	fmt.Println(rbig.Solve(2020, times))
}

type op struct {
	code, arg int
}

func parse(input string) []op {
	var ops []op
	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		switch {
		case line == "":
		case line == "deal into new stack":
			ops = append(ops, op{})
		case strings.HasPrefix(line, "cut "):
			rem := strings.TrimPrefix(line, "cut ")
			ops = append(ops, op{code: 1, arg: mustInt(rem)})
		case strings.HasPrefix(line, "deal with increment "):
			rem := strings.TrimPrefix(line, "deal with increment ")
			ops = append(ops, op{code: 2, arg: mustInt(rem)})
		default:
			panic(line)
		}
	}

	return ops
}

func mustInt(s string) int {
	if v, err := strconv.Atoi(s); err != nil {
		panic(fmt.Sprint(s, err))
	} else {
		return v
	}
}
