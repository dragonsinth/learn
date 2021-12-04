package main

import (
	"fmt"
	"strconv"
	"strings"
)

var sample = `
00100
11110
10110
10111
10101
01111
00111
11100
10000
11001
00010
01010
`

const width = 5

func main1() {
	zeros := make([]int, width)
	ones := make([]int, width)
	for _, line := range strings.Split(sample, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if len(line) != width {
			panic(len(line))
		}

		for i := 0; i < width; i++ {
			switch line[i] {
			case '0':
				// nothing
				zeros[i]++
			case '1':
				ones[i]++
			default:
				panic(line)
			}
		}
	}

	gamma, epsilon := 0, 0
	for i := 0; i < width; i++ {
		gamma <<= 1
		epsilon <<= 1
		if ones[i] > zeros[i] {
			gamma++
		} else if ones[i] < zeros[i] {
			epsilon++
		} else {
			panic(i)
		}
	}

	fmt.Println(zeros, ones, gamma, epsilon, gamma*epsilon)
}

func main() {
	var o2set, co2set []string
	for _, line := range strings.Split(sample, "\n") {
		if line == "" {
			continue
		}
		if len(line) != width {
			panic(len(line))
		}
		o2set = append(o2set, line)
		co2set = append(co2set, line)
	}

	// o2set
	for i := 0; i < width; i++ {
		// calc
		zero, ones := 0, 0
		for _, s := range o2set {
			switch s[i] {
			case '0':
				zero++
			case '1':
				ones++
			default:
				panic(s)
			}
		}

		// filter
		search := byte('0')
		if ones >= zero {
			search = byte('1')
		}
		var newSet []string
		for _, s := range o2set {
			if s[i] == search {
				newSet = append(newSet, s)
			}
		}

		o2set = newSet
		if len(o2set) == 1 {
			break
		}
	}

	// co2set
	for i := 0; i < width; i++ {
		// calc
		zero, ones := 0, 0
		for _, s := range co2set {
			switch s[i] {
			case '0':
				zero++
			case '1':
				ones++
			default:
				panic(s)
			}
		}

		// filter
		search := byte('1')
		if ones >= zero {
			search = byte('0')
		}
		var newSet []string
		for _, s := range co2set {
			if s[i] == search {
				newSet = append(newSet, s)
			}
		}

		co2set = newSet
		if len(co2set) == 1 {
			break
		}
	}

	o2str, co2str := o2set[0], co2set[0]
	o2, err := strconv.ParseInt(o2str, 2, 64)
	if err != nil {
		panic(err)
	}
	co2, err := strconv.ParseInt(co2str, 2, 64)
	if err != nil {
		panic(err)
	}
	fmt.Println(o2, co2, o2 * co2)
}
