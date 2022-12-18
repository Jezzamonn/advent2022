package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Edge struct {
	DestName string
	Cost     int
}

type Node struct {
	Name  string
	Value int
	Edges map[string]*Edge
}

const startNodeName = "AA"

func parseFile(filename string) (nodes map[string]*Node) {
	// Read input
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// Parse the input into a graph
	nodes = make(map[string]*Node)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		// Parse an example line like this:
		// Valve GJ has flow rate=14; tunnels lead to valves UV, AO, MM, UD, GM
		// Using a regex to parse the line
		re := regexp.MustCompile(`Valve (\w+) has flow rate=(\d+); tunnel(s?) lead(s?) to valve(s?) (\w+(, \w+)*)`)
		matches := re.FindStringSubmatch(line)
		if matches == nil {
			panic("Invalid line: " + line)
		}

		// Create the node
		value, err := strconv.Atoi(matches[2])
		if err != nil {
			panic(err)
		}
		node := &Node{
			Name:  matches[1],
			Value: value,
			Edges: make(map[string]*Edge),
		}
		// Create the edges
		for _, destName := range strings.Split(matches[6], ", ") {
			node.Edges[destName] = &Edge{
				DestName: destName,
				Cost:     1,
			}
		}
		nodes[node.Name] = node
	}
	return nodes
}

func solve(filename string) {
	fmt.Println("Solving", filename)

	nodes := parseFile(filename)
	// For the moment, just print the nodes?
	for _, node := range nodes {
		fmt.Println(node)
	}

	// outputAsDotFile(nodes, "16/graph.dot")
	removeZeroValueNodes(nodes)

	// outputAsDotFile(nodes, "16/opt.dot")

	// Now, search through different paths to find the one with the best value.
	// We'll use a DFS to find the best path, treating each 'node' of the search as a possible path.
	// This may be too slow, so we'll see how it goes.

	// Right now this can revisit nodes which probably isn't going to lead to the optimal solution.
	minutesLeft := 30
	pt1 := getMostFlowPossible(nodes, make(map[string]struct{}), startNodeName, 0, minutesLeft, 0)

	fmt.Println("Part 1:", pt1)
}

// Recursive approach, too slow.
func getMostFlowPossible(nodes map[string]*Node, openedValves map[string]struct{}, startNode string, flowSoFar, minutesLeft, depth int) int {
	if minutesLeft <= 0 {
		// No time left to open any valves
		return flowSoFar
	}

	flowOptions := make([]int, 0)

	// If this valve isn't opened, we can try opening it.
	if _, ok := openedValves[startNode]; !ok {
		newOpenedValves := make(map[string]struct{})
		for k, v := range openedValves {
			newOpenedValves[k] = v
		}
		newOpenedValves[startNode] = struct{}{}

		if minutesLeft > 20 {
			fmt.Println(strings.Repeat(" ", depth), depth, "Opening valve", startNode, "at", minutesLeft, "minutes left. Flow so far:", flowSoFar)
		}

		flow := getMostFlowPossible(
			nodes,
			newOpenedValves,
			startNode,
			flowSoFar+nodes[startNode].Value*(minutesLeft-1),
			minutesLeft-1,
			depth+1)
		flowOptions = append(flowOptions, flow)
	}

	// Try going through each edge
	node := nodes[startNode]
	for _, edge := range node.Edges {
		if minutesLeft > 20 {
			fmt.Println(strings.Repeat(" ", depth), depth, "Going through", edge.DestName, "at", minutesLeft, "minutes left. Flow so far:", flowSoFar)
		}

		flow := getMostFlowPossible(
			nodes,
			openedValves,
			edge.DestName,
			flowSoFar,
			minutesLeft-edge.Cost,
			depth+1)
		flowOptions = append(flowOptions, flow)
	}

	// Now, find the best flow option
	bestFlow := 0
	for _, flow := range flowOptions {
		if flow > bestFlow {
			bestFlow = flow
		}
	}
	return bestFlow
}

func removeZeroValueNodes(nodes map[string]*Node) {
	// Optimize the graph by removing nodes with value of 0.
	// First get all the keys because we're going to modify the map
	// while iterating over it.
	nodeNames := make([]string, len(nodes))
	i := 0
	for nodeName := range nodes {
		nodeNames[i] = nodeName
		i++
	}

	for _, nodeName := range nodeNames {
		node := nodes[nodeName]
		if node.Value != 0 {
			continue
		}
		// Special case the start node because we need it for the solution.
		if node.Name == startNodeName {
			continue
		}
		deleteNode(nodes, nodeName)
	}
}

func deleteNode(nodes map[string]*Node, nodeName string) {
	node := nodes[nodeName]
	// Remove this node
	delete(nodes, nodeName)
	// Remove this node from the other node's edges
	for _, edge := range node.Edges {
		otherNode := nodes[edge.DestName]
		delete(otherNode.Edges, nodeName)
	}

	// Connect all the nodes that pass through this node
	// to each other.
	for _, firstEdge := range node.Edges {
		for _, secondEdge := range node.Edges {
			if firstEdge == secondEdge {
				continue
			}

			// Connect the two nodes
			firstNode := nodes[firstEdge.DestName]
			secondNode := nodes[secondEdge.DestName]

			// TODO: Check if the edge already exists, and maybe update the cost
			firstNode.Edges[secondNode.Name] = &Edge{
				DestName: secondNode.Name,
				Cost:     firstEdge.Cost + secondEdge.Cost,
			}
			secondNode.Edges[firstNode.Name] = &Edge{
				DestName: firstNode.Name,
				Cost:     firstEdge.Cost + secondEdge.Cost,
			}
		}
	}
}

func outputAsDotFile(nodes map[string]*Node, filename string) {
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// The graph is not directed.
	f.WriteString("digraph {\n")
	for _, node := range nodes {
		for _, edge := range node.Edges {
			f.WriteString(fmt.Sprintf("\t%s -> %s [label=%d];\n", node.Name, edge.DestName, edge.Cost))
		}
	}
	f.WriteString("}\n")
}

func main() {
	solve("16/input.txt")
}
