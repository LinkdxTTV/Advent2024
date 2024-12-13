package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

type machine struct {
	ax int
	ay int
	bx int
	by int
	px int
	py int
}

func main() {
	file, err := os.Open("./input")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	buttonRegexp := regexp.MustCompile(`Button (.): X\+([0-9]+), Y\+([0-9]+)`)
	priceRegexp := regexp.MustCompile(`Prize: X=([0-9]+), Y=([0-9]+)`)
	machines := []machine{}

	var ax, ay, bx, by, px, py int = 0, 0, 0, 0, 0, 0
	var Machine machine
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			ax, ay, bx, by, px, py = 0, 0, 0, 0, 0, 0
			continue
		}
		finds := buttonRegexp.FindAllStringSubmatch(line, -1)
		if len(finds) == 1 {
			if finds[0][1] == "A" {
				Machine = machine{}
				ax, err = strconv.Atoi(finds[0][2])
				if err != nil {
					panic(err)
				}
				ay, err = strconv.Atoi(finds[0][3])
				if err != nil {
					panic(err)
				}
				Machine.ax = ax
				Machine.ay = ay
			} else {
				bx, err = strconv.Atoi(finds[0][2])
				if err != nil {
					panic(err)
				}
				by, err = strconv.Atoi(finds[0][3])
				if err != nil {
					panic(err)
				}
				Machine.bx = bx
				Machine.by = by
			}
		} else {
			finds = priceRegexp.FindAllStringSubmatch(line, -1)
			if len(finds) == 1 {
				px, err = strconv.Atoi(finds[0][1])
				if err != nil {
					panic(err)
				}
				py, err = strconv.Atoi(finds[0][2])
				if err != nil {
					panic(err)
				}
				Machine.px = px
				Machine.py = py
				machines = append(machines, Machine)
			}
		}
	}
	// fmt.Println(machines)

	sum := 0
	for _, mach := range machines {
		possibleTokens := []int{}
		for i := 0; i < 100; i++ {
			for j := 0; j < 100; j++ {
				if i*mach.ax+j*mach.bx > mach.px || i*mach.ay+j*mach.by > mach.py {
					break
				}
				if i*mach.ax+j*mach.bx == mach.px && i*mach.ay+j*mach.by == mach.py {
					possibleTokens = append(possibleTokens, 3*i+j)
					break
				}
			}
		}
		if len(possibleTokens) == 1 {
			sum += possibleTokens[0]
		}
	}

	fmt.Println(sum)

	// Part 2.. Need to think of a new algorithm to solve this.
	tokens := []int{}

	for _, mach := range machines {
		mach.px += 10000000000000
		mach.py += 10000000000000
		remainder := (mach.px*mach.ay - mach.ax*mach.py) % (mach.bx*mach.ay - mach.ax*mach.by)
		if remainder != 0 {
			continue
		}
		b := (mach.px*mach.ay - mach.ax*mach.py) / (mach.bx*mach.ay - mach.ax*mach.by)
		remainder = (mach.px - b*mach.bx) % mach.ax
		if remainder != 0 {
			continue
		}
		a := (mach.px - b*mach.bx) / mach.ax

		tokens = append(tokens, 3*a+b)
	}

	// fmt.Println(tokens)
	sum = 0
	for _, token := range tokens {
		sum += token
	}
	fmt.Println(sum)
}
