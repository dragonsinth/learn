package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var sample = `
ecl:gry pid:860033327 eyr:2020 hcl:#fffffd
byr:1937 iyr:2017 cid:147 hgt:183cm

iyr:2013 ecl:amb cid:350 eyr:2023 pid:028048884
hcl:#cfa07d byr:1929

hcl:#ae17e1 iyr:2013
eyr:2024
ecl:brn pid:760753108 byr:1931
hgt:179cm

hcl:#cfa07d eyr:2025 pid:166559648
iyr:2011 ecl:brn hgt:59in
`

var (
	input = sample
)

func main() {
	cur := map[string]string{}
	valid := 0

	check := func() {
		fmt.Println(cur)
		if !validYear(cur["byr"], 1920, 2002) {
			return
		}
		if !validYear(cur["iyr"], 2010, 2020) {
			return
		}
		if !validYear(cur["eyr"], 2020, 2030) {
			return
		}
		if !validHeight(cur["hgt"]) {
			return
		}
		if !validHair(cur["hcl"]) {
			return
		}
		if !validEye(cur["ecl"]) {
			return
		}
		if !validPid(cur["pid"]) {
			return
		}
		valid++
	}

	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			check()
			cur = map[string]string{}
			continue
		}

		parts := strings.Split(line, " ")
		for _, p := range parts {
			kv := strings.Split(p, ":")
			cur[kv[0]] = kv[1]
		}

	}
	check()
	fmt.Println(valid)
}

func validYear(s string, min, max int) bool {
	v, err := strconv.Atoi(s)
	return err == nil && v >= min && v <= max
}

func validHeight(s string) bool {
	if strings.HasSuffix(s, "cm") {
		v, err := strconv.Atoi(strings.TrimSuffix(s, "cm"))
		return err == nil && v >= 150 && v <= 193
	} else if strings.HasSuffix(s, "in") {
		v, err := strconv.Atoi(strings.TrimSuffix(s, "in"))
		return err == nil && v >= 59 && v <= 76
	} else {
		return false
	}
}

func validHair(s string) bool {
	return regexp.MustCompile(`^#[0-9a-f]{6}$`).MatchString(s)
}

func validEye(s string) bool {
	return map[string]bool{
		"amb": true,
		"blu": true,
		"brn": true,
		"gry": true,
		"grn": true,
		"hzl": true,
		"oth": true,
	}[s]

}

func validPid(s string) bool {
	return regexp.MustCompile(`^[0-9]{9}$`).MatchString(s)
}
