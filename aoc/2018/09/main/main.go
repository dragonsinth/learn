package main

import "fmt"

type game struct {
	p, n int
}

var samples = []game{
	{9, 25},
	{10, 1618},
	{13, 7999},
	{17, 1104},
	{21, 6111},
	{30, 5807},
}

func main() {
	for _, g := range samples {
		play(g)
	}
}

type node struct {
	val        int
	prev, next *node
}

func (n *node) Left(i int) *node {
	for i > 0 {
		n = n.prev
		i--
	}
	return n
}

func (n *node) InsertRight(r *node) {
	r.prev = n
	r.next = n.next
	n.next = r
	r.next.prev = r
}

func (n *node) Remove() {
	n.prev.next = n.next
	n.next.prev = n.prev
	n.prev, n.next = n, n
}

func (n *node) Walk() []int {
	ret := []int{n.val}
	for c := n.next; c != n; c = c.next {
		ret = append(ret, c.val)
	}
	return ret
}

func play(g game) {
	zero := &node{}
	zero.prev = zero
	zero.next = zero
	//fmt.Println(0, 0, zero.Walk())

	scores := make([]int, g.p)

	cur := zero
	player := 0
	for i := 1; i <= g.n; i++ {
		if i%23 == 0 {
			scores[player] += i
			target := cur.Left(7)
			cur = target.next
			target.Remove()
			scores[player] += target.val
		} else {
			n := &node{val: i}
			cur.next.InsertRight(n)
			cur = n
		}

		//fmt.Println(i, cur.val, zero.Walk())
		player = (player + 1) % g.p
	}

	best := 0
	for _, s := range scores {
		if s > best {
			best = s
		}
	}
	fmt.Println(best, scores)
}
