package main

import "testing"

func TestIntersection(t *testing.T) {
	for i := 0; i < 101*103; i++ {
		if (i%101 == wRepeat) && (i%103 == hRepeat) {
			t.Log(i)
			if i != p2solution {
				t.Error("wrong")
			}
			return
		}
	}
	t.Errorf("not found")
}
