package main

import (
	"fmt"
	"strconv"
	"strings"
)

var sample = `
1
2
-3
3
-2
0
4
`

type node struct {
	ord int
	val int
}

func main() {
	run(sample, 1, 1)
	run(sample, 811589153, 10)
}

func run(input string, mult int, count int) {
	var buf []node

	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		buf = append(buf, node{
			ord: len(buf),
			val: mustInt(line) * mult,
		})
	}

	sz := len(buf)
	for i := 0; i < count; i++ {
		for ord := 0; ord < sz; ord++ {
			// fmt.Println(ord, ":", str(buf))
			// fmt.Println("move:", buf[findByOrd(buf, ord)].val)
			mix(buf, ord)
		}
	}
	// fmt.Println("* :", str(buf))

	for ord := 0; ord < sz; ord++ {
		findByOrd(buf, ord)
	}

	pos := findByVal(buf, 0)
	a := buf[(pos+1000)%sz].val
	b := buf[(pos+2000)%sz].val
	c := buf[(pos+3000)%sz].val
	fmt.Println(a+b+c, ":", a, b, c) // -924 : -8606 3938 3744 wrong
}

func mix(in []node, ord int) {
	sz := len(in)
	srcPos := findByOrd(in, ord)
	srcNode := in[srcPos]

	mod := sz - 1
	dstPos := (srcPos + srcNode.val) % mod
	if dstPos < 0 {
		dstPos += mod
	}

	if srcPos < dstPos {
		// shift left, put node at end
		copy(in[srcPos:dstPos], in[srcPos+1:dstPos+1])
	} else if srcPos > dstPos {
		// shift right, put node at start
		copy(in[dstPos+1:srcPos+1], in[dstPos:srcPos])
	}
	in[dstPos] = srcNode
}

func findByVal(in []node, val int) int {
	for i, v := range in {
		if v.val == val {
			return i
		}
	}
	panic("here")
}

func findByOrd(in []node, ord int) int {
	for i, v := range in {
		if v.ord == ord {
			return i
		}
	}
	panic("here")
}

func str(in []node) string {
	var out strings.Builder
	for i, v := range in {
		if i > 0 {
			out.WriteString(", ")
		}
		out.WriteString(strconv.Itoa(v.val))
	}
	return out.String()
}

func mustInt(s string) int {
	if v, err := strconv.Atoi(s); err != nil {
		panic(fmt.Sprint(s, err))
	} else {
		return v
	}
}
