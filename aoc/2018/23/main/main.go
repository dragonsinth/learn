package main

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
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

type dim int

const (
	X = dim(iota)
	Y
	Z
)

var (
	re = regexp.MustCompile(`^pos=<(-?\d+),(-?\d+),(-?\d+)>, r=(\d+)$`)
)

type pos [3]int

func (p pos) rad() int {
	return p[X] + p[Y] + p[Z]
}

func (p pos) String() string {
	return fmt.Sprintf("pos{%d, %d, %d}", p[X], p[Y], p[Z])
}

func main() {
	fmt.Println(findPointsInBest(parse(sample)))

	rand.Seed(time.Now().UnixMicro())
	bots := parse(sample2)
	search := pos{10, 10, 10}
	bestScore := calcScore(search, bots)
	fmt.Println(bestScore, search.rad(), search)

	for {
		score, match := randomMatch(search, pos{0, 2000, 2000}, 100000, bots)
		fmt.Println(score, match.rad(), match)
		if score > bestScore || (score == bestScore && match.rad() < search.rad()) {
			search = match
			bestScore = score
		}

		score, match = exhaustiveMatch(search, 25, bots)
		fmt.Println(score, match.rad(), match)
		if score > bestScore || (score == bestScore && match.rad() < search.rad()) {
			search = match
			bestScore = score
		}
	}
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
			rad: mustInt(m[4]),
		}.normalize())
	}
	return bots
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
