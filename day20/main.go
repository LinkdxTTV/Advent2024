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

var directions []vector = []vector{
	{-1, 0},
	{0, 1},
	{1, 0},
	{0, -1},
}

type bigMap [][]string

type player struct {
	location         point
	usedCheat        bool
	inCheatMode      bool
	cheatCounter     int
	distance         int
	lastLocation     point
	firstWallTouched point
	firstOutOfWall   point
}

type cheat struct {
	firstWall      point
	firstOutOfWall point
}

const (
	wallStr  string = "#"
	emptyStr string = "."
	startStr string = "S"
	endStr   string = "E"
)

func main() {
	file, err := os.Open("./input")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	bigMap := bigMap{}
	var start, end point

	y := 0
	for scanner.Scan() {
		newRow := []string{}
		line := scanner.Text()
		for x, char := range line {
			newRow = append(newRow, string(char))
			if string(char) == startStr {
				start = point{x, y}
			}
			if string(char) == endStr {
				end = point{x, y}
			}
		}
		bigMap = append(bigMap, newRow)
		y++
	}

	xMax := len(bigMap[0])
	yMax := len(bigMap)

	fmt.Println(start, end)
	bigMap.print()

	// Start with a cheatless player to find the base distance
	playerQueue := []player{{
		location:     start,
		usedCheat:    false,
		inCheatMode:  false,
		cheatCounter: 2,
		distance:     0,
		lastLocation: point{0, 0},
	}}

	finishedPlayers := []player{}
	noCheatsPath := map[point]int{}

	for len(playerQueue) > 0 {
		currentPlayer := playerQueue[0]
		playerQueue = playerQueue[1:]

		if bigMap.at(currentPlayer.location) == endStr {
			finishedPlayers = append(finishedPlayers, currentPlayer)
			if !currentPlayer.usedCheat {
				noCheatsPath[currentPlayer.location] = currentPlayer.distance
			}
			continue
		}

		if currentPlayer.inCheatMode {
			currentPlayer.cheatCounter--
			if currentPlayer.cheatCounter == 0 {
				currentPlayer.inCheatMode = false
			}
		}

		if !currentPlayer.inCheatMode && bigMap.at(currentPlayer.location) == wallStr {
			// Cannot come out in a wall
			continue
		}

		if bigMap.at(currentPlayer.lastLocation) == wallStr && bigMap.at(currentPlayer.location) != wallStr {
			currentPlayer.firstOutOfWall = currentPlayer.location
		}

		if !currentPlayer.usedCheat {
			noCheatsPath[currentPlayer.location] = currentPlayer.distance
		}

		for _, direction := range directions {
			newPoint := point{currentPlayer.location.x + direction.dx, currentPlayer.location.y + direction.dy}
			if newPoint.x < 0 || newPoint.x >= xMax || newPoint.y < 0 || newPoint.y >= yMax {
				continue
			}
			// Dont go back
			if newPoint == currentPlayer.lastLocation {
				continue
			}
			if bigMap.at(newPoint) == wallStr {
				if currentPlayer.usedCheat == false || currentPlayer.inCheatMode {
					// Send in a cheat
					playerQueue = append(playerQueue, player{
						location:         newPoint,
						usedCheat:        true,
						inCheatMode:      true,
						cheatCounter:     currentPlayer.cheatCounter,
						distance:         currentPlayer.distance + 1,
						lastLocation:     currentPlayer.location,
						firstWallTouched: newPoint,
					})
				} else {
					continue
				}
			} else {
				playerQueue = append(playerQueue, player{
					location:         newPoint,
					usedCheat:        currentPlayer.usedCheat,
					inCheatMode:      currentPlayer.inCheatMode,
					cheatCounter:     currentPlayer.cheatCounter,
					distance:         currentPlayer.distance + 1,
					lastLocation:     currentPlayer.location,
					firstWallTouched: currentPlayer.firstWallTouched,
					firstOutOfWall:   currentPlayer.firstOutOfWall,
				})
			}
		}
	}

	// fmt.Println(finishedPlayers)

	// Dedupe cheats
	dedupedFinishedPlayers := map[player]bool{}
	for _, player := range finishedPlayers {
		dedupedFinishedPlayers[player] = true
	}

	baseline := 0
	for player := range dedupedFinishedPlayers {
		if player.usedCheat == false {
			baseline = player.distance
		}
	}

	fmt.Println("baseline", baseline)
	threshold := 100
	cheatsThatSavedOverTheshold := 0
	cheatSavingsMap := map[int]int{}

	for player := range dedupedFinishedPlayers {
		cheatSavingsMap[baseline-player.distance]++
		if baseline-player.distance >= threshold {
			cheatsThatSavedOverTheshold++
		}
	}

	// for k, v := range cheatSavingsMap {
	// 	if k > 0 {
	// 		fmt.Println(v, k)
	// 	}
	// }
	fmt.Println("Part 1:")
	fmt.Println(cheatsThatSavedOverTheshold)

	// Part 1 check: using part 2 method:
	validPart1Cheats := 0
	// Do a quick retraversal of the path
	// fmt.Println(noCheatsPath)
	for p1, d1 := range noCheatsPath {
		for p2, d2 := range noCheatsPath {
			if d2-d1-p2.manhattanDistance(p1) >= 100 {
				if p2.manhattanDistance(p1) <= 2 {
					validPart1Cheats++
				}
			}
		}
	}
	fmt.Println("Check with part 2 approach:", validPart1Cheats)
	// Part 2
	// Cross every two points on the existing path. If their distance (time in this problem) value is >= 100
	// and they are within 20 manhattan distance of eachother, mark it as a valid cheat path
	validPart2Cheats := 0
	// Do a quick retraversal of the path
	// fmt.Println(noCheatsPath)
	for p1, d1 := range noCheatsPath {
		for p2, d2 := range noCheatsPath {
			if d2-d1-p2.manhattanDistance(p1) >= 100 {
				if p2.manhattanDistance(p1) <= 20 {
					validPart2Cheats++
				}
			}
		}
	}
	fmt.Println(validPart2Cheats)
}

func (b bigMap) at(p point) string {
	return b[p.y][p.x]
}

func (b bigMap) print() {
	out := ""
	for _, row := range b {
		for _, char := range row {
			out += char
		}
		out += "\n"
	}

	fmt.Print(out)
}

func (v vector) add(v2 vector) vector {
	return vector{v.dx + v2.dx, v.dy + v2.dy}
}

func (p point) manhattanDistance(end point) int {
	return abs(p.x-end.x) + abs(p.y-end.y)
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}
