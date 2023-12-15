package main

import (
	"fmt"
	"golang.org/x/exp/slices"
	"strings"
)

var sample = `rn=1,cm-,qp=3,cm=2,qp-,pc=4,ot=9,ab=5,pc-,pc=6,ot=7`

func main() {
	run(sample, true)
}

func run(input string, debug bool) {
	parts := strings.Split(input, ",")
	sum := 0
	var boxes [256]box
	for _, inst := range parts {
		h := hash(inst)
		sum += h

		l := parseInst(inst)
		which := hash(l.name)
		boxes[which].apply(l)
		if debug {
			fmt.Println(inst, "=", h)
			for i := range boxes {
				if !boxes[i].empty() {
					fmt.Printf("Box %d: %s\n", i, boxes[i].String())
				}
			}
		}
	}
	fmt.Println("sum", sum)

	pow := 0
	for i, b := range boxes {
		if !b.empty() {
			p := b.pow()
			if debug {
				fmt.Printf("Box %d: %d\n", i, p)
			}
			pow += (i + 1) * p
		}
	}

	fmt.Println("pow", pow)
}

func hash(p string) int {
	var sum int32
	for _, c := range p {
		sum += c
		sum *= 17
		sum %= 256
	}
	return int(sum)
}

func parseInst(inst string) lens {
	if strings.HasSuffix(inst, "-") {
		return lens{
			name: inst[:len(inst)-1],
			pow:  0,
		}
	}
	parts := strings.Split(inst, "=")
	if len(parts) != 2 {
		panic(parts)
	}
	if len(parts[1]) != 1 {
		panic(parts[1])
	}
	return lens{
		name: parts[0],
		pow:  int(parts[1][0] - '0'),
	}
}

type box struct {
	lenses []lens
}

func (b *box) String() string {
	var s strings.Builder
	for i, l := range b.lenses {
		if i > 0 {
			s.WriteRune(' ')
		}
		s.WriteRune('[')
		s.WriteString(l.name)
		s.WriteRune(' ')
		s.WriteRune(rune(l.pow + '0'))
		s.WriteRune(']')
	}
	return s.String()
}

func (b *box) apply(l lens) {
	if l.pow == 0 {
		b.lenses = slices.DeleteFunc(b.lenses, func(test lens) bool {
			return l.name == test.name
		})
	} else {
		pos := slices.IndexFunc(b.lenses, func(test lens) bool {
			return l.name == test.name
		})
		if pos < 0 {
			b.lenses = append(b.lenses, l)
		} else {
			b.lenses[pos].pow = l.pow
		}
	}
}

func (b *box) empty() bool {
	return len(b.lenses) == 0
}

func (b *box) pow() int {
	ret := 0
	for i, l := range b.lenses {
		ret += (i + 1) * l.pow
	}
	return ret
}

type lens struct {
	name string
	pow  int
}
