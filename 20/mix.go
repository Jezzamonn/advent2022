package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

type Node struct {
	Value int
	Prev  *Node
	Next  *Node
}

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func nodesToString(head *Node) string {
	node := head
	str := ""
	for {
		str += fmt.Sprintf("%d, ", node.Value)
		node = node.Next
		if node == head {
			break
		}
	}
	return str
}

func nodesToStringBackwards(head *Node) string {
	node := head
	str := ""
	for {
		str += fmt.Sprintf("%d, ", node.Value)
		node = node.Prev
		if node == head {
			break
		}
	}
	return str
}

func validateLength(head *Node, expectedLength int) {
	forwardLength := getForwardLength(head)
	backwardLength := getBackwardLength(head)
	if forwardLength != backwardLength {
		panic(fmt.Sprintf("Forward length is %d, backward length is %d", forwardLength, backwardLength))
	}
	if forwardLength != expectedLength {
		panic(fmt.Sprintf("Length is %d, expected %d", forwardLength, expectedLength))
	}
}

func validatePointers(head *Node) {
	node := head
	for {
		if node.Next.Prev != node {
			panic("Next pointer is not correct")
		}
		if node.Prev.Next != node {
			panic("Prev pointer is not correct")
		}
		node = node.Next
		if node == head {
			break
		}
	}
}

func getForwardLength(head *Node) int {
	node := head
	length := 0
	for {
		length++
		node = node.Next
		if node == head {
			break
		}
		if length > 10000 {
			// Invalid loop
			return -1
		}
	}
	return length
}

func getBackwardLength(head *Node) int {
	node := head
	length := 0
	for {
		length++
		node = node.Prev
		if node == head {
			break
		}
		if length > 10000 {
			// Invalid loop
			return -1
		}
	}
	return length
}

func solvePt1(filename string) {
	fmt.Println("Solving", filename, "part 1...")

	nodes := parseToList(filename)
	head := nodes[0]

	// fmt.Println("Initial configuration:")
	// fmt.Println(nodesToString(head))
	// fmt.Println(nodesToStringBackwards(head.Prev))
	// fmt.Println("Length:", len(nodes))

	validateLength(head, len(nodes))
	validatePointers(head)

	mix(nodes)

	value := getListValue(nodes)
	fmt.Println(value)
}

func solvePt2(filename string) {
	fmt.Println("Solving", filename, "part 2...")

	nodes := parseToList(filename)
	head := nodes[0]

	// fmt.Println("Initial configuration:")
	// fmt.Println(nodesToString(head))
	// fmt.Println(nodesToStringBackwards(head.Prev))
	// fmt.Println("Length:", len(nodes))

	validateLength(head, len(nodes))
	validatePointers(head)

	for _, node := range nodes {
		node.Value *= 811589153
	}

	for i := 0; i < 10; i++ {
		mix(nodes)
	}

	value := getListValue(nodes)
	fmt.Println(value)
}

func parseToList(filename string) []*Node {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// Create a doubly linked list
	var head *Node
	var tail *Node
	nodes := make([]*Node, 0)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		value, err := strconv.Atoi(line)
		if err != nil {
			panic(err)
		}
		node := &Node{Value: value}
		if head == nil {
			head = node
			tail = node
		} else {
			tail.Next = node
			node.Prev = tail
			tail = node
		}
		nodes = append(nodes, node)
	}
	// Join the start and the end.
	head.Prev = tail
	tail.Next = head

	return nodes
}

func mix(nodes []*Node) {
	// Do one mix
	for _, node := range nodes {
		// Move the node forward or backward.
		shuffleAmount := Abs(node.Value) % (len(nodes) - 1)
		// fmt.Println("Shuffle amount:", shuffleAmount)

		for i := 0; i < shuffleAmount; i++ {
			if node.Value > 0 {
				oldPrev := node.Prev
				oldNext := node.Next

				newPrev := node.Next
				newNext := node.Next.Next

				// Remove node from old position
				oldPrev.Next = oldNext
				oldNext.Prev = oldPrev

				// Insert node into new position
				newPrev.Next = node
				node.Prev = newPrev
				newNext.Prev = node
				node.Next = newNext
			} else {
				oldPrev := node.Prev
				oldNext := node.Next

				newPrev := node.Prev.Prev
				newNext := node.Prev

				// Remove node from old position
				oldPrev.Next = oldNext
				oldNext.Prev = oldPrev

				// Insert node into new position
				newPrev.Next = node
				node.Prev = newPrev
				newNext.Prev = node
				node.Next = newNext
			}
		}

		// fmt.Println(node.Value, "moves between", node.Prev.Value, "and", node.Next.Value)
		// fmt.Println(nodesToString(head))
		// fmt.Println()

		head := nodes[0]
		validateLength(head, len(nodes))
		validatePointers(head)
	}
}

func getListValue(nodes []*Node) int {
	// For the solution, we need to find the 0 node.
	var zeroNode *Node
	for _, node := range nodes {
		if node.Value == 0 {
			zeroNode = node
			break
		}
	}

	sum := 0
	for _, dist := range []int{1000, 2000, 3000} {
		dist %= len(nodes)
		// Move dist nodes forward from zero, add to sum.
		node := zeroNode
		for i := 0; i < dist; i++ {
			node = node.Next
		}
		fmt.Println("Value:", node.Value)
		sum += node.Value
	}
	return sum
}

func main() {
	solvePt1("20/input.txt")
	solvePt2("20/demo.txt")
	solvePt2("20/input.txt")
}
