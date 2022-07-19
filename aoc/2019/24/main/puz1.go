package main

import "fmt"

func run1(input string) {
	p := puz{data: parse(input)}
	seen := map[string]bool{
		p.String(): true,
	}
	for {
		p = p.next()
		key := p.String()
		if seen[key] {
			fmt.Println(p.Score())
			break
		}
		seen[key] = true
	}
}

type puz struct {
	data []byte
}

func (p *puz) val(x, y int) int {
	if x < 0 || x >= w || y < 0 || y >= h {
		return 0
	}
	if p.data[y*w+x] == '#' {
		return 1
	}
	return 0
}

func (p *puz) next() puz {
	newData := make([]byte, sz)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			v := byte('.')
			sum := p.val(x-1, y) + p.val(x+1, y) + p.val(x, y-1) + p.val(x, y+1)
			switch p.val(x, y) {
			case 0:
				if sum == 1 || sum == 2 {
					v = '#'
				}
			case 1:
				if sum == 1 {
					v = '#'
				}
			}
			newData[y*w+x] = v

		}
	}
	return puz{data: newData}
}

func (p *puz) String() string {
	return string(p.data)
}

func (p *puz) Score() int {
	sum := 0
	pwr := 1
	for _, v := range p.data {
		if v == '#' {
			sum += pwr
		}
		pwr <<= 1
	}
	return sum
}
