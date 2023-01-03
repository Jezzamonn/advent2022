package main

import (
	"container/heap"
	"fmt"
	"os"
	"time"
)

type Point struct {
	X, Y int
}

func (p Point) DistTo(o Point) int {
	return Abs(p.X-o.X) + Abs(p.Y-o.Y)
}

func (p Point) Add(o Point) Point {
	return Point{p.X + o.X, p.Y + o.Y}
}

type Position struct {
	P    Point
	T    int
	Prev *Position
}

type PositionsByDistanceToEnd struct {
	Positions []Position
	End       Point
}

func (p PositionsByDistanceToEnd) Len() int {
	return len(p.Positions)
}

func (p PositionsByDistanceToEnd) Less(i, j int) bool {
	return p.Positions[i].AStarValue(p.End) < p.Positions[j].AStarValue(p.End)
}

func (p PositionsByDistanceToEnd) Swap(i, j int) {
	p.Positions[i], p.Positions[j] = p.Positions[j], p.Positions[i]
}

func (p *PositionsByDistanceToEnd) Push(x interface{}) {
	item := x.(Position)
	p.Positions = append(p.Positions, item)
}

func (p *PositionsByDistanceToEnd) Pop() interface{} {
	old := p.Positions
	n := len(old)
	item := old[n-1]
	p.Positions = old[0 : n-1]
	return item
}

// A delta for moving each direction, and one for staying still.
var deltas = []Position{
	{P: Point{X: 0, Y: -1}, T: 1},
	{P: Point{X: 1, Y: 0}, T: 1},
	{P: Point{X: 0, Y: 1}, T: 1},
	{P: Point{X: -1, Y: 0}, T: 1},
	{P: Point{X: 0, Y: 0}, T: 1},
}

func (p Position) AStarValue(End Point) int {
	return p.P.DistTo(End) + p.T
}

func animate(b *Blizzard) {
	// :)
	for i := 0; true; i = (i + 1) % b.Period() {
		// Clear the console
		fmt.Print("\033[H\033[2J")
		fmt.Print(b.StringAtTime(i))
		// Sleep a little
		time.Sleep(100 * time.Millisecond)
	}
}

// Do an A* search
func search(b *Blizzard, start Position, end Point) Position {
	// Just use a regular list, we can do a priority queue later if needed.
	visited := make(map[Position]bool)
	// For debugging
	// visitedPoints := make(map[Point]bool)

	toVisit := PositionsByDistanceToEnd{Positions: []Position{start}, End: end}
	i := 0
	for toVisit.Len() > 0 {
		// Pop the first element
		p := heap.Pop(&toVisit).(Position)

		// Check if this is a valid position
		if !b.IsEmptyAtTime(p.P.X, p.P.Y, p.T) {
			continue
		}

		// If we've already visited this position, skip it.
		// To avoid getting caught in a loop, create a new position mod the board period.
		pModPeriod := p
		pModPeriod.T = p.T % b.Period()
		// Oh and we have to clear the prev node
		pModPeriod.Prev = nil
		if visited[pModPeriod] {
			continue
		}
		visited[pModPeriod] = true
		// visitedPoints[p.P] = true
		i++

		// if (i % 100) == 0 {
		// 	printSearchStateAnimation(b, p)
		// }

		// If we've reached the end, we're done
		if p.P == end {
			return p
		}

		// Add the substates to the list of states to visit
		for _, delta := range deltas {
			newPosition := Position{
				P:    p.P.Add(delta.P),
				T:    p.T + delta.T,
				Prev: &p,
			}
			heap.Push(&toVisit, newPosition)
		}
	}
	return Position{T: -1}
}

func printSearchStateAnimation(b *Blizzard, currentPosition Position) {
	// Clear the console
	printSearchState(b, currentPosition)
	// Sleep a little
	// time.Sleep(10 * time.Millisecond)
}

// Print a map of the blizzard.
func printSearchState(b *Blizzard, currentPosition Position) {
	// Move the cursor to the top left of the console.
	outStr := "\033[H\033[2J"
	for y := 0; y < b.Height; y++ {
		for x := 0; x < b.Width; x++ {
			p := Point{X: x, Y: y}
			// Use console colors to represent whether the point is on the path or not.
			// If on the path, color the square with bright yellow.
			// Otherwise use gray
			colorString := "\033[0;90m"

			path := &currentPosition
			for path != nil {
				if path.P == p {
					colorString = "\033[1;33m"
					break
				}
				path = path.Prev
			}

			var s string
			if p == currentPosition.P {
				s = "E"
			} else {
				s = b.GetTileAtTime(x, y, currentPosition.T).String()
			}
			outStr += colorString + s
		}
		outStr += "\n"
	}
	// Reset color
	outStr += "\033[0m"
	// Also print the current time
	outStr += fmt.Sprintf("Time: %d", currentPosition.T)
	fmt.Print(outStr)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: 24 <input file>")
		os.Exit(1)
	}
	b := ParseBlizzardFromFile(os.Args[1])

	start := Point{X: 1, Y: 0}
	end := Point{X: b.Width - 2, Y: b.Height - 1}

	// animate(b)
	p1 := search(b, Position{P: start, T: 0}, end)
	// 308 is too high?

	printSearchStateAnimation(b, p1)
	fmt.Println("Part 1:", p1.T)

	p1.Prev = nil
	p2 := search(b, p1, start)

	p2.Prev = nil
	p3 := search(b, p2, end)

	printSearchStateAnimation(b, p1)

	fmt.Println("Part 2:", p3.T)

}
