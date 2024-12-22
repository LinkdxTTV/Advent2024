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

	maxIterations := 25

	for _, direction := range directions {
		fmt.Println(direction)
		robotDirections := GenerateDirectionsOnPad(direction, "A", numpad, maxIterations, 0)
		// fmt.Println(robotDirections)
		for i := 0; i < maxIterations; i++ {
			robotDirections = GenerateDirectionsOnPad(robotDirections, "A", keypad, maxIterations, i)
			// fmt.Println(robotDirections)
		}

		sequenceLength := len(robotDirections)
		numericCode := calculateNumericCode(direction)

		sum += sequenceLength * numericCode
		fmt.Println(sequenceLength, numericCode)
	}

	fmt.Println(sum)

	// for i, direction := range []string{">v", "v>", "^>", ">^", "<^", "^<", "<v", "v<"} {
	// 	if i%2 == 0 {
	// 		fmt.Println()
	// 	}
	// 	robot := direction
	// 	for i := 0; i < 1; i++ {
	// 		robot = GenerateDirectionsOnPad(robot, "A", keypad, 1, i)
	// 	}
	// 	fmt.Println(direction, len(robot))
	// }
}

func calculateNumericCode(input string) int {
	justNumbers := strings.TrimSuffix(input, "A")
	asInt, err := strconv.Atoi(justNumbers)
	if err != nil {
		panic(err)
	}
	return asInt
}

func GenerateDirectionsOnPad(input string, start string, pad [][]string, maxIterations, iteration int) string {
	if input == "" {
		return ""
	}
	outString := GeneratePath(start, string(input[0]), pad, maxIterations, iteration)
	outString += "A"
	for i := 0; i < len(input)-1; i++ {
		j := i + 1
		outString += GeneratePath(string(input[i]), string(input[j]), pad, maxIterations, iteration)
		outString += "A"
	}
	return outString
}

var memo = map[string]string{}

// Does not include A's
func GeneratePath(start, end string, pad [][]string, maxIterations, iteration int) string {
	iterationsLeft := maxIterations - iteration
	last2Iterations := iterationsLeft < 2
	if out, ok := memo[fmt.Sprintf("%s.%s.%v", start, end, last2Iterations)]; ok {
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
			if iterationsLeft <= 1 {
				outString += outStringX
				outString += outStringY
			} else {
				outString += outStringY
				outString += outStringX
			}
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
			if iterationsLeft <= 1 {
				outString += outStringY
				outString += outStringX
			} else {
				outString += outStringX
				outString += outStringY
			}
		}

	} else if diffX <= 0 && diffY >= 0 { // Left and Down || Go left first
		if end == "<" {
			outString += outStringY
			outString += outStringX
		} else {
			if iterationsLeft <= 1 {
				outString += outStringY
				outString += outStringX
			} else {
				outString += outStringX
				outString += outStringY
			}
		}

	}
	memo[fmt.Sprintf("%s.%s.%v", start, end, last2Iterations)] = outString
	return outString
}

// 242484
// 294209504640384
