package main

import (
	"fmt"
	"regexp"
	"strings"
)

var sample1 = `
RL

AAA = (BBB, CCC)
BBB = (DDD, EEE)
CCC = (ZZZ, GGG)
DDD = (DDD, DDD)
EEE = (EEE, EEE)
GGG = (GGG, GGG)
ZZZ = (ZZZ, ZZZ)
`

var sample2 = `
LLR

AAA = (BBB, BBB)
BBB = (AAA, ZZZ)
ZZZ = (ZZZ, ZZZ)
`

var sample3 = `
LR

11A = (11B, XXX)
11B = (XXX, 11Z)
11Z = (11B, XXX)
22A = (22B, XXX)
22B = (22C, 22C)
22C = (22Z, 22Z)
22Z = (22B, 22B)
XXX = (XXX, XXX)
`

var re1 = regexp.MustCompile(`^[RL]+$`)
var re2 = regexp.MustCompile(`^([A-Z0-9]{3}) = \(([A-Z0-9]{3}), ([A-Z0-9]{3})\)$`)

func main() {
	part1(parse(sample1))
	part1(parse(sample2))
	part2(parse(sample3))
}

func parse(input string) ([]byte, map[string]*node) {
	nodes := map[string]*node{}

	var path []byte
	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if path == nil {
			if !re1.MatchString(line) {
				panic(line)
			}
			path = []byte(line)
		} else {
			if !re2.MatchString(line) {
				panic(line)
			}
			matches := re2.FindStringSubmatch(line)
			src, lf, rt := matches[1], matches[2], matches[3]
			nodes[src] = &node{
				src:   src,
				lf:    lf,
				rt:    rt,
				left:  nil,
				right: nil,
			}
		}
	}

	for _, v := range nodes {
		v.left = nodes[v.lf]
		if v.left == nil {
			panic(v.lf)
		}
		v.right = nodes[v.rt]
		if v.right == nil {
			panic(v.right)
		}
	}

	return path, nodes
}

type node struct {
	src, lf, rt string
	left, right *node
}

func (cur *node) next(wat byte) *node {
	switch wat {
	case 'L':
		return cur.left
	case 'R':
		return cur.right
	default:
		panic(wat)
	}
}

func part1(path []byte, nodes map[string]*node) {
	cur := nodes["AAA"]
	step := 0
	for cur.src != "ZZZ" {
		pos := step % len(path)
		cur = cur.next(path[pos])
		step++
	}
	fmt.Println(step)
}

type state struct {
	src string
	pos int
}

func part2(path []byte, nodes map[string]*node) {
	var factors []int
	for _, v := range nodes {
		if v.src[2] == 'A' {
			// Compute a chain that loops.
			cur := v
			seen := map[state]int{}
			step := 0
			zpos := 0
			for {
				if cur.src[2] == 'Z' {
					zpos = step
				}

				pos := step % len(path)
				st := state{
					src: cur.src,
					pos: pos,
				}
				if last, ok := seen[st]; ok {
					fmt.Printf("src=%s, pos=%s, seen=%d, repeat=%d, zpos=%d\n", v.src, st.src, last, step, zpos)
					if zpos != step-last {
						panic(zpos)
					}
					factors = append(factors, zpos)
					break
				}
				seen[st] = step
				cur = cur.next(path[pos])
				step++
			}
		}
	}

	fmt.Println(factors)
	m := 1
	for _, v := range factors {
		m = lcm(m, v)
	}
	fmt.Println(m)
}

func gcd(a, b int) int {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

func lcm(a, b int) int {
	return a * (b / gcd(a, b))
}
