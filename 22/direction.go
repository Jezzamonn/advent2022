package main

type Direction int

const (
	East  Direction = 0
	South Direction = 1
	West  Direction = 2
	North Direction = 3
)

func (d Direction) TurnLeft() Direction {
	return (d + 3) % 4
}

func (d Direction) TurnRight() Direction {
	return (d + 1) % 4
}

func (d Direction) TurnAround() Direction {
	return (d + 2) % 4
}

func (d Direction) Subtract(other Direction) int {
	return (int(d) - int(other) + 4) % 4
}

func (d Direction) Deltas() (int, int) {
	switch d {
	case North:
		return 0, -1
	case East:
		return 1, 0
	case South:
		return 0, 1
	case West:
		return -1, 0
	default:
		panic("Unknown direction")
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

// Returns ">", "v", "<", or "^".
func (d Direction) ArrowString() string {
	switch d {
	case North:
		return "^"
	case East:
		return ">"
	case South:
		return "v"
	case West:
		return "<"
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
