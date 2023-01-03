package main

import (
	"bufio"
	"fmt"
	"os"
)

type Tile int

const (
	Empty Tile = iota
	Wall
	NorthWind
	EastWind
	SouthWind
	WestWind
	MultipleWinds
)

var winds = []Tile{NorthWind, EastWind, SouthWind, WestWind}

func TileFromString(s string) Tile {
	switch s {
	case ".":
		return Empty
	case "#":
		return Wall
	case "^":
		return NorthWind
	case ">":
		return EastWind
	case "v":
		return SouthWind
	case "<":
		return WestWind
	}
	return Empty
}

func (t Tile) String() string {
	switch t {
	case Empty:
		return "."
	case Wall:
		return "#"
	case NorthWind:
		return "^"
	case EastWind:
		return ">"
	case SouthWind:
		return "v"
	case WestWind:
		return "<"
	case MultipleWinds:
		return "X"
	}
	return "?"
}

func (t Tile) Delta() (int, int) {
	switch t {
	case NorthWind:
		return 0, -1
	case EastWind:
		return 1, 0
	case SouthWind:
		return 0, 1
	case WestWind:
		return -1, 0
	}
	return 0, 0
}

func (t Tile) IsWind() bool {
	return t >= NorthWind && t <= WestWind
}

type Blizzard struct {
	board  []Tile
	Width  int
	Height int
}

func ParseBlizzardFromFile(filename string) *Blizzard {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var b Blizzard
	for scanner.Scan() {
		line := scanner.Text()
		if b.Width == 0 {
			b.Width = len(line)
		} else if b.Width != len(line) {
			panic("Blizzard line length mismatch")
		}
		b.Height++
		for _, c := range line {
			b.board = append(b.board, TileFromString(string(c)))
		}
	}
	return &b
}

func (b *Blizzard) Period() int {
	// The LCM of (width - 2) and (height - 2), to account for the walls.
	return LCM(b.Width-2, b.Height-2)
}

func (b *Blizzard) Get(x, y int) Tile {
	// Outside the bounds = wall
	if x < 0 || x >= b.Width || y < 0 || y >= b.Height {
		return Wall
	}
	return b.board[y*b.Width+x]
}

func (b *Blizzard) Set(x, y int, t Tile) {
	// Outside the bounds = panic
	if x < 0 || x >= b.Width || y < 0 || y >= b.Height {
		panic(fmt.Sprintf("Set(%d, %d, %d) out of bounds", x, y, t))
	}
	b.board[y*b.Width+x] = t
}

func (b *Blizzard) GetTileAtTime(x, y, t int) Tile {
	// Winds don't exist on the edges of the map, so we can just use the regular GetTile method
	if x <= 0 || x >= b.Width-1 || y <= 0 || y >= b.Height-1 {
		return b.Get(x, y)
	}

	// For each direction, check if a wind tile from that direction has moved to this square.
	lastWind := Empty
	for _, wind := range winds {
		dx, dy := wind.Delta()
		// Modoulo to account for the wind looping. We have to also exclude the walls around the map.
		relX := x - 1
		relY := y - 1
		newRelX := PosMod((relX - dx*t), (b.Width - 2))
		newRelY := PosMod((relY - dy*t), (b.Height - 2))
		if b.Get(newRelX+1, newRelY+1) == wind {
			if lastWind == Empty {
				lastWind = wind
			} else {
				return MultipleWinds
			}
		}
	}
	return lastWind
}

func (b *Blizzard) IsEmptyAtTime(x, y, t int) bool {
	return b.GetTileAtTime(x, y, t) == Empty
}

func (b *Blizzard) StringAtTime(t int) string {
	s := ""
	for y := 0; y < b.Height; y++ {
		for x := 0; x < b.Width; x++ {
			s += b.GetTileAtTime(x, y, t).String()
		}
		s += "\n"
	}
	return s
}

func (b *Blizzard) String() string {
	return b.StringAtTime(0)
}
