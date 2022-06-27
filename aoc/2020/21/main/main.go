package main

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

const sample = `
mxmxvkd kfcds sqjhc nhms (contains dairy, fish)
trh fvjkl sbzzf mxmxvkd (contains dairy)
sqjhc fvjkl (contains soy)
sqjhc mxmxvkd sbzzf (contains fish)
`

var (
	input = sample

	lineRe = regexp.MustCompile(`^([a-z]+)(\s[a-z]+)* \(contains ([a-z]+)(, [a-z]+)*\)$`)
)

type ingredient string
type allergen string

type puzzle struct {
	allIngs map[ingredient]bool
	allAlgs map[allergen]bool

	entries []entry

	solved map[ingredient]allergen
}

func (p puzzle) String() string {
	var buf strings.Builder
	for _, e := range p.entries {
		fmt.Fprintf(&buf, "%+s\n", e)
	}
	for k, v := range p.solved {
		fmt.Fprintf(&buf, "%s => %s\n", k, v)
	}
	return buf.String()
}

type entry struct {
	ings []ingredient
	algs []allergen
}

func (e entry) HasAllergen(alg allergen) bool {
	for _, v := range e.algs {
		if v == alg {
			return true
		}
	}
	return false
}

func (e entry) Without(ing ingredient, alg allergen) entry {
	var ret entry
	for _, v := range e.ings {
		if v != ing {
			ret.ings = append(ret.ings, v)
		}
	}
	for _, v := range e.algs {
		if v != alg {
			ret.algs = append(ret.algs, v)
		}
	}
	return ret
}

func main() {
	p := puzzle{
		allIngs: map[ingredient]bool{},
		allAlgs: map[allergen]bool{},
		solved:  map[ingredient]allergen{},
	}

	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if lineRe.MatchString(line) {
			split := strings.Split(line, " (contains ")
			ings := strings.Split(split[0], " ")
			algs := strings.Split(split[1][:len(split[1])-1], ", ")
			sort.Strings(ings)
			sort.Strings(algs)

			var e entry
			for _, v := range ings {
				v := ingredient(v)
				e.ings = append(e.ings, v)
				p.allIngs[v] = true
			}
			for _, v := range algs {
				v := allergen(v)
				e.algs = append(e.algs, v)
				p.allAlgs[v] = true
			}
			p.entries = append(p.entries, e)
		} else {
			panic(line)
		}
	}

	for {
		// Find one we can solve.
		// For each allergen, find all the rules for that allergen.
		// See if the intersection of ingredients devolves to a single case.
		// Map and simplify.
		found := false
		for alg := range p.allAlgs {
			var entries []entry
			for _, e := range p.entries {
				if e.HasAllergen(alg) {
					entries = append(entries, e)
				}
			}

			// Retain only the ingr common to all rules
			var common []ingredient
			for i, e := range entries {
				if i == 0 {
					common = e.ings
				} else {
					common = intersect(common, e.ings)
				}
			}

			if len(common) == 1 {
				// Found a match; emit and simplify the puzzle
				found = true
				ing := common[0]
				p.solved[ing] = alg

				// Remove both the ing and alg from every rule.
				var newEntries []entry
				for _, e := range p.entries {
					e = e.Without(ing, alg)
					if len(e.ings) > 0 {
						newEntries = append(newEntries, e)
					}
				}
				p.entries = newEntries

				break
			}
		}

		if !found {
			break
		}
	}

	if len(p.allAlgs) != len(p.solved) {
		panic("unsolved")
	}

	// Count the ings remaining
	sum := 0
	for _, e := range p.entries {
		sum += len(e.ings)
	}
	fmt.Println(p.String())
	fmt.Println(sum)

	type danger struct {
		ing ingredient
		alg allergen
	}

	var ret []danger
	for k, v := range p.solved {
		ret = append(ret, danger{
			ing: k,
			alg: v,
		})
	}
	sort.Slice(ret, func(i, j int) bool {
		return ret[i].alg < ret[j].alg
	})
	for i, d := range ret {
		if i > 0 {
			fmt.Print(",")
		}
		fmt.Print(d.ing)
	}
}

func intersect(a []ingredient, b []ingredient) []ingredient {
	var ret []ingredient
	for ai, bi := 0, 0; ai < len(a) && bi < len(b); {
		if a[ai] == b[bi] {
			ret = append(ret, a[ai])
			ai++
			bi++
		} else if a[ai] < b[bi] {
			ai++
		} else if a[ai] > b[bi] {
			bi++
		} else {
			panic("nope")
		}
	}
	return ret
}
