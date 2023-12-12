package main

import (
	"testing"
)

func TestState(t *testing.T) {
	in := `?###????????`
	expect := []int{3, 2, 1}
	st := newState(in, expect)
	t.Log(st.Count(map[state]int{}))

	count := 0
	work := []state{st}
	for i := 0; i < len(work); i++ {
		cur := work[i]
		if cur.terminal() {
			count++
			continue
		}
		valid := cur.nextValid()
		t.Log(cur, valid)
		work = append(work, valid...)
	}
	t.Log(count)
}
