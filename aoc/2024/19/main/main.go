package main

import (
	"fmt"
	"regexp"
	"strings"
)

const sample = `
r, wr, b, g, bwu, rb, gb, br

brwrr
bggr
gbbr
rrbgbr
ubwu
bwurrg
brgr
bbrgwb
`

func main() {
	run(sample, true)
}

func run(input string, debug bool) {
	var rules []string
	var tests []string
	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if rules == nil {
			rules = strings.Split(line, ", ")
			continue
		}
		tests = append(tests, line)
	}
	re := buildRe(rules)

	sum := 0
	bigSum := 0
	dict := map[string]int{}
	for _, t := range tests {
		match := re.MatchString(t)
		perms := countPermsDict(dict, t, rules)
		if debug {
			fmt.Println(match, perms, t)
		}
		if match {
			sum++
		}
		bigSum += perms
	}
	fmt.Println(sum, bigSum)
}

func buildRe(alt []string) *regexp.Regexp {
	var sb strings.Builder
	sb.WriteString("^(")
	for i, e := range alt {
		if i > 0 {
			sb.WriteRune('|')
		}
		sb.WriteString(e)
	}
	sb.WriteString(")+$")
	return regexp.MustCompile(sb.String())
}

func countPerms(t string, rules []string) int {
	if t == "" {
		return 1
	}
	sum := 0
	for i, r := range rules {
		_ = i
		if strings.HasPrefix(t, r) {
			next := strings.TrimPrefix(t, r)
			sum += countPerms(next, rules)
		}
	}
	return sum
}

func countPermsDict(dict map[string]int, t string, rules []string) int {
	if t == "" {
		return 1
	}
	if v, ok := dict[t]; ok {
		return v
	}
	sum := 0
	for _, r := range rules {
		if strings.HasPrefix(t, r) {
			sum += countPermsDict(dict, strings.TrimPrefix(t, r), rules)
		}
	}
	dict[t] = sum
	return sum
}
