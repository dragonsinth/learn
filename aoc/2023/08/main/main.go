package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var sample = `
px{a<2006:qkq,m>2090:A,rfg}
pv{a>1716:R,A}
lnx{m>1548:A,A}
rfg{s<537:gd,x>2440:R,A}
qs{s>3448:A,lnx}
qkq{x<1416:A,crn}
crn{x>2662:A,R}
in{s<1351:px,qqz}
qqz{s>2770:qs,m<1801:hdj,R}
gd{a>3333:R,R}
hdj{m>838:A,pv}

{x=787,m=2655,a=1222,s=2876}
{x=1679,m=44,a=2067,s=496}
{x=2036,m=264,a=79,s=2244}
{x=2461,m=1339,a=466,s=291}
{x=2127,m=1623,a=2188,s=1013}`

var reRule = regexp.MustCompile(`^([a-z]+){((?:(?:[xmas][<>]\d+:[a-zAR]+),|[AR],)+)([a-zAR]+)}$`)
var rePart = regexp.MustCompile(`^{x=(\d+),m=(\d+),a=(\d+),s=(\d+)}$`)

func main() {
	p, parts := parse(sample)
	p.part1(parts)
	p.part2()
}

func parse(input string) (puz, []part) {
	p := puz{
		rules: map[target]rule{},
	}
	var parts []part
	state := 0
	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			if len(p.rules) > 0 {
				state = 1
			}
			continue
		}

		if state == 0 {
			if !reRule.MatchString(line) {
				panic(line)
			}
			matches := reRule.FindStringSubmatch(line)
			t := parseTarget(matches[1])
			p.rules[t] = rule{
				name:  t,
				steps: parseSteps(matches[2]),
				def:   parseTarget(matches[3]),
			}
		} else {
			if !rePart.MatchString(line) {
				panic(line)
			}
			matches := rePart.FindStringSubmatch(line)
			parts = append(parts, part{
				X: mustInt(matches[1]),
				M: mustInt(matches[2]),
				A: mustInt(matches[3]),
				S: mustInt(matches[4]),
			})
		}
	}
	return p, parts
}

var reStep = regexp.MustCompile(`^([xmas])([<>])(\d+):([a-zAR]+)$`)

func parseSteps(in string) []step {
	var steps []step
	for _, p := range strings.Split(in, ",") {
		if p == "" {
			continue
		}
		if !reStep.MatchString(p) {
			panic(p)
		}
		matches := reStep.FindStringSubmatch(p)
		steps = append(steps, step{
			q:   parseQuality(matches[1][0]),
			sym: parseSym(matches[2][0]),
			val: mustInt(matches[3]),
			t:   parseTarget(matches[4]),
		})
	}
	return steps
}

func parseQuality(s byte) quality {
	switch s {
	case 'x':
		return X
	case 'm':
		return M
	case 'a':
		return A
	case 's':
		return S
	default:
		panic(s)
	}
}

func parseSym(s byte) sym {
	if s != '<' && s != '>' {
		panic(s)
	}
	return sym(s)
}

func parseTarget(s string) target {
	return target(s)
}

type puz struct {
	rules map[target]rule
}
type rule struct {
	name  target
	steps []step
	def   target
}

type step struct {
	q   quality
	sym sym
	val int
	t   target
}

type sym byte

type target string

type quality byte

const (
	X = quality(iota)
	M
	A
	S
)

func mustInt(s string) int {
	if v, err := strconv.Atoi(s); err != nil {
		panic(fmt.Sprint(s, err))
	} else {
		return v
	}
}
