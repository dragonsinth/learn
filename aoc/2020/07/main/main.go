package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var sample = `
light red bags contain 1 bright white bag, 2 muted yellow bags.
dark orange bags contain 3 bright white bags, 4 muted yellow bags.
bright white bags contain 1 shiny gold bag.
muted yellow bags contain 2 shiny gold bags, 9 faded blue bags.
shiny gold bags contain 1 dark olive bag, 2 vibrant plum bags.
dark olive bags contain 3 faded blue bags, 4 dotted black bags.
vibrant plum bags contain 5 faded blue bags, 6 dotted black bags.
faded blue bags contain no other bags.
dotted black bags contain no other bags.
`

var sample2 = `
shiny gold bags contain 2 dark red bags.
dark red bags contain 2 dark orange bags.
dark orange bags contain 2 dark yellow bags.
dark yellow bags contain 2 dark green bags.
dark green bags contain 2 dark blue bags.
dark blue bags contain 2 dark violet bags.
dark violet bags contain no other bags.
`

var (
	input = sample
)

var (
	lineParse = regexp.MustCompile(`^(.+) bags contain (.+)\.$`)
	itemParse = regexp.MustCompile(`^(\d+) (.+) bags?$`)
)

type stuff struct {
	count int
	rhs   string
}

func main() {
	isContained := map[string][]string{}
	contains := map[string][]stuff{}

	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if strings.HasSuffix(line, "no other bags.") {
			continue
		}

		if !lineParse.MatchString(line) {
			panic(line)
		}

		parts := lineParse.FindStringSubmatch(line)
		lhs := parts[1]
		for _, item := range strings.Split(parts[2], ",") {
			item = strings.TrimSpace(item)
			if item == "" {
				continue
			}

			if !itemParse.MatchString(item) {
				panic(item)
			}

			itemParts := itemParse.FindStringSubmatch(item)
			//fmt.Println(itemParts)
			count := mustInt(itemParts[1])
			rhs := itemParts[2]
			fmt.Printf("%s => %d %s\n", lhs, count, rhs)
			isContained[rhs] = append(isContained[rhs], lhs)
			contains[lhs] = append(contains[lhs], stuff{
				count: count,
				rhs:   rhs,
			})
		}
	}

	seen := rmap(isContained, "shiny gold")
	fmt.Println(seen)
	fmt.Println(len(seen))

	fmt.Println(fmap(contains, "shiny gold"))
}

func fmap(contains map[string][]stuff, in string) int {
	sum := 0
	for _, it := range contains[in] {
		sum += it.count
		sum += it.count * fmap(contains, it.rhs)
	}
	return sum
}

func rmap(isContained map[string][]string, in string) map[string]bool {
	seen := map[string]bool{}
	work := append([]string{}, isContained[in]...)

	for i := 0; i < len(work); i++ {
		item := work[i]
		if seen[item] {
			continue
		}

		seen[item] = true
		work = append(work, isContained[item]...)
	}

	return seen
}

func mustInt(s string) int {
	if v, err := strconv.Atoi(s); err != nil {
		panic(fmt.Sprint(s, err))
	} else {
		return v
	}
}
