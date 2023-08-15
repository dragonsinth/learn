package main

import "testing"

func TestVerticesOcta(t *testing.T) {
	baseExpect := [6][3]int{
		{-4, 0, 0}, {4, 0, 0}, {0, -4, 0}, {0, 4, 0}, {0, 0, -4}, {0, 0, 4},
	}

	tcs := [][3]int{
		{0, 0, 0},
		{5, 0, 0},
		{0, 5, 0},
		{0, 0, 5},
		{3, 4, 5},
		{-3, 4, 5},
	}

	for _, tc := range tcs {
		octa := toOctas([]bot{{p: tc, rad: 1}})[0]

		expect := baseExpect
		for i := range expect {
			for d := range expect[i] {
				expect[i][d] += tc[d] * 4
			}
		}

		vertices := octa.vertices()
		t.Log(vertices)
		if vertices != expect {
			t.Errorf("expect=%+v, actual=%+v", expect, vertices)
		}
	}
}
