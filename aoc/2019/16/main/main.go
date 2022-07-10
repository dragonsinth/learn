package main

import (
	"fmt"
	"strconv"
)

var (
	samples1 = []string{
		"12345678",
		"80871224585914546619083218645595",
		"19617804207202209144916044189917",
		"69317163492948606335995924319873",
	}

	samples2 = []string{
		"03036732577212944063491565474664",
		"02935109699940807407585447034323",
		"03081770884921959731165446850517",
	}
)

var (
	pattern = []int{0, 1, 0, -1}
)

func main() {
	for _, input := range samples1 {
		var vals []byte
		for _, c := range input {
			vals = append(vals, byte(c-'0'))
		}

		next := make([]byte, len(vals))
		for i := 0; i < 100; i++ {
			computeNext(vals, next)
			vals, next = next, vals
		}
		fmt.Println(toString(vals[:8]))
	}
	fmt.Println()

	for _, input := range samples2 {
		var vals []byte
		for _, c := range input {
			vals = append(vals, byte(c-'0'))
		}

		pos := mustInt(input[:7])
		targetSize := len(vals) * 10000
		targetLen := targetSize - pos
		fmt.Println(pos, targetSize, targetLen)

		{
			small := make([]byte, 0, targetLen)
			small = append(small, vals[pos%len(vals):]...)
			for len(small) < targetLen {
				small = append(small, vals...)
			}
			vals = small
		}

		next := make([]byte, len(vals))
		for i := 0; i < 100; i++ {
			computeFast(vals, next)
			vals, next = next, vals
		}
		fmt.Println(toString(vals[:8]))
	}
}

func toString(vals []byte) string {
	out := make([]byte, len(vals))
	for i, v := range vals {
		out[i] = v + '0'
	}
	return string(out)
}

func computeFast(vals []byte, next []byte) {
	// Compute from tail -> head using a running total.
	sum := 0
	for i := len(vals) - 1; i >= 0; i-- {
		sum += int(vals[i])
		next[i] = byte(sum % 10)
	}
}

func computeNext(vals []byte, next []byte) {
	for i := 0; i < len(vals); i++ {
		next[i] = compute(vals, i+1)
	}
}

func compute(vals []byte, reps int) byte {
	patternIdx := 0
	repsLeft := reps - 1
	if reps == 1 {
		patternIdx = 1
		repsLeft = 1
	}

	sum := 0
	idx, maxLen := 0, len(vals)
	for idx < maxLen {
		switch pattern[patternIdx] {
		case -1:
			for repsLeft > 0 && idx < maxLen {
				sum -= int(vals[idx])
				idx++
				repsLeft--
			}
			if sum < -10000000 {
				sum = sum % 100000
				fmt.Println(sum)
			}
		case 0:
			idx += repsLeft
		case 1:
			for repsLeft > 0 && idx < maxLen {
				sum += int(vals[idx])
				idx++
				repsLeft--
			}
			if sum > 10000000 {
				sum = sum % 100000
				fmt.Println(sum)
			}
		default:
			panic("here")
		}

		patternIdx = (patternIdx + 1) % len(pattern)
		repsLeft = reps
	}

	// Last digit only
	if sum < 0 {
		sum = -sum
	}
	return byte(sum % 10)
}

func mustInt(s string) int {
	if v, err := strconv.Atoi(s); err != nil {
		panic(fmt.Sprint(s, err))
	} else {
		return v
	}
}
