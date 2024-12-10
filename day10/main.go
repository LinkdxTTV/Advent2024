package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

type point struct {
	x int
	y int
}

type vector struct {
	dx int
	dy int
}

var checkDirections []vector = []vector{
	{1, 0}, {-1, 0}, {0, 1}, {0, -1},
}

var maxX, maxY int = 0, 0

func main() {
	file, err := os.Open("./input")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	topography := [][]int{}
	startPoints := []point{}

	y := 0
	for scanner.Scan() {
		line := scanner.Text()
		newLine := []int{}
		for x, char := range line {
			charAsInt, err := strconv.Atoi(string(char))
			if err != nil {
				panic(err)
			}
			newLine = append(newLine, charAsInt)
			if charAsInt == 0 {
				startPoints = append(startPoints, point{x, y})
			}
		}
		topography = append(topography, newLine)
		y++
	}

	// Assuming a square
	maxY = len(topography) - 1
	maxX = len(topography[0]) - 1

	// fmt.Println(topography)

	trailsMap := map[point]map[point]int{} // Map from starting point -> map of reachable points -> Unique paths to that end point
	for _, point := range startPoints {
		recursivelyMoveUpwards(topography, point, point, trailsMap)
	}

	// fmt.Println(trailsMap)
	sum := 0
	for _, trails := range trailsMap {
		sum += len(trails)
	}
	fmt.Println(sum)

	// Part 2
	sum = 0
	for _, trails := range trailsMap {
		for _, uniquePaths := range trails {
			sum += uniquePaths
		}
	}
	fmt.Println(sum)

}

func recursivelyMoveUpwards(topography [][]int, currentPlace point, start point, trailsMap map[point]map[point]int) {
	currentValue := topography[currentPlace.y][currentPlace.x]
	if currentValue == 9 {
		// We made it
		list, ok := trailsMap[start]
		if !ok {
			list = map[point]int{}
			trailsMap[start] = list
		}
		list[currentPlace]++
		trailsMap[start] = list // not sure if needed
	}

	for _, direction := range checkDirections {
		newX := currentPlace.x + direction.dx
		newY := currentPlace.y + direction.dy

		if newX < 0 || newX > maxX || newY < 0 || newY > maxY {
			continue
		}
		if topography[newY][newX] == currentValue+1 {
			// Go up
			recursivelyMoveUpwards(topography, point{newX, newY}, start, trailsMap)
		}
	}
}
