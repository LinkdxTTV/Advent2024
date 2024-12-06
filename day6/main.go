package main

import (
	"bufio"
	"fmt"
	"os"
)

type point struct {
	x int
	y int
}

type vector struct {
	dx int
	dy int
}

type pointvector struct {
	x  int
	y  int
	dx int
	dy int
}

var directions []vector = []vector{
	{0, -1}, {1, 0}, {0, 1}, {-1, 0},
}

var distinctPositions map[point]bool = map[point]bool{}

const obstruction string = "#"

func main() {
	file, err := os.Open("./input")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	bigMap := []string{}

	for scanner.Scan() {
		line := scanner.Text()

		bigMap = append(bigMap, line)
	}

	fmt.Println(bigMap)
	xMax := len(bigMap[0]) - 1
	yMax := len(bigMap) - 1

	fmt.Println(xMax, yMax)

	// Part 1

	currentPos := findInitialPosition(bigMap)
	directionIndex := 0

	for {
		distinctPositions[currentPos] = true
		// Look ahead
		direction := directions[directionIndex]
		nextPoint := currentPos.move(direction)

		// Exit the map
		if nextPoint.x < 0 || nextPoint.x > xMax || nextPoint.y < 0 || nextPoint.y > yMax {
			break
		}

		if nextPoint.character(bigMap) == obstruction {
			directionIndex++
			if directionIndex > 3 {
				directionIndex = 0
			}
		} else {
			currentPos = nextPoint
		}
	}

	fmt.Println("Part 1")
	fmt.Println(len(distinctPositions))
	// Part 2
	numLoops := 0

	// We can only place obstructions at places we've been as they are the only places that could affect the trajectory.
	// Dont forget to remove the starting location.
	delete(distinctPositions, findInitialPosition(bigMap))
	i := 0

	for point := range distinctPositions {
		if placingObstructionAtPointWouldCauseLoop(point, bigMap, xMax, yMax) {
			numLoops++
		}
		i++
		fmt.Printf("Progress: %d / %d \n", i, len(distinctPositions))
	}

	fmt.Println("Part 2")
	fmt.Println(numLoops)
}

func findInitialPosition(bigMap []string) point {
	for y, row := range bigMap {
		for x, char := range row {
			if string(char) == "^" {
				return point{x, y}
			}
		}
	}
	return point{}
}

func (p point) move(v vector) point {
	return point{p.x + v.dx, p.y + v.dy}
}

func (p point) character(b []string) string {
	return string(b[p.y][p.x])
}

// Part 2

func placingObstructionAtPointWouldCauseLoop(p point, bigMap []string, xMax, yMax int) bool {
	distinctPositionAndDirection := map[pointvector]bool{}
	currentPos := findInitialPosition(bigMap)
	directionIndex := 0

	bigMap = copyOfMapWithObstruction(bigMap, p)

	for {
		// fmt.Println(currentPos)
		direction := directions[directionIndex]
		currentPointVector := pointvector{
			currentPos.x, currentPos.y, direction.dx, direction.dy,
		}

		if distinctPositionAndDirection[currentPointVector] {
			// We've been at this position with this vector before
			return true
		}

		distinctPositionAndDirection[currentPointVector] = true
		// Look ahead
		nextPoint := currentPos.move(direction)

		// Exit the map
		if nextPoint.x < 0 || nextPoint.x > xMax || nextPoint.y < 0 || nextPoint.y > yMax {
			return false
		}

		if nextPoint.character(bigMap) == obstruction {
			directionIndex++
			if directionIndex > 3 {
				directionIndex = 0
			}
		} else {
			currentPos = nextPoint
		}
	}
}

func copyOfMapWithObstruction(bigMap []string, p point) []string {
	out := []string{}

	for y, row := range bigMap {
		if y != p.y {
			out = append(out, row)
		} else {
			newRow := ""
			for x, char := range row {
				if p.x != x {
					newRow += string(char)
				} else {
					newRow += "#"
				}
			}
			out = append(out, newRow)
		}
	}

	return out
}

func printMap(bigMap []string) {
	for _, row := range bigMap {
		fmt.Println(row)
	}
}
