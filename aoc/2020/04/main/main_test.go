package main

import "testing"

func TestHair(t* testing.T) {
	t.Log(validHair(`#123abc`))
	t.Log(validHair(`#123abz`))
	t.Log(validHair(`123abc`))
	t.Log(validHair(`#123ab`))
}
