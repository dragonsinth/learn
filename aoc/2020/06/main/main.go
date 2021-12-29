package main

import (
	"fmt"
	"strings"
)

var sample = `
abc

a
b
c

ab
ac

a
a
a
a

b
`

var (
	input = sample
)

func main() {
	sumUnion := 0
	sumInter := 0

	union := [26]bool{}
	inter := [26]int{}
	groupCount := 0
	finishGroup := func() {
		uCount := 0
		for _, v := range union {
			if v {
				uCount++
			}
		}
		if uCount == 0 {
			return
		}

		iCount := 0
		for _, v := range inter {
			if v == groupCount {
				iCount++
			}
		}

		fmt.Println(uCount, iCount)
		sumUnion += uCount
		sumInter += iCount
		union = [26]bool{}
		inter = [26]int{}
		groupCount = 0
	}

	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			finishGroup()
			continue
		}

		for _, b := range []byte(line) {
			union[b-'a'] = true
			inter[b-'a']++
		}
		groupCount++
	}
	finishGroup()
	fmt.Println(sumUnion, sumInter)
}
