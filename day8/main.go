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

func main() {
	file, err := os.Open("./input")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	antennaes := map[string][]point{}

	y := 0
	yMax, xMax := 0, 0
	for scanner.Scan() {
		line := scanner.Text()
		xMax = len(line)

		for x, char := range line {
			if string(char) != "." {
				points, ok := antennaes[string(char)]
				if ok {
					points = append(points, point{x, y})
					antennaes[string(char)] = points
				} else {
					antennaes[string(char)] = []point{{x, y}}
				}
			}
		}
		y++
	}
	yMax = y

	uniqueAntinodes := map[point]bool{}

	for _, points := range antennaes {

		for i, point1 := range points {
			for j, point2 := range points {
				if i == j {
					// Points dont reflect themselves
					continue
				}
				antinode := point{2*point1.x - point2.x, 2*point1.y - point2.y}
				if antinode.x >= 0 && antinode.x < xMax && antinode.y >= 0 && antinode.y < yMax {
					uniqueAntinodes[antinode] = true
				}
			}
		}
	}

	fmt.Println(len(uniqueAntinodes))

	// Part 2
	uniqueAntinodes2 := map[point]bool{}

	for _, points := range antennaes {

		for i, point1 := range points {
			for j, point2 := range points {
				if i == j {
					// Points dont reflect themselves
					continue
				}
				antinode := point2
				dx, dy := point2.x-point1.x, point2.y-point1.y
				for antinode.x >= 0 && antinode.x < xMax && antinode.y >= 0 && antinode.y < yMax {
					uniqueAntinodes2[antinode] = true
					antinode.x += dx
					antinode.y += dy
				}
			}
		}
	}

	fmt.Println(len(uniqueAntinodes2))
}
