package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	file, err := os.Open("./input")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	type equation struct {
		target int
		input  []int
	}

	equations := []equation{}

	for scanner.Scan() {
		line := scanner.Text()

		split := strings.Split(line, ": ")

		target, err := strconv.Atoi(split[0])
		if err != nil {
			panic(err)
		}

		inputs := strings.Split(split[1], " ")
		inputsInt := []int{}
		for _, char := range inputs {
			asInt, err := strconv.Atoi(char)
			if err != nil {
				panic(err)
			}
			inputsInt = append(inputsInt, asInt)
		}

		equations = append(equations, equation{
			target: target,
			input:  inputsInt,
		})
	}

	sum := 0
	for _, equation := range equations {
		start := equation.input[0]
		isValid := false
		recursiveOperation(&isValid, equation.target, start, 1, equation.input)

		if isValid {
			sum += equation.target
		}
	}

	fmt.Println(sum)

	sum = 0
	for _, equation := range equations {
		start := equation.input[0]
		isValid := false
		recursiveOperationPart2(&isValid, equation.target, start, 1, equation.input)

		if isValid {
			sum += equation.target
		}
	}
	fmt.Println(sum)
}

func recursiveOperation(isValid *bool, target, current, i int, input []int) {
	if i == len(input) {
		if current == target {
			*isValid = true
		}
		return
	}

	recursiveOperation(isValid, target, current+input[i], i+1, input)
	recursiveOperation(isValid, target, current*input[i], i+1, input)
}

func recursiveOperationPart2(isValid *bool, target, current, i int, input []int) {
	if i == len(input) {
		if current == target {
			*isValid = true
		}
		return
	}

	recursiveOperationPart2(isValid, target, current+input[i], i+1, input)
	recursiveOperationPart2(isValid, target, current*input[i], i+1, input)
	recursiveOperationPart2(isValid, target, concatenateTwoInts(current, input[i]), i+1, input)
}

func concatenateTwoInts(a, b int) int {
	together := fmt.Sprintf("%d%d", a, b)
	asInt, err := strconv.Atoi(together)
	if err != nil {
		fmt.Println(together)
		panic(err)
	}
	return asInt
}
