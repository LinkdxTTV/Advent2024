package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
)

type fileBlock struct {
	id        int
	files     int
	freeSpace int
}

func main() {
	file, err := os.Open("./input")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	scanner.Scan()
	line := scanner.Text()
	line += "0" // Fix the weird length encoding

	encoding := []string{}
	part2Encoding := []fileBlock{}
	id := 0

	for i := 0; i < len(line); i += 2 {
		num1, err := strconv.Atoi(string(line[i]))
		if err != nil {
			panic(err)
		}
		num2, err := strconv.Atoi(string(line[i+1]))
		if err != nil {
			panic(err)
		}

		for j := 0; j < num1; j++ {
			encoding = append(encoding, fmt.Sprintf("%d", id))
		}
		for j := 0; j < num2; j++ {
			encoding = append(encoding, ".")
		}

		part2Encoding = append(part2Encoding, fileBlock{
			id:        id,
			files:     num1,
			freeSpace: num2,
		})

		id++
	}

	// Two pointer swap approach
	i := 0
	j := len(encoding) - 1

	for i < j {
		if encoding[i] != "." {
			i++
			continue
		}
		if encoding[j] == "." {
			j--
			continue
		}
		encoding[i] = encoding[j]
		encoding[j] = "."
	}

	checksum := 0
	for i, numStr := range encoding {
		if numStr == "." {
			break
		}
		numInt, err := strconv.Atoi(numStr)
		if err != nil {
			panic(err)
		}
		checksum += i * numInt
	}

	fmt.Println(checksum)

	// Part 2
	// fmt.Println(part2Encoding)

	// Resort part 2
	j = len(part2Encoding) - 1
	for {
		if j == 0 {
			break
		}
		moved := false
		for i := 0; i < j; i++ {
			if part2Encoding[i].freeSpace >= part2Encoding[j].files {
				// We are going to move file j to where i's empty space is.
				// First add empty space to where j used to be
				temp := part2Encoding[j]
				part2Encoding[j-1].freeSpace += temp.files + temp.freeSpace

				part2Encoding = slices.Delete(part2Encoding, j, j+1)

				// Redo empty space on inserted file
				temp.freeSpace = part2Encoding[i].freeSpace - temp.files
				// Remove empty space
				part2Encoding[i].freeSpace = 0

				// Slice shenanigans
				part2Encoding = slices.Insert(part2Encoding, i+1, temp)

				// part2String := printEncoding(part2Encoding)
				// fmt.Println(part2String)
				moved = true
				break
			}
		}
		if !moved {
			j--
		}
	}

	// After resort
	// fmt.Println(part2Encoding)

	part2String := printEncoding(part2Encoding)

	checksum = 0
	for i, numStr := range part2String {
		if numStr == "." {
			continue
		}
		numInt, err := strconv.Atoi(numStr)
		if err != nil {
			panic(err)
		}
		checksum += i * numInt
	}

	fmt.Println(checksum)
}

// 10048198664969 was wrong

func printEncoding(part2Encoding []fileBlock) []string {
	part2String := []string{}
	for _, fileBlock := range part2Encoding {
		for i := 0; i < fileBlock.files; i++ {
			part2String = append(part2String, fmt.Sprintf("%d", fileBlock.id))
		}
		for i := 0; i < fileBlock.freeSpace; i++ {
			part2String = append(part2String, ".")
		}
	}

	return part2String
}
