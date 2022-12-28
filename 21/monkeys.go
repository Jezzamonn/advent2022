package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

type HasValue interface {
	GetValue(values map[string]HasValue) float64
	GetGraphvizRepresentation(values map[string]HasValue) string
}

type Leaf struct {
	Name  string
	Value float64
}

func (l *Leaf) GetValue(values map[string]HasValue) float64 {
	return l.Value
}

func (l *Leaf) GetGraphvizRepresentation(values map[string]HasValue) string {
	// Label this node with it's value
	return fmt.Sprintln(l.Name, "[label=\"", l.Name, " (", l.Value, ")\"];")
}

type Node struct {
	Name       string
	Child1Name string
	Child2Name string
	Operation  string
	// HasCachedValue bool
	// CachedValue    int
}

func (n *Node) GetValue(values map[string]HasValue) float64 {
	// // If we've already calculated this value, return it
	// if n.HasCachedValue {
	// 	return n.CachedValue
	// }
	value1 := values[n.Child1Name].GetValue(values)
	value2 := values[n.Child2Name].GetValue(values)
	return n.calculateValue(value1, value2)
	// n.CachedValue = n.calculateValue(value1, value2)
	// n.HasCachedValue = true
	// return n.CachedValue
}

func (n *Node) calculateValue(value1 float64, value2 float64) float64 {
	switch n.Operation {
	case "+":
		return value1 + value2
	case "-":
		return value1 - value2
	case "*":
		return value1 * value2
	case "/":
		return value1 / value2
	default:
		panic("Unknown operation: " + n.Operation)
	}
}

func (n *Node) GetGraphvizRepresentation(values map[string]HasValue) string {
	// Return a label for this node with it's operation and it's value, and edges to it's children
	// Warning: Now that I got rid of the cache, this is SLOW.
	return fmt.Sprintf("%s [label=\"%s (%s) = %d\"];\n", n.Name, n.Name, n.Operation, n.GetValue(values)) +
		fmt.Sprintf("%s -> %s;\n", n.Name, n.Child1Name) +
		fmt.Sprintf("%s -> %s;\n", n.Name, n.Child2Name)
}

func solve(filename string) {
	fmt.Println("Solving", filename, "...")

	values := parseInput(filename)

	outputGraphviz(filename, values)

	// Part 1:
	// Print the value of the root node
	fmt.Println("Part 1:")
	fmt.Println(values["root"].GetValue(values))

	// Part 2:
	//
	// Lazy approach: Modify the graph, knowing the solution will be linear. We
	// can then modify humn, and use that to figure out what value for humn is
	// needed to make root zero.

	// Ah, but it fails because of integer division :(

	values["root"].(*Node).Operation = "-"
	values["humn"].(*Leaf).Value = 0

	whenHumnIsZero := values["root"].GetValue(values)
	fmt.Println("whenHumnIsZero:", whenHumnIsZero)

	values["humn"].(*Leaf).Value = 1
	whenHumnIsOne := values["root"].GetValue(values)
	fmt.Println("whenHumnIsOne:", whenHumnIsOne)

	// Equation of form y = mx + b
	m := whenHumnIsOne - whenHumnIsZero
	b := whenHumnIsZero

	// Solve for y = 0
	humn := -b / m
	fmt.Println("Part 2:")
	fmt.Println(humn)
	fmt.Println(int(humn))
}

// Print the graph in Graphviz format
func outputGraphviz(filename string, values map[string]HasValue) {
	graphvizFilename := filename + ".dot"
	graphvizFile, err := os.Create(graphvizFilename)
	if err != nil {
		panic(err)
	}
	defer graphvizFile.Close()
	graphvizFile.WriteString("digraph G {\n")
	for _, value := range values {
		graphvizFile.WriteString(value.GetGraphvizRepresentation(values))
	}
	graphvizFile.WriteString("}\n")
}

func parseInput(filename string) map[string]HasValue {
	// Read input
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// Parse the input into a graph
	values := make(map[string]HasValue)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		// Parse the line with regex.
		// Example lines:
		// `root: pppw + sjmn`
		// `sjmn: 2`
		leafRe := regexp.MustCompile(`^([a-z]+): ([0-9]+)$`)
		nodeRe := regexp.MustCompile(`^([a-z]+): ([a-z]+) ([+*-/]) ([a-z]+)$`)
		nodeMatch := nodeRe.FindStringSubmatch(line)
		if nodeMatch != nil {
			node := &Node{
				Name:       nodeMatch[1],
				Child1Name: nodeMatch[2],
				Child2Name: nodeMatch[4],
				Operation:  nodeMatch[3],
			}
			values[nodeMatch[1]] = node
			continue
		}
		leafMatch := leafRe.FindStringSubmatch(line)
		if leafMatch != nil {
			value, err := strconv.ParseFloat(leafMatch[2], 64)
			if err != nil {
				panic(err)
			}
			values[leafMatch[1]] = &Leaf{
				Name:  leafMatch[1],
				Value: value,
			}
			continue
		}
		panic("Unknown line")
	}
	return values
}

func main() {
	solve("21/demo.txt")
	solve("21/input.txt")
}
