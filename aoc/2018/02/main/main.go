package main

import (
	"fmt"
	"sort"
	"strings"
)

const (
	sample1 = `
abcdef
bababc
abbcde
abcccd
aabcdd
abcdee
ababab
`
	sample2 = `
abcde
fghij
klmno
pqrst
fguij
axcye
wvxyz
`
)

func main() {
	var words []string
	lines := strings.Split(sample2, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		words = append(words, line)
	}

	sort.Strings(words)
	for i := range words {
		if i == 0 {
			continue
		}
		a := words[i-1]
		b := words[i]
		if diff(a, b) == 1 {
			fmt.Println(a, b, same(a, b))
			return
		}
	}
}

func diff(a string, b string) int {
	ret := 0
	for i := range a {
		if a[i] != b[i] {
			ret++
		}
	}
	return ret
}

func same(a string, b string) string {
	var ret string
	for i := range a {
		if a[i] == b[i] {
			ret = ret + string(a[i])
		}
	}
	return ret
}

func main1() {
	exact := make([]int, 10)
	lines := strings.Split(sample1, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		letterCounts := make([]int, 26)
		for _, c := range line {
			letterCounts[c-'a']++
		}
		var hasTwo, hasThree bool
		for _, n := range letterCounts {
			hasTwo = hasTwo || n == 2
			hasThree = hasThree || n == 3
		}
		if hasTwo {
			exact[2]++
		}
		if hasThree {
			exact[3]++
		}
	}
	fmt.Println(exact[2], exact[3], exact[2]*exact[3])
}
