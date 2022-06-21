package main

func (rs *ruleSet) tryMatch(state int, in string) bool {
	for _, v := range rs.rules[state].resolved {
		if v == in {
			return true
		}
	}
	return false
}

/*
0: 8 11
8: 42
11: 42 31

8: 42 | 42 8
11: 42 31 | 42 11 31

0
8 11
42 42 31

0
8 11
42+ (42 31)+

42(n), 31(m) where n > m
*/

func (rs *ruleSet) tryMatchCustom(in string) bool {
	sz := rs.rules[42].minLen
	if len(in)%sz != 0 {
		return false
	}
	// Strip every 42 from the head, every 31 from tail.
	nHead, nTail := 0, 0
	for len(in) > 0 {
		if rs.tryMatch(42, in[:sz]) {
			in = in[sz:]
			nHead++
		} else if rs.tryMatch(31, in[len(in)-sz:]) {
			in = in[:len(in)-sz]
			nTail++
		} else {
			return false
		}
	}
	return nHead > 1 && nTail > 0 && nHead > nTail
}
