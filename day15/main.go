package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
)

type point struct {
	x int
	y int
}

type vector struct {
	dx int
	dy int
}

const (
	boxStr      string = "O"
	robotStr    string = "@"
	wallStr     string = "#"
	emptyStr    string = "."
	leftBoxStr  string = "["
	rightBoxStr string = "]"
)

var (
	left  vector = vector{-1, 0}
	up    vector = vector{0, -1} // Reversed based on our grid
	right vector = vector{1, 0}
	down  vector = vector{0, 1}
)

var directionMap map[string]vector = map[string]vector{
	"<": left,
	"^": up,
	">": right,
	"v": down,
}

var translateWASD map[string]string = map[string]string{
	"w": "^",
	"a": "<",
	"s": "v",
	"d": ">",
}

// READ ME

// Change this to false to allow the simulation to automatically run
// If this is true, you can use "wasd" and the terminal to move the character to test whatever you want.
const interactive bool = true

func main() {
	file, err := os.Open("./input")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	bigMap := [][]string{}
	bigMap2 := [][]string{}
	moves := ""
	robotStart := point{}
	robotStart2 := point{}

	parsingMoves := false
	y := 0
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			parsingMoves = true
		}

		if !parsingMoves {
			newRow := []string{}
			newRow2 := []string{}
			for x, char := range line {
				newRow = append(newRow, string(char))
				if string(char) == robotStr {
					robotStart = point{x, y}
					robotStart2 = point{2 * x, y}
					newRow2 = append(newRow2, robotStr)
					newRow2 = append(newRow2, emptyStr)
				} else if string(char) == boxStr {
					newRow2 = append(newRow2, leftBoxStr)
					newRow2 = append(newRow2, rightBoxStr)
				} else {
					newRow2 = append(newRow2, string(char))
					newRow2 = append(newRow2, string(char))
				}
			}
			bigMap = append(bigMap, newRow)
			bigMap2 = append(bigMap2, newRow2)
			y++
		}
		if parsingMoves {
			moves += line
		}
	}

	// printMap(bigMap)
	// fmt.Println(moves)

	currentRobotPosition := robotStart
	for _, move := range moves {
		// printMap(bigMap)
		// fmt.Println(i, string(move))
		direction := directionMap[string(move)]
		nextRobotPosition := point{currentRobotPosition.x + direction.dx, currentRobotPosition.y + direction.dy}
		if bigMap[nextRobotPosition.y][nextRobotPosition.x] == wallStr {
			// Do nothing
			continue
		}
		if bigMap[nextRobotPosition.y][nextRobotPosition.x] == emptyStr {
			bigMap[nextRobotPosition.y][nextRobotPosition.x] = robotStr
			bigMap[currentRobotPosition.y][currentRobotPosition.x] = emptyStr
			currentRobotPosition = nextRobotPosition
			continue
		}
		if bigMap[nextRobotPosition.y][nextRobotPosition.x] == boxStr {
			if tryToMoveBoxes(bigMap, nextRobotPosition, direction) {
				bigMap[nextRobotPosition.y][nextRobotPosition.x] = robotStr
				bigMap[currentRobotPosition.y][currentRobotPosition.x] = emptyStr
				currentRobotPosition = nextRobotPosition
			} else {
				// Cant move, continue
				continue
			}
		}
	}

	// printMap(bigMap)

	fmt.Println(scoreForMap(bigMap, boxStr))

	// Part 2
	ch := make(chan string)
	go func(ch chan string) {
		// disable input buffering
		exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
		// do not display entered characters on the screen
		exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
		var b []byte = make([]byte, 1)
		for {
			os.Stdin.Read(b)
			ch <- string(b)
		}
	}(ch)

	currentRobotPosition = robotStart2
	for i, move := range moves {
		moveStr := string(move)
		if interactive {
			printMap(bigMap2)
			fmt.Println("Ctrl+C to stop")
			input := <-ch
			moveStr = translateWASD[input]
			fmt.Println(i, moveStr)
			moves = "v" // Ensures infinite loop
		}

		direction := directionMap[moveStr]
		nextRobotPosition := point{currentRobotPosition.x + direction.dx, currentRobotPosition.y + direction.dy}
		if bigMap2[nextRobotPosition.y][nextRobotPosition.x] == wallStr {
			// Do nothing
			continue
		}
		if bigMap2[nextRobotPosition.y][nextRobotPosition.x] == emptyStr {
			bigMap2[nextRobotPosition.y][nextRobotPosition.x] = robotStr
			bigMap2[currentRobotPosition.y][currentRobotPosition.x] = emptyStr
			currentRobotPosition = nextRobotPosition
			continue
		}
		if bigMap2[nextRobotPosition.y][nextRobotPosition.x] == leftBoxStr || bigMap2[nextRobotPosition.y][nextRobotPosition.x] == rightBoxStr {
			if canMoveBiggerBoxes(bigMap2, nextRobotPosition, direction) {
				doMoveBiggerBoxes(bigMap2, nextRobotPosition, direction)
				bigMap2[nextRobotPosition.y][nextRobotPosition.x] = robotStr
				bigMap2[currentRobotPosition.y][currentRobotPosition.x] = emptyStr
				currentRobotPosition = nextRobotPosition
			} else {
				// Cant move, continue
				continue
			}
		}
	}
	printMap(bigMap2)
	fmt.Println(scoreForMap(bigMap2, leftBoxStr))

	// 1468997 was wrong
}

func tryToMoveBoxes(bigMap [][]string, box point, direction vector) bool {
	if bigMap[box.y][box.x] == emptyStr {
		return true
	}
	if bigMap[box.y][box.x] == wallStr {
		return false
	}

	nextPoint := point{box.x + direction.dx, box.y + direction.dy}
	if tryToMoveBoxes(bigMap, nextPoint, direction) {
		bigMap[nextPoint.y][nextPoint.x] = boxStr
		return true
	}
	return false
}

func doMoveBiggerBoxes(bigMap [][]string, box point, direction vector) bool {
	if direction == right {
		consideredStr := bigMap[box.y][box.x]
		if consideredStr == emptyStr {
			return true
		}
		if consideredStr == wallStr {
			return false
		}
		if consideredStr == leftBoxStr {
			nextPoint := point{box.x + 2, box.y}
			if doMoveBiggerBoxes(bigMap, nextPoint, direction) {
				bigMap[box.y][box.x] = emptyStr
				bigMap[box.y][box.x+1] = leftBoxStr
				bigMap[box.y][box.x+2] = rightBoxStr
				return true
			}
		}
		return false
	}
	if direction == left {
		consideredStr := bigMap[box.y][box.x]
		if consideredStr == emptyStr {
			return true
		}
		if consideredStr == wallStr {
			return false
		}
		if consideredStr == rightBoxStr {
			nextPoint := point{box.x - 2, box.y}
			if doMoveBiggerBoxes(bigMap, nextPoint, direction) {
				bigMap[box.y][box.x] = emptyStr
				bigMap[box.y][box.x-2] = leftBoxStr
				bigMap[box.y][box.x-1] = rightBoxStr
				return true
			}
		}
		return false
	}
	if direction == up {
		consideredStr := bigMap[box.y][box.x]
		if consideredStr == emptyStr {
			return true
		}
		if consideredStr == wallStr {
			return false
		}
		if consideredStr == leftBoxStr {
			nextPoint := point{box.x, box.y - 1}
			otherPoint := point{box.x + 1, box.y - 1}
			if doMoveBiggerBoxes(bigMap, nextPoint, direction) && doMoveBiggerBoxes(bigMap, otherPoint, direction) {
				// Move the whole box up
				bigMap[nextPoint.y][nextPoint.x] = leftBoxStr
				bigMap[otherPoint.y][otherPoint.x] = rightBoxStr
				// Delete current box
				bigMap[box.y][box.x] = emptyStr
				bigMap[box.y][box.x+1] = emptyStr
				return true
			}
		}
		if consideredStr == rightBoxStr {
			nextPoint := point{box.x, box.y - 1}
			otherPoint := point{box.x - 1, box.y - 1}
			if doMoveBiggerBoxes(bigMap, nextPoint, direction) && doMoveBiggerBoxes(bigMap, otherPoint, direction) {
				// Move the whole box up
				bigMap[nextPoint.y][nextPoint.x] = rightBoxStr
				bigMap[otherPoint.y][otherPoint.x] = leftBoxStr
				// Delete current box
				bigMap[box.y][box.x] = emptyStr
				bigMap[box.y][box.x-1] = emptyStr
				return true
			}
		}
		return false
	}
	if direction == down {
		consideredStr := bigMap[box.y][box.x]
		if consideredStr == emptyStr {
			return true
		}
		if consideredStr == wallStr {
			return false
		}
		if consideredStr == leftBoxStr {
			nextPoint := point{box.x, box.y + 1}
			otherPoint := point{box.x + 1, box.y + 1}
			if doMoveBiggerBoxes(bigMap, nextPoint, direction) && doMoveBiggerBoxes(bigMap, otherPoint, direction) {
				// Move the whole box down
				bigMap[nextPoint.y][nextPoint.x] = leftBoxStr
				bigMap[otherPoint.y][otherPoint.x] = rightBoxStr
				// Delete current box
				bigMap[box.y][box.x] = emptyStr
				bigMap[box.y][box.x+1] = emptyStr
				return true
			}
		}
		if consideredStr == rightBoxStr {
			nextPoint := point{box.x, box.y + 1}
			otherPoint := point{box.x - 1, box.y + 1}
			if doMoveBiggerBoxes(bigMap, nextPoint, direction) && doMoveBiggerBoxes(bigMap, otherPoint, direction) {
				// Move the whole box down
				bigMap[nextPoint.y][nextPoint.x] = rightBoxStr
				bigMap[otherPoint.y][otherPoint.x] = leftBoxStr
				// Delete current box
				bigMap[box.y][box.x] = emptyStr
				bigMap[box.y][box.x-1] = emptyStr
				return true
			}
		}
		return false
	}

	return false
}

func canMoveBiggerBoxes(bigMap [][]string, box point, direction vector) bool {
	if direction == right {
		consideredStr := bigMap[box.y][box.x]
		if consideredStr == emptyStr {
			return true
		}
		if consideredStr == wallStr {
			return false
		}
		if consideredStr == leftBoxStr {
			nextPoint := point{box.x + 2, box.y}
			if canMoveBiggerBoxes(bigMap, nextPoint, direction) {
				return true
			}
		}
		return false
	}
	if direction == left {
		consideredStr := bigMap[box.y][box.x]
		if consideredStr == emptyStr {
			return true
		}
		if consideredStr == wallStr {
			return false
		}
		if consideredStr == rightBoxStr {
			nextPoint := point{box.x - 2, box.y}
			if canMoveBiggerBoxes(bigMap, nextPoint, direction) {
				return true
			}
		}
		return false
	}
	if direction == up {
		consideredStr := bigMap[box.y][box.x]
		if consideredStr == emptyStr {
			return true
		}
		if consideredStr == wallStr {
			return false
		}
		if consideredStr == leftBoxStr {
			nextPoint := point{box.x, box.y - 1}
			otherPoint := point{box.x + 1, box.y - 1}
			if canMoveBiggerBoxes(bigMap, nextPoint, direction) && canMoveBiggerBoxes(bigMap, otherPoint, direction) {
				return true
			}
		}
		if consideredStr == rightBoxStr {
			nextPoint := point{box.x, box.y - 1}
			otherPoint := point{box.x - 1, box.y - 1}
			if canMoveBiggerBoxes(bigMap, nextPoint, direction) && canMoveBiggerBoxes(bigMap, otherPoint, direction) {
				return true
			}
		}
		return false
	}
	if direction == down {
		consideredStr := bigMap[box.y][box.x]
		if consideredStr == emptyStr {
			return true
		}
		if consideredStr == wallStr {
			return false
		}
		if consideredStr == leftBoxStr {
			nextPoint := point{box.x, box.y + 1}
			otherPoint := point{box.x + 1, box.y + 1}
			if canMoveBiggerBoxes(bigMap, nextPoint, direction) && canMoveBiggerBoxes(bigMap, otherPoint, direction) {
				return true
			}
		}
		if consideredStr == rightBoxStr {
			nextPoint := point{box.x, box.y + 1}
			otherPoint := point{box.x - 1, box.y + 1}
			if canMoveBiggerBoxes(bigMap, nextPoint, direction) && canMoveBiggerBoxes(bigMap, otherPoint, direction) {
				return true
			}
		}
		return false
	}

	return false
}

func printMap(bigMap [][]string) {
	bigStr := ""
	for _, row := range bigMap {
		for _, subStr := range row {
			bigStr += subStr
		}
		bigStr += "\n"
	}
	fmt.Print(bigStr)
}

func scoreForMap(bigMap [][]string, idChar string) int {
	score := 0
	for y, row := range bigMap {
		for x, character := range row {
			if character == idChar {
				score += 100*y + x
			}
		}
	}
	return score
}
