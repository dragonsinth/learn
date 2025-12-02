package main

import (
	"fmt"
	"strconv"
	"strings"
)

var samples = []string{`
0,0,0,0
 3,0,0,0
 0,3,0,0
 0,0,3,0
 0,0,0,3
 0,0,0,6
 9,0,0,0
12,0,0,0
`, `
-1,2,2,0
0,0,2,-2
0,0,0,-2
-1,2,0,0
-2,-2,-2,2
3,0,2,-1
-1,3,2,2
-1,0,-1,0
0,2,1,-2
3,0,0,0
`, `
1,-1,0,1
2,0,-1,0
3,2,-1,0
0,0,3,1
0,0,-1,-1
2,3,-2,0
-2,2,0,0
2,-2,0,-1
1,-1,0,-1
3,2,0,2
`, `
1,-1,-1,-2
-2,-2,0,1
0,2,1,3
-2,3,-2,1
0,2,3,-2
-1,-1,1,-2
0,-2,-1,0
-2,2,3,-1
1,2,2,0
-1,-2,0,-2
`,
}

func main() {
	for _, sample := range samples {
		run(sample)
	}
}

func run(input string) {
	pts := parse(input)

	// build a graph of "close enough" points
	edges := map[int][]int{}

	// O(N^2) measure all distances, plot all edges
	for i := range pts {
		for j := range pts {
			if i < j {
				if dist(pts[i], pts[j]) <= 3 {
					edges[i] = append(edges[i], j)
					edges[j] = append(edges[j], i)
				}
			}
		}
	}

	assignments := make([]int, len(pts))
	constellation := 0
	for {
		// find a start node
		start := -1
		for i, v := range assignments {
			if v == 0 {
				start = i
				break
			}
		}

		if start == -1 {
			// no more nodes, we're done
			fmt.Println(constellation)
			return
		}

		constellation++

		// form a complete graph based on the edges
		walk(start, edges, assignments, constellation)
	}
}

func walk(start int, edges map[int][]int, assignments []int, constellation int) {
	if assignments[start] == constellation {
		return // already seen
	}
	if assignments[start] != 0 {
		panic("here") // reached an unreachable?
	}
	assignments[start] = constellation

	for _, e := range edges[start] {
		walk(e, edges, assignments, constellation)
	}
}

func dist(a point, b point) int {
	return abs(a[0]-b[0]) + abs(a[1]-b[1]) + abs(a[2]-b[2]) + abs(a[3]-b[3])
}

type point [4]int

func parse(input string) []point {
	var ret []point

	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		ret = append(ret, parseNums(line))
	}

	return ret
}

func parseNums(s string) [4]int {
	var ret []int
	for _, p := range strings.Split(s, ",") {
		ret = append(ret, mustInt(p))
	}
	if len(ret) != 4 {
		panic(ret)
	}
	return [4]int(ret)
}

func mustInt(s string) int {
	if v, err := strconv.Atoi(s); err != nil {
		panic(fmt.Sprint(s, err))
	} else {
		return v
	}
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
