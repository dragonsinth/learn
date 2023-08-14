package main

import (
	"bytes"
	"fmt"
	"golang.org/x/exp/slices"
)

func collapseRegionsLinear(dx int, bots []bot) region {
	iter := 1
	regions := toRegions(bots)

	for i := range regions {
		r := &regions[i]
		for d := range (pos{}) {
			if d != dx {
				r.min[d] = 0
				r.max[d] = 0
			}
		}
	}
	orig := regions

	fmt.Printf("iter=%d, len=%d, vol=%d\n", iter, len(regions), volume(regions))
	printLinear(dx, regions)
	for len(regions) > 1 {
		regions = intersectAll(regions)
		regions = dedup(regions)
		iter++
		fmt.Printf("iter=%d, len=%d, vol=%d\n", iter, len(regions), volume(regions))
		printLinear(dx, regions)
	}

	// validate!
	final := regions[0]
	expectInRange := map[int]bool{}
	for _, id := range final.ids {
		expectInRange[id] = true
	}

	for x := final.min[dx]; x <= final.max[dx]; x++ {
		for id, r := range orig {
			p := r.min
			p[dx] = x
			if expectInRange[id] != r.inRange(p) {
				panic(fmt.Sprintf("expectInRange(%v) != r.inRange(%v), id=%d", expectInRange[id], r.inRange(p), id))
			}
		}
	}

	return final
}

func printLinear(dx int, regions []region) {
	slices.SortFunc(regions, func(a, b region) bool {
		return a.min[dx] < b.min[dx]
	})
	for _, r := range regions {
		var buf bytes.Buffer
		if r.min[dx] >= 0 && r.max[dx] <= 80 {
			buf.Write(bytes.Repeat([]byte{'.'}, r.min[dx]))
			buf.Write(bytes.Repeat([]byte{'*'}, r.max[dx]-r.min[dx]+1))
			buf.Write(bytes.Repeat([]byte{'.'}, 80-r.max[dx]-1))
		} else {
			buf.Write(bytes.Repeat([]byte{'?'}, 80))
		}
		buf.WriteString(fmt.Sprintf(" [%d->%d]", r.min[dx], r.max[dx]))
		buf.WriteString(fmt.Sprintf(" %+v", r.ids))
		fmt.Println(buf.String())
	}
}
