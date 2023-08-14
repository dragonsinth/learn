package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const sample = `
pos=<0,0,0>, r=4
pos=<1,0,0>, r=1
pos=<4,0,0>, r=3
pos=<0,2,0>, r=1
pos=<0,5,0>, r=3
pos=<0,0,3>, r=1
pos=<1,1,1>, r=1
pos=<1,1,2>, r=1
pos=<1,3,1>, r=1
`

const sample2 = `
pos=<10,12,12>, r=2
pos=<12,14,12>, r=2
pos=<16,12,12>, r=4
pos=<14,14,14>, r=6
pos=<50,50,50>, r=200
pos=<10,10,10>, r=5
`

const (
	dims = 3
)

var (
	re = regexp.MustCompile(`^pos=<(-?\d+),(-?\d+),(-?\d+)>, r=(\d+)$`)
)

type pos = [dims]int

func main() {
	// fmt.Println(findPointsInBest(parse(sample)))

	s2 := parse(sample2)
	fmt.Println("collapse X")
	collapseRegionsLinear(0, s2)
	fmt.Println("collapse Y")
	collapseRegionsLinear(1, s2)
	fmt.Println("collapse Z")
	collapseRegionsLinear(2, s2)

	fmt.Println("collapse X/Y")
	collapseRegionsPlanar(0, 1, s2)
	collapseRegionsPlanar(1, 2, s2)
	collapseRegionsPlanar(2, 0, s2)

	collapseRegionsCubic(s2)

	//closest = findClosest(parse(data))
	//fmt.Println(closest, dist(pos{}, closest))
}

func parse(input string) []bot {
	var bots []bot

	lines := strings.Split(input, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if !re.MatchString(line) {
			panic(line)
		}
		m := re.FindStringSubmatch(line)
		bots = append(bots, bot{
			p: pos{
				mustInt(m[1]),
				mustInt(m[2]),
				mustInt(m[3]),
			},
			r: mustInt(m[4]),
		})
	}
	return bots
}

func toRegions(bots []bot) []region {
	var ret []region
	for i, bot := range bots {
		r := bot.toRegion()
		r.ids = []int{i}
		ret = append(ret, r)
	}
	return ret
}

func mustInt(s string) int {
	if v, err := strconv.Atoi(s); err != nil {
		panic(fmt.Sprint(s, err))
	} else {
		return v
	}
}

func abs(i int) int {
	if i < 0 {
		return -i
	}
	return i
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
