package main

import (
	"bufio"
	"fmt"
	"os"
)

func solve(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	sum := 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		bq := BalancedQuinary(scanner.Text())
		sum += bq.Int()
	}

	fmt.Println("Part 1:", BalancedQuinaryFromInt(sum))
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide a filename")
		return
	}
	solve(os.Args[1])
}
