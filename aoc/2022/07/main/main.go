package main

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
)

var sample = `
$ cd /
$ ls
dir a
14848514 b.txt
8504156 c.dat
dir d
$ cd a
$ ls
dir e
29116 f
2557 g
62596 h.lst
$ cd e
$ ls
584 i
$ cd ..
$ cd ..
$ cd d
$ ls
4060174 j
8033020 d.log
5626152 d.ext
7214296 k
`

func main() {
	dirs := map[string]bool{}
	files := map[string]int{}
	path := "/"
	dirs[path] = true
	for _, line := range strings.Split(sample, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "$ cd ") {
			where := strings.TrimPrefix(line, "$ cd ")
			if where == "/" {
				path = "/"
			} else if where == ".." {
				path = filepath.Dir(path)
			} else {
				path = filepath.Join(path, where)
			}
			dirs[path] = true
		} else if line == "$ ls" {

		} else if strings.HasPrefix(line, "dir ") {

		} else {
			parts := strings.Split(line, " ")
			if len(parts) != 2 {
				panic(parts)
			}

			sz, name := mustInt(parts[0]), parts[1]
			files[filepath.Join(path, name)] = sz
		}
	}

	// For each dir, sum files
	outer := 0
	used := 0
	for d := range dirs {
		sum := 0
		for f, sz := range files {
			if strings.HasPrefix(f, d) {
				sum += sz
			}
		}
		fmt.Println(d, sum)
		if d == "/" {
			used = sum
		}
		if sum < 100000 {
			outer += sum
		}
	}
	fmt.Println("part 1:", outer)

	need := 30000000 - (70000000 - used)
	fmt.Println("need:", need)

	bestSize := 70000000
	for d := range dirs {
		sum := 0
		for f, sz := range files {
			if strings.HasPrefix(f, d) {
				sum += sz
			}
		}
		if sum > need && sum < bestSize {
			bestSize = sum
		}
	}
	fmt.Println(bestSize)
}

func mustInt(s string) int {
	if v, err := strconv.Atoi(s); err != nil {
		panic(fmt.Sprint(s, err))
	} else {
		return v
	}
}
