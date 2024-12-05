package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

// Just gonna make a guess that I'll need a linked list or some dynamic programming after this.
// Maybe I can just generate it better later.. dont preoptimize lel.
type pageRule struct {
	value        int
	mustBeAfter  map[int]bool
	mustBeBefore map[int]bool
}

var bigMapofRules map[int]*pageRule

func main() {
	file, err := os.Open("./input")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	rules := []string{}
	manuals := []string{}

	manualScan := false

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" { // Swap to manual page numbers after.
			manualScan = true
			continue
		}

		if !manualScan {
			rules = append(rules, line)
		} else {
			manuals = append(manuals, line)
		}
	}

	bigMapofRules = map[int]*pageRule{}

	for _, rule := range rules {
		split := strings.Split(rule, "|")
		pg1, err := strconv.Atoi(split[0])
		if err != nil {
			panic(err)
		}
		pg2, err := strconv.Atoi(split[1])
		if err != nil {
			panic(err)
		}

		rule, ok := bigMapofRules[pg1]
		if !ok {
			bigMapofRules[pg1] = &pageRule{
				value:       pg1,
				mustBeAfter: map[int]bool{},
				mustBeBefore: map[int]bool{
					pg2: true,
				},
			}
		} else {
			rule.mustBeBefore[pg2] = true
		}

		// Other way too
		rule, ok = bigMapofRules[pg2]
		if !ok {
			bigMapofRules[pg2] = &pageRule{
				value: pg2,
				mustBeAfter: map[int]bool{
					pg1: true,
				},
				mustBeBefore: map[int]bool{},
			}
		} else {
			rule.mustBeAfter[pg1] = true
		}
	}

	// fmt.Println(bigMapofRules)

	// Convert manual page ordering to an int array
	manualsAsInts := [][]int{}
	for _, manual := range manuals {
		manualsAsInts = append(manualsAsInts, convertManualStringArrayToIntArray(manual))
	}

	numValid := 0
	sum := 0

	// Part 2 init
	invalidManuals := [][]int{}

	for _, manual := range manualsAsInts {
		valid := isManualValid(manual)
		if valid {
			numValid++
			sum += manual[(len(manual)-1)/2] // Grab the middle number as per problem statement
		} else {
			invalidManuals = append(invalidManuals, manual)
		}
	}

	fmt.Println(sum)

	sum2 := 0
	fixedManuals := [][]int{}

	// fmt.Println(invalidManuals)
	for _, manual := range invalidManuals {
		fixedManuals = append(fixedManuals, fixInvalidManual(manual))
	}

	for _, manual := range fixedManuals {
		sum2 += manual[(len(manual)-1)/2]
	}

	fmt.Println(sum2)
	// for k, v := range bigMapofRules {
	// 	fmt.Println(k, *&v.mustBeBefore)
	// }
}

func isManualValid(manual []int) bool {
	valid := true
	for i, num := range manual {
		for _, num2 := range manual[i+1:] { // Check the forward list
			rule, ok := bigMapofRules[num2]
			if ok {
				if rule.mustBeBefore[num] {
					valid = false
					break
				}
			}
		}
		if !valid {
			break
		}
	}
	return valid
}

func convertManualStringArrayToIntArray(manual string) []int {
	split := strings.Split(manual, ",")
	out := []int{}
	for _, entry := range split {
		num, err := strconv.Atoi(entry)
		if err != nil {
			panic(err)
		}
		out = append(out, num)
	}
	return out
}

func fixInvalidManual(invalidManual []int) []int {
	mustBeBefores := map[int]int{}
	for i, num1 := range invalidManual {
		mustBeBefore := 0
		for j, num2 := range invalidManual {
			if i == j {
				continue
			}
			rule, ok := bigMapofRules[num1]
			if ok {
				if rule.mustBeBefore[num2] {
					mustBeBefore++
				}
			}
		}
		mustBeBefores[num1] = mustBeBefore
	}

	// fmt.Println(mustBeBefores)

	sort.Slice(invalidManual, func(i, j int) bool {
		return mustBeBefores[invalidManual[i]] > mustBeBefores[invalidManual[j]]
	})

	return invalidManual
}
