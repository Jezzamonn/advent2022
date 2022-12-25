package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Point struct {
	X int
	Y int
	Z int
}

type Slice3D struct {
	Points []bool
	MaxX   int
	MaxY   int
	MaxZ   int
}

var faceDirections = []Point{
	{1, 0, 0},
	{-1, 0, 0},
	{0, 1, 0},
	{0, -1, 0},
	{0, 0, 1},
	{0, 0, -1},
}

func NewSlice3D(maxX, maxY, maxZ int) *Slice3D {
	return &Slice3D{
		Points: make([]bool, maxX*maxY*maxZ),
		MaxX:   maxX,
		MaxY:   maxY,
		MaxZ:   maxZ,
	}
}

func (s *Slice3D) Get(x, y, z int) bool {
	// Return false if the point is out of bounds
	if x < 0 || x >= s.MaxX || y < 0 || y >= s.MaxY || z < 0 || z >= s.MaxZ {
		return false
	}
	return s.Points[x+y*s.MaxX+z*s.MaxX*s.MaxY]
}

func (s *Slice3D) Set(x, y, z int, val bool) {
	// Panic if out of bounds
	if x < 0 || x >= s.MaxX || y < 0 || y >= s.MaxY || z < 0 || z >= s.MaxZ {
		panic("Out of bounds")
	}
	s.Points[x+y*s.MaxX+z*s.MaxX*s.MaxY] = val
}

func CreatePointsAndSlice3D(filename string) ([]Point, *Slice3D) {
	// Read the input file into a string
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// Read each line of x,y,z coordinates into a slice of points
	points := make([]Point, 0)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		// Split the line into a slice of strings
		lineParts := strings.Split(line, ",")

		// Convert the strings to ints
		x, err := strconv.Atoi(lineParts[0])
		if err != nil {
			panic(err)
		}
		y, err := strconv.Atoi(lineParts[1])
		if err != nil {
			panic(err)
		}
		z, err := strconv.Atoi(lineParts[2])
		if err != nil {
			panic(err)
		}

		// Add the point to the slice
		points = append(points, Point{x, y, z})
	}
	// Find the maximum x, y, and z values
	maxX := 0
	maxY := 0
	maxZ := 0
	for _, point := range points {
		if point.X > maxX {
			maxX = point.X
		}
		if point.Y > maxY {
			maxY = point.Y
		}
		if point.Z > maxZ {
			maxZ = point.Z
		}
	}

	// Create a 3D slice of bools
	slice := NewSlice3D(maxX+1, maxY+1, maxZ+1)
	// Add each point
	for _, point := range points {
		slice.Set(point.X, point.Y, point.Z, true)
	}

	return points, slice
}

func solvePt1(filename string) {
	// Create the points and slice
	points, slice := CreatePointsAndSlice3D(filename)

	// Iterate over each point, count empty neighbours.
	surfaceArea := 0
	for _, point := range points {
		emptyNeighbours := 0
		for _, direction := range faceDirections {
			// Check if the neighbour is empty
			if !slice.Get(point.X+direction.X, point.Y+direction.Y, point.Z+direction.Z) {
				emptyNeighbours++
			}
		}
		surfaceArea += emptyNeighbours
	}

	// Print the result
	fmt.Println(surfaceArea)
}

func solvePt2(filename string) {
	points, slice := CreatePointsAndSlice3D(filename)

	// Create a duplicate slice, and do a fill from the outside.
	// Make it a little larger than the original slice
	padding := 1
	emptySpace := NewSlice3D(
		slice.MaxX+2*padding,
		slice.MaxY+2*padding,
		slice.MaxZ+2*padding)
	explored := make(map[Point]struct{})
	toBeExplored := make([]Point, 0)
	// Start at the top left corner
	toBeExplored = append(toBeExplored, Point{0, 0, 0})
	for len(toBeExplored) > 0 {
		next := toBeExplored[0]
		toBeExplored = toBeExplored[1:]

		// Check if the point is out of bounds
		if next.X < 0 || next.X >= slice.MaxX+2*padding ||
			next.Y < 0 || next.Y >= slice.MaxY+2*padding ||
			next.Z < 0 || next.Z >= slice.MaxZ+2*padding {
			continue
		}

		// Check if the point has already been explored
		if _, ok := explored[next]; ok {
			continue
		}

		// Mark the point as explored
		explored[next] = struct{}{}

		// Ensure the point is empty
		if slice.Get(next.X-padding, next.Y-padding, next.Z-padding) {
			continue
		}

		// Set the point in the second slice
		emptySpace.Set(next.X, next.Y, next.Z, true)

		// Add the neighbours to the toBeExplored slice
		for _, direction := range faceDirections {
			toBeExplored = append(toBeExplored, Point{
				next.X + direction.X,
				next.Y + direction.Y,
				next.Z + direction.Z})
		}
	}

	// Count surface area
	surfaceArea := 0
	for _, point := range points {
		emptyNeighbours := 0
		for _, direction := range faceDirections {
			// Check if the neighbour is empty, using the emptySpace slice.
			if emptySpace.Get(
				point.X+direction.X+padding,
				point.Y+direction.Y+padding,
				point.Z+direction.Z+padding) {

				emptyNeighbours++
			}
		}
		surfaceArea += emptyNeighbours
	}

	// Print the result
	fmt.Println("Pt2", surfaceArea)
}

func main() {
	solvePt1("18/input.txt")
	solvePt2("18/input.txt")
}
