package main

type lex struct {
	in  []byte
	pos int
}

func (l *lex) nextInt() int {
	i := l.pos
	for i < len(l.in) {
		c := l.in[i]
		if c >= '0' && c <= '9' {
			i++
		} else {
			break
		}
	}
	ret := mustInt(string(l.in[l.pos:i]))
	l.pos = i
	return ret
}

func (l *lex) nextDir() turn {
	if l.pos >= len(l.in) {
		return 0
	}
	switch c := l.in[l.pos]; c {
	case 'L', 'R':
		l.pos++
		return turn(c)
	default:
		panic(c)
	}
}
