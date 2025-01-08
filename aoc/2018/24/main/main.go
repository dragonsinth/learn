package main

import (
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

const sample = `
Immune System:
17 units each with 5390 hit points (weak to radiation, bludgeoning) with an attack that does 4507 fire damage at initiative 2
989 units each with 1274 hit points (immune to fire; weak to bludgeoning, slashing) with an attack that does 25 slashing damage at initiative 3

Infection:
801 units each with 4706 hit points (weak to radiation) with an attack that does 116 bludgeoning damage at initiative 1
4485 units each with 2961 hit points (immune to radiation; weak to fire, cold) with an attack that does 12 slashing damage at initiative 4
`

var (
	re = regexp.MustCompile(`(\d+) units each with (\d+) hit points (\([^)]+\) )?with an attack that does (\d+) (\w+) damage at initiative (\d+)`)
)

func main() {
	run(sample, 0, true)
	//run(data, 0, false)

	part2(sample, 0)
	// part2(data, 34)
}

func part2(input string, deadlock int) {
	for boost := 1; true; boost <<= 1 {
		if run(input, boost, false) {
			fmt.Printf("boost %d wins\n", boost)
			// immune system finally won, now binary search
			lo := boost / 2
			hi := boost
			for {
				mid := (lo + hi) / 2
				if mid != deadlock && run(input, mid, false) {
					// immune system won; lower the hi mark
					fmt.Printf("boost %d wins\n", mid)
					hi = mid
				} else {
					// immune system lost; raise the lo mark
					fmt.Printf("boost %d loses\n", mid)
					lo = mid
				}
				if lo+1 == hi {
					// done, hi is the smallest boost
					fmt.Println("smallest boost:", hi)
					run(input, hi, true)
					break
				}
			}
			break
		}
		fmt.Printf("boost %d loses\n", boost)
	}
}

func run(input string, boost int, debug bool) bool {
	imm, inf := parse(input)
	for _, g := range imm.groups {
		g.dmg += boost
	}
	battle(imm, inf, debug)
	return len(imm.groups) > 0
}

func battle(imm *army, inf *army, debug bool) {
	for len(imm.groups) > 0 && len(inf.groups) > 0 {
		if debug {
			imm.debugPrint()
			inf.debugPrint()
			fmt.Println()
		}
		fight(imm, inf, debug)
	}
	if debug {
		imm.debugPrint()
		inf.debugPrint()
	}
}

func fight(imm *army, inf *army, debug bool) {
	t1 := selectTargets(imm, inf, debug)
	t2 := selectTargets(inf, imm, debug)
	allTargets := append(t1, t2...)
	slices.SortFunc(allTargets, func(a, b target) int {
		return b.attacker.init - a.attacker.init
	})

	if debug {
		fmt.Println()
	}
	for _, t := range allTargets {
		killed := t.attacker.damage(t.defender) / t.defender.hp
		killed = min(killed, t.defender.count)
		if debug {
			fmt.Printf("%s group %d attacks defending group %d, killing %d units\n", t.attacker.army.name, t.attacker.id, t.defender.id, killed)
		}
		t.defender.count -= killed
	}
	imm.prune()
	inf.prune()
}

type target struct {
	attacker *group
	defender *group
}

func selectTargets(atk *army, def *army, debug bool) []target {
	slices.SortFunc(atk.groups, func(a, b *group) int {
		if aPow, bPow := a.power(), b.power(); aPow != bPow {
			return bPow - aPow
		}
		return b.init - a.init
	})

	var targets []target
	defGroups := append([]*group{}, def.groups...)
	for _, g := range atk.groups {
		if len(defGroups) == 0 {
			break
		}

		slices.SortFunc(defGroups, func(a, b *group) int {
			if aDmg, bDmg := g.damage(a), g.damage(b); aDmg != bDmg {
				return bDmg - aDmg
			}
			if aPow, bPow := a.power(), b.power(); aPow != bPow {
				return bPow - aPow
			}
			return b.init - a.init
		})

		best := defGroups[0]
		if dmg := g.damage(best); dmg > 0 {
			if debug {
				fmt.Printf("%s group %d would deal defending group %d %d damage\n", g.army.name, g.id, best.id, dmg)
			}
			targets = append(targets, target{attacker: g, defender: best})
			defGroups = defGroups[1:]
		}
	}

	return targets
}

type army struct {
	name   string
	groups []*group
}

func (a *army) debugPrint() {
	fmt.Println(a.name + ":")
	sum := 0
	for _, g := range a.groups {
		fmt.Println(g.desc())
		sum += g.count
	}
	fmt.Println(sum, "total units")
}

func (a *army) addGroup(g group) {
	g.id = len(a.groups) + 1
	g.army = a
	a.groups = append(a.groups, &g)
}

func (a *army) prune() {
	a.groups = slices.DeleteFunc(a.groups, func(g *group) bool {
		return g.count == 0
	})
}

type group struct {
	id    int
	army  *army
	count int

	hp   int
	dmg  int
	typ  string
	init int

	strong []string
	weak   []string
}

func (g *group) desc() string {
	return fmt.Sprintf("Group %d contains %d units", g.id, g.count)
}

func (g *group) power() int {
	return g.dmg * g.count
}

func (g *group) damage(def *group) int {
	if slices.Contains(def.strong, g.typ) {
		return 0
	}
	if slices.Contains(def.weak, g.typ) {
		return g.power() * 2
	}
	return g.power()
}

func parse(input string) (*army, *army) {
	var imm, inf army
	var a *army

	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if line == "Immune System:" {
			a = &imm
			a.name = "Immune System"
		} else if line == "Infection:" {
			a = &inf
			a.name = "Infection"
		} else if re.MatchString(line) {
			matches := re.FindStringSubmatch(line)
			count := mustInt(matches[1])
			hp := mustInt(matches[2])
			strong, weak := parseAttributes(matches[3])
			dmg := mustInt(matches[4])
			typ := matches[5]
			init := mustInt(matches[6])
			a.addGroup(group{
				count:  count,
				hp:     hp,
				dmg:    dmg,
				typ:    typ,
				init:   init,
				strong: strong,
				weak:   weak,
			})
		} else {
			panic(line)
		}
	}

	return &imm, &inf
}

func parseAttributes(s string) ([]string, []string) {
	if s == "" {
		return nil, nil
	}
	s = s[1 : len(s)-2]
	var strong, weak []string
	for _, clause := range strings.Split(s, ";") {
		clause = strings.TrimSpace(clause)
		if strings.HasPrefix(clause, "weak to ") {
			clause = strings.TrimPrefix(clause, "weak to ")
			weak = parseList(clause)
		} else if strings.HasPrefix(clause, "immune to ") {
			clause = strings.TrimPrefix(clause, "immune to ")
			strong = parseList(clause)
		} else {
			panic(clause)
		}
	}
	return strong, weak
}

func parseList(l string) []string {
	return strings.Split(l, ", ")
}

func mustInt(s string) int {
	if v, err := strconv.Atoi(s); err != nil {
		panic(fmt.Sprint(s, err))
	} else {
		return v
	}
}
