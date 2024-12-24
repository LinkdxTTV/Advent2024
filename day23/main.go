package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"
)

type node struct {
	name        string
	connections map[string]*node
}

func main() {
	file, err := os.Open("./input")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	connections := []string{}

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		connections = append(connections, line)
	}

	nodeMap := map[string]*node{}
	// Make nodes
	for _, connection := range connections {
		split := strings.Split(connection, "-")

		// Ensure nodes exist
		_, ok := nodeMap[split[0]]
		if !ok {
			nodeMap[split[0]] = &node{
				name:        split[0],
				connections: map[string]*node{},
			}
		}
		_, ok = nodeMap[split[1]]
		if !ok {
			nodeMap[split[1]] = &node{
				name:        split[1],
				connections: map[string]*node{},
			}
		}

		// Connect them
		nodeMap[split[0]].connections[split[1]] = nodeMap[split[1]]
		nodeMap[split[1]].connections[split[0]] = nodeMap[split[0]]
	}

	mapOf3Networks := map[string][]string{}
	// All size 3 networks
	for name, _ := range nodeMap {
		for _, net := range exploreAndReturnNetworksOfThree(name, nodeMap) {
			slices.Sort(net)
			mapOf3Networks[strings.Join(net, ",")] = net
		}
	}

	size3NetworksWithT := 0
	// fmt.Println(mapOf3Networks)
	for _, v := range mapOf3Networks {
		for _, computer := range v {
			if strings.HasPrefix(computer, "t") {
				// fmt.Println(v)
				size3NetworksWithT++
				break
			}
		}
	}

	fmt.Println(size3NetworksWithT)

	// Part 2, find largest totally connected network
	largest := []string{}
	for comp := range nodeMap {
		network := returnLargestNetwork(comp, nodeMap)
		if len(network) > len(largest) {
			largest = network
		}
	}

	slices.Sort(largest)
	fmt.Println(strings.Join(largest, ","))
}

func exploreAndReturnNetworksOfThree(in string, nodeMap map[string]*node) [][]string {
	networksOfThree := [][]string{}

	for connName, conn := range nodeMap[in].connections {
		connected := []string{in, connName}
		for connName2, conn2 := range conn.connections {
			if connName2 == connName {
				continue
			}
			connected2 := slices.Concat(connected, []string{connName2})
			for connName3, _ := range conn2.connections {
				if connName3 == connName2 {
					continue
				}
				if connName3 == in {
					networksOfThree = append(networksOfThree, connected2)
				}
			}
		}
	}

	return networksOfThree
}

// Do I need a new type here
type nodeSearch struct {
	node    string
	visited []string
}

// Should return largest network that (in string) is in
func returnLargestNetwork(in string, nodeMap map[string]*node) []string {
	// fmt.Println("starting at", in, nodeMap[in])
	largest := []string{}

	exploreQueue := []nodeSearch{
		{node: in,
			visited: []string{in},
		},
	}
	alreadyChecked := map[string]bool{}

	for len(exploreQueue) > 0 {
		currentNode := exploreQueue[0]
		exploreQueue = exploreQueue[1:]
		// fmt.Println("considering", currentNode)
		if alreadyChecked[currentNode.node] {
			// fmt.Println("skipping", currentNode.node)
			continue
		}
		alreadyChecked[currentNode.node] = true

		for connName, conn := range nodeMap[currentNode.node].connections {
			if slices.Contains(currentNode.visited, connName) {
				continue // We've already been here
			}

			isFullyConnected := true
			for _, visited := range currentNode.visited {
				_, ok := conn.connections[visited]
				if !ok {
					isFullyConnected = false
				}
			}
			if !isFullyConnected {
				continue
			}

			newVisited := slices.Clone(currentNode.visited)
			newVisited = append(newVisited, connName)

			if len(newVisited) > len(largest) {
				largest = newVisited
			}

			// exploreQueue = append(exploreQueue, nodeSearch{
			// 	node:    connName,
			// 	visited: newVisited,
			// })
			exploreQueue = slices.Concat([]nodeSearch{{
				node:    connName,
				visited: newVisited,
			}}, exploreQueue)
		}
	}

	return largest

}
