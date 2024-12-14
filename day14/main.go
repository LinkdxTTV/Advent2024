package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
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

type robot struct {
	position point
	velocity vector
}

const (
	maxX    int = 101 // 11 // 101
	maxY    int = 103 /// 7 // 103
	seconds int = 1   // 100
)

func main() {
	file, err := os.Open("./input")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	regex := regexp.MustCompile(`p=(-?\d+),(-?\d+) v=(-?\d+),(-?\d+)`)

	robots := []robot{}
	for scanner.Scan() {
		line := scanner.Text()
		finds := regex.FindAllStringSubmatch(line, -1)

		position := point{}
		velocity := vector{}
		position.x, err = strconv.Atoi(finds[0][1])
		if err != nil {
			panic(err)
		}
		position.y, err = strconv.Atoi(finds[0][2])
		if err != nil {
			panic(err)
		}
		velocity.dx, err = strconv.Atoi(finds[0][3])
		if err != nil {
			panic(err)
		}
		velocity.dy, err = strconv.Atoi(finds[0][4])
		if err != nil {
			panic(err)
		}
		robots = append(robots, robot{position, velocity})
	}

	fmt.Println(robots)
	robotPositions := map[point]int{}
	for _, robot := range robots {
		// Afer 100 iterations
		x := mod(robot.position.x+robot.velocity.dx*seconds, maxX)
		y := mod(robot.position.y+robot.velocity.dy*seconds, maxY)

		robotPositions[point{x, y}]++
	}

	// fmt.Println(robotPositions)

	quadrants := map[int]int{}

	for position, robots := range robotPositions {
		if position.x < (maxX-1)/2 && position.y < (maxY-1)/2 {
			quadrants[1] += robots
		}
		if position.x > (maxX-1)/2 && position.y < (maxY-1)/2 {
			quadrants[2] += robots
		}
		if position.x > (maxX-1)/2 && position.y > (maxY-1)/2 {
			quadrants[3] += robots
		}
		if position.x < (maxX-1)/2 && position.y > (maxY-1)/2 {
			quadrants[4] += robots
		}
	}

	product := 1
	for _, bots := range quadrants {
		product = product * bots
	}

	fmt.Println("Part 1", product)
	drawMap(robotPositions)
	i := 0
	// Part 2
	for {
		shouldWait := true
		robotPositions := map[point]int{}
		for _, robot := range robots {
			// Afer 100 iterations
			x := mod(robot.position.x+robot.velocity.dx*i, maxX)
			y := mod(robot.position.y+robot.velocity.dy*i, maxY)

			robotPositions[point{x, y}]++
		}
		for _, numRobots := range robotPositions {
			if numRobots > 1 {
				shouldWait = false
			}
		}
		if shouldWait {
			drawMap(robotPositions)
			fmt.Println(i, "Ctrl-C to exit, enter to continue") // Happened at 8149, cycle was 18552-8149 = 10403
			fmt.Scanln()
		}
		i++
	}
}

func mod(a, b int) int {
	return (a%b + b) % b
}

func drawMap(robots map[point]int) {
	baseMap := [][]string{}

	for i := 0; i < maxY; i++ {
		newString := []string{}
		for j := 0; j < maxX; j++ {
			newString = append(newString, ".")
		}
		baseMap = append(baseMap, newString)
	}

	for robotPos := range robots {
		baseMap[robotPos.y][robotPos.x] = "X"
	}

	for _, row := range baseMap {
		for _, column := range row {
			fmt.Print(column)
		}
		fmt.Println()
	}
}
