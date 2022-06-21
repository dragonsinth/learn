package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var sample = `
42: 9 14 | 10 1
9: 14 27 | 1 26
10: 23 14 | 28 1
1: "a"
11: 42 31
5: 1 14 | 15 1
19: 14 1 | 14 14
12: 24 14 | 19 1
16: 15 1 | 14 14
31: 14 17 | 1 13
6: 14 14 | 1 14
2: 1 24 | 14 4
0: 8 11
13: 14 3 | 1 12
15: 1 | 14
17: 14 2 | 1 7
23: 25 1 | 22 14
28: 16 1
4: 1 1
20: 14 14 | 1 15
3: 5 14 | 16 1
27: 1 6 | 14 18
14: "b"
21: 14 1 | 1 14
25: 1 1 | 1 14
22: 14 14
8: 42
26: 14 22 | 1 20
18: 15 15
7: 14 5 | 1 21
24: 14 1

abbbbbabbbaaaababbaabbbbabababbbabbbbbbabaaaa
bbabbbbaabaabba
babbbbaabbbbbabbbbbbaabaaabaaa
aaabbbbbbaaaabaababaabababbabaaabbababababaaa
bbbbbbbaaaabbbbaaabbabaaa
bbbababbbbaaaaaaaabbababaaababaabab
ababaaaaaabaaab
ababaaaaabbbaba
baabbaaaabbaaaababbaababb
abbbbabbbbaaaababbbbbbaaaababb
aaaaabbaabaaaaababaa
aaaabbaaaabbaaa
aaaabbaabbaaaaaaabbbabbbaaabbaabaaa
babaaabbbaaabaababbaabababaaab
aabbbbbaabbbaaaaaabbbbbababaaaaabbaaabba
`

var (
	input = sample

	ruleRe = regexp.MustCompile(`^(\d+): (.*)$`)
	msgRe  = regexp.MustCompile(`^[ab]+$`)
)

type char byte

const (
	A = char(254)
	B = char(255)
)

type seq []char

type rule struct {
	raw      []seq
	resolved []string
	minLen   int
	maxLen   int
}

func main() {
	// Support up to 256 rules
	rules := make([]rule, 256)
	var inputs []string
	maxRule := -1
	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if ruleRe.MatchString(line) {
			rNum, r := parseRule(line)
			rules[rNum] = r
			if rNum > maxRule {
				maxRule = rNum
			}
		} else if msgRe.MatchString(line) {
			inputs = append(inputs, line)
		} else {
			panic(line)
		}
	}

	rules = rules[:maxRule+1]
	rs := resolveRules(rules)

	r42 := rs.rules[42]
	r31 := rs.rules[31]
	for _, v := range []int{r42.minLen, r42.maxLen, r31.minLen, r31.maxLen} {
		if v != r42.minLen {
			panic(v)
		}
	}

	// Ensure no chunks in 42 and 31 match.
	for _, v1 := range r42.resolved {
		for _, v2 := range r31.resolved {
			if v1 == v2 {
				panic(v1)
			}
		}
	}

	count := 0
	for _, in := range inputs {
		if rs.tryMatchCustom(in) {
			fmt.Println("+++: ", in)
			count++
		} else {
			fmt.Println("---: ", in)
		}
	}

	fmt.Println(count)
}

func parseRule(line string) (int, rule) {
	matches := ruleRe.FindStringSubmatch(line)
	if len(matches) != 3 {
		panic(matches)
	}

	ruleNum := mustInt(matches[1])
	expr := matches[2]
	switch expr {
	case `"a"`:
		return ruleNum, rule{raw: []seq{{A}}}
	case `"b"`:
		return ruleNum, rule{raw: []seq{{B}}}
	}

	// Parse it
	var ret rule
	var curSeq seq
	tokens := strings.Split(expr, " ")
	for _, token := range tokens {
		if token == "|" {
			if len(curSeq) == 0 {
				panic("here")
			}
			ret.raw = append(ret.raw, curSeq)
			curSeq = nil
		} else {
			curSeq = append(curSeq, char(mustInt(token)))
		}
	}
	if len(curSeq) == 0 {
		panic("here")
	}
	ret.raw = append(ret.raw, curSeq)
	return ruleNum, ret
}

func mustInt(s string) int {
	if v, err := strconv.Atoi(s); err != nil {
		panic(fmt.Sprint(s, err))
	} else {
		return v
	}
}
