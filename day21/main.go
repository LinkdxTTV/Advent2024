package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type point struct {
	x int
	y int
}

var numpad [][]string = [][]string{
	{"7", "8", "9"},
	{"4", "5", "6"},
	{"1", "2", "3"},
	{"X", "0", "A"},
}

var keypad [][]string = [][]string{
	{"X", "^", "A"},
	{"<", "v", ">"},
}

func main() {
	file, err := os.Open("./input")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	directions := []string{}

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		directions = append(directions, line)
	}

	fmt.Println(directions)

	sum := 0

	for _, direction := range directions {
		fmt.Println(direction)
		robotDirections := GenerateDirectionsOnPad(direction, "A", numpad)
		fmt.Println(robotDirections)
		for i := 0; i < 2; i++ {
			robotDirections = GenerateDirectionsOnPad(robotDirections, "A", keypad)
			fmt.Println(robotDirections)
		}

		sequenceLength := len(robotDirections)
		numericCode := calculateNumericCode(direction)

		sum += sequenceLength * numericCode
		fmt.Println(sequenceLength, numericCode)
	}

	fmt.Println(sum)

	for i, direction := range []string{">>v", "v>>", "^>>", ">>^", "<<^", "^<<", "<<v", "v<<"} {
		if i%2 == 0 {
			fmt.Println()
		}
		robot := direction
		for i := 0; i < 6; i++ {
			robot = GenerateDirectionsOnPad(robot, "A", keypad)
		}
		fmt.Println(direction, len(robot))
	}
}

func calculateNumericCode(input string) int {
	justNumbers := strings.TrimSuffix(input, "A")
	asInt, err := strconv.Atoi(justNumbers)
	if err != nil {
		panic(err)
	}
	return asInt
}

func GenerateDirectionsOnPad(input string, start string, pad [][]string) string {
	if input == "" {
		return ""
	}
	outString := GeneratePath(start, string(input[0]), pad)
	outString += "A"
	for i := 0; i < len(input)-1; i++ {
		j := i + 1
		outString += GeneratePath(string(input[i]), string(input[j]), pad)
		outString += "A"
	}
	return outString
}

// Does not include A's
func GeneratePath(start, end string, pad [][]string) string {
	// Find startPoint and endPoint
	var startPoint, endPoint = point{}, point{}
	for y, row := range pad {
		for x, char := range row {
			if char == start {
				startPoint = point{x, y}
			}
			if char == end {
				endPoint = point{x, y}
			}
		}
	}

	// Diff them
	diffX := endPoint.x - startPoint.x
	diffY := endPoint.y - startPoint.y

	// If moving right, we do rights first. If moving left, we do verticals first
	outStringX := ""
	outStringY := ""

	if diffX >= 0 {
		for i := 0; i < diffX; i++ {
			outStringX += ">"
		}
	} else {
		for i := 0; i < -diffX; i++ {
			outStringX += "<"
		}
	}

	if diffY >= 0 {
		for i := 0; i < diffY; i++ {
			outStringY += "v"
		}
	} else {
		for i := 0; i < -diffY; i++ {
			outStringY += "^"
		}
	}
	outString := ""

	// This is a bit tricky. You should always prefer the ordering in the following order: right, up, down, left
	// But we need to make sure we never hit the panic space
	if diffX >= 0 && diffY >= 0 { // Right and Down // Go right first
		outString += outStringX
		outString += outStringY
	} else if diffX >= 0 && diffY <= 0 { // Right and Up// Go up first
		if start == "<" {
			outString += outStringX
			outString += outStringY
		} else {
			outString += outStringY
			outString += outStringX
		}
	} else if diffX <= 0 && diffY >= 0 { // Left and Down || Go left first
		if end == "<" {
			outString += outStringY
			outString += outStringX
		} else {
			outString += outStringX
			outString += outStringY
		}
	} else if diffX <= 0 && diffY <= 0 { // Left and Up || Go up first
		outString += outStringY
		outString += outStringX

	}

	return outString
}
