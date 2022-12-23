package main

import "fmt"

type pos struct {
	x, y int
}

func (pt pos) advance(d dir) pos {
	switch d {
	case NORTH:
		pt.y--
	case SOUTH:
		pt.y++
	case WEST:
		pt.x--
	case EAST:
		pt.x++
	default:
		panic(d)
	}
	return pt
}

func (pt pos) String() string {
	return fmt.Sprintf("(%d,%d)", pt.x, pt.y)
}

func (pt pos) adjacent(d dir) []pos {
	switch d {
	case NORTH:
		return []pos{
			{pt.x - 1, pt.y - 1},
			{pt.x + 0, pt.y - 1},
			{pt.x + 1, pt.y - 1},
		}
	case SOUTH:
		return []pos{
			{pt.x - 1, pt.y + 1},
			{pt.x + 0, pt.y + 1},
			{pt.x + 1, pt.y + 1},
		}
	case WEST:
		return []pos{
			{pt.x - 1, pt.y - 1},
			{pt.x - 1, pt.y + 0},
			{pt.x - 1, pt.y + 1},
		}
	case EAST:
		return []pos{
			{pt.x + 1, pt.y - 1},
			{pt.x + 1, pt.y + 0},
			{pt.x + 1, pt.y + 1},
		}
	default:
		panic(d)
	}
}

func (pt pos) allAdjacent() []pos {
	return []pos{
		{pt.x - 1, pt.y},
		{pt.x + 1, pt.y},
		{pt.x, pt.y - 1},
		{pt.x, pt.y + 1},
		{pt.x - 1, pt.y - 1},
		{pt.x + 1, pt.y + 1},
		{pt.x + 1, pt.y - 1},
		{pt.x - 1, pt.y + 1},
	}
}

type dir int

const (
	NORTH = dir(iota)
	SOUTH
	WEST
	EAST
)

func (d dir) byte() byte {
	return "^v<>"[d]
}
