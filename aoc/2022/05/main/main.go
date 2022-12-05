package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type input struct {
	stacks []string
	moves  string
}

var (
	sample = input{
		stacks: []string{
			"ZN",
			"MCD",
			"P",
		},
		moves: `
move 1 from 2 to 1
move 3 from 1 to 3
move 2 from 2 to 1
move 1 from 1 to 2
`,
	}
)

var (
	lineRe = regexp.MustCompile(`^move (\d+) from (\d+) to (\d+)$`)
)

func main() {
	main1()

	var stacks [][]byte
	for _, st := range sample.stacks {
		stacks = append(stacks, []byte(st))
	}

	for _, line := range strings.Split(sample.moves, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if !lineRe.MatchString(line) {
			panic(line)
		}

		matches := lineRe.FindStringSubmatch(line)
		if len(matches) != 4 {
			panic(matches)
		}

		count, from, to := mustInt(matches[1]), mustInt(matches[2]), mustInt(matches[3])

		// 0 base
		from--
		to--

		// whole stack
		pos := len(stacks[from]) - count
		moved := stacks[from][pos:]
		stacks[from] = stacks[from][:pos]
		stacks[to] = append(stacks[to], moved...)
	}

	for _, st := range stacks {
		fmt.Print(string(rune(st[len(st)-1])))
	}
	fmt.Println()
}

func main1() {
	var stacks [][]byte
	for _, st := range sample.stacks {
		stacks = append(stacks, []byte(st))
	}

	for _, line := range strings.Split(sample.moves, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if !lineRe.MatchString(line) {
			panic(line)
		}

		matches := lineRe.FindStringSubmatch(line)
		if len(matches) != 4 {
			panic(matches)
		}

		count, from, to := mustInt(matches[1]), mustInt(matches[2]), mustInt(matches[3])

		// 0 base
		from--
		to--
		for i := 0; i < count; i++ {
			tail := len(stacks[from]) - 1
			c := stacks[from][tail]
			stacks[from] = stacks[from][:tail]
			stacks[to] = append(stacks[to], c)
		}
	}

	for _, st := range stacks {
		fmt.Print(string(rune(st[len(st)-1])))
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
