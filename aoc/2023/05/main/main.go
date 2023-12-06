package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/dragonsinth/learn/aoc/sliceheap"
)

var sample = `
seeds: 79 14 55 13

seed-to-soil map:
50 98 2
52 50 48

soil-to-fertilizer map:
0 15 37
37 52 2
39 0 15

fertilizer-to-water map:
49 53 8
0 11 42
42 0 7
57 7 4

water-to-light map:
88 18 7
18 25 70

light-to-temperature map:
45 77 23
81 45 19
68 64 13

temperature-to-humidity map:
0 69 1
1 0 69

humidity-to-location map:
60 56 37
56 93 4
`

var rulesetNames = []string{
	"seed-to-soil",
	"soil-to-fertilizer",
	"fertilizer-to-water",
	"water-to-light",
	"light-to-temperature",
	"temperature-to-humidity",
	"humidity-to-location",
}

type ruleset []rule

func (rs ruleset) apply(v int) int {
	for _, r := range rs {
		// 98, 2 => 98, 99
		if r.src <= v && v < r.src+r.len {
			// in range of this rule, do the transform
			return v - r.src + r.dst
		}
	}
	// unchanged
	return v
}

func (rs ruleset) applySegments(segs []segment) []segment {
	var ret []segment
	work := sliceheap.New(segmentSort)
	for _, seg := range segs {
		work.Push(seg)
	}

	match := func(seg segment) bool {
		for _, r := range rs {
			if r.src <= seg.src && seg.end() <= r.end() {
				// rule fully contains segment, just transform it
				ret = append(ret, segment{
					src: seg.src - r.src + r.dst,
					len: seg.len,
				})
				return true
			}

			if seg.isBroken(r.src) {
				// start of rule breaks segment; break _to the left_ of the rule src
				pt := r.src
				a := segment{
					src: seg.src,
					len: pt - seg.src,
				}
				assert(a.end() == pt-1)
				assert(!a.isBroken(pt))
				b := segment{
					src: pt,
					len: seg.len - a.len,
				}
				assert(a.len+b.len == seg.len)
				assert(!b.isBroken(pt))
				assert(b.end() == seg.end())
				work.Push(a)
				work.Push(b)
				return true
			}

			if seg.isBroken(r.end()) {
				// end of rule breaks segment; break _to the right_ of the rule end
				pt := r.end()
				a := segment{
					src: seg.src,
					len: pt - seg.src + 1,
				}
				assert(a.end() == pt)
				assert(!a.isBroken(pt))
				b := segment{
					src: pt + 1,
					len: seg.len - a.len,
				}
				assert(a.len+b.len == seg.len)
				assert(!b.isBroken(pt))
				assert(b.end() == seg.end())
				work.Push(a)
				work.Push(b)
				return true
			}

			// disjoint, continue
		}
		return false
	}

	// break the incoming segment into as many segments as needed to cleanly apply
	for work.Len() > 0 {
		seg := work.Pop()
		if !match(seg) {
			// no rule matched
			ret = append(ret, seg)
		}
	}
	return ret
}

func assert(b bool) {
	if !b {
		panic(b)
	}
}

type segment struct {
	src int
	len int
}

func (s segment) end() int {
	if s.len < 1 {
		panic(s.len)
	}
	return s.src + s.len - 1
}

func (s segment) isBroken(pt int) bool {
	return pt > s.src && pt < s.end()
}

type rule struct {
	dst int
	segment
}

func main() {
	run(sample)
}

func run(input string) {
	var state string
	var seeds []int
	nextRuleset := 0
	var ruleSets []ruleset
	var curRuleSet ruleset
	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			if state != "" {
				ruleSets = append(ruleSets, curRuleSet)
				curRuleSet = nil
				state = ""
				nextRuleset++
			}
			continue
		}

		switch {
		case state == "" && strings.HasPrefix(line, "seeds: "):
			seeds = parseNums(strings.TrimPrefix(line, "seeds: "))
		case state == "" && strings.HasSuffix(line, " map:"):
			state = strings.TrimSuffix(line, " map:")
			if state != rulesetNames[nextRuleset] {
				panic(state)
			}
		default:
			vals := parseNums(line)
			if len(vals) != 3 {
				panic(vals)
			}
			curRuleSet = append(curRuleSet, rule{
				dst: vals[0],
				segment: segment{
					src: vals[1],
					len: vals[2],
				},
			})
		}
	}
	ruleSets = append(ruleSets, curRuleSet)

	// part 1
	bestSeed, bestLoc := 0, math.MaxInt
	for _, s := range seeds {
		v := s
		for _, rs := range ruleSets {
			v = rs.apply(v)
		}
		if v < bestLoc {
			bestSeed, bestLoc = s, v
		}
	}
	fmt.Println(bestSeed, bestLoc)

	// part 2
	var segs []segment
	for i := 0; i < len(seeds); i += 2 {
		segs = append(segs, segment{
			src: seeds[i],
			len: seeds[i+1],
		})
	}

	v := segs
	for _, rs := range ruleSets {
		v = rs.applySegments(v)
	}
	bestSeg := segment{src: math.MaxInt, len: 0}
	for _, seg := range v {
		if seg.src < bestSeg.src {
			bestSeg = seg
		}
	}
	fmt.Println(bestSeg)
}

func parseNums(s string) []int {
	var ret []int
	for _, p := range strings.Fields(s) {
		ret = append(ret, mustInt(p))
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

func segmentSort(a, b segment) bool {
	return a.src < b.src
}
