package main

import (
	"fmt"
	"testing"
)

func TestExpand(t *testing.T) {
	{
		robotPad := computeShortestPaths(codePad, humanCosts{})
		debugPrint(expand(robotPad, "029A"))
		debugPrint("<A^A>^^AvvvA")
	}
	{
		aPad := computeShortestPaths(directionalPad, humanCosts{})
		bPad := computeShortestPaths(codePad, robotCosts{aPad})
		debugPrint(expand(bPad, "029A"))
		debugPrint("v<<A>>^A<A>AvA<^AA>A<vAAA>^A")
	}
	{
		aPad := computeShortestPaths(directionalPad, humanCosts{})
		bPad := computeShortestPaths(directionalPad, robotCosts{aPad})
		cPad := computeShortestPaths(codePad, robotCosts{bPad})
		debugPrint(expand(cPad, "029A"))
		debugPrint("<vA<AA>>^AvAA<^A>A<v<A>>^AvA^A<vA>^A<v<A>^A>AAvA^A<v<A>A>^AAAvA<^A>A")
	}
}

func debugPrint(v string) {
	fmt.Println(len(v), v)
}

func TestReduce(t *testing.T) {
	reduce("<v<A>>^A<A>AvA<^AA>A<vAAA>^A", directionalPad)
	reduce("v<<A>>^A<A>AvA<^AA>A<vAAA>^A", directionalPad)
	reduce("<v<A>>^A<A>A<A>vA<^A>A<vAAA>^A", directionalPad)

	reduce("<A^A>^^AvvvA", codePad)
	reduce("<A^A^>^AvvvA", codePad)
}

func reduce(in string, pad map[byte]pos) {
	var ret string
	posMap := map[pos]byte{}
	for k, pt := range pad {
		posMap[pt] = k
	}

	pt := pad['A']
	for _, c := range in {
		switch c {
		case '<':
			pt.x--
		case '>':
			pt.x++
		case '^':
			pt.y--
		case 'v':
			pt.y++
		case 'A':
			ret += string(rune(posMap[pt]))
		}
	}

	fmt.Println(ret)
}
