package main

import "testing"

func TestEncodeDecode(t *testing.T) {
	for _, tc := range []struct {
		decimal int
		snafu   string
	}{
		{1, "1"},
		{2, "2"},
		{3, "1="},
		{4, "1-"},
		{5, "10"},
		{6, "11"},
		{7, "12"},
		{8, "2="},
		{9, "2-"},
		{10, "20"},
		{15, "1=0"},
		{20, "1-0"},
		{2022, "1=11-2"},
		{12345, "1-0---0"},
		{314159265, "1121-1110-1=0"},
	} {
		{
			got := decode(tc.snafu)
			if tc.decimal != got {
				t.Errorf("want=%d, got=%d", tc.decimal, got)
			}
		}
		{
			got := encode(tc.decimal)
			if tc.snafu != got {
				t.Errorf("want=%s, got=%s", tc.snafu, got)
			}
		}
	}
}
