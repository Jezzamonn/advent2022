package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

type HasValue interface {
	GetValue(values map[string]HasValue) (int, error)
	GetGraphvizRepresentation(values map[string]HasValue) string
}

type Leaf struct {
	Name  string
	Value int
}

func (l *Leaf) GetValue(values map[string]HasValue) (int, error) {
	return l.Value, nil
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

func (n *Node) GetValue(values map[string]HasValue) (int, error) {
	// // If we've already calculated this value, return it
	// if n.HasCachedValue {
	// 	return n.CachedValue
	// }
	child1, ok := values[n.Child1Name]
	if !ok {
		return 0, fmt.Errorf("Node %s: Child not found: %s", n, n.Child1Name)
	}

	if n.Operation == "eq" {
		// Lazy node type: just return the value of the first child
		return child1.GetValue(values)
	}

	child2, ok := values[n.Child2Name]
	if !ok {
		return 0, fmt.Errorf("Node %s: Child not found: %s", n, n.Child2Name)
	}
	value1, err1 := child1.GetValue(values)
	if err1 != nil {
		return 0, err1
	}
	value2, err2 := child2.GetValue(values)
	if err2 != nil {
		return 0, err2
	}
	return n.calculateValue(value1, value2)
	// n.CachedValue = n.calculateValue(value1, value2)
	// n.HasCachedValue = true
	// return n.CachedValue
}

func (n *Node) calculateValue(value1 int, value2 int) (int, error) {
	switch n.Operation {
	case "+":
		return value1 + value2, nil
	case "-":
		return value1 - value2, nil
	case "*":
		return value1 * value2, nil
	case "/":
		if value2 == 0 {
			return 0, fmt.Errorf("Division by zero: %d / %d", value1, value2)
		}
		if value1%value2 != 0 {
			return 0, fmt.Errorf("Division is not an integer: %d / %d", value1, value2)
		}
		return value1 / value2, nil
	default:
		panic("Unknown operation: " + n.Operation)
	}
}

func (n *Node) GetGraphvizRepresentation(values map[string]HasValue) string {
	// Return a label for this node with it's operation and it's value, and edges to it's children
	// value, _ := n.GetValue(values)
	out := fmt.Sprintf("%s [label=\"%s (%s)\"];\n", n.Name, n.Name, n.Operation)
	out += fmt.Sprintf("%s -> %s;\n", n.Name, n.Child1Name)
	if n.Operation != "eq" {
		out += fmt.Sprintf("%s -> %s;\n", n.Name, n.Child2Name)
	}
	return out
}

func (n *Node) Invert(focus string) {
	var other string
	if n.Child1Name == focus {
		other = n.Child2Name
	} else if n.Child2Name == focus {
		other = n.Child1Name
	} else {
		panic("Child not found: " + focus)
	}

	switch n.Operation {
	case "+":
		// y = x + a -> x = y - a
		n.Name, n.Child1Name, n.Child2Name = focus, n.Name, other
		n.Operation = "-"
	case "-":
		// Two cases:
		if n.Child1Name == focus {
			// y = x - a -> x = y + a
			n.Name, n.Child1Name, n.Child2Name = focus, n.Name, other
			n.Operation = "+"
		} else {
			// y = a - x -> x = a - y
			n.Name, n.Child1Name, n.Child2Name = focus, other, n.Name
			n.Operation = "-"
		}
	case "*":
		// y = x * a -> x = y / a
		n.Name, n.Child1Name, n.Child2Name = focus, n.Name, other
		n.Operation = "/"
	case "/":
		// Two cases:
		if n.Child1Name == focus {
			// y = x / a -> x = y * a
			n.Name, n.Child1Name, n.Child2Name = focus, n.Name, other
			n.Operation = "*"
		} else {
			// y = a / x -> x = a / y
			n.Name, n.Child1Name, n.Child2Name = focus, other, n.Name
			n.Operation = "/"
		}
	case "eq":
		// This is a special case. Make this return the other child.
		n.Name, n.Child1Name, n.Child2Name = focus, other, ""
		n.Operation = "eq"
	default:
		panic("Unknown operation: " + n.Operation)
	}
}

func (n *Node) String() string {
	return fmt.Sprintf("%s: %s %s %s", n.Name, n.Child1Name, n.Operation, n.Child2Name)
}

func solve(filename string) {
	fmt.Println("Solving", filename, "...")

	values := parseInput(filename)

	// outputGraphviz(filename, values)

	// Part 1:
	// Print the value of the root node
	fmt.Println("Part 1:")
	fmt.Println(values["root"].GetValue(values))

	// Part 2:
	root := values["root"].(*Node)
	root.Operation = "eq"

	fmt.Println("inverting...")

	nextToInvert := "humn"
	nodesToInvert := make([]*Node, 0)
	for {
		var n *Node
		for _, value := range values {
			nn, ok := value.(*Node)
			if !ok {
				continue
			}
			if nn.Child1Name == nextToInvert || nn.Child2Name == nextToInvert {
				if n != nil {
					panic("Found two nodes to invert")
				}
				n = nn
			}
		}
		if n == nil {
			panic("Node not found")
		}
		nodesToInvert = append(nodesToInvert, n)
		nextToInvert = n.Name
		if n.Name == "root" {
			break
		}
	}

	nextToInvert = "humn"
	for _, n := range nodesToInvert {
		fmt.Println("Inverting", n.Name, "to", nextToInvert)

		oldName := n.Name
		delete(values, oldName)
		n.Invert(nextToInvert)
		values[n.Name] = n
		nextToInvert = oldName

		if nextToInvert == "root" {
			break
		}
	}
	// Lets just output as dot file and see what we did.
	// outputGraphviz(filename+"2", values)

	// Print the value of the humn node
	fmt.Println("Part 2:")
	fmt.Println(values["humn"].GetValue(values))
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
			value, err := strconv.Atoi(leafMatch[2])
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
