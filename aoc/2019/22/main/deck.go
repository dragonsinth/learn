package main

func newDeck(count int) deck {
	d := deck{}
	for i := 0; i < count; i++ {
		d = append(d, i)
	}
	return d
}

type deck []int

func (d deck) run(ops []op) deck {
	for _, o := range ops {
		switch o.code {
		case 0:
			d = d.reverse()
		case 1:
			d = d.cut(o.arg)
		case 2:
			d = d.increment(o.arg)
		}
	}
	return d
}

func (d deck) reverse() deck {
	r := make(deck, len(d))
	for i := 0; i < len(d); i++ {
		r[len(d)-i-1] = d[i]
	}
	return r
}

func (d deck) cut(v int) deck {
	if v < 0 {
		v = len(d) + v
	}

	r := make(deck, len(d))
	copy(r, d[v:])
	copy(r[len(d)-v:], d)
	return r
}

func (d deck) increment(v int) deck {
	r := make(deck, len(d))
	for i := 0; i < len(d); i++ {
		dst := (i * v) % len(d)
		r[dst] = d[i]
	}
	return r
}
