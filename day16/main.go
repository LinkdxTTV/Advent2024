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

var fourDirections []vector = []vector{
	{-1, 0}, {0, 1}, {1, 0}, {0, -1},
}

var zeroVector vector = vector{0, 0}

type Reindeer struct {
	p point
	v vector
}

type ReindeerWithPath struct {
	r    Reindeer
	path []point
}

type PathWithScore *[]point

const (
	wallStr  string = "#"
	emptyStr string = "."
	startStr string = "S"
	endStr   string = "E"
)

type bigMap [][]string

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

	// bigMap.print()
	// fmt.Println(start, end)

	startReindeer := Reindeer{start, vector{1, 0}}
	startReindeerWithPath := ReindeerWithPath{startReindeer, []point{start}}

	bestPaths := map[PathWithScore]int{}
	reindeerMapToScore := map[Reindeer]int{
		startReindeer: 0,
	}
	// Begin BFS.. We start facing east on Start tile
	reindeerQueue := []ReindeerWithPath{startReindeerWithPath}
	bestScoreSoFar := 99999999999999

	for len(reindeerQueue) > 0 {
		currentReindeer := reindeerQueue[0]
		reindeerQueue = reindeerQueue[1:]

		if currentReindeer.r.p == end {
			if reindeerMapToScore[currentReindeer.r] <= bestScoreSoFar {
				bestScoreSoFar = reindeerMapToScore[currentReindeer.r]
				bestPaths[&currentReindeer.path] = bestScoreSoFar
			}
			continue
		}
		for _, direction := range fourDirections {
			// Dont go backwards
			if currentReindeer.r.v.add(direction) == zeroVector {
				continue
			}
			nextPoint := point{currentReindeer.r.p.x + direction.dx, currentReindeer.r.p.y + direction.dy}
			if bigMap.at(nextPoint) == wallStr {
				continue
			}
			var addedScore int
			if currentReindeer.r.v == direction {
				addedScore = 1
			} else {
				addedScore = 1001
			}
			nextReindeer := Reindeer{nextPoint, direction}
			newScore := reindeerMapToScore[currentReindeer.r] + addedScore
			if newScore > bestScoreSoFar {
				continue
			}
			existingScore, ok := reindeerMapToScore[nextReindeer]
			if !ok {
				reindeerMapToScore[nextReindeer] = newScore
			} else {
				if existingScore < newScore {
					continue
				}
				// We found a better score
				reindeerMapToScore[nextReindeer] = newScore
			}
			newPath := copyPath(currentReindeer.path)
			newPath = append(newPath, nextPoint)
			reindeerQueue = append(reindeerQueue, ReindeerWithPath{nextReindeer, newPath})
		}
	}

	// fmt.Println(reindeerMapToScore)
	lowestScore := 99999999999999
	for deer, score := range reindeerMapToScore {
		if deer.p == end {
			if score < lowestScore {
				lowestScore = score
			}
		}
	}

	fmt.Println(lowestScore)
	uniquePointsInAnyBestPath := map[point]bool{}
	for pathList, score := range bestPaths {

		if score == lowestScore {
			for _, point := range *pathList {
				uniquePointsInAnyBestPath[point] = true
			}
		}

	}
	// Illustration purposes only

	// for point := range uniquePointsInAnyBestPath {
	// 	bigMap[point.y][point.x] = "O"
	// }
	// bigMap.print()
	fmt.Println(len(uniquePointsInAnyBestPath))
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

func copyPath(in []point) []point {
	out := []point{}
	out = append(out, in...)
	return out
}
