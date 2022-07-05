package main

import (
	"fmt"
	"strconv"
)

func main() {
	fmt.Println(111111, check(111111))
	fmt.Println(223450, check(223450))
	fmt.Println(123789, check(123789))

	fmt.Println(112233, check(112233))
	fmt.Println(123444, check(123444))
	fmt.Println(111122, check(111122))

	sum := 0
	for i := 183564; i <= 657474; i++ {
		if check(i) {
			sum++
		}
	}
	fmt.Println(sum)
}

func check(n int) bool {
	v := strconv.Itoa(n)
	hasDouble := false
	tally := 1
	for i := 1; i < 6; i++ {
		if v[i] < v[i-1] {
			return false
		}
		if v[i] == v[i-1] {
			tally++
		} else {
			hasDouble = hasDouble || tally == 2
			tally = 1
		}
	}
	return hasDouble || tally == 2
}
