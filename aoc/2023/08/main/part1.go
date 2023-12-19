package main

import "fmt"

func (p puz) part1(parts []part) {
	var accepted []part
	for _, pt := range parts {
		state, ok := p.rules["in"]
		if !ok {
			panic("in")
		}
	loop:
		for {
			t := state.def
			for _, st := range state.steps {
				val := pt[st.q]
				pass := val < st.val && st.sym == '<'
				pass = pass || val > st.val && st.sym == '>'
				if pass {
					t = st.t
					break
				}
			}

			switch t {
			case "A":
				accepted = append(accepted, pt)
				break loop
			case "R":
				break loop
			default:
				next, ok := p.rules[t]
				if !ok {
					panic(t)
				}
				state = next
			}
		}
	}
	fmt.Println(len(accepted))
	sum := 0
	for _, pt := range accepted {
		sum += pt[X]
		sum += pt[M]
		sum += pt[A]
		sum += pt[S]
	}
	fmt.Println(sum)
}

type part [4]int
