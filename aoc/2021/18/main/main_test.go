package main

import (
	"testing"
)

func TestExplode(t *testing.T) {
	tcs := []struct {
		input  string
		expect string
	}{
		{`[[[[[9,8],1],2],3],4]`, `[[[[0,9],2],3],4]`},
		{`[7,[6,[5,[4,[3,2]]]]]`, `[7,[6,[5,[7,0]]]]`},
		{`[[6,[5,[4,[3,2]]]],1]`, `[[6,[5,[7,0]]],3]`},
		{`[[3,[2,[1,[7,3]]]],[6,[5,[4,[3,2]]]]]`, `[[3,[2,[8,0]]],[9,[5,[4,[3,2]]]]]`},
		{`[[3,[2,[8,0]]],[9,[5,[4,[3,2]]]]]`, `[[3,[2,[8,0]]],[9,[5,[7,0]]]]`},
	}

	for i, tc := range tcs {
		n := parse(tc.input)
		n = explode(n)
		actual := n.String()
		if tc.expect != actual {
			t.Errorf("case %d: want=%q, got=%q", i, tc.expect, actual)
		}
	}
}

func TestSplit(t *testing.T) {
	tcs := []struct {
		input  *Number
		expect string
	}{
		{literal(15), `[7,8]`},
		{literal(14), `[7,7]`},
		{pair(literal(0), literal(17)), `[0,[8,9]]`},
	}

	for i, tc := range tcs {
		n := split(tc.input)
		actual := n.String()
		if tc.expect != actual {
			t.Errorf("case %d: want=%q, got=%q", i, tc.expect, actual)
		}
	}
}

func TestReduce(t *testing.T) {
	a := parse(`[[[[4,3],4],4],[7,[[8,4],9]]]`)
	b := parse(`[1,1]`)
	n := pair(a, b)
	n = reduce(n)
	const want = `[[[[0,7],4],[[7,8],[6,0]]],[8,1]]`
	if got := n.String(); want != got {
		t.Errorf("want=%q, got=%q", want, got)
	}

}
