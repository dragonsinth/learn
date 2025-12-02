package main

import (
	"fmt"
	"strconv"
	"strings"
)

const sample = `11-22,95-115,998-1012,1188511880-1188511890,222220-222224,1698522-1698528,446443-446449,38593856-38593862,565653-565659,824824821-824824827,2121212118-2121212124`

type rng struct {
	lo, hi int
}

func main() {
	rs := parse(sample)
	run(rs, true, invalid1)
	run(rs, true, invalid2)
}

func parse(input string) []rng {
	var ret []rng
	for _, line := range strings.Split(input, ",") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Split(line, "-")
		if len(parts) != 2 {
			panic(parts)
		}
		ret = append(ret, rng{
			lo: mustInt(parts[0]),
			hi: mustInt(parts[1]),
		})
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

func run(rs []rng, debug bool, invalid func(i int) bool) {
	sum := 0
	for _, r := range rs {
		if debug {
			fmt.Println(r.lo, r.hi)
		}
		for i := r.lo; i <= r.hi; i++ {
			if invalid(i) {
				if debug {
					fmt.Println("invalid index:", i)
				}
				sum += i
			}
		}
	}
	fmt.Println("sum:", sum)
}

func invalid1(i int) bool {
	str := strconv.Itoa(i)
	if len(str)%2 != 0 {
		return false
	}
	return str[0:len(str)/2] == str[len(str)/2:]
}

func invalid2(i int) bool {
	str := strconv.Itoa(i)
	for i := 1; i <= len(str)/2; i++ {
		if isInvalid(str, i) {
			return true
		}
	}
	return false
}

func isInvalid(str string, sz int) bool {
	if len(str)%sz != 0 {
		return false
	}

	for i := sz; i < len(str); i += sz {
		// compare str[0] to str[i]
		if str[0:sz] != str[i:i+sz] {
			return false // found a mismatch
		}
	}

	return true
}
