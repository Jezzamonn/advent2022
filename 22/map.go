package main

import (
	"fmt"
	"strings"
)

type Map [][]Tile

type Tile byte

const (
	Outside Tile = iota
	Empty
	Wall
)

func TileFromRune(r rune) Tile {
	switch r {
	case ' ':
		return Outside
	case '.':
		return Empty
	case '#':
		return Wall
	default:
		panic(fmt.Sprintf("Unknown character in map: %c", r))
	}
}

func (t Tile) String() string {
	switch t {
	case Outside:
		return " "
	case Empty:
		return "."
	case Wall:
		return "#"
	default:
		panic(fmt.Sprintf("Unknown tile: %d", t))
	}
}

func ParseMap(s string) Map {
	// Split the string into lines.
	lines := strings.Split(s, "\n")

	width := 0
	for _, line := range lines {
		if len(line) > width {
			width = len(line)
		}
	}
	// Create a 2D byte array.
	m := make(Map, len(lines))
	for i, line := range lines {
		m[i] = make([]Tile, width)
		for j, c := range line {
			m[i][j] = TileFromRune(c)
		}
	}
	return m
}

func (m Map) String() string {
	var sb strings.Builder
	for _, row := range m {
		for _, tile := range row {
			sb.WriteString(tile.String())
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

func (m Map) ToStringWithPosition(x, y int, d Direction) string {
	var sb strings.Builder
	for i, row := range m {
		for j, tile := range row {
			if i == y && j == x {
				// Prints the arrow in a bold, yellow color in the terminal.
				sb.WriteString(fmt.Sprintf("\x1b[33;1m%s\x1b[0m", d.String()))
			} else {
				sb.WriteString(tile.String())
			}
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

func (m Map) ToStringWithVisited(visited [][]byte) string {
	var sb strings.Builder
	for i, row := range m {
		for j, tile := range row {
			if visited[i][j] > 0 {
				d := Direction(visited[i][j] - 1)
				// Prints the visited position in a bold, yellow color in the terminal.
				sb.WriteString(fmt.Sprintf("\x1b[33;1m%s\x1b[0m", d.ArrowString()))
			} else {
				sb.WriteString(tile.String())
			}
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

func (m Map) Width() int {
	return len(m[0])
}

func (m Map) Height() int {
	return len(m)
}
