package main

type Direction int

type Point struct {
	X, Y int
}

func (p Point) Add(o Point) Point {
	return Point{p.X + o.X, p.Y + o.Y}
}

// Also includes diagonal neighbors
func (p Point) Adjacent() []Point {
	return []Point{
		Point{p.X + 0, p.Y - 1},
		Point{p.X + 1, p.Y + 0},
		Point{p.X + 0, p.Y + 1},
		Point{p.X - 1, p.Y + 0},

		Point{p.X - 1, p.Y - 1},
		Point{p.X + 1, p.Y - 1},
		Point{p.X + 1, p.Y + 1},
		Point{p.X - 1, p.Y + 1},
	}
}

const (
	East  Direction = 0
	South Direction = 1
	West  Direction = 2
	North Direction = 3
)

func (d Direction) Delta() Point {
	switch d {
	case North:
		return Point{0, -1}
	case East:
		return Point{1, 0}
	case South:
		return Point{0, 1}
	case West:
		return Point{-1, 0}
	default:
		panic("Unknown direction")
	}
}

func (d Direction) DeltasToCheck() []Point {
	delta := d.Delta()
	return []Point{
		Point{delta.X, delta.Y},
		Point{delta.X + delta.Y, delta.Y - delta.X},
		Point{delta.X - delta.Y, delta.Y + delta.X},
	}
}

func (d Direction) String() string {
	switch d {
	case North:
		return "N"
	case East:
		return "E"
	case South:
		return "S"
	case West:
		return "W"
	default:
		panic("Unknown direction")
	}
}

func DirectionFromString(s string) Direction {
	switch s {
	case "N":
		return North
	case "E":
		return East
	case "S":
		return South
	case "W":
		return West
	default:
		panic("Unknown direction")
	}
}
