package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var sample = `
root: pppw + sjmn
dbpl: 5
cczh: sllz + lgvd
zczc: 2
ptdq: humn - dvpt
dvpt: 3
lfqf: 4
humn: 5
ljgn: 2
sjmn: drzm * dbpl
sllz: 4
pppw: cczh / lfqf
lgvd: ljgn * ptdq
drzm: hmdt - zczc
hmdt: 32`

var (
	re  = regexp.MustCompile(`^([a-z]{4}): (\d+|[a-z]{4} [+\-*/] [a-z]{4})$`)
	re2 = regexp.MustCompile(`^([a-z]{4}) ([+\-*/]) ([a-z]{4})$`)
)

type unresolvedExpr struct {
	key string
	lhs string
	op  string
	rhs string
}

func main() {
	refs := map[string]int{}
	unresolved := map[string]any{}

	for _, line := range strings.Split(sample, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if !re.MatchString(line) {
			panic(line)
		}

		matches := re.FindStringSubmatch(line)
		key := matches[1]
		if re2.MatchString(matches[2]) {
			matches := re2.FindStringSubmatch(matches[2])
			op := unresolvedExpr{key, matches[1], matches[2], matches[3]}
			unresolved[key] = op
			refs[op.lhs]++
			refs[op.rhs]++
			if refs[op.lhs] > 1 {
				panic(op.lhs)
			}
			if refs[op.rhs] > 1 {
				panic(op.rhs)
			}
		} else {
			unresolved[key] = literal{key, mustInt(matches[2])}
		}
	}

	resolved := map[string]expr{}
	root := resolve(unresolved, "root", resolved)
	root.Code(os.Stdout)
	fmt.Println(root.Eval())

	part2 := root.Simplify().(operation)
	part2.op = '='
	part2.Code(os.Stdout)
	fmt.Println(part2.Reverse(1))
}

func resolve(unresolved map[string]any, key string, resolved map[string]expr) expr {
	if expr, ok := resolved[key]; ok {
		return expr
	}
	switch t := unresolved[key].(type) {
	case literal:
		resolved[key] = t
		return t
	case unresolvedExpr:
		lhs := resolve(unresolved, t.lhs, resolved)
		rhs := resolve(unresolved, t.rhs, resolved)
		op := operation{
			key: key,
			lhs: lhs,
			op:  t.op[0],
			rhs: rhs,
		}
		resolved[key] = op
		return op
	default:
		panic(t)
	}
}

func mustInt(s string) int {
	if v, err := strconv.Atoi(s); err != nil {
		panic(fmt.Sprint(s, err))
	} else {
		return v
	}
}
