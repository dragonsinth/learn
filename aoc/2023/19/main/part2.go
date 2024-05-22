package main

import "fmt"

func (p puz) part2() {
	ruleStart, ok := p.rules["in"]
	if !ok {
		panic("in")
	}
	voxelStart := voxel{
		X: {1, 4000},
		M: {1, 4000},
		A: {1, 4000},
		S: {1, 4000},
	}

	var sum int64
	p.traverseRule(voxelStart, ruleStart, func(v voxel) {
		sum += v.volume()
	})
	fmt.Println(sum)
}

func (p puz) traverseRule(v voxel, r rule, accept func(v voxel)) {
	// Split the voxel by each step in the rule.
	for _, st := range r.steps {
		vrem, vspl := v, v
		switch st.sym {
		case '<':
			// splitting 1..4000 <1000 yields split=1..999, rem=1000..4000
			vspl[st.q].max = min(v[st.q].max, st.val-1)
			vrem[st.q].min = max(v[st.q].min, st.val)
		case '>':
			// splitting 1..4000 >1000 yields rem=1..1000, split=1001..4000
			vrem[st.q].max = min(v[st.q].max, st.val)
			vspl[st.q].min = max(v[st.q].min, st.val+1)
		default:
			panic(st.sym)
		}
		if vspl.valid() {
			p.traverseTarget(vspl, st.t, accept)
		}
		if !vrem.valid() {
			return
		}
		v = vrem
	}
	p.traverseTarget(v, r.def, accept)
}

func (p puz) traverseTarget(v voxel, t target, accept func(v voxel)) {
	switch t {
	case "A":
		accept(v)
	case "R":
		return
	default:
		nextRule, ok := p.rules[t]
		if !ok {
			panic(t)
		}
		p.traverseRule(v, nextRule, accept)
	}
}

type voxel [4]dimRange

func (v voxel) valid() bool {
	return v[X].min <= v[X].max &&
		v[M].min <= v[M].max &&
		v[A].min <= v[A].max &&
		v[S].min <= v[S].max
}

func (v voxel) volume() int64 {
	ret := int64(1)
	ret *= int64(v[X].max - v[X].min + 1)
	ret *= int64(v[M].max - v[M].min + 1)
	ret *= int64(v[A].max - v[A].min + 1)
	ret *= int64(v[S].max - v[S].min + 1)
	return ret
}

type dimRange struct {
	min, max int
}
