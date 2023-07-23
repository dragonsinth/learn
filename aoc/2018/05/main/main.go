package main

import (
	"fmt"
)

const sample = `dabAcCaCBAcCcaDA`

const diff = 'a' - 'A'

func main() {
	s := reduce(sample)
	fmt.Println(len(s), s)

	best := s
	for r := byte('A'); r <= byte('Z'); r++ {
		f := filter(sample, r)
		f = reduce(f)
		if len(f) < len(best) {
			best = f
		}
	}
	fmt.Println(len(best), best)
}

func filter(s string, r byte) string {
	r2 := r + diff
	buf := []byte(s)
	wIdx := 0
	for _, c := range buf {
		if c != r && c != r2 {
			buf[wIdx] = c
			wIdx++
		}
	}
	return string(buf[:wIdx])
}

func reduce(in string) string {
	buf := []byte(in)
	for {
		changed := false
		for i := 0; i < len(buf)-1; i++ {
			a, b := buf[i], buf[i+1]
			if df(a, b) == diff {
				buf = append(buf[:i], buf[i+2:]...)
				changed = true
			}
		}
		if !changed {
			break
		}
	}
	return string(buf)
}

func df(a byte, b byte) int {
	u := int(a) - int(b)
	if u < 0 {
		return -u
	}
	return u
}
