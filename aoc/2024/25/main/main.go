package main

import (
	"fmt"
	"strings"
)

const sample = `
#####
.####
.####
.####
.#.#.
.#...
.....

#####
##.##
.#.##
...##
...#.
...#.
.....

.....
#....
#....
#...#
#.#.#
#.###
#####

.....
.....
#.#..
###..
###.#
###.#
#####

.....
.....
.....
#....
#.#..
#.#.#
#####
`

func main() {
	run(sample, true)
}

func run(input string, debug bool) {
	keys, locks := parse(input)
	if debug {
		for _, l := range locks {
			fmt.Println(toString(l))
		}
		fmt.Println()
		for _, k := range keys {
			fmt.Println(toString(k))
		}
		fmt.Println()
	}

	sum := 0
	for _, k := range keys {
		for _, l := range locks {
			if fit(k, l) {
				sum++
			}
		}
	}
	fmt.Println(sum)
}

func fit(k key, l lock) bool {
	for i := 0; i < 5; i++ {
		if k[i]+l[i] > 5 {
			return false
		}
	}
	return true
}

type key [5]byte

type lock [5]byte

func toString(in [5]byte) string {
	out := []byte("0,0,0,0,0")
	for i := 0; i < 5; i++ {
		out[i*2] = in[i] + '0'
	}
	return string(out)
}

func parse(sample string) ([]key, []lock) {
	var keys []key
	var locks []lock
	var lines []string
	for _, line := range strings.Split(sample, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		lines = append(lines, line)
		if len(lines) == 7 {
			if lines[0] == "#####" {
				// new lock
				if lines[6] != "....." {
					panic(lines)
				}

				var l lock
				for i := 0; i < 5; i++ {
					// from the top, position of last '#' denotes size
					for row, rowLine := range lines {
						if rowLine[i] == '.' {
							l[i] = byte(row) - 1
							break
						}
					}
				}
				locks = append(locks, l)
			} else if lines[6] == "#####" {
				// new key
				if lines[0] != "....." {
					panic(lines)
				}
				var k key
				for i := 0; i < 5; i++ {
					// from the top, position of first '#' denotes size
					for row, rowLine := range lines {
						if rowLine[i] == '#' {
							k[i] = 6 - byte(row)
							break
						}
					}
				}
				keys = append(keys, k)
			} else {
				panic(lines)
			}
			lines = nil
		}
	}
	if lines != nil {
		panic("non-empty")
	}
	return keys, locks
}

//
//func distance(as []int, bs []int) int {
//	sum := 0
//	for i := range as {
//		sum = sum + abs(as[i]-bs[i])
//	}
//	return sum
//}
//
//func similarity(as []int, bs []int) int {
//	ret := 0
//	for _, a := range as {
//		for _, b := range bs {
//			if a == b {
//				ret += a
//			}
//		}
//	}
//	return ret
//}
//
//func mustInt(s string) int {
//	if v, err := strconv.Atoi(s); err != nil {
//		panic(fmt.Sprint(s, err))
//	} else {
//		return v
//	}
//}
//
//func abs(a int) int {
//	if a < 0 {
//		return -a
//	}
//	return a
//}
