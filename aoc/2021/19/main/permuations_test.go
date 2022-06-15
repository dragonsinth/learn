package main

import "testing"

func TestPermuations(t *testing.T) {
	const input = `
		--- scanner 0 ---
		-1,-1,1
		-2,-2,2
		-3,-3,3
		-2,-3,1
		5,6,-4
		8,0,7
	`

	scanners := parseScanners(input)
	sc := scanners[0]
	scPerms := permutations(sc.points)
	for _, scPerm := range scPerms {
		t.Log(scPerm)
	}
}
