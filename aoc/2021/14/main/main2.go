package main

import (
	"fmt"
	"sort"
	"strings"
)

type cacheKey struct {
	a, b  byte
	depth int
}

// map of letter -> score
type scoreCard [26]int64

var (
	cache = map[cacheKey]scoreCard{}
	pairs = map[pair]byte{}
)

func main() {
	var buf []byte
	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if strings.Contains(line, " -> ") {
			parts := strings.Split(line, " -> ")
			pairs[pair{
				a: parts[0][0],
				b: parts[0][1],
			}] = parts[1][0]
		} else {
			buf = []byte(line)
		}
	}

	sc := computeScores(buf, iters)

	sorted := make([]score, 0, len(sc))
	for k, v := range sc {
		if v > 0 {
			sorted = append(sorted, score{
				c:     byte(k) + 'A',
				score: v,
			})
		}
	}
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].score < sorted[j].score
	})
	min, max := sorted[0].score, sorted[len(sorted)-1].score
	fmt.Println(max, " - ", min, " = ", max-min)
	for _, sc := range sorted {
		fmt.Println(string(sc.c), ": ", sc.score)
	}
}

func computeScores(buf []byte, depth int) scoreCard {
	ret := scoreCard{}
	for i := 0; i < len(buf)-1; i++ {
		a, b := buf[i], buf[i+1]
		r := computeScore(a, b, depth)
		for i := range ret {
			ret[i] += r[i]
		}
		ret[b-'A']-- // don't count the right edge twice
	}
	ret[buf[len(buf)-1]-'A']++ // count the final right edge
	return ret
}

func computeScore(a, b byte, depth int) scoreCard {
	key := cacheKey{
		a:     a,
		b:     b,
		depth: depth,
	}
	if r, ok := cache[key]; ok {
		return r
	}

	ret := scoreCard{}
	if depth == 0 {
		ret[a-'A']++
		ret[b-'A']++
	} else {
		// Sum each subtree.
		p := pair{a: a, b: b}
		c, ok := pairs[p]
		if !ok {
			panic(p)
		}

		r1 := computeScore(a, c, depth-1)
		r2 := computeScore(c, b, depth-1)
		for i := range ret {
			ret[i] = r1[i] + r2[i]
		}
		// The middle value is added in twice so subtract it out.
		ret[c-'A']--
	}

	cache[key] = ret
	return ret
}
