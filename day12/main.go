package main

import (
	"bufio"
	"fmt"
	"os"
)

var xMax int
var yMax int

type point struct {
	x int
	y int
}

type vector struct {
	dx int
	dy int
}

var directions []vector = []vector{{-1, 0}, {0, 1}, {1, 0}, {0, -1}}

type Region struct {
	id               string
	areaPoints       map[point]bool
	perimeterPoints  map[point]int
	allTouchedPoints map[point]bool
}

// Globals to simplify
var allRegions []*Region = []*Region{}

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

	// printMap(bigMap)
	xMax = len(bigMap[0])
	yMax = len(bigMap)

	seenPoints := map[point]bool{}

	for y := 0; y < yMax; y++ {
		for x := 0; x < xMax; x++ {
			point := point{x, y}
			recursivelyScoutIsland(bigMap, seenPoints, point, nil, charAtPoint(bigMap, point))
		}
	}

	fmt.Println(allRegions)
	totalPrice := 0
	totalPrice2 := 0
	for _, region := range allRegions {
		fmt.Println(region.id)
		fmt.Println("Area", len(region.areaPoints))
		fmt.Println("Perimeter", calcPerimeter(region.perimeterPoints))
		sides := calculateNumberOfSidesFromAreaMap(bigMap, *region)
		fmt.Println("Sides", sides)
		totalPrice += len(region.areaPoints) * calcPerimeter(region.perimeterPoints)
		totalPrice2 += len(region.areaPoints) * sides
	}
	fmt.Println(totalPrice)
	fmt.Println(totalPrice2)
}

func recursivelyScoutIsland(bigMap []string, seenPoints map[point]bool, currentPoint point, region *Region, id string) {
	if currentPoint.x < 0 || currentPoint.x >= xMax || currentPoint.y < 0 || currentPoint.y >= yMax {
		region.perimeterPoints[currentPoint]++
		return
	}
	// We will do two different things depending on if we are looking for area or perimeter
	if charAtPoint(bigMap, currentPoint) == id {
		// Add to area
		_, ok := seenPoints[currentPoint]
		if ok {
			// We dont need to go to points twice
			return
		}
		if region == nil {
			// This is a new region
			fmt.Println("new region!")
			region = &Region{
				id:              id,
				areaPoints:      map[point]bool{},
				perimeterPoints: map[point]int{},
			}

			allRegions = append(allRegions, region)
		}
		region.areaPoints[currentPoint] = true
		seenPoints[currentPoint] = true

		// Recursive Forward
		for _, direction := range directions {
			newX := currentPoint.x + direction.dx
			newY := currentPoint.y + direction.dy
			newPoint := point{newX, newY}
			recursivelyScoutIsland(bigMap, seenPoints, newPoint, region, id)
		}
	} else {
		// Add to perimeter
		region.perimeterPoints[currentPoint]++
	}
}

func printMap(bigMap []string) {
	for _, row := range bigMap {
		fmt.Println(row)
	}
}

func charAtPoint(bigMap []string, p point) string {
	if p.x < 0 || p.x >= xMax || p.y < 0 || p.y >= yMax {
		return "."
	}
	return string(bigMap[p.y][p.x])
}

func calcPerimeter(pmap map[point]int) int {
	sum := 0
	for _, v := range pmap {
		sum += v
	}
	return sum
}

func calculateNumberOfSidesFromAreaMap(bigMap []string, region Region) int {
	corners := 0
	for point := range region.areaPoints {
		corners += numCornersFromPoint(bigMap, point, region.id)
	}
	return corners
}

func numCornersFromPoint(bigMap []string, p point, char string) int {
	corners := 0
	// 4 Cases

	// Up Left
	if charAtPoint(bigMap, point{p.x, p.y + 1}) == char && charAtPoint(bigMap, point{p.x - 1, p.y}) == char && charAtPoint(bigMap, point{p.x - 1, p.y + 1}) != char {
		corners++
	}
	if charAtPoint(bigMap, point{p.x, p.y + 1}) != char && charAtPoint(bigMap, point{p.x - 1, p.y}) != char {
		corners++
	}

	// Up Right
	if charAtPoint(bigMap, point{p.x, p.y + 1}) == char && charAtPoint(bigMap, point{p.x + 1, p.y}) == char && charAtPoint(bigMap, point{p.x + 1, p.y + 1}) != char {
		corners++
	}
	if charAtPoint(bigMap, point{p.x, p.y + 1}) != char && charAtPoint(bigMap, point{p.x + 1, p.y}) != char {
		corners++
	}

	// Bottom Left
	if charAtPoint(bigMap, point{p.x, p.y - 1}) == char && charAtPoint(bigMap, point{p.x - 1, p.y}) == char && charAtPoint(bigMap, point{p.x - 1, p.y - 1}) != char {
		corners++
	}
	if charAtPoint(bigMap, point{p.x, p.y - 1}) != char && charAtPoint(bigMap, point{p.x - 1, p.y}) != char {
		corners++
	}

	// Bottom Right
	if charAtPoint(bigMap, point{p.x, p.y - 1}) == char && charAtPoint(bigMap, point{p.x + 1, p.y}) == char && charAtPoint(bigMap, point{p.x + 1, p.y - 1}) != char {
		corners++
	}
	if charAtPoint(bigMap, point{p.x, p.y - 1}) != char && charAtPoint(bigMap, point{p.x + 1, p.y}) != char {
		corners++
	}

	return corners
}
