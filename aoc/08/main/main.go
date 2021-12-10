package main

import (
	"fmt"
	"sort"
	"strings"
)

var sample = `
be cfbegad cbdgef fgaecd cgeb fdcge agebfd fecdb fabcd edb | fdgacbe cefdb cefbgd gcbe
edbfga begcd cbg gc gcadebf fbgde acbgfd abcde gfcbed gfec | fcgedb cgb dgebacf gc
fgaebd cg bdaec gdafb agbcfd gdcbef bgcad gfac gcb cdgabef | cg cg fdcagb cbg
fbegcd cbd adcefb dageb afcb bc aefdc ecdab fgdeca fcdbega | efabcd cedba gadfec cb
aecbfdg fbg gf bafeg dbefa fcge gcbea fcaegb dgceab fcbdga | gecf egdcabf bgf bfgea
fgeab ca afcebg bdacfeg cfaedg gcfdb baec bfadeg bafgc acf | gebdcfa ecba ca fadegcb
dbcfg fgd bdegcaf fgec aegbdf ecdfab fbedc dacgb gdcebf gf | cefg dcbef fcge gbcadfe
bdfegc cbegaf gecbf dfcage bdacg ed bedf ced adcbefg gebcd | ed bcgafe cdgba cbgef
egadfb cdbfeg cegd fecab cgb gbdefca cg fgcdab egfdb bfceg | gbdfcae bgc cg cgb
gcafb gcf dcaebfg ecagb gf abcdeg gaef cafbge fdbac fegbdc | fgae cfgab fg bagce
`

var input = sample

//  0000
// 1    2
// 1    2
//  3333
// 4    5
// 4    5
//  6666

var digitsToSegments = [10][]int{
	0: {0, 1, 2, 4, 5, 6},
	1: {2, 5},
	2: {0, 2, 3, 4, 6},
	3: {0, 2, 3, 5, 6},
	4: {1, 2, 3, 5},
	5: {0, 1, 3, 5, 6},
	6: {0, 1, 3, 4, 5, 6},
	7: {0, 2, 5},
	8: {0, 1, 2, 3, 4, 5, 6},
	9: {0, 1, 2, 3, 5, 6},
}

type permuation struct {
	codesByDigit [10]string
	digitsByCode map[string]int

	key string // lexically sorted codes
}

func main() {
	// Generate all possible input permutations.
	permsByInput := map[string]*permuation{}
	permsByKey := map[string]*permuation{}
	generatePerms("abcdefg", func(s string) {
		p := createPermuation(s)
		permsByInput[s] = p
		permsByKey[p.key] = p
	})

	// Analyze the input.
	sum := 0
	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.Split(line, " | ")
		if len(parts) != 2 {
			panic(parts)
		}

		signals := strings.Split(parts[0], " ")
		if len(signals) != 10 {
			panic(signals)
		}

		// Map each signal to a permutation.
		for i := range signals {
			signals[i] = lexicalSort(signals[i])
		}
		key := asKey(signals)
		p, ok := permsByKey[key]
		if !ok {
			panic(key)
		}

		// Now that we know what perm we're on, we can compute the digits.
		digitCodes := strings.Split(parts[1], " ")
		if len(digitCodes) != 4 {
			panic(digitCodes)
		}

		val := 0
		for i := 0; i < 4; i++ {
			val *= 10
			digit := p.digitsByCode[lexicalSort(digitCodes[i])]
			val += digit
		}
		fmt.Println(val)
		sum += val
	}

	fmt.Println(sum)
}

func generatePerms(s string, accept func(string)) {
	in := []byte(s)
	generatePermsIn(in, 0, len(in)-1, accept)
}
func generatePermsIn(in []byte, left, right int, accept func(string)) {
	if left == right {
		accept(string(in))
	} else {
		for i := left; i <= right; i++ {
			in[left], in[i] = in[i], in[left]
			generatePermsIn(in, left+1, right, accept)
			in[left], in[i] = in[i], in[left]
		}
	}
}

func createPermuation(s string) *permuation {
	ret := &permuation{
		codesByDigit: [10]string{},
		digitsByCode: map[string]int{},
		key:          "",
	}

	letters := []byte(s)
	var keys []string
	for i := 0; i < 10; i++ {
		// Compute the code for this digit.
		var codeBytes []byte
		for _, seg := range digitsToSegments[i] {
			codeBytes = append(codeBytes, letters[seg])
		}

		code := lexicalSort(string(codeBytes))
		ret.codesByDigit[i] = code
		ret.digitsByCode[code] = i

		keys = append(keys, code)
	}

	ret.key = asKey(keys)
	return ret
}

func lexicalSort(s string) string {
	// Sort the code lexically.
	b := []byte(s)
	sort.Slice(b, func(i, j int) bool {
		return b[i] < b[j]
	})
	return string(b)
}

func asKey(s []string) string {
	sort.Strings(s)
	return strings.Join(s, " ")
}

func main1() {
	count := 0
	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.Split(line, " | ")
		if len(parts) != 2 {
			panic(parts)
		}

		outputs := strings.Split(parts[1], " ")
		if len(outputs) != 4 {
			panic(outputs)
		}

		for _, output := range outputs {
			switch len(output) {
			case 2:
				//fmt.Println(1)
				count++
			case 3:
				//fmt.Println(7)
				count++
			case 4:
				//fmt.Println(4)
				count++
			case 7:
				//fmt.Println(8)
				count++
			}
		}
	}

	fmt.Println(count)
}
