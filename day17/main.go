package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"
)

type instruction func(registers []int, operand int, inPtr *int) (bool, int)

func main() {
	file, err := os.Open("./input")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	registers := []int{}
	program := []int{}
	parsingProgram := false
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			parsingProgram = true
			continue
		}

		if !parsingProgram {
			split := strings.Split(line, ": ")
			asInt, err := strconv.Atoi(split[1])
			if err != nil {
				panic(err)
			}
			registers = append(registers, asInt)
		} else {
			split := strings.Split(line, ": ")
			split2 := strings.Split(split[1], ",")
			for _, num := range split2 {
				asInt, err := strconv.Atoi(num)
				if err != nil {
					panic(err)
				}
				program = append(program, asInt)
			}
		}
	}

	originalRegisters := slices.Clone(registers)
	fmt.Println(registers)
	fmt.Println(program)
	inPtr := 0
	output := []int{}

	for inPtr != len(program) {
		opcode := program[inPtr]
		operand := program[inPtr+1]

		inPtrAdvances, out := instructionSet[opcode](registers, operand, &inPtr)
		if opcode == 5 {
			output = append(output, out)
		}
		if inPtrAdvances {
			inPtr += 2
		}
	}

	fmt.Println(output)

	outPutAsStr := []string{}
	for _, num := range output {
		outPutAsStr = append(outPutAsStr, fmt.Sprintf("%d", num))
	}

	fmt.Println(strings.Join(outPutAsStr, ","))

	// Part 2 //////////////////////////////////////
	for a := 0; a < 8; a++ {
		for b := 0; b < 8; b++ {
			for c := 0; c < 8; c++ {
				for d := 0; d < 8; d++ {
					for e := 0; e < 8; e++ {
						octal := []int{1, 0, 3, 5, 5, 1, 0, 0, 0, 5, 1, a, b, c, d, e} // Got these numbers by manually testing
						output = []int{}

						octalAsStr := []string{}
						for _, num := range octal {
							octalAsStr = append(octalAsStr, fmt.Sprintf("%d", num))
						}

						asOneString := strings.Join(octalAsStr, "")
						asDec, err := strconv.ParseInt(asOneString, 8, 64)
						if err != nil {
							panic(err)
						}
						i := int(asDec)
						// fmt.Println(i)

						inPtr = 0
						output = []int{}
						outputCounter := 0
						shouldBreak := false
						registers = slices.Clone(originalRegisters)
						registers[0] = i
						// fmt.Println("trying register A as", i)

						// fmt.Println(fmt.Sprintf("%o", i))

						// fmt.Println(program)

						for inPtr != len(program) {
							opcode := program[inPtr]
							operand := program[inPtr+1]

							inPtrAdvances, out := instructionSet[opcode](registers, operand, &inPtr)
							if opcode == 5 {
								output = append(output, out)
								if out != program[outputCounter] {
									// shouldBreak = true
									// fmt.Println("no match", program, out)
								}
								outputCounter++
							}
							if inPtrAdvances {
								inPtr += 2
							}
							if shouldBreak {
								break
							}
						}
						// fmt.Println("output", output)
						if slices.Equal(program, output) {
							fmt.Println(i)
						}
					}
				}
			}
		}
	}

	// fmt.Println(output)
}

func GetOperand(registers []int, combo bool, operand int) int {
	if !combo {
		return operand
	}

	if operand < 4 {
		return operand
	}

	if operand == 7 {
		panic("combo operand 7 is illegal")
	}

	return registers[operand-4]
}

// Instructions // Return value indicates if we should move forward 2 instructions
// Return int indicates an output
func adv(registers []int, operand int, inPtr *int) (bool, int) {
	numerator := float64(registers[0])
	denominator := math.Pow(2, float64(GetOperand(registers, true, operand)))
	division := int(numerator / denominator)
	registers[0] = division
	return true, 0
}

func bxl(registers []int, operand int, inPtr *int) (bool, int) {
	left := registers[1]
	right := GetOperand(registers, false, operand)
	xor := left ^ right
	registers[1] = xor
	return true, 0
}

func bst(registers []int, operand int, inPtr *int) (bool, int) {
	registers[1] = GetOperand(registers, true, operand) % 8
	return true, 0
}

func jnz(registers []int, operand int, inPtr *int) (bool, int) {
	if registers[0] == 0 {
		return true, 0
	}
	*inPtr = GetOperand(registers, false, operand)
	return false, 0
}

func bxc(registers []int, operand int, inPtr *int) (bool, int) {
	registers[1] = registers[1] ^ registers[2]
	return true, 0
}

func out(registers []int, operand int, inPtr *int) (bool, int) {
	output := GetOperand(registers, true, operand) % 8
	return true, output
}

func bdv(registers []int, operand int, inPtr *int) (bool, int) {
	numerator := float64(registers[0])
	denominator := math.Pow(2, float64(GetOperand(registers, true, operand)))
	division := int(numerator / denominator)
	registers[1] = division
	return true, 0
}

func cdv(registers []int, operand int, inPtr *int) (bool, int) {
	numerator := float64(registers[0])
	denominator := math.Pow(2, float64(GetOperand(registers, true, operand)))
	division := int(numerator / denominator)
	registers[2] = division
	return true, 0
}

var instructionSet []instruction = []instruction{
	adv, bxl, bst, jnz, bxc, out, bdv, cdv,
}
