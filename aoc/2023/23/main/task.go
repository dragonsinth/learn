package main

type task struct {
	x, y  int
	d     dir
	route []taskKey
}

func (t task) pos() pos {
	return pos{x: t.x, y: t.y}
}

func (t task) next() task {
	switch t.d {
	case N:
		t.y--
	case E:
		t.x++
	case S:
		t.y++
	case W:
		t.x--
	default:
		panic(t.d)
	}
	route := append([]taskKey{}, t.route...)
	t.route = append(route, t.key())
	return t
}

func (t task) left() task {
	t.d = t.d.left()
	return t.next()
}

func (t task) right() task {
	t.d = t.d.right()
	return t.next()
}

type taskKey struct {
	x, y int
	d    dir
}

func (t task) key() taskKey {
	return taskKey{
		x: t.x,
		y: t.y,
		d: t.d,
	}
}
