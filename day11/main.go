package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type linkedNode struct {
	value int
	left  *linkedNode
	right *linkedNode
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

	start := []int{}
	split := strings.Split(line, " ")
	for _, num := range split {
		numInt, err := strconv.Atoi(num)
		if err != nil {
			panic(err)
		}
		start = append(start, numInt)
	}

	fmt.Println(start)
	stoneCounter := len(start)
	fmt.Println("stone counter:", stoneCounter)

	root := linkedNode{value: start[0]}
	var lastNode *linkedNode = &root
	// Assemble the linked list because I feel like we will need it
	for i, num := range start {
		if i == 0 {
			continue
		}
		newNode := linkedNode{
			value: num,
			left:  lastNode,
			right: nil,
		}
		lastNode.right = &newNode
		lastNode = &newNode
	}

	for blink := 0; blink < 25; blink++ {
		fmt.Println("blink", blink)
		node := &root
		// Iterate over the list
		for node != nil {
			// Do something
			if node.value == 0 {
				node.value = 1
			} else if hasEvenNumberOfDigits(node.value) {
				// Split the node
				leftNum, rightNum := splitNumberIntoTwo(node.value)
				node.value = leftNum

				// Insert new node
				oldRight := node.right
				newNode := &linkedNode{
					value: rightNum,
					left:  node,
					right: oldRight,
				}
				node.right = newNode
				// Incase this was one of the last elements
				if oldRight != nil {
					oldRight.left = newNode
				}
				stoneCounter++
				// Do not process this next node, it is simultaneous, so we skip it
				node = node.right
			} else {
				node.value = node.value * 2024
			}

			node = node.right
		}
		// PrintList(root)
	}

	fmt.Println(stoneCounter)

	// Part 2 thonkulating
	numberMap := map[int]int{}
	for _, num := range start {
		numberMap[num] = 1
	}

	fmt.Println(numberMap)

	for blink := 0; blink < 75; blink++ {
		tempMap := map[int]int{}
		for k, v := range numberMap {
			if k == 0 {
				tempMap[1] += v
			} else if hasEvenNumberOfDigits(k) {
				left, right := splitNumberIntoTwo(k)
				tempMap[left] += v
				tempMap[right] += v
			} else {
				tempMap[k*2024] += v
			}
		}
		numberMap = tempMap
	}
	sum := 0
	for _, v := range numberMap {
		sum += v
	}

	fmt.Println(sum)
}

func hasEvenNumberOfDigits(num int) bool {
	return len(fmt.Sprintf("%d", num))%2 == 0
}

func splitNumberIntoTwo(num int) (int, int) {
	asString := fmt.Sprintf("%d", num)
	length := len(asString)
	left := asString[:length/2]
	right := asString[length/2:]
	// fmt.Println("split", num, left, right)

	leftInt, err := strconv.Atoi(left)
	if err != nil {
		panic(err)
	}
	rightInt, err := strconv.Atoi(right)
	if err != nil {
		panic(err)
	}

	return leftInt, rightInt
}

func PrintList(root linkedNode) {
	node := &root
	for node != nil {
		fmt.Print(node.value, " ")
		node = node.right
	}
	fmt.Println()
}
