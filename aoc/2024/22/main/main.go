package main

import (
	"fmt"
	"strconv"
	"strings"
)

const sample1 = `
1
10
100
2024
`

const sample2 = `
1
2
3
2024
`

func main() {
	fmt.Println(calc(123, 10, true))

	part1(sample1, true)

	bigMap := map[sol]int{}
	addSequences(bigMap, 123, 9, true)
	fmt.Println(bigMap)

	part2(sample2, false)
}

func part1(input string, debug bool) {
	sum := 0
	for _, val := range parse(input) {
		res := calc(val, 2000, false)
		if debug {
			fmt.Println(val, res)
		}
		sum += res
	}
	fmt.Println(sum)
}

type sol [4]int8

func part2(input string, debug bool) {
	bigMap := map[sol]int{}

	for _, val := range parse(input) {
		addSequences(bigMap, val, 2000, debug)
	}
	var bestSol sol
	var bestScr int
	for kSol, vScr := range bigMap {
		if vScr > bestScr {
			bestSol, bestScr = kSol, vScr
		}
	}
	fmt.Println(bestSol, bestScr)
}

func addSequences(bigMap map[sol]int, secret int, steps int, debug bool) {
	seen := make(map[sol]bool, steps)

	last := int8(secret % 10)
	if debug {
		fmt.Println(secret, last)
	}

	var rolling sol
	for i := 0; i < steps; i++ {
		secret = nextSecret(secret)
		next := int8(secret % 10)
		rolling[0], rolling[1], rolling[2] = rolling[1], rolling[2], rolling[3]
		rolling[3] = next - last
		if debug {
			fmt.Println(secret, next, int(rolling[3]))
		}
		if i > 2 {
			if !seen[rolling] {
				seen[rolling] = true
				bigMap[rolling] += int(next)
			}
		}
		last = next
	}
}

func calc(secret int, steps int, debug bool) int {
	for i := 0; i < steps; i++ {
		secret = nextSecret(secret)
		if debug {
			fmt.Println(i, secret)
		}
	}
	return secret
}

func nextSecret(secret int) int {
	secret = prune(mix(secret*64, secret))
	secret = prune(mix(secret/32, secret))
	secret = prune(mix(secret*2048, secret))
	return secret
}

func mix(result int, secret int) int {
	return result ^ secret
}

func prune(v int) int {
	return v % 16777216
}

func parse(input string) []int {
	var ret []int
	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		ret = append(ret, mustInt(line))
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
