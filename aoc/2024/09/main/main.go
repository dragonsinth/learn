package main

import (
	"fmt"
	"slices"
)

func main() {
	part1("12345", true)
	part1("2333133121414131402", true)

	part2("2333133121414131402", true)
}

type block struct {
	free bool
	sz   int
	id   int
}

func parse(in string) []block {
	var ret []block
	for i, c := range in {
		sz := int(c - '0')
		if i%2 == 0 {
			id := i / 2
			if sz < 1 {
				panic(id)
			}
			ret = append(ret, block{
				free: false,
				sz:   sz,
				id:   id,
			})
		} else {
			if sz > 0 {
				ret = append(ret, block{
					free: true,
					sz:   sz,
					id:   0,
				})
			}
		}
	}
	return ret
}

func sum(in []block) int {
	pos := 0
	ret := 0
	for _, b := range in {
		if b.free {
			pos += b.sz
			continue
		}
		// add each pos
		for sz := b.sz; sz > 0; sz-- {
			ret += pos * b.id
			pos++
		}
	}
	return ret
}

func printBlocks(in []block) {
	for _, b := range in {
		if b.sz < 1 {
			panic(b)
		}
		if b.free {
			for i := 0; i < b.sz; i++ {
				fmt.Print(".")
			}
		} else {
			s := string(rune(b.id + '0'))
			for i := 0; i < b.sz; i++ {
				fmt.Print(s)
			}
		}
	}
	fmt.Println()
}

func part1(input string, debug bool) {
	in := parse(input)

	var out []block

	// merge copy from both ends
	li := 0
	ri := len(in) - 1
	for {
		if debug {
			printBlocks(out)
		}

		if li == ri {
			// this is the last block, copy only if it has data
			b := &in[li]
			if !b.free {
				out = append(out, *b)
			}
			break // all done
		}

		lb := &in[li]
		rb := &in[ri]

		// existing files on the left stay put
		if !lb.free {
			out = append(out, *lb)
			li++
			continue
		}

		// skip free blocks on the right
		if rb.free {
			ri--
			continue
		}

		// need to move right to left... 3 cases!
		if lb.sz == rb.sz {
			// exact fit!
			out = append(out, *rb)
			li++
			ri--
		} else if lb.sz < rb.sz {
			// move part of the file, but not all
			out = append(out, block{
				free: false,
				sz:   lb.sz,
				id:   rb.id,
			})
			li++
			rb.sz -= lb.sz
		} else {
			// move entire file, but don't consume all the free space
			out = append(out, *rb)
			ri--
			lb.sz -= rb.sz
		}
	}
	if debug {
		printBlocks(out)
	}

	fmt.Println(sum(out))
}

func part2(input string, debug bool) {
	in := parse(input)

	for ri := len(in) - 1; ri >= 0; ri-- {
		if debug {
			printBlocks(in)
		}

		rb := &in[ri]
		if rb.free {
			continue
		}

		// find the leftmost spot this could fit
		for li := 0; li < ri; li++ {
			lb := &in[li]
			if !lb.free {
				continue
			}
			if lb.sz < rb.sz {
				continue // can't fit
			} else if lb.sz == rb.sz {
				// perfect fit
				lb.free = false
				lb.id = rb.id
				rb.free = true
				rb.id = 0
				break
			} else { // lb.sz > rb.sz
				// must splice!
				rem := lb.sz - rb.sz
				lb.free = false
				lb.id = rb.id
				lb.sz = rb.sz

				rb.free = true
				rb.id = 0
				in = slices.Insert(in, li+1, block{
					free: true,
					sz:   rem,
					id:   0,
				})
				ri++ // shift our current index right to account for the insert
				break
			}
		}
	}

	fmt.Println(sum(in))
}
