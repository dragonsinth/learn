package main

import (
	"bufio"
	"fmt"
	"github.com/dragonsinth/learn/aoc/2019/intcode"
	"os"
)

const script = `north
north
take sand
south
south
south
west
take wreath
south
south
take pointer
north
north
east
north
west
south
take planetoid
north
west
south
west
inv
north
`

func main() {
	codes := intcode.Parse(`99`)

	scanner := bufio.NewScanner(os.Stdin)

	inbuf := []byte(script)
	var outbuf []byte
	m := intcode.NewIntMachine(codes, func() int {
		if len(inbuf) == 0 {
			if scanner.Scan() {
				inbuf = scanner.Bytes()
				inbuf = append(inbuf, '\n')
			} else if err := scanner.Err(); err != nil {
				panic(err)
			} else {
				os.Exit(0)
			}
		}
		r := inbuf[0]
		inbuf = inbuf[1:]
		return int(r)
	}, func(v int) {
		if v == '\n' {
			fmt.Println(string(outbuf))
			outbuf = nil
		} else {
			outbuf = append(outbuf, byte(v))
		}
	})
	m.Run()
}
