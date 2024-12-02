package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	file, err := os.Open("./input")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	reports := [][]int{}
	for scanner.Scan() {
		report := []int{}
		line := scanner.Text()
		split := strings.Split(line, " ")
		for _, num := range split {
			if num == "" {
				continue
			}
			asInt, err := strconv.Atoi(num)
			if err != nil {
				panic(err)
			}
			report = append(report, asInt)
		}

		reports = append(reports, report)
	}

	fmt.Println(reports)
	safeReports := 0

	for _, report := range reports {

		if reportIsSafe(report) {
			safeReports++
		} else {
			// just try a bunch
			for i := 0; i < len(report); i++ {
				newlySafe := reportIsSafe(returnReportWithIndexMissing(report, i))
				if newlySafe {
					safeReports++
					break
				}
			}
		}
	}

	fmt.Println(safeReports)
}

func abs(a int) int {
	if a > 0 {
		return a
	}
	return -a
}

// if not safe, reports the index that was not safe.
func reportIsSafe(report []int) bool {
	safe := true
	increasing := true
	if report[1]-report[0] < 0 {
		increasing = false
	}
	for i := 0; i < len(report)-1; i++ {
		if report[i+1] == report[i] {
			safe = false
			break
		}
		if abs(report[i+1]-report[i]) > 3 {
			safe = false
			break
		}
		if increasing {
			if report[i+1]-report[i] <= 0 {
				safe = false
				break
			}
		} else {
			if report[i+1]-report[i] >= 0 {
				safe = false
				break
			}
		}
	}
	return safe
}

func returnReportWithIndexMissing(report []int, index int) []int {
	out := []int{}
	out = append(out, report[:index]...)
	out = append(out, report[index+1:]...)
	return out
}
