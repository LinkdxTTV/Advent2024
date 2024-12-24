package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"slices"
	"strconv"
	"strings"
)

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

	for _, gate := range gates {
		registerDependencyMap[gate.outputRegister] = []string{gate.leftRegister, gate.rightRegister}
	}

	registerWasWrong := map[int]int{}
	cleanRegisterCopy := makeRegisterCopy(register)

	for i := 0; i < 1000; i++ {
		register = makeRegisterCopy(cleanRegisterCopy)
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
		// fmt.Println("x  ", x)
		// fmt.Println("y  ", y)
		// fmt.Println("x+y", x+y)
		// fmt.Println("z  ", z)
		// fmt.Println("-- as binary --")
		// fmt.Println("x  ", fmt.Sprintf("%b", x))
		// fmt.Println("y  ", fmt.Sprintf("%b", y))
		// fmt.Println("x+y", fmt.Sprintf("%b", x+y))
		// fmt.Println("z  ", fmt.Sprintf("%b", z))

		// Find out which registers are wrong
		for i, char := range fmt.Sprintf("%b", x+y) {
			if string(fmt.Sprintf("%b", z)[i]) != string(char) {
				registerWasWrong[i]++
			}
		}
	}
	type wrongRegister struct {
		register int
		wrong    int
	}

	wrongRegisters := []wrongRegister{}
	for k, v := range registerWasWrong {
		wrongRegisters = append(wrongRegisters, wrongRegister{register: k, wrong: v})
	}

	slices.SortFunc(wrongRegisters, func(a, b wrongRegister) int {
		return b.wrong - a.wrong
	})

	fmt.Println(wrongRegisters)

	for _, wrongRegister := range wrongRegisters {
		for i := 0; i < wrongRegister.wrong; i++ {
			recursivelyAddBadRegistersToMap(fmt.Sprintf("z%d", wrongRegister.register), registerDependencyMap, 1)
		}
	}

	type MentionedRegister struct {
		register string
		mentions int
	}
	mentionedRegistersList := []MentionedRegister{}
	for k, v := range mentionedRegister {
		mentionedRegistersList = append(mentionedRegistersList, MentionedRegister{
			register: k,
			mentions: v,
		})
	}

	slices.SortFunc(mentionedRegistersList, func(a, b MentionedRegister) int {
		return b.mentions - a.mentions
	})

	registersTaken := 0
	swappableList := []string{}
	for _, poss := range mentionedRegistersList {
		if registersTaken == 12 {
			break
		}
		if strings.HasPrefix(poss.register, "x") || strings.HasPrefix(poss.register, "y") {
			continue
		}
		fmt.Println(poss.register, poss.mentions)
		swappableList = append(swappableList, poss.register)
		registersTaken++
	}

	fmt.Println(swappableList)
	indexToGateName := map[int]string{}
	swappableListAsGateIndex := []int{}
	for i, gate := range gates {
		if slices.Contains(swappableList, gate.outputRegister) {
			swappableListAsGateIndex = append(swappableListAsGateIndex, i)
			indexToGateName[i] = gate.outputRegister
		}
	}
	fmt.Println(swappableListAsGateIndex)

	for a, order := range permutations(swappableListAsGateIndex) {
		if a%100 == 0 {
			fmt.Println(a)
		}
		swaps := order[:8]
		gateCopy := slices.Clone(gates)

		for n := 0; n < 4; n++ {
			gateCopy[swaps[2*n]], gateCopy[swaps[2*n+1]] = swapGates(gates[swaps[2*n]], gates[swaps[2*n+1]])
		}
		register = makeRegisterCopy(cleanRegisterCopy)
		works := TestGateConfigurationGood(gateCopy)
		if works {
			output := []string{}
			for _, num := range swaps {
				output = append(output, indexToGateName[num])
			}

			slices.Sort(output)
			fmt.Println(strings.Join(output, ","))
		}
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

		for i, char := range fmt.Sprintf("%b", x+y) {
			if string(fmt.Sprintf("%b", z)[i]) == string(char) {
				score++
			}
		}
	}
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
