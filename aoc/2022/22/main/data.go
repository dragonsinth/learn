package main

var data = input{
	width: 150,
	rsize: 50,
	// YMMV
	regions: [6]region{
		{id: 0, x: 50, y: 0, translate: [4]translation{{1, EAST}, {2, SOUTH}, {3, EAST}, {5, EAST}}},
		{id: 1, x: 100, y: 0, translate: [4]translation{{4, WEST}, {2, WEST}, {0, WEST}, {5, NORTH}}},
		{id: 2, x: 50, y: 50, translate: [4]translation{{1, NORTH}, {4, SOUTH}, {3, SOUTH}, {0, NORTH}}},
		{id: 3, x: 0, y: 100, translate: [4]translation{{4, EAST}, {5, SOUTH}, {0, EAST}, {2, EAST}}},
		{id: 4, x: 50, y: 100, translate: [4]translation{{1, WEST}, {5, WEST}, {3, WEST}, {2, NORTH}}},
		{id: 5, x: 0, y: 150, translate: [4]translation{{4, NORTH}, {1, SOUTH}, {0, SOUTH}, {3, NORTH}}},
	},
}
