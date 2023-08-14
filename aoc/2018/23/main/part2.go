package main

import (
	"golang.org/x/exp/slices"
)

type region struct {
	min, max pos
	ids      []int
}

type key struct {
	min, max pos
}

func (r region) key() key {
	return key{
		min: r.min,
		max: r.max,
	}
}

func (r region) inRange(p pos) bool {
	for d := range (pos{}) {
		if p[d] < r.min[d] || p[d] > r.max[d] {
			return false
		}
	}
	return true
}

func (r region) volume() int {
	vol := 1
	for d := range (pos{}) {
		vol *= r.max[d] - r.min[d] + 1
	}
	return vol
}

func intersectAll(regions []region) []region {
	var ret []region
	for i := 0; i < len(regions); i++ {
		for j := i + 1; j < len(regions); j++ {
			if r, ok := intersectRegions(regions[i], regions[j]); ok {
				ret = append(ret, r)
			}
		}
	}
	return ret
}

func intersectRegions(r1 region, r2 region) (region, bool) {
	// compute the intersection
	var ret region
	for i := range (pos{}) {
		ret.min[i] = max(r1.min[i], r2.min[i])
		ret.max[i] = min(r1.max[i], r2.max[i])
		if ret.max[i] < ret.min[i] {
			return region{}, false
		}
	}
	ret.ids = zipperMerge(r1.ids, r2.ids)
	return ret, true
}

func zipperMerge(a []int, b []int) []int {
	r := make([]int, len(a)+len(b))
	ai, bi, ri := 0, 0, 0
	for ai < len(a) && bi < len(b) {
		av := a[ai]
		bv := b[bi]
		switch {
		case av < bv:
			r[ri] = av
			ri++
			ai++
		case bv < av:
			r[ri] = bv
			ri++
			bi++
		case av == bv:
			r[ri] = av
			ri++
			ai++
			bi++
		default:
			panic(av)
		}
	}
	for ai < len(a) {
		r[ri] = a[ai]
		ri++
		ai++
	}
	for bi < len(b) {
		r[ri] = b[bi]
		ri++
		bi++
	}
	return r[:ri]
}

func dedup(regions []region) []region {
	slices.SortFunc(regions, func(a, b region) bool {
		for i := range (pos{}) {
			if a.min[i] != b.min[i] {
				return a.min[i] < b.min[i]
			}
			if a.max[i] != b.max[i] {
				return a.max[i] < b.max[i]
			}
		}
		return false
	})

	wIdx := 1
	last := &regions[0]
	for i := 1; i < len(regions); i++ {
		if regions[i].key() == last.key() {
			last.ids = zipperMerge(last.ids, regions[i].ids)
		} else {
			regions[wIdx] = regions[i]
			last = &regions[wIdx]
			wIdx++
		}
	}
	return regions[:wIdx]
}

func volume(regions []region) int {
	sum := 0
	for _, r := range regions {
		sum += r.volume()
	}
	return sum
}
