package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
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

	left := []int{}
	right := []int{}
	for scanner.Scan() {
		line := scanner.Text()
		split := strings.Split(line, "   ")
		int1, err := strconv.Atoi(split[0])
		if err != nil {
			panic(err)
		}
		int2, err := strconv.Atoi(split[1])
		if err != nil {
			panic(err)
		}

		left = append(left, int1)
		right = append(right, int2)
	}

	sort.Ints(left)
	sort.Ints(right)

	sum := 0
	for i := 0; i < len(left); i++ {
		sum += abs(left[i] - right[i])
	}

	fmt.Println(sum)

	// Part 2
	rightMap := map[int]int{}
	for _, num := range right {
		rightMap[num]++
	}

	sum2 := 0
	for _, num := range left {
		occurences, ok := rightMap[num]
		if ok {
			sum2 += occurences * num
		}
	}
	fmt.Println(sum2)
}

func abs(a int) int {
	if a > 0 {
		return a
	}
	return -a
}
