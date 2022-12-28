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
	value1, err1 := values[n.Child1Name].GetValue(values)
	if err1 != nil {
		return 0, err1
	}
	value2, err2 := values[n.Child2Name].GetValue(values)
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
	// Warning: Now that I got rid of the cache, this is SLOW.
	value, _ := n.GetValue(values)
	return fmt.Sprintf("%s [label=\"%s (%s) = %d\"];\n", n.Name, n.Name, n.Operation, value) +
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

	values["root"].(*Node).Operation = "-"

	var firstValidInputValue int = -1
	var firstValidOutputValue int
	var secondValidInputValue int = -1
	var secondValidOutputValue int
	for i := 0; secondValidInputValue == -1; i++ {
		values["humn"].(*Leaf).Value = i
		outputValue, err := values["root"].GetValue(values)
		if err != nil {
			continue
		}
		if firstValidInputValue == -1 {
			firstValidInputValue = i
			firstValidOutputValue = outputValue
			fmt.Println("firstValidInputValue:", firstValidInputValue, "firstValidOutputValue:", firstValidOutputValue)
		} else {
			secondValidInputValue = i
			secondValidOutputValue = outputValue
			fmt.Println("secondValidInputValue:", secondValidInputValue, "secondValidOutputValue:", secondValidOutputValue)
		}
	}

	// Solve the linear equation to find input value when the output is 0.
	// y = (p / q)x + b
	p := firstValidOutputValue - secondValidOutputValue
	q := firstValidInputValue - secondValidInputValue
	b := firstValidOutputValue - (p * firstValidInputValue / q)
	var humn int
	if p < q {
		humn = (0 - b) * q / p
	} else {
		humn = (0 - b) / (p / q)
	}

	fmt.Println("Part 2:")
	fmt.Println(humn)

	// Verify the solution.
	values["humn"].(*Leaf).Value = humn
	outputValue, err := values["root"].GetValue(values)
	if err != nil {
		panic(err)
	}
	if outputValue != 0 {
		panic("Solution is not correct!")
	}
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
