package main

import (
	"math"
)

func resolveRules(rules []rule) *ruleSet {
	rs := &ruleSet{rules: rules}
	rs.resolve(0)
	return rs
}

type ruleSet struct {
	rules []rule
}

func (rs *ruleSet) resolve(i int) []string {
	r := &rs.rules[i]
	if len(r.resolved) > 0 {
		return r.resolved
	}

	// Must resolve each piece.
	for _, s := range r.raw {
		r.resolved = append(r.resolved, rs.resolveSeq(s)...)
	}
	min, max := math.MaxInt32, 0
	for _, r := range r.resolved {
		if len(r) < min {
			min = len(r)
		}
		if len(r) > max {
			max = len(r)
		}
	}
	r.minLen, r.maxLen = min, max
	return r.resolved
}

func (rs *ruleSet) resolveSeq(s seq) []string {
	// For each piece, collect all possibilities.
	ans := make([][]string, len(s))
	for i, c := range s {
		switch c {
		case A:
			ans[i] = []string{"a"}
		case B:
			ans[i] = []string{"b"}
		default:
			ans[i] = rs.resolve(int(c))
		}
	}

	// Now cross multiply.
	var ret []string
	for i := range ans {
		// Base case
		if i == 0 {
			for _, v := range ans[0] {
				ret = append(ret, v)
			}
		} else {
			// Concat each existing answer with each alternative here.
			iter := ret
			ret = nil
			for _, pre := range iter {
				for _, post := range ans[i] {
					ret = append(ret, pre+post)
				}
			}
		}
	}

	return ret
}
