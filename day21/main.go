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

	fmt.Println("====Part 1====")
	maxDepth := 3
	sum := 0

	for _, direction := range directions {
		length := GetSequenceLength(direction, maxDepth, maxDepth)
		fmt.Println(length, calculateNumericCode(direction))
		sum += length * calculateNumericCode(direction)
	}

	fmt.Println(sum)

	// Part 2
	fmt.Println("====Part 2====")
	maxDepth = 26
	sum = 0

	for _, direction := range directions {
		length := GetSequenceLength(direction, maxDepth, maxDepth)
		fmt.Println(length, calculateNumericCode(direction))
		sum += length * calculateNumericCode(direction)
	}

	fmt.Println(sum)

}

type sequenceDepth struct {
	sequence string
	depth    int
}

var sequenceMemo = map[sequenceDepth]int{}

func GetSequenceLength(sequence string, depth, maxDepth int) int {
	if out, ok := sequenceMemo[sequenceDepth{sequence, depth}]; ok {
		return out
	}

	pad := [][]string{}
	if depth == maxDepth {
		pad = numpad
	} else {
		pad = keypad
	}

	length := 0
	if depth == 0 {
		length = len(sequence)
	} else {
		currentLetter := "A"
		for _, nextLetter := range sequence {
			lengthOfMove := getMoves(currentLetter, string(nextLetter), pad, depth, maxDepth)
			currentLetter = string(nextLetter)
			length += lengthOfMove
		}
	}

	sequenceMemo[sequenceDepth{sequence, depth}] = length
	return length
}

func getMoves(currentLetter, nextLetter string, pad [][]string, depth, maxDepth int) int {
	if currentLetter == nextLetter {
		return 1 // No need to move if the letters are the same, just press again
	}

	return GetSequenceLength(GeneratePath(currentLetter, nextLetter, pad), depth-1, maxDepth)
}

func calculateNumericCode(input string) int {
	justNumbers := strings.TrimSuffix(input, "A")
	asInt, err := strconv.Atoi(justNumbers)
	if err != nil {
		panic(err)
	}
	return asInt
}

var pathMemo = map[string]string{}

// Generate the string that would get us from start to end, including the A at the end
func GeneratePath(start, end string, pad [][]string) string {
	if out, ok := pathMemo[fmt.Sprintf("%s%s", start, end)]; ok {
		return out
	}
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

	// This is a bit tricky. Depending on how far you are from the human input, the preference changes:

	// If there are >= 2 iterations left,
	// But we need to make sure we never hit the panic space
	if diffX >= 0 && diffY >= 0 { // Right and Down // Go down first
		if (start == "1" || start == "4" || start == "7") && (end == "0" || end == "A") { // Go right first
			outString += outStringX
			outString += outStringY
		} else {
			outString += outStringY
			outString += outStringX

		}
	} else if diffX >= 0 && diffY <= 0 { // Right and Up// Go up first
		if start == "<" {
			outString += outStringX
			outString += outStringY
		} else {
			outString += outStringY
			outString += outStringX
		}
	} else if diffX <= 0 && diffY <= 0 { // Left and Up || Go left first
		if (start == "0" || start == "A") && (end == "1" || end == "4" || end == "7") { // Gotta go up first
			outString += outStringY
			outString += outStringX
		} else {
			outString += outStringX
			outString += outStringY
		}

	} else if diffX <= 0 && diffY >= 0 { // Left and Down || Go left first
		if end == "<" {
			outString += outStringY
			outString += outStringX
		} else {
			outString += outStringX
			outString += outStringY
		}

	}
	outString += "A"
	pathMemo[fmt.Sprintf("%s%s", start, end)] = outString
	return outString
}
