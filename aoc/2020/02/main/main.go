package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var sample = `
1-3 a: abcde
1-3 b: cdefg
2-9 c: ccccccccc
`

var input = sample

var parseLine = regexp.MustCompile(`^(\d+)-(\d+) ([a-z]): ([a-z]+)$`)

func main() {
	valid := 0
	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if !parseLine.MatchString(line) {
			panic(line)
		}

		subs := parseLine.FindStringSubmatch(line)
		if len(subs) != 5 {
			panic(subs)
		}

		first, second := mustInt(subs[1]), mustInt(subs[2])
		letter, password := subs[3][0], subs[4]

		if (password[first-1] == letter) != (password[second-1] == letter) {
			valid++
		}
	}
	fmt.Println(valid)
}

func main1() {
	valid := 0
	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if !parseLine.MatchString(line) {
			panic(line)
		}

		subs := parseLine.FindStringSubmatch(line)
		if len(subs) != 5 {
			panic(subs)
		}

		min, max := mustInt(subs[1]), mustInt(subs[2])
		letter, password := subs[3][0], subs[4]

		count := 0
		for _, b := range []byte(password) {
			if b == letter {
				count++
			}
		}
		if count >= min && count <= max {
			valid++
		}
	}
	fmt.Println(valid)
}

func mustInt(s string) int {
	if v, err := strconv.Atoi(s); err != nil {
		panic(fmt.Sprint(s, err))
	} else {
		return v
	}
}
