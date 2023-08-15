package main

import (
	"bytes"
	"fmt"
	"golang.org/x/exp/slices"
)

type rng struct {
	min, max int
}

type segment struct {
	r   rng
	ids []int
}

func (s segment) length() int {
	return s.r.max - s.r.min
}

func (s segment) inRange(p int) bool {
	return p >= s.r.min && p <= s.r.max
}

func (s segment) intersect(o segment) (segment, bool) {
	// compute the intersection
	ret := segment{
		r: rng{
			min: max(s.r.min, o.r.min),
			max: min(s.r.max, o.r.max),
		},
		ids: nil,
	}
	if ret.r.max < ret.r.min {
		return segment{}, false
	}
	ret.ids = zipperMerge(s.ids, o.ids)
	return ret, true
}

func (s segment) key() rng {
	return s.r
}

func (s segment) loc() int {
	return mean(s.r.min, s.r.max)
}

func (s segment) rads() int {
	return rads(s.r.max, s.r.min)
}

func collapseRegionsLinear(dx int, bots []bot) segment {
	iter := 1
	segments := toSegments(dx, bots)
	orig := segments

	fmt.Printf("iter=%d, len=%d, length=%d\n", iter, len(segments), length(segments))
	printLinear(segments)
	for len(segments) > 1 {
		segments = intersectAllSegments(iter, segments)
		segments = dedupSegments(segments)
		iter++
		fmt.Printf("iter=%d, len=%d, length=%d\n", iter, len(segments), length(segments))
		printLinear(segments)
	}

	// validate!
	final := segments[0]
	expectInRange := map[int]bool{}
	for _, id := range final.ids {
		expectInRange[id] = true
	}

	for x := final.r.min; x <= final.r.max; x++ {
		for id, seg := range orig {
			p := seg.r.min
			p = x
			if expectInRange[id] != seg.inRange(p) {
				panic(fmt.Sprintf("expectInRange(%v) != seg.inRange(%v), id=%d", expectInRange[id], seg.inRange(p), id))
			}
		}
	}

	return final
}

func intersectAllSegments(minLen int, segs []segment) []segment {
	var ret []segment
	for i, seg := range segs {
		found := false
		for j := i + 1; j < len(segs); j++ {
			if newSeg, ok := seg.intersect(segs[j]); ok {
				if len(newSeg.ids) >= minLen {
					ret = append(ret, newSeg)
					found = true
				}
			}
		}
		if !found && len(seg.ids) >= minLen {
			ret = append(ret, seg)
		}
	}
	return ret
}

func dedupSegments(segs []segment) []segment {
	slices.SortFunc(segs, func(a, b segment) bool {
		if a.r.min != b.r.min {
			return a.r.min < b.r.min
		}
		if a.r.max != b.r.max {
			return a.r.max < b.r.max
		}
		return false
	})

	wIdx := 1
	last := &segs[0]
	for i := 1; i < len(segs); i++ {
		if segs[i].key() == last.key() {
			last.ids = zipperMerge(last.ids, segs[i].ids)
		} else {
			segs[wIdx] = segs[i]
			last = &segs[wIdx]
			wIdx++
		}
	}
	return segs[:wIdx]
}

func toSegments(dx int, bots []bot) []segment {
	var ret []segment
	for i, b := range bots {
		x := b.p[dx]
		r := b.rad

		// quadruple the coordinate system for precise math
		x *= 4
		r *= 4

		seg := segment{
			r: rng{
				min: x - r,
				max: x + r,
			},
			ids: []int{i},
		}
		if seg.loc() != x {
			panic("wrong loc")
		}

		ret = append(ret, seg)
	}
	return ret
}

func printLinear(segs []segment) {
	const width = 100
	slices.SortFunc(segs, func(a, b segment) bool {
		return a.r.min < b.r.min
	})
	for _, seg := range segs {
		var buf bytes.Buffer
		if seg.r.min >= 0 && seg.r.max < width {
			buf.Write(bytes.Repeat([]byte{'.'}, seg.r.min))
			buf.Write(bytes.Repeat([]byte{'*'}, seg.r.max-seg.r.min+1))
			buf.Write(bytes.Repeat([]byte{'.'}, width-seg.r.max-1))
		} else {
			buf.Write(bytes.Repeat([]byte{'?'}, width))
		}
		buf.WriteString(fmt.Sprintf(" [%d->%d]", seg.r.min, seg.r.max))
		buf.WriteString(fmt.Sprintf(" [loc=%d, rads=%d]", seg.loc(), seg.rads()))
		buf.WriteString(fmt.Sprintf(" %+v", seg.ids))
		fmt.Println(buf.String())
	}
}

func length(segs []segment) int {
	sum := 0
	for _, seg := range segs {
		sum += seg.length()
	}
	return sum
}
