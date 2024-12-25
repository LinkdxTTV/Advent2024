package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"slices"
	"strconv"
	"strings"

	"gonum.org/v1/gonum/stat/combin"
)

// READ ME

// For part 2, I did the problem on pen and paper by analyzing how to form Z01.. Z02.. ETC. They form an adder circuit which has a predictable pattern. From there, it is easy to see where things go wrong
/*
		Z##
	aaa XOR bbb
 X XOR Y   ccc      OR      ddd
		 X-1 XOR Y-1  Z-1 AND Z-1 (These are the previous two children of Z-1)

		 etc
*/

type gate struct {
	leftRegister   string
	rightRegister  string
	instruction    string
	outputRegister string
}

var register = map[string]*int{}
var registerDependencyMap = map[string][]string{}

var one int = 1
var zero int = 0

var instructionMap = map[string]func(a, b *int) *int{
	"AND": func(a, b *int) *int {
		if a == nil || b == nil {
			return nil
		}
		if *a == 1 && *b == 1 {
			return &one
		}
		return &zero
	},
	"OR": func(a, b *int) *int {
		if a == nil || b == nil {
			return nil
		}
		if *a == 0 && *b == 0 {
			return &zero
		}
		return &one
	},
	"XOR": func(a, b *int) *int {
		if a == nil || b == nil {
			return nil
		}
		if *a == *b {
			return &zero
		} else {
			return &one
		}
	},
}

func main() {
	file, err := os.Open("./input")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	parsingGates := false

	gates := []gate{}

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			parsingGates = true
			continue
		}
		if !parsingGates {
			split := strings.Split(line, ": ")
			asInt, err := strconv.Atoi(split[1])
			if err != nil {
				panic(err)
			}
			register[split[0]] = &asInt
		} else {
			split := strings.Split(line, " -> ")
			leftSide := strings.Split(split[0], " ")
			// Right side
			if _, ok := register[split[1]]; !ok {
				register[split[1]] = nil
			}
			// Left side
			if _, ok := register[leftSide[0]]; !ok {
				register[leftSide[0]] = nil
			}
			if _, ok := register[leftSide[2]]; !ok {
				register[leftSide[2]] = nil
			}
			gates = append(gates, gate{
				leftRegister:   leftSide[0],
				rightRegister:  leftSide[2],
				instruction:    leftSide[1],
				outputRegister: split[1],
			})
		}
	}

	mapFromIndexToGateString := map[int]string{}
	for i, gate := range gates {
		mapFromIndexToGateString[i] = gate.outputRegister
	}

	type similarity struct {
		swap  []int
		score int
	}

	scores := []similarity{}
	for i, swap := range combin.Combinations(len(gates), 2) {
		if i%1000 == 0 {
			fmt.Println(i)
		}
		gateCopy := slices.Clone(gates)
		gateCopy[swap[0]], gateCopy[swap[1]] = swapGates(gates[swap[0]], gates[swap[1]])

		score := GetGateConfigurationScore(gateCopy)
		scores = append(scores, similarity{
			swap:  swap,
			score: score,
		})
	}

	slices.SortFunc(scores, func(a, b similarity) int {
		return b.score - a.score
	})

	fmt.Println(scores[:30])
	gateCopy := slices.Clone(gates)
	for _, poss := range scores[0:4] {
		swap := poss.swap
		gateCopy[swap[0]], gateCopy[swap[1]] = swapGates(gates[swap[0]], gates[swap[1]])
	}
	works := TestGateConfigurationGood(gateCopy)
	if works {
		out := []string{}
		for _, poss := range scores[0:4] {
			out = append(out, mapFromIndexToGateString[poss.swap[0]])
			out = append(out, mapFromIndexToGateString[poss.swap[1]])
		}
		slices.Sort(out)
		fmt.Println(strings.Join(out, ","))
	}
}

var mentionedRegister = map[string]int{}

func recursivelyAddBadRegistersToMap(register string, dependencyMap map[string][]string, modifier int) {
	mentionedRegister[register] += modifier
	for _, child := range dependencyMap[register] {
		recursivelyAddBadRegistersToMap(child, dependencyMap, modifier)
	}
}

func giveRegisterString(x int, prefix string) string {
	if x >= 10 {
		return fmt.Sprintf("%s%d", prefix, x)
	} else {
		return fmt.Sprintf("%s0%d", prefix, x)
	}
}

func returnNumberFromRegistersStartingWith(prefix string) int {
	bits := ""
	for z := 100; z >= 0; z-- {
		out, ok := register[giveRegisterString(z, prefix)]
		if ok {
			if out != nil {
				bits += fmt.Sprintf("%d", *out) // Will panic if z register not set
			} else {
				bits += "0"
			}
		}
	}

	decimal, err := strconv.ParseInt(bits, 2, 64)
	if err != nil {
		panic(err)
	}
	return int(decimal)
}

func returnBinaryFromRegistersStartingWith(prefix string) string {
	bits := ""
	for z := 100; z >= 0; z-- {
		out, ok := register[giveRegisterString(z, prefix)]
		if ok {
			bits += fmt.Sprintf("%d", *out) // Will panic if z register not set
		}
	}

	return bits
}

func swapGates(gate1, gate2 gate) (gate, gate) {
	temp := gate1
	gate1.outputRegister = gate2.outputRegister
	gate2.outputRegister = temp.outputRegister
	return gate1, gate2
}

func overwriteRegistersThatStartWithPrefixToValue(prefix string, value *int) {
	for z := 0; z <= 100; z++ {
		_, ok := register[giveRegisterString(z, prefix)]
		if ok {
			register[giveRegisterString(z, prefix)] = value
		}
	}
}

func makeRegisterCopy(in map[string]*int) map[string]*int {
	out := map[string]*int{}
	for k, v := range in {
		out[k] = v
	}
	return out
}

func randomizeRegistersThatStartWithPrefix(prefix string) {
	for z := 0; z <= 100; z++ {
		_, ok := register[giveRegisterString(z, prefix)]
		if ok {

			if rand.Intn(2) == 1 {
				register[giveRegisterString(z, prefix)] = &one
			} else {
				register[giveRegisterString(z, prefix)] = &zero
			}
		}
	}
}

func TestGateConfigurationGood(gates []gate) bool {
	for i := 0; i < 10; i++ {
		randomizeRegistersThatStartWithPrefix("x")
		randomizeRegistersThatStartWithPrefix("y")
		for i := 0; i < len(gates); i++ {
			for _, gate := range gates {
				register[gate.outputRegister] = instructionMap[gate.instruction](register[gate.leftRegister], register[gate.rightRegister])
			}
		}
		x := returnNumberFromRegistersStartingWith("x")
		y := returnNumberFromRegistersStartingWith("y")
		z := returnNumberFromRegistersStartingWith("z")
		if z != x+y {
			return false
		}
	}
	return true
}

func GetGateConfigurationScore(gates []gate) int {
	score := 0
	for i := 0; i < 5; i++ {
		randomizeRegistersThatStartWithPrefix("x")
		randomizeRegistersThatStartWithPrefix("y")
		for i := 0; i < len(gates); i++ {
			for _, gate := range gates {
				register[gate.outputRegister] = instructionMap[gate.instruction](register[gate.leftRegister], register[gate.rightRegister])
			}
		}
		x := returnNumberFromRegistersStartingWith("x")
		y := returnNumberFromRegistersStartingWith("y")
		z := returnNumberFromRegistersStartingWith("z")

		xyAsBin := fmt.Sprintf("%b", x+y)
		zAsBin := fmt.Sprintf("%b", z)
		// fmt.Println(xyAsBin, zAsBin, len(xyAsBin), len(zAsBin))
		if len(xyAsBin) != len(zAsBin) {
			continue
		}
		for i, char := range zAsBin {
			if string(xyAsBin[i]) == string(char) {
				score++
			}
		}
	}
	// fmt.Println(score)
	return score
}

func permutations(n, r int) [][]int {
	if r > n {
		return [][]int{}
	}

	result := [][]int{}
	elements := make([]int, n)
	for i := 0; i < n; i++ {
		elements[i] = i + 1
	}

	generatePermutations(elements, r, 0, []int{}, &result)
	return result
}

func generatePermutations(elements []int, r, index int, current []int, result *[][]int) {
	if index == r {
		*result = append(*result, append([]int{}, current...))
		return
	}

	for i := 0; i < len(elements); i++ {
		if !contains(current, elements[i]) {
			generatePermutations(elements, r, index+1, append(current, elements[i]), result)
		}
	}
}

func contains(arr []int, val int) bool {
	for _, v := range arr {
		if v == val {
			return true
		}
	}
	return false
}
