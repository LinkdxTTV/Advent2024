package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"
)

type design struct {
	target    string
	remaining string
	order     []string
}

func main() {
	file, err := os.Open("./input")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	availableTowels := []string{}
	desiredDesigns := []string{}
	parsingDesigns := false

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			parsingDesigns = true
			continue
		}

		if !parsingDesigns {
			split := strings.Split(line, ", ")
			availableTowels = append(availableTowels, split...)
		} else {
			desiredDesigns = append(desiredDesigns, line)
		}
	}

	validDesigns := 0
	sum := 0
	for _, desiredDesign := range desiredDesigns {
		ways := waysToMakeTowel(desiredDesign, availableTowels)
		if ways != 0 {
			validDesigns++
		}
		sum += ways
	}
	fmt.Println(validDesigns)
	fmt.Println(sum)
}

// use a quick global in recursion
var memo map[string]int = map[string]int{}

// Tells you how many ways you can make "target" which could be a towel or anything
func waysToMakeTowel(target string, towels []string) int {
	ways := 0
	if slices.Contains(towels, target) {
		ways += 1
	}

	if count, ok := memo[target]; ok {
		return count
	}

	for _, towel := range towels {
		if strings.HasPrefix(target, towel) {
			ways += waysToMakeTowel(strings.TrimPrefix(target, towel), towels)
		}
	}

	memo[target] = ways
	return ways
}
