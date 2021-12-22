package main

import (
	"fmt"
	"sort"
	"strings"
)

var sample = `
NNCB

CH -> B
HH -> N
CB -> H
NH -> C
HB -> C
HC -> B
HN -> C
NN -> C
BH -> H
NC -> B
NB -> B
BN -> B
BB -> N
BC -> B
CC -> N
CN -> C
`

var input = sample

const iters = 10

type pair struct {
	a, b byte
}

type score struct {
	c     byte
	score int64
}

func main1() {
	pairs := map[pair]byte{}
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

	fmt.Println(string(buf))

	for step := 1; step <= iters; step++ {
		newBuf := make([]byte, 0, len(buf)*2)
		for i := 0; i < len(buf)-1; i++ {
			p := pair{
				a: buf[i],
				b: buf[i+1],
			}
			newBuf = append(newBuf, p.a)
			if c, ok := pairs[p]; ok {
				newBuf = append(newBuf, c)
			} else {
				fmt.Println("no pair for: ", string(buf[i:i+2]))
			}
		}
		newBuf = append(newBuf, buf[len(buf)-1])
		buf = newBuf
		if step < 6 {
			fmt.Println(step, len(buf), string(buf))
		} else {
			fmt.Println(step, len(buf), string(buf))
		}
	}

	counts := map[byte]int64{}
	for _, v := range buf {
		counts[v]++
	}

	sorted := make([]score, 0, len(counts))
	for k, v := range counts {
		sorted = append(sorted, score{
			c:     k,
			score: v,
		})
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
