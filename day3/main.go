package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
)

func main() {
	file, err := os.Open("./input")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	stringBuffer := ""

	for scanner.Scan() {
		stringBuffer += scanner.Text()
	}

	regex, err := regexp.Compile(`mul\(([0-9]+)\,([0-9]+)\)`)
	if err != nil {
		panic(err)
	}

	matches := regex.FindAllStringSubmatch(stringBuffer, -1)

	sum := 0

	for _, match := range matches {

		sum += getProductFromMatch(match)
	}

	fmt.Println(sum)

	// Part 2

	matchIndices := regex.FindAllStringIndex(stringBuffer, -1)

	doMatch := regexp.MustCompile(`do\(\)`)
	dontMatch := regexp.MustCompile(`don\'t\(\)`)

	doMatches := doMatch.FindAllStringIndex(stringBuffer, -1)
	dontMatches := dontMatch.FindAllStringIndex(stringBuffer, -1)

	// Stack them
	type Command struct {
		command string
		index   int
		product int
	}

	commands := []Command{}

	// Muls
	for i := 0; i < len(matches); i++ {
		command := Command{
			command: "mul",
			index:   matchIndices[i][0],
			product: getProductFromMatch(matches[i]),
		}
		commands = append(commands, command)
	}

	// Dos
	for _, doMatch := range doMatches {
		commands = append(commands, Command{
			command: "do",
			index:   doMatch[0],
		})
	}

	for _, dontMatch := range dontMatches {
		commands = append(commands, Command{
			command: "dont",
			index:   dontMatch[0],
		})
	}

	sort.Slice(commands, func(i, j int) bool {
		return commands[i].index < commands[j].index
	})

	enabled := true
	sum = 0
	for _, command := range commands {
		if command.command == "do" {
			enabled = true
			continue
		}
		if command.command == "dont" {
			enabled = false
			continue
		}
		if enabled && command.command == "mul" {
			sum += command.product
		}
	}

	fmt.Println(sum)
}

func getProductFromMatch(match []string) int {
	num1, err := strconv.Atoi(match[1])
	if err != nil {
		panic(err)
	}
	num2, err := strconv.Atoi(match[2])
	if err != nil {
		panic(err)
	}
	return num1 * num2
}
