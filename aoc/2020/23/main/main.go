package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"time"
)

var (
	fInput  = flag.String("in", `389125467`, "input")
	fRounds = flag.Int("rounds", 100, "rounds")
	fSz     = flag.Int("size", 9, "size")

	w io.Writer
)

func main() {
	flag.Parse()
	start := time.Now()
	defer func() {
		fmt.Println(time.Since(start))
	}()

	//run2(*fInput, *fRounds, *fSz)
	run3(*fInput, *fRounds, *fSz)
}

type node struct {
	v int32
	n *node
}

func run3(input string, rounds int, size int) {
	if rounds*size < 10000 {
		w = os.Stdout
	}
	nodes := make([]node, size)
	nodesByVal := make([]*node, size+1)

	for i, c := range input {
		v := c - '0'
		nodes[i].v = v
		nodesByVal[v] = &nodes[i]
	}
	for i := len(input); i < size; i++ {
		v := int32(i + 1)
		nodes[i].v = v
		nodesByVal[v] = &nodes[i]
	}

	for i := 0; i < size-1; i++ {
		nodes[i].n = &nodes[i+1]
	}
	nodes[size-1].n = &nodes[0]
	sz := int32(size)

	head := &nodes[0]
	for move := 0; move < rounds; move++ {
		fmtPrintf("-- move %d --\n", move+1)
		fmtPrintf("cups: %+v\n", render(head))

		a := head.n
		b := a.n
		c := b.n
		head.n = c.n

		fmtPrintf("pick up: %+v\n", []int32{a.v, b.v, c.v})

		tgt := head.v - 1
		for ; true; tgt = (tgt + sz) % (sz + 1) {
			if tgt == 0 || tgt == a.v || tgt == b.v || tgt == c.v {
				continue
			}
			break
		}

		// insert after
		t := nodesByVal[tgt]
		c.n = t.n
		t.n = a

		head = head.n
		fmtPrintln()
	}

	fmtPrintln(render(head))

	one := nodesByVal[1]
	a := one.n.v
	b := one.n.n.v
	fmt.Println(a, b, int(a)*int(b))
}

func render(head *node) []int32 {
	if w == nil {
		return nil
	}
	ret := []int32{head.v}
	for t := head.n; t != head; t = t.n {
		ret = append(ret, t.v)
	}
	return ret
}

func fmtPrintln(args ...interface{}) {
	if w != nil {
		fmt.Fprintln(w, args...)
	}
}

func fmtPrintf(s string, args ...interface{}) {
	if w != nil {
		fmt.Fprintf(w, s, args...)
	}
}
