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
	part1(p, examples)
	part2(p, examples)
}

func parse(input string) (puz, [][]int) {
	p := puz{
		rules: map[int]*rule{},
	}
	examples := [][]int{}

	lines := strings.Split(input, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if strings.Contains(line, "|") {
			parts := strings.Split(line, "|")
			prev, next := mustInt(parts[0]), mustInt(parts[1])

			{
				r := p.rules[next]
				if r == nil {
					r = &rule{
						page: next,
					}
					p.rules[next] = r
				}
				r.mustComeAfter = append(r.mustComeAfter, prev)
			}

			{
				r := p.rules[prev]
				if r == nil {
					r = &rule{
						page: prev,
					}
					p.rules[prev] = r
				}
				r.mustComeBefore = append(r.mustComeBefore, next)
			}

		} else if strings.Contains(line, ",") {
			nums := parseNums(line)
			examples = append(examples, nums)
		} else {
			panic(line)
		}
	}
	return p, examples
}

func part1(p puz, examples [][]int) {
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

func part2(p puz, examples [][]int) {
	sum := 0
	for _, nums := range examples {
		if !p.violates(nums) {
			continue
		}

		slices.SortFunc(nums, func(a, b int) int {
			r := p.rules[a]
			if r == nil {
				panic("rule not found")
			}
			if slices.Contains(r.mustComeBefore, b) {
				return -1
			} else if slices.Contains(r.mustComeAfter, b) {
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

type rule struct {
	page           int
	mustComeAfter  []int
	mustComeBefore []int
}

type puz struct {
	rules map[int]*rule
}

func (p puz) violates(nums []int) bool {
	if len(nums) < 2 {
		return false
	}
	first, rest := nums[0], nums[1:]
	firstRule := p.rules[first]
	for _, num := range rest {
		if slices.Contains(firstRule.mustComeAfter, num) {
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
