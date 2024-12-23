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

	initialSecretNumbers := []int{}

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		asInt, err := strconv.Atoi(line)
		if err != nil {
			panic(err)
		}
		initialSecretNumbers = append(initialSecretNumbers, asInt)

	}
	// fmt.Println(initialSecretNumbers)
	diffs := []string{}
	allPrices := []string{}
	lastSecretNumber := []int{}

	for _, secretNumber := range initialSecretNumbers {
		diff := ""
		prices := ""
		for i := 0; i < 2000; i++ {
			lastNum := secretNumber
			secretNumber = GetNextSecretNumber(secretNumber)
			diffAsInt := GetLastDigit(secretNumber) - GetLastDigit(lastNum)
			if diffAsInt < 0 {
				diff += fmt.Sprintf("%d", diffAsInt)
			} else {
				diff += fmt.Sprintf("0%d", diffAsInt)
			}
			prices += fmt.Sprintf("%d", GetLastDigit(secretNumber))
		}
		// fmt.Println(diff)
		lastSecretNumber = append(lastSecretNumber, secretNumber)
		diffs = append(diffs, diff)
		allPrices = append(allPrices, prices)
	}
	sum := 0

	for _, secretNumber := range lastSecretNumber {
		sum += secretNumber
	}
	fmt.Println("Part 1:", sum)

	// Part 2
	// Just brute force it?
	sequenceToProfit := map[string]int{}

	for i := -9; i < 10; i++ {
		fmt.Println(i)
		for j := -9; j < 10; j++ {
			for k := -9; k < 10; k++ {
				for l := -9; l < 10; l++ {
					sequenceProfit := 0
					searchString := intToString(i) + intToString(j) + intToString(k) + intToString(l)
					for x, diff := range diffs {
						if n := strings.Index(diff, searchString); n != -1 {
							sequenceProfit += ezStrToInt(string(allPrices[x][n/2+3]))
						}
					}
					sequenceToProfit[searchString] = sequenceProfit
				}
			}
		}
	}

	maxSoFar := 0
	for k, v := range sequenceToProfit {
		if v > maxSoFar {
			fmt.Println(k, v)
			maxSoFar = v
		}
	}
}

func ezStrToInt(in string) int {
	asInt, err := strconv.Atoi(in)
	if err != nil {
		panic(err)
	}
	return asInt
}

func GetNextSecretNumber(in int) int {
	temp := (in) ^ (in * 64)
	temp = Prune(temp)

	temp = (temp / 32) ^ temp
	temp = Prune(temp)

	temp = (temp) ^ (temp * 2048)
	temp = Prune(temp)
	return temp
}

func Prune(in int) int {
	return in % 16777216
}

func GetLastDigit(in int) int {
	return in % 10
}

func intToString(in int) string {
	if in < 0 {
		return fmt.Sprintf("%d", in)
	} else {
		return fmt.Sprintf("0%d", in)
	}
}
