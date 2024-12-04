package main

import (
	"bufio"
	"fmt"
	"os"
)

type direction struct {
	dx int
	dy int
}

var xMax, yMax int = 0, 0

var directions []direction = []direction{
	{-1, 1}, {0, 1}, {1, 1}, {-1, 0}, {1, 0}, {-1, -1}, {0, -1}, {1, -1},
}

func main() {
	file, err := os.Open("./input")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	grid := []string{}

	for scanner.Scan() {
		grid = append(grid, scanner.Text())
	}

	yMax = len(grid) - 1
	xMax = len(grid[0]) - 1

	total := 0
	total2 := 0
	for y := 0; y <= yMax; y++ {
		for x := 0; x <= xMax; x++ {
			total += xmasFoundFromThisPointLinear(grid, x, y)
			total2 += xmasFoundFromThisPointX(grid, x, y)
		}
	}

	fmt.Println(total)
	fmt.Println(total2)
}

const MAS = "MAS"

func xmasFoundFromThisPointLinear(grid []string, x, y int) int {
	found := 0
	if string(grid[y][x]) != "X" {
		return found
	}
	originalX, originalY := x, y
	for _, direction := range directions {
		x, y = originalX, originalY
		xmasFound := true
		for _, letter := range MAS {
			x += direction.dx
			y += direction.dy

			if x < 0 || x > xMax || y < 0 || y > yMax || string(grid[y][x]) != string(letter) {
				xmasFound = false
				break
			}
		}
		if xmasFound {
			found++
		}
	}
	return found
}

var leftCross []direction = []direction{{-1, 1}, {1, -1}}
var rightCross []direction = []direction{{1, 1}, {-1, -1}}

func xmasFoundFromThisPointX(grid []string, x, y int) int {
	found := 0
	if string(grid[y][x]) != "A" {
		return found
	}

	leftCrossLetters := map[string]int{}
	for _, direction := range leftCross {
		newY := y + direction.dy
		newX := x + direction.dx
		if newY < 0 || newY > yMax || newX < 0 || newX > xMax {
			break
		}
		leftCrossLetters[string(grid[newY][newX])]++
	}

	rightCrossLetters := map[string]int{}
	for _, direction := range rightCross {
		newY := y + direction.dy
		newX := x + direction.dx
		if newY < 0 || newY > yMax || newX < 0 || newX > xMax {
			break
		}
		rightCrossLetters[string(grid[newY][newX])]++
	}
	if rightCrossLetters["S"] == 1 && rightCrossLetters["M"] == 1 && leftCrossLetters["S"] == 1 && leftCrossLetters["M"] == 1 {
		found++
	}
	return found
}
