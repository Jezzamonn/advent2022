package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

type Elf struct {
	Cur     Point
	New     Point
	Clashed bool
}

func parse(filename string) ([]*Elf, map[Point]*Elf) {
	f, err := os.Open(filename)
	if err != nil {
		f, err = os.Open("../" + filename)
		if err != nil {
			panic(err)
		}
	}

	elves := make([]*Elf, 0)
	elfMap := make(map[Point]*Elf)

	scanner := bufio.NewScanner(f)
	y := 0
	for scanner.Scan() {
		line := scanner.Text()
		for x, c := range line {
			if c == '#' {
				elf := Elf{Point{x, y}, Point{x, y}, false}
				elves = append(elves, &elf)
				elfMap[elf.Cur] = &elf
			}
		}
		y++
	}
	return elves, elfMap
}

func simulateStep(elves []*Elf, elfMap map[Point]*Elf, directions []Direction) int {
	// Calculate where to move to.
	for _, elf := range elves {
		// Check if any adjacent cells are occupied
		hasNeighbour := false
		for _, d := range elf.Cur.Adjacent() {
			if _, ok := elfMap[d]; ok {
				hasNeighbour = true
				break
			}
		}

		if !hasNeighbour {
			continue
		}

		// We have a neighbour, now we need to check each direction in order.
		for _, d := range directions {
			allEmpty := true
			for _, p := range d.DeltasToCheck() {
				if _, ok := elfMap[elf.Cur.Add(p)]; ok {
					allEmpty = false
					break
				}
			}
			if !allEmpty {
				continue
			}
			elf.New = elf.Cur.Add(d.Delta())
			break
		}
	}
	// Write the new positions to the map, but leave the old positions for the moment. So we can detect clashes.
	for _, elf := range elves {
		if elf.New == elf.Cur {
			continue
		}
		if clashed, ok := elfMap[elf.New]; ok {
			elf.Clashed = true
			clashed.Clashed = true
			continue
		}
		elfMap[elf.New] = elf
	}
	elvesMoved := 0
	// Remove the old positions from the map.
	for _, elf := range elves {
		if elf.New == elf.Cur {
			continue
		}
		if elf.Clashed {
			delete(elfMap, elf.New)
			elf.Clashed = false
			elf.New = elf.Cur
		} else {
			delete(elfMap, elf.Cur)
			elf.Cur = elf.New
			elvesMoved++
		}
	}
	return elvesMoved
}

func getBounds(elves []*Elf) (Point, Point) {
	min := Point{0, 0}
	max := Point{0, 0}
	for _, elf := range elves {
		if elf.Cur.X < min.X {
			min.X = elf.Cur.X
		}
		if elf.Cur.X > max.X {
			max.X = elf.Cur.X
		}
		if elf.Cur.Y < min.Y {
			min.Y = elf.Cur.Y
		}
		if elf.Cur.Y > max.Y {
			max.Y = elf.Cur.Y
		}
	}
	return min, max
}

func print(elves []*Elf, elfMap map[Point]*Elf) {
	min, max := getBounds(elves)
	// Round the min and max to the nearest 10.
	min.X = min.X/10*10 - 10
	min.Y = min.Y/10*10 - 10
	max.X = max.X/10*10 + 10
	max.Y = max.Y/10*10 + 10

	for y := min.Y; y <= max.Y; y++ {
		for x := min.X; x <= max.X; x++ {
			if elf, ok := elfMap[Point{x, y}]; ok {
				if elf.Clashed {
					fmt.Print("X")
				} else {
					fmt.Print("#")
				}
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
}

func printAnimationFrame(elves []*Elf, elfMap map[Point]*Elf, step int) {
	// Clear the screen.
	fmt.Print("\033[H\033[2J")
	print(elves, elfMap)
	fmt.Println("Step", step)
	// Sleep for a bit.
	time.Sleep(100 * time.Millisecond)
}

func solve(filename string) {
	fmt.Println("Solving", filename)
	elves, elfMap := parse(filename)

	animate := false

	directions := []Direction{North, South, West, East}
	step := 0
	for {
		if animate {
			printAnimationFrame(elves, elfMap, step)
		}
		numMoved := simulateStep(elves, elfMap, directions)
		step++
		if numMoved == 0 {
			break
		}
		// Move the first direction to the end.
		directions = append(directions[1:], directions[0])
	}
	if animate {
		printAnimationFrame(elves, elfMap, step)
	}

	min, max := getBounds(elves)
	totalArea := (max.X - min.X + 1) * (max.Y - min.Y + 1)
	emptySpaces := totalArea - len(elves)
	fmt.Println("Empty spaces:", emptySpaces)
	fmt.Println("Num steps:", step)

	if animate {
		// Sleep a little so we can see the result
		time.Sleep(3000 * time.Millisecond)
	}
}

func main() {
	solve("23/demo2.txt")
	solve("23/demo.txt")
	solve("23/input.txt")
}
