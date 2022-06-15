package main

import (
	"fmt"
)

var (
	minEnergy = 100000
	maxDepth  = 0

	moveToHall     int64
	moveToRoom     int64
	moveRoomToRoom int64

	stepBuffer = [256]step{}
)

type step struct {
	p    puzzle
	desc string
	cost int
}

func main() {
	stepBuffer[0].p = sample
	stepBuffer[0].desc = "start"
	stepBuffer[0].cost = 0

	findLowestEnergy(&stepBuffer[0].p, 0, 0)
	fmt.Println(moveToHall, moveToRoom, moveRoomToRoom)
}

var sample = puzzle{
	rooms: [4][2]what{
		{A, B},
		{D, C},
		{C, B},
		{A, D},
	},
}

const roomDepth = 2

//var sample = puzzle{
//	rooms: [4][4]what{
//		{A, D, D, B},
//		{D, B, C, C},
//		{C, A, B, B},
//		{A, C, A, D},
//	},
//}
//
//const roomDepth = 4

var (
	costs = []int{
		A: 1,
		B: 10,
		C: 100,
		D: 1000,
	}

	roomToJoinPoint = []int{2, 4, 6, 8}

	isJoinPoint = [hallMax]bool{2: true, 4: true, 6: true, 8: true}
)

type what int

const (
	EMPTY = what(iota)
	A
	B
	C
	D
)

const (
	roomMax = 4
	hallMax = 11
)

type puzzle struct {
	rooms [roomMax][roomDepth]what
	hall  [hallMax]what
}

func (p *puzzle) isSolved() bool {
	for _, v := range p.hall {
		if v != EMPTY {
			return false
		}
	}

	return true
}

func (p *puzzle) canPull(room int) (int, bool) {
	if p.rooms[room][0] == EMPTY {
		return 0, false
	}

	// If we find any aliens, we can pull.
	desiredWat := what(room + 1)
	for i := 0; i < roomDepth; i++ {
		v := p.rooms[room][i]
		if v == EMPTY {
			return 0, false // no aliens
		}
		if v != desiredWat {
			// Found an alien, return the highest non-empty space
			for j := i + 1; j < roomDepth; j++ {
				if p.rooms[room][j] == EMPTY {
					return j - 1, true
				}
			}
			return roomDepth - 1, true
		}
	}

	return 0, false
}

func (p *puzzle) canPush(room int) (int, bool) {
	desiredWat := what(room + 1)

	// Must have an empty spot, and no undesirables
	for i, v := range p.rooms[room] {
		if v == EMPTY {
			return i, true
		}
		if v != desiredWat {
			return 0, false
		}
	}
	panic("wat")
}

func (p *puzzle) top(room int) (what, int) {
	for i := roomDepth - 1; i >= 0; i-- {
		if p.rooms[room][i] != EMPTY {
			return p.rooms[room][i], i
		}
	}
	return EMPTY, 0
}

type result int

const (
	INVALID = result(iota)
	NO
	YES
)

func findLowestEnergy(p *puzzle, cost int, depth int) result {
	if cost > 0 && p.isSolved() {
		if cost >= minEnergy {
			return NO
		}

		fmt.Println(cost, depth)
		minEnergy = cost
		for i := 0; i <= depth; i++ {
			for _, v := range stepBuffer[i].p.hall {
				switch v {
				case EMPTY:
					fmt.Print(".")
				default:
					fmt.Print(string(rune(v + 'A' - 1)))
				}
			}
			fmt.Print(" ")
			fmt.Print(stepBuffer[i].desc)
			fmt.Print(" ")
			fmt.Print(stepBuffer[i].cost)
			fmt.Println()

			// Print room states.
			for r := roomDepth - 1; r >= 0; r-- {
				fmt.Print("  ")
				for j := 0; j < roomMax; j++ {
					v := stepBuffer[i].p.rooms[j][r]
					if v == EMPTY {
						fmt.Print(".")
					} else {
						fmt.Print(string(rune(v + 'A' - 1)))
					}
					fmt.Print(" ")
				}
				fmt.Println()
			}
		}

		return YES
	}

	// Try moving each pod in the hallway to its destination room.
	for hall := 0; hall < hallMax; hall++ {
		if isJoinPoint[hall] {
			continue
		}
		wat := p.hall[hall]
		if wat == EMPTY {
			continue
		}

		dstRoom := int(wat - 1)
		to, ok := p.canPush(dstRoom)
		if !ok {
			continue
		}

		res := tryMoveToRoom(p, wat, hall, to, dstRoom, cost, depth)
		if res == INVALID {
			continue
		}
		// A successful move is definitive
		return res
	}

	// Try moving from each room to its direct destination, cutting off a ton of intermediate states.
	for srcRoom := 0; srcRoom < roomMax; srcRoom++ {
		from, ok := p.canPull(srcRoom)
		if !ok {
			continue
		}
		wat := p.rooms[srcRoom][from]

		dstRoom := int(wat - 1)
		if dstRoom == srcRoom {
			continue
		}
		to, ok := p.canPush(dstRoom)
		if !ok {
			continue
		}

		res := tryMoveRoomToRoom(p, wat, from, srcRoom, to, dstRoom, cost, depth)
		if res == INVALID {
			continue
		}
		// A successful move is definitive
		return res
	}

	// If all else fails, try moving one pod from each room to every legal hall position.
	var bestRes result
	for srcRoom := 0; srcRoom < roomMax; srcRoom++ {
		from, ok := p.canPull(srcRoom)
		if !ok {
			continue
		}
		wat := p.rooms[srcRoom][from]

		for hall := 0; hall < hallMax; hall++ {
			if isJoinPoint[hall] {
				continue
			}
			if p.hall[hall] != EMPTY {
				continue
			}

			res := tryMoveToHall(p, wat, from, srcRoom, hall, cost, depth)
			if res > bestRes {
				bestRes = res
			}
		}
	}

	return bestRes
}

func tryMoveToHall(p *puzzle, wat what, from int, srcRoom int, hall int, cost int, depth int) result {
	nSpaces := checkPath(p, roomToJoinPoint[srcRoom], hall)
	if nSpaces == 0 {
		return INVALID
	}
	moveToHall++
	// Add spaces to get out of the room.
	nSpaces += roomDepth - from // (2 - 0) = 2

	newCost := costs[wat] * nSpaces
	cost += newCost
	if cost >= minEnergy {
		return NO // no better solution this way
	}

	// Do the move
	depth++
	if depth > maxDepth {
		maxDepth = depth
		fmt.Println(maxDepth)
	}
	stepBuffer[depth].p = puzzle{
		rooms: p.rooms,
		hall:  p.hall,
	}
	stepBuffer[depth].desc = "room-to-hall"
	stepBuffer[depth].cost = newCost
	newP := &stepBuffer[depth].p
	newP.rooms[srcRoom][from] = EMPTY
	newP.hall[hall] = wat

	// Recurse
	return findLowestEnergy(newP, cost, depth)
}

func tryMoveToRoom(p *puzzle, wat what, hall int, to, dstRoom int, cost int, depth int) result {
	nSpaces := checkPath(p, hall, roomToJoinPoint[dstRoom])
	if nSpaces == 0 {
		return INVALID
	}
	moveToRoom++
	// Add spaces to get into the room.
	nSpaces += roomDepth - to // (2 - 0)

	newCost := costs[wat] * nSpaces
	cost += newCost
	if cost >= minEnergy {
		return NO // no better solution this way
	}

	// Do the move
	depth++
	if depth > maxDepth {
		maxDepth = depth
		fmt.Println(maxDepth)
	}
	stepBuffer[depth].p = puzzle{
		rooms: p.rooms,
		hall:  p.hall,
	}
	stepBuffer[depth].desc = "hall-to-room"
	stepBuffer[depth].cost = newCost
	newP := &stepBuffer[depth].p
	newP.rooms[dstRoom][to] = wat
	newP.hall[hall] = EMPTY

	// Recurse
	return findLowestEnergy(newP, cost, depth)
}

func tryMoveRoomToRoom(p *puzzle, wat what, from, srcRoom, to, dstRoom int, cost int, depth int) result {
	nSpaces := checkPath(p, roomToJoinPoint[srcRoom], roomToJoinPoint[dstRoom])
	if nSpaces == 0 {
		return INVALID
	}
	moveRoomToRoom++

	// Add spaces to get out of the room.
	nSpaces += roomDepth - from // (2 - 0) = 2

	// Add spaces to get into the room.
	nSpaces += roomDepth - to // (2 - 0) = 2

	newCost := costs[wat] * nSpaces
	cost += newCost
	if cost >= minEnergy {
		return NO // no better solution this way
	}

	// Do the move
	depth++
	if depth > maxDepth {
		maxDepth = depth
		fmt.Println(maxDepth)
	}
	stepBuffer[depth].p = puzzle{
		rooms: p.rooms,
		hall:  p.hall,
	}
	stepBuffer[depth].desc = "room-to-room"
	stepBuffer[depth].cost = newCost
	newP := &stepBuffer[depth].p
	newP.rooms[srcRoom][from] = EMPTY
	newP.rooms[dstRoom][to] = wat

	// Recurse
	return findLowestEnergy(newP, cost, depth)
}

func checkPath(p *puzzle, src int, dst int) int {
	// Don't actually check src, just make sure all target spots are empty
	if src < dst {
		for i := src + 1; i <= dst; i++ {
			if p.hall[i] != EMPTY {
				return 0
			}
		}
		return dst - src
	} else if src > dst {
		for i := src - 1; i >= dst; i-- {
			if p.hall[i] != EMPTY {
				return 0
			}
		}
		return src - dst
	} else {
		panic("wat")
	}
}
