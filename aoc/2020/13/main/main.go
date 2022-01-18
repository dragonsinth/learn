package main

import (
	"fmt"
	"strings"
)

var samples = []string{
	`7,13,x,x,59,x,31,19`,
	`17,x,13,19`,
	`67,7,59,61`,
	`67,x,7,59,61`,
	`67,7,x,59,61`,
	`1789,37,47,1889`,
}

type bus struct {
	offset int64
	when   int64
}

func main() {
	for _, input := range samples {
		var busses []bus
		parts := strings.Split(input, ",")
		for offset, p := range parts {
			if p == "x" {
				continue
			}
			busses = append(busses, bus{
				offset: int64(offset),
				when:   int64(mustInt(p)),
			})
		}
		fmt.Println(busses)

		for len(busses) > 1 {
			first := busses[0]
			second := busses[1]
			comb := combine(first, second)
			busses = append([]bus{comb}, busses[2:]...)
		}

		fmt.Println(busses[0].offset)
	}

	//count := 0
	//for i := 0; count < 3; i += 7 {
	//	if i%13 == 12 {
	//		fmt.Println(i)
	//		count++
	//	}
	//}

	// At 77, bus 7 leaves, and bus 13 leaves in 1 minute
	// At 77 + 91(n)
	// Replace with bus{77, 91}
}

func combine(first bus, second bus) bus {
	var target int64
	if second.offset >= second.when {
		target = second.when - (second.offset % second.when)
	} else {
		target = second.when - second.offset
	}

	for i := first.offset; true; i += first.when {
		if i%second.when == target {
			return bus{
				offset: i,
				when:   first.when * second.when,
			}
		}
	}
	panic("should not get here")
}
