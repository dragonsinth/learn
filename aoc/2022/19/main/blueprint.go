package main

import "fmt"

type blueprint struct {
	index             int
	oreCostOre        int
	clayCostOre       int
	obsidianCostOre   int
	obsidianCostClay  int
	geodeCostOre      int
	geodeCostObsidian int
}

func (bp *blueprint) String() string {
	return fmt.Sprintf(`Blueprint %d: Each ore robot costs %d ore. Each clay robot costs %d ore. Each obsidian robot costs %d ore and %d clay. Each geode robot costs %d ore and %d obsidian.`,
		bp.index,
		bp.oreCostOre,
		bp.clayCostOre,
		bp.obsidianCostOre,
		bp.obsidianCostClay,
		bp.geodeCostOre,
		bp.geodeCostObsidian,
	)
}

func (bp *blueprint) toCostSheet() *costsheet {
	return &costsheet{
		{bp.oreCostOre, 0, 0, 0},
		{bp.clayCostOre, 0, 0, 0},
		{bp.obsidianCostOre, bp.obsidianCostClay, 0, 0},
		{bp.geodeCostOre, 0, bp.geodeCostObsidian, 0},
	}
}
