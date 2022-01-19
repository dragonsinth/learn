package main

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

var (
	memLine = regexp.MustCompile(`^mem\[(\d+)\] = (\d+)$`)
)

func main() {
	const sample = `
mask = 000000000000000000000000000000X1001X
mem[42] = 100
mask = 00000000000000000000000000000000X0XX
mem[26] = 1
`

	mem := map[uint64]uint64{}

	floatMask := uint64(0)
	orMask := uint64(0)

	for _, line := range strings.Split(sample, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "mask = ") {
			maskStr := strings.TrimPrefix(line, "mask = ")
			if len(maskStr) != 36 {
				panic(maskStr)
			}

			floatMask = uint64(0)
			orMask = uint64(0)
			for _, c := range maskStr {
				floatMask <<= 1
				orMask <<= 1

				switch c {
				case '0':
					// no-op both
					floatMask |= 0
					orMask |= 0
				case '1':
					// force to 1
					floatMask |= 0
					orMask |= 1
				case 'X':
					// float mask
					floatMask |= 1
					orMask |= 0
				}
			}
		} else if memLine.MatchString(line) {
			matches := memLine.FindStringSubmatch(line)
			if len(matches) != 3 {
				panic(matches)
			}

			addr := mustUInt(matches[1])
			val := mustUInt(matches[2])
			apply(mem, addr, val, orMask, math.MaxUint64, floatMask)

		} else {
			panic(line)
		}
	}

	var sum uint64
	for _, v := range mem {
		sum += v
	}
	fmt.Println(sum)

}

func apply(mem map[uint64]uint64, addr uint64, val, orMask, andMask, floatMask uint64) {
	if floatMask == 0 {
		// End recursion
		addr = (addr & andMask) | orMask
		fmt.Println(addr, "=", val)
		mem[addr] = val
		return
	}

	// Find the first non-zero bit in floatMask.
	for i, c := uint64(1), uint64(1<<36); i < c; i <<= 1 {
		if floatMask&i != 0 {
			// Both combinations: force a 0 and force a 1
			apply(mem, addr, val, orMask, andMask^i, floatMask^i)
			apply(mem, addr, val, orMask|i, andMask, floatMask^i)
			return
		}
	}
	panic("should not get here")
}

func main1() {
	const sample = `
mask = XXXXXXXXXXXXXXXXXXXXXXXXXXXXX1XXXX0X
mem[8] = 11
mem[7] = 101
mem[8] = 0
`

	mem := map[int]uint64{}

	andMask := uint64(math.MaxUint64)
	orMask := uint64(0)

	for _, line := range strings.Split(sample, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "mask = ") {
			maskStr := strings.TrimPrefix(line, "mask = ")
			if len(maskStr) != 36 {
				panic(maskStr)
			}

			andMask = uint64(math.MaxUint64)
			orMask = uint64(0)
			for _, c := range maskStr {
				andMask <<= 1
				orMask <<= 1

				switch c {
				case '0':
					// force to 0
					andMask |= 0
					orMask |= 0
				case '1':
					// force to 1
					andMask |= 1
					orMask |= 1
				case 'X':
					// do nothing
					andMask |= 1
					orMask |= 0
				}
			}
		} else if memLine.MatchString(line) {
			matches := memLine.FindStringSubmatch(line)
			if len(matches) != 3 {
				panic(matches)
			}

			addr := mustInt(matches[1])
			val, err := strconv.ParseUint(matches[2], 10, 64)
			if err != nil {
				panic(err)
			}

			maskedVal := (val & andMask) | orMask
			mem[addr] = maskedVal
		} else {
			panic(line)
		}
	}

	var sum uint64
	for _, v := range mem {
		sum += v
	}
	fmt.Println(sum)
}

func toBin(n uint64) string {
	return strconv.FormatUint(n, 2)
}

func mustInt(s string) int {
	if v, err := strconv.Atoi(s); err != nil {
		panic(fmt.Sprint(s, err))
	} else {
		return v
	}
}

func mustUInt(s string) uint64 {
	if v, err := strconv.ParseUint(s, 10, 64); err != nil {
		panic(fmt.Sprint(s, err))
	} else {
		return v
	}
}
