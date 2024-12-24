package main

import (
	"testing"
)

func TestTrial(t *testing.T) {
	const hi = 3       // 127
	const shiftMax = 1 // 44 - 8

	for shift := 0; shift < shiftMax; shift++ {
		for i := 0; i < hi; i++ {
			for j := 0; j < hi; j++ {
				a := i << shift
				b := j << shift
				sum := runTrial(a, b)
				if sum != a+b {
					t.Errorf("sum(%d,%d)=%d, want %d, shift=%d", a, b, sum, a+b, shift)
					t.FailNow()
				}
			}
		}
	}
}

func runTrial(a int, b int) int {
	x := encode(a)
	y := encode(b)
	z := trial(x, y)
	return decode(z)
}

func encode(n int) [63]bool {
	var ret [63]bool
	for i := 0; i < 63; i++ {
		ret[i] = n&1 == 1
		n = n >> 1
	}
	return ret
}

func decode(z [63]bool) int {
	var ret int
	for i := 62; i >= 0; i-- {
		ret = ret << 1
		if z[i] {
			ret |= 1
		}
	}
	return ret
}
