package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type point struct {
	x int
	y int
}

type vector struct {
	dx int
	dy int
}

var directions []vector = []vector{
	{-1, 0},
	{0, 1},
	{1, 0},
	{0, -1},
}

const size int = 71 // 7 or 71

func main() {
	file, err := os.Open("./input")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	fallingRocks := []point{}

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		split := strings.Split(line, ",")
		x, err := strconv.Atoi(split[0])
		if err != nil {
			panic(err)
		}
		y, err := strconv.Atoi(split[1])
		if err != nil {
			panic(err)
		}

		fallingRocks = append(fallingRocks, point{x, y})
	}

	grid := generateEmptyGrid(size)

	numOfFallenRocks := 1024
	for i, rock := range fallingRocks {
		if i == numOfFallenRocks {
			break
		}
		grid[rock.y][rock.x] = "#"
	}

	printGrid(grid)

	start := point{0, 0}
	end := point{size - 1, size - 1}

	scoreMap := map[point]int{start: 0}

	pointQueue := []point{start}

	for len(pointQueue) > 0 {
		currentPoint := pointQueue[0]
		pointQueue = pointQueue[1:]

		for _, direction := range directions {
			newPoint := point{currentPoint.x + direction.dx, currentPoint.y + direction.dy}
			if newPoint.x < 0 || newPoint.x >= size || newPoint.y < 0 || newPoint.y >= size {
				continue
			}
			if grid[newPoint.y][newPoint.x] == "#" {
				continue
			}
			newScore := scoreMap[currentPoint] + 1
			existingScore, ok := scoreMap[newPoint]
			if !ok {
				existingScore = math.MaxInt
			}

			if newScore >= existingScore {
				continue
			}
			scoreMap[newPoint] = newScore

			pointQueue = append(pointQueue, newPoint)
		}
	}

	fmt.Println(scoreMap[end])
	lastPointBeforeCutoff := 0
	// Part 2
	for i := 0; i < len(fallingRocks); i++ {
		grid = generateEmptyGrid(size)
		for j, rock := range fallingRocks {
			grid[rock.y][rock.x] = "#"
			if j == i {
				break
			}
		}
		start := point{0, 0}
		end := point{size - 1, size - 1}

		scoreMap := map[point]int{start: 0}

		pointQueue := []point{start}

		for len(pointQueue) > 0 {
			currentPoint := pointQueue[0]
			pointQueue = pointQueue[1:]

			for _, direction := range directions {
				newPoint := point{currentPoint.x + direction.dx, currentPoint.y + direction.dy}
				if newPoint.x < 0 || newPoint.x >= size || newPoint.y < 0 || newPoint.y >= size {
					continue
				}
				if grid[newPoint.y][newPoint.x] == "#" {
					continue
				}
				newScore := scoreMap[currentPoint] + 1
				existingScore, ok := scoreMap[newPoint]
				if !ok {
					existingScore = math.MaxInt
				}

				if newScore >= existingScore {
					continue
				}
				scoreMap[newPoint] = newScore

				pointQueue = append(pointQueue, newPoint)
			}
		}

		_, ok := scoreMap[end]
		if ok {
			lastPointBeforeCutoff = i
		} else {
			break
		}
	}

	fmt.Println(fallingRocks[lastPointBeforeCutoff+1])
}

func generateEmptyGrid(size int) [][]string {
	output := [][]string{}
	for i := 0; i < size; i++ {
		newRow := []string{}
		for j := 0; j < size; j++ {
			newRow = append(newRow, ".")
		}
		output = append(output, newRow)
	}
	return output
}

func printGrid(grid [][]string) {
	printString := ""
	for _, row := range grid {
		for _, char := range row {
			printString += char
		}
		printString += "\n"
	}
	fmt.Println(printString)
}
