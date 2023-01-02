package main

import (
	"bufio"
	"container/heap"
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

// Parse an example string like this:
// Valve GJ has flow rate=14; tunnels lead to valves UV, AO, MM, UD, GM
func NodeFromString(s string) *Node {
	// Using a regex to parse the line
	re := regexp.MustCompile(`Valve (\w+) has flow rate=(\d+); tunnel(s?) lead(s?) to valve(s?) (\w+(, \w+)*)`)
	matches := re.FindStringSubmatch(s)
	if matches == nil {
		panic("Invalid line: " + s)
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
	return node
}

type NamePair struct {
	Name1 string
	Name2 string
}

type Graph struct {
	Nodes     map[string]*Node
	Distances map[NamePair]int
}

func (g *Graph) RemoveZeroValueNodes() {
	// Optimize the graph by removing nodes with value of 0.
	// First get all the keys because we're going to modify the map
	// while iterating over it.
	nodeNames := make([]string, len(g.Nodes))
	i := 0
	for nodeName := range g.Nodes {
		nodeNames[i] = nodeName
		i++
	}

	for _, nodeName := range nodeNames {
		node := g.Nodes[nodeName]
		if node.Value != 0 {
			continue
		}
		// Special case the start node because we need it for the solution.
		if node.Name == startNodeName {
			continue
		}
		g.DeleteNode(nodeName)
	}
}

func (g *Graph) DeleteNode(nodeName string) {
	node := g.Nodes[nodeName]
	// Remove this node
	delete(g.Nodes, nodeName)
	// Remove this node from the other node's edges
	for _, edge := range node.Edges {
		otherNode := g.Nodes[edge.DestName]
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
			firstNode := g.Nodes[firstEdge.DestName]
			secondNode := g.Nodes[secondEdge.DestName]

			newCost := firstEdge.Cost + secondEdge.Cost

			// Check if the edge already exists, and maybe update the cost
			if _, ok := firstNode.Edges[secondNode.Name]; ok {
				firstNode.Edges[secondNode.Name].Cost = intMin(firstNode.Edges[secondNode.Name].Cost, newCost)
				secondNode.Edges[firstNode.Name].Cost = intMin(secondNode.Edges[firstNode.Name].Cost, newCost)
				continue
			}

			firstNode.Edges[secondNode.Name] = &Edge{
				DestName: secondNode.Name,
				Cost:     newCost,
			}
			secondNode.Edges[firstNode.Name] = &Edge{
				DestName: firstNode.Name,
				Cost:     newCost,
			}
		}
	}
}

func (g *Graph) CalculateDistancesToEachNode() {
	distances := make(map[NamePair]int)
	// Do a search from each node
	for _, node := range g.Nodes {
		visited := make(map[string]bool)

		toVisit := NodePriorityQueue{}

		heap.Push(&toVisit, NodeDistPair{Node: node, Dist: 0})

		for len(toVisit) > 0 {
			// Pop the first node
			currentNode := toVisit[0].Node
			currentDist := toVisit[0].Dist
			toVisit = toVisit[1:]

			// Check if we've already visited this node
			if visited[currentNode.Name] {
				continue
			}

			distances[NamePair{
				Name1: node.Name,
				Name2: currentNode.Name,
			}] = currentDist
			visited[currentNode.Name] = true

			// Add the edges to the list of nodes to visit
			for _, edge := range currentNode.Edges {
				heap.Push(&toVisit, NodeDistPair{
					Node: g.Nodes[edge.DestName],
					Dist: currentDist + edge.Cost})
			}
		}
	}
	g.Distances = distances
}

func (g *Graph) PrintNodes() {
	for _, node := range g.Nodes {
		fmt.Println(node)
	}
}

func (g *Graph) OutputAsDotFile(filename string) {
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// The graph is not directed.
	f.WriteString("digraph {\n")
	for _, node := range g.Nodes {
		for _, edge := range node.Edges {
			f.WriteString(fmt.Sprintf("\t%s -> %s [label=%d];\n", node.Name, edge.DestName, edge.Cost))
		}
	}
	f.WriteString("}\n")
}

func GraphFromFile(filename string) *Graph {
	// Read input
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// Parse the input into a graph
	graph := &Graph{
		Nodes: make(map[string]*Node),
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		node := NodeFromString(line)
		graph.Nodes[node.Name] = node
	}

	// Optimize the graph
	graph.RemoveZeroValueNodes()

	// Calculate the distances between each node
	graph.CalculateDistancesToEachNode()

	return graph
}
