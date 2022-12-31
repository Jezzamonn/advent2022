package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const animate = false

func followInstructions(m Map, s string) (x, y int, d Direction) {
	// Find the starting position (the first empty space of the first row).
	for i, b := range m[0] {
		if b == 1 {
			x = i
			break
		}
	}
	d = East

	w, h := m.Width(), m.Height()

	// 2D array of visited positions, along with direction. I guess we'll store this as 1 more than the direction byte.
	visited := make([][]byte, h)
	for i := range visited {
		visited[i] = make([]byte, w)
	}

	// Parse the instructions, using regex!
	// Example instructions: '15L3R31L27'
	instructionsRegex := regexp.MustCompile(`[LR]|(?:[0-9]+)`)
	instructions := instructionsRegex.FindAllStringSubmatch(s, -1)
	for _, instruction := range instructions {
		// If it's a number, go forward that many steps.
		switch instruction[0] {
		case "L":
			d = d.TurnLeft()
			visited[y][x] = 1 + byte(d)

			PrintAnimationFrame(m, visited)

		case "R":
			d = d.TurnRight()
			visited[y][x] = 1 + byte(d)

			PrintAnimationFrame(m, visited)

		default:
			dx, dy := d.Deltas()
			steps, err := strconv.Atoi(instruction[0])
			if err != nil {
				panic(err)
			}
			for i := 0; i < steps; i++ {
				// Starting positions
				sx := x
				sy := y

				// Move forward, repeating if needed if we are off the map.
				for {
					y = (y + dy + h) % h
					x = (x + dx + w) % w

					if m[y][x] != Outside {
						break
					}
				}

				// Can't move onto a wall, return to where we started.
				// Also no point trying to move forward more.
				if m[y][x] == 2 {
					x = sx
					y = sy
					break
				}

				visited[y][x] = 1 + byte(d)

				PrintAnimationFrame(m, visited)
			}
		}
	}
	return x, y, d
}

func followInstructionsOnCube(m Map, s string) (x, y int, d Direction) {
	// Find the starting position (the first empty space of the first row).
	for i, b := range m[0] {
		if b == 1 {
			x = i
			break
		}
	}
	d = East

	w, h := m.Width(), m.Height()

	// 2D array of visited positions, along with direction. I guess we'll store this as 1 more than the direction byte.
	visited := make([][]byte, h)
	for i := range visited {
		visited[i] = make([]byte, w)
	}

	// Parse the instructions, using regex!
	// Example instructions: '15L3R31L27'
	instructionsRegex := regexp.MustCompile(`[LR]|(?:[0-9]+)`)
	instructions := instructionsRegex.FindAllStringSubmatch(s, -1)
	for _, instruction := range instructions {
		// If it's a number, go forward that many steps.
		switch instruction[0] {
		case "L":
			d = d.TurnLeft()
			visited[y][x] = 1 + byte(d)

			PrintAnimationFrame(m, visited)

		case "R":
			d = d.TurnRight()
			visited[y][x] = 1 + byte(d)

			PrintAnimationFrame(m, visited)

		default:
			steps, err := strconv.Atoi(instruction[0])
			if err != nil {
				panic(err)
			}
			for i := 0; i < steps; i++ {
				dx, dy := d.Deltas()
				// Starting positions
				sx := x
				sy := y
				sd := d
				// Starting face
				f := GetFaceIndex(x, y)

				// Move forward, moving onto the other face of the cube if appropriate.
				y += dy
				x += dx

				if x < 0 || x >= w || y < 0 || y >= h || m[y][x] == Outside {
					// Move onto a different face.
					fe := FaceEdge{
						f, d,
					}
					// New face
					nf, ok := FaceConnections[fe]
					if !ok {
						panic("Don't know where this face connects to")
					}
					// New position on the face
					fx := x - f.X() - FaceSize*dx
					fy := y - f.Y() - FaceSize*dy

					// How much we need to rotate
					r := d.Subtract(nf.Edge.TurnAround())
					// New direction is just the opposite direction of the edge of the new face
					d = nf.Edge.TurnAround()
					// New position on the new face
					nx, ny := RotatePosition(fx, fy, r)
					// New position on the map
					x = nx + nf.Index.X()
					y = ny + nf.Index.Y()
				}

				// Can't move onto a wall, return to where we started.
				// Also no point trying to move forward more.
				if m[y][x] == 2 {
					x = sx
					y = sy
					d = sd
					break
				}

				visited[y][x] = 1 + byte(d)

				PrintAnimationFrame(m, visited)
			}
		}
	}
	return x, y, d
}

func RotatePosition(x, y, r int) (int, int) {
	switch r {
	case 0:
		return x, y
	case 1:
		return y, (FaceSize - 1) - x
	case 2:
		return (FaceSize - 1) - x, (FaceSize - 1) - y
	case 3:
		return (FaceSize - 1) - y, x
	default:
		panic("Invalid rotation")
	}
}

func PrintAnimationFrame(m Map, visited [][]byte) {
	if !animate {
		return
	}
	// Clear the console.
	fmt.Print("\033[H\033[2J")
	fmt.Println(m.ToStringWithVisited(visited))
	// Sleep so we can see the animation.
	time.Sleep(100 * time.Millisecond)
}

func solve(filename string) {
	fmt.Println("Solving", filename)
	// Read the full file into a string.
	s, err := ioutil.ReadFile(filename)
	if err != nil {
		s, err = ioutil.ReadFile("../" + filename)
		if err != nil {
			panic(err)
		}
	}
	// Split the string into the map and the instructions.
	split := strings.Split(string(s), "\n\n")

	m := ParseMap(split[0])

	// x, y, d := followInstructions(m, split[1])
	x, y, d := followInstructionsOnCube(m, split[1])
	// Convert them to 1-based
	x++
	y++
	fmt.Printf("Part 1: row=%d, col=%d, dir=%s (%d)\n", y, x, d, d)
	fmt.Println(1000*y + 4*x + int(d))
}

func main() {
	// solve("22/demo.txt")
	solve("22/input.txt")
}
