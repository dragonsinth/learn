package main

import (
	"bytes"
	"fmt"
)

const (
	size = 200_000_000
)

func main() {
	in := make([]byte, 0, size)
	in = append(in, 3, 7)
	i, j := 0, 1
	for len(in) < size {
		a, b := in[i], in[j]
		sum := a + b
		if sum >= 10 {
			in = append(in, sum/10)
			in = append(in, sum%10)
		} else {
			in = append(in, sum)
		}
		i = (i + 1 + int(a)) % len(in)
		j = (j + 1 + int(b)) % len(in)
	}
	doPrint(in, 5)
	doPrint(in, 9)
	doPrint(in, 18)
	doPrint(in, 2018)

	search(in, "01245")
	search(in, "51589")
	search(in, "92510")
	search(in, "59414")
}

func doPrint(in []byte, pos int) {
	buf := make([]byte, 10)
	for i := 0; i < 10; i++ {
		buf[i] = in[pos+i] + '0'
	}
	fmt.Println(pos, string(buf))
}

func search(in []byte, s string) {
	match := []byte(s)
	for i := range match {
		match[i] -= '0'
	}

	sz := len(match)
	for i, c := 0, len(in)-sz; i < c; i++ {
		if bytes.Equal(in[i:i+sz], match) {
			fmt.Println(i)
			return
		}
	}
	fmt.Println("not found")
}
