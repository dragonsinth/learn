package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var sample = `
Blueprint 1: Each ore robot costs 4 ore. Each clay robot costs 2 ore. Each obsidian robot costs 3 ore and 14 clay. Each geode robot costs 2 ore and 7 obsidian.
Blueprint 2: Each ore robot costs 2 ore. Each clay robot costs 3 ore. Each obsidian robot costs 3 ore and 8 clay. Each geode robot costs 3 ore and 12 obsidian.
`

var (
	re = regexp.MustCompile(`^Blueprint (\d+): Each ore robot costs (\d+) ore. Each clay robot costs (\d+) ore. Each obsidian robot costs (\d+) ore and (\d+) clay. Each geode robot costs (\d+) ore and (\d+) obsidian.$`)
)

func main() {
	main1()

	prod := 1
	for _, line := range strings.Split(sample, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if !re.MatchString(line) {
			panic(line)
		}

		matches := re.FindStringSubmatch(line)

		bp := &blueprint{
			index:             mustInt(matches[1]),
			oreCostOre:        mustInt(matches[2]),
			clayCostOre:       mustInt(matches[3]),
			obsidianCostOre:   mustInt(matches[4]),
			obsidianCostClay:  mustInt(matches[5]),
			geodeCostOre:      mustInt(matches[6]),
			geodeCostObsidian: mustInt(matches[7]),
		}

		if bp.index > 3 {
			continue
		}

		fmt.Println(bp.String())

		g := &game{
			cs:      bp.toCostSheet(),
			maxTurn: 32,
			bots:    [4]int{1, 0, 0, 0},
		}

		g = g.run()
		g.summary()

		fmt.Println(turns, done)
		turns, done = 0, 0

		prod *= g.resources[GEODE]
	}
	fmt.Println(prod)
}

func main1() {
	sum := 0
	for _, line := range strings.Split(sample, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if !re.MatchString(line) {
			panic(line)
		}

		matches := re.FindStringSubmatch(line)

		bp := &blueprint{
			index:             mustInt(matches[1]),
			oreCostOre:        mustInt(matches[2]),
			clayCostOre:       mustInt(matches[3]),
			obsidianCostOre:   mustInt(matches[4]),
			obsidianCostClay:  mustInt(matches[5]),
			geodeCostOre:      mustInt(matches[6]),
			geodeCostObsidian: mustInt(matches[7]),
		}

		fmt.Println(bp.String())

		g := &game{
			cs:      bp.toCostSheet(),
			maxTurn: 24,
			bots:    [4]int{1, 0, 0, 0},
		}

		g = g.run()
		g.summary()

		fmt.Println(turns, done)
		turns, done = 0, 0

		sum += bp.index * g.resources[GEODE]
	}
	fmt.Println(sum)
}

func mustInt(s string) int {
	if v, err := strconv.Atoi(s); err != nil {
		panic(fmt.Sprint(s, err))
	} else {
		return v
	}
}
