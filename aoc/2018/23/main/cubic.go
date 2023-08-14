package main

import (
	"bytes"
	"fmt"
)

func collapseRegionsCubic(bots []bot) region {
	iter := 1
	regions := toRegions(bots)
	orig := regions

	fmt.Printf("iter=%d, len=%d, vol=%d\n", iter, len(regions), volume(regions))
	for len(regions) > 1 {
		printCubic(regions)
		regions = intersectAll(regions)
		regions = dedup(regions)
		iter++
		fmt.Printf("iter=%d, len=%d, vol=%d\n", iter, len(regions), volume(regions))
		printCubic(regions)
	}

	// validate!
	final := regions[0]
	expectInRange := map[int]bool{}
	for _, id := range final.ids {
		expectInRange[id] = true
	}

	for z := final.min[2]; z <= final.max[2]; z++ {
		for y := final.min[1]; y <= final.max[1]; y++ {
			for x := final.min[0]; x <= final.max[0]; x++ {
				p := pos{x, y, z}
				fmt.Println(p)
				for id, r := range orig {
					fmt.Println(id, r, expectInRange[id], r.inRange(p))
					if expectInRange[id] != r.inRange(p) {
						panic(fmt.Sprintf("expectInRange(%t) != r.inRange(%t), id=%d", expectInRange[id], r.inRange(p), id))
					}
				}
			}
		}
	}

	return final
}

func printCubic(regions []region) {
	for _, r := range regions {
		var buf bytes.Buffer
		buf.WriteString(fmt.Sprintf(" [%+v->%+v]", r.min, r.max))
		buf.WriteString(fmt.Sprintf(" %+v", r.ids))
		fmt.Println(buf.String())
	}
}
