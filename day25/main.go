package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	file, err := os.Open("./input")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	locksStr := [][]string{}
	keysStr := [][]string{}

	buffer := []string{}
	// Might need to add a newline to the input to make this clean
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			if buffer[0] == "#####" {
				locksStr = append(locksStr, buffer)
			} else {
				keysStr = append(keysStr, buffer)
			}
			buffer = []string{}
			continue
		}
		buffer = append(buffer, line)
	}

	keys := [][]int{}
	for _, key := range keysStr {
		heights := parseKeyToColumnHeight(key)
		keys = append(keys, heights)
	}

	locks := [][]int{}
	for _, lock := range locksStr {
		heights := parseLockToColumnHeight(lock)
		locks = append(locks, heights)
	}

	fits := 0
	for _, key := range keys {
		for _, lock := range locks {
			if doesKeyFitInLock(key, lock) {
				fits++
			}
		}
	}

	fmt.Println(fits)
}

func parseLockToColumnHeight(in []string) []int {
	out := []int{}
	for x := 0; x < 5; x++ {
		for y := 0; y < 7; y++ {
			if string(in[y][x]) == "." {
				out = append(out, y-1)
				break
			}
		}
	}
	return out
}

func parseKeyToColumnHeight(in []string) []int {
	out := []int{}
	for x := 0; x < 5; x++ {
		for y := 0; y < 7; y++ {
			if string(in[y][x]) == "#" {
				out = append(out, 6-y)
				break
			}
		}
	}
	return out
}

func doesKeyFitInLock(key []int, lock []int) bool {
	for i := 0; i < 5; i++ {
		if key[i]+lock[i] >= 6 {
			return false
		}
	}
	return true
}
