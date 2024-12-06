package main

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
)

const sample = `
47|53
97|13
97|61
97|47
75|29
61|13
75|53
29|13
97|29
53|29
61|53
97|53
61|29
47|13
75|47
97|75
47|61
75|61
47|29
75|13
53|13

75,47,61,53,29
97,61,53,29,13
75,29,13
75,97,47,61,53
61,13,29
97,13,75,29,47
`

func main() {
	p, examples := parse(sample)
	p.part1(examples)
	p.part2(examples)
}

func parse(input string) (*puz, [][]int) {
	p := &puz{
		rules: map[int]*rule{},
	}
	var examples [][]int

	lines := strings.Split(input, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if strings.Contains(line, "|") {
			parts := strings.Split(line, "|")
			prev, next := mustInt(parts[0]), mustInt(parts[1])

			nextRule := p.getRule(next)
			nextRule.mustComeAfter = append(nextRule.mustComeAfter, prev)

			prevRule := p.getRule(prev)
			prevRule.mustComeBefore = append(prevRule.mustComeBefore, next)
		} else if strings.Contains(line, ",") {
			nums := parseNums(line)
			examples = append(examples, nums)
		} else {
			panic(line)
		}
	}
	return p, examples
}

type rule struct {
	page           int
	mustComeAfter  []int
	mustComeBefore []int
}

func (r rule) isBefore(n int) bool {
	return slices.Contains(r.mustComeBefore, n)
}

func (r rule) isAfter(n int) bool {
	return slices.Contains(r.mustComeAfter, n)
}

type puz struct {
	rules map[int]*rule
}

func (p *puz) getRule(num int) *rule {
	if p.rules[num] == nil {
		p.rules[num] = &rule{
			page: num,
		}
	}
	return p.rules[num]
}

func (p *puz) part1(examples [][]int) {
	sum := 0

	for _, nums := range examples {
		if p.violates(nums) {
		} else {
			v := nums[len(nums)/2]
			sum += v
		}
	}

	fmt.Println(sum)
}

func (p *puz) part2(examples [][]int) {
	sum := 0
	for _, nums := range examples {
		if !p.violates(nums) {
			continue
		}

		slices.SortFunc(nums, func(a, b int) int {
			aRule := p.rules[a]
			if aRule == nil {
				panic("rule not found")
			}
			if aRule.isBefore(b) {
				return -1
			} else if aRule.isAfter(b) {
				return 1
			}
			panic("don't know")
		})

		if p.violates(nums) {
			panic("here")
		}
		v := nums[len(nums)/2]
		sum += v
	}

	fmt.Println(sum)
}

func (p *puz) violates(nums []int) bool {
	if len(nums) < 2 {
		return false
	}
	first, rest := nums[0], nums[1:]
	firstRule := p.rules[first]
	for _, num := range rest {
		if firstRule.isAfter(num) {
			return true
		}
	}

	return p.violates(rest)
}

func mustInt(s string) int {
	if v, err := strconv.Atoi(s); err != nil {
		panic(fmt.Sprint(s, err))
	} else {
		return v
	}
}

func parseNums(s string) []int {
	var ret []int
	for _, p := range strings.Split(s, ",") {
		ret = append(ret, mustInt(p))
	}
	return ret
}
