package main

import (
	"fmt"
	"math"
	"os"
)

const boardWidth = 7

// These look flipped compared to the actual shapes, which is because of the bit
// representation we're using for the board/pieces.
var pieceShapes = [][]byte{
	// Horizontal line
	{0b1111},
	// Plus shape
	{0b010, 0b111, 0b010},
	// L shape
	{0b111, 0b100, 0b100},
	// Vertical line
	{0b1, 0b1, 0b1, 0b1},
	// Block
	{0b11, 0b11},
}

type Piece struct {
	// A byte array representation of this piece.
	Shape []byte
	// Position of the bottom left corner of this piece.
	X      int
	Y      int
	Width  int
	Height int
}

// Gets the value of the piece at the position relative to the bottom left corner.
// False if the position is out of bounds.
func (p Piece) Get(x, y int) bool {
	if x < 0 || x >= p.Width || y < 0 || y >= p.Height {
		return false
	}
	return p.Shape[y]&(1<<x) != 0
}

func NewPiece(shape []byte) Piece {
	// Find the maximum width of any row.
	width := 0
	for _, row := range shape {
		rowWidth := 0
		for i := 0; i < boardWidth; i++ {
			if row&(1<<i) != 0 {
				rowWidth = i + 1
			}
		}
		if rowWidth > width {
			width = rowWidth
		}
	}
	return Piece{
		Shape:  shape,
		Width:  width,
		Height: len(shape),
	}
}

// Represents a 2D board as bits in a byte array.
// The first byte is the bottom of the board, and each byte is higher up.
// The lowest bit of each byte is the leftmost column of the board.
type Board []byte

// Gets the value of the board at the position.
// If x is out of bounds, treat it as a wall, return true.
// If y is too low, treat it as the floor and return true.
// If y is too high, treat it as an open ceiling and return false.
func (b Board) Get(x, y int) bool {
	if x < 0 || x >= boardWidth || y < 0 {
		return true
	}
	if y >= len(b) {
		return false
	}
	return b[y]&(1<<x) != 0
}

func (b Board) Set(x, y int, val bool) {
	// Panic if out of bounds
	if x < 0 || x >= boardWidth || y < 0 || y >= len(b) {
		panic(fmt.Sprintf("Out of bounds: %d, %d", x, y))
	}
	if val {
		// fmt.Println("Setting ", x, y, " to true")
		// fmt.Println("Before: ", b[y])
		b[y] |= (1 << x)
		// fmt.Println("After: ", b[y])
		// fmt.Println("1 << x: ", 1<<x)
	} else {
		b[y] &= ^(1 << x)
	}
}

func (b Board) Height() int {
	return len(b)
}

func (b Board) PieceIsColiding(p Piece) bool {
	for y := 0; y < p.Height; y++ {
		for x := 0; x < p.Width; x++ {
			if p.Get(x, y) && b.Get(x+p.X, y+p.Y) {
				return true
			}
		}
	}
	return false
}

func (b Board) PlacePiece(p Piece) Board {
	// Place the piece on the board, growing the size of the slice if needed.
	// fmt.Println("Placing piece at ", p.X, p.Y, " board size is ", len(b))
	for y := 0; y < p.Height; y++ {
		// Grow the board if needed.
		for y+p.Y >= len(b) {
			b = append(b, 0)
		}
		for x := 0; x < p.Width; x++ {
			if p.Get(x, y) {
				// fmt.Println("Setting ", x+p.X, y+p.Y, " to true")
				b.Set(x+p.X, y+p.Y, true)
			}
		}
	}
	// fmt.Println("Placed piece at ", p.X, p.Y, " board size is ", len(b))
	return b
}

// Prints the board as a string, with an ascii boarder for the walls and the bottom, but not the ceiling.
func (b Board) String() string {
	str := ""
	// Print the board from top to bottom
	for y := len(b) - 1; y >= 0; y-- {
		str += "|"
		for x := 0; x < boardWidth; x++ {
			if b.Get(x, y) {
				str += "#"
			} else {
				str += "."
			}
		}
		str += "|\n"
	}
	str += "+-------+"
	return str
}

func (b Board) ToStringWithPiece(p Piece) string {
	str := ""
	// Print the board and piece, from top to bottom.
	// We start at 5 above the top of the board, to fit the piece.
	for y := len(b) + 6; y >= 0; y-- {
		str += "|"
		for x := 0; x < boardWidth; x++ {
			if b.Get(x, y) {
				str += "#"
			} else if p.Get(x-p.X, y-p.Y) {
				str += "O"
			} else {
				str += "."
			}
		}
		str += "|\n"
	}
	str += "+-------+"
	return str
}

// movements is a string with characters '>' and '<' representing the moves to make.
func simulatePieces(movements string, numMoves int, numPieces int) (int, int) {
	// Create the board. Start with some arbitrary height.
	board := make(Board, 0, 10)
	// Index of the current piece
	pieceIndex := 0
	// Create the first piece
	piece := NewPiece(pieceShapes[pieceIndex%len(pieceShapes)])
	// Place the piece at the top.
	piece.X = 2
	piece.Y = board.Height() + 3

	// // Clear the console and move the cursor to the top left
	// fmt.Print("\033[H\033[2J")
	// // Print the state of the board and piece before the simulation.
	// fmt.Println(board.ToStringWithPiece(piece))

	lastHeight := board.Height()
	lastPieceIndex := pieceIndex

	for i := 0; i < numMoves && pieceIndex < numPieces; i++ {
		if i%len(movements) == 0 {
			heightDiff := board.Height() - lastHeight
			pieceDiff := pieceIndex - lastPieceIndex
			fmt.Println("Starting loop ", i/len(movements), " height diff ", heightDiff, " piece diff ", pieceDiff)
			lastHeight = board.Height()
			lastPieceIndex = pieceIndex
		}

		// Check if the pattern has looped
		if (i != 0) && (i%len(movements) == 0) && (pieceIndex%len(pieceShapes) == 0) {
			fmt.Println("Looped after ", i, " moves, ", pieceIndex, " blocks")
			break
		}
		// Try move the piece. If it collides, move it back.
		if movements[i%len(movements)] == '>' {
			piece.X++
			if board.PieceIsColiding(piece) {
				piece.X--
			}
		} else {
			piece.X--
			if board.PieceIsColiding(piece) {
				piece.X++
			}
		}
		// Move the piece one down.
		piece.Y--
		// If the piece collides, place it on the board and create a new piece.
		if board.PieceIsColiding(piece) {
			piece.Y++
			board = board.PlacePiece(piece)
			pieceIndex++
			piece = NewPiece(pieceShapes[pieceIndex%len(pieceShapes)])
			piece.X = 2
			piece.Y = board.Height() + 3
			// fmt.Println(board.ToStringWithPiece(piece))
			// fmt.Println()
		}

		// // Clear the console and move the cursor to the top left
		// fmt.Print("\033[H\033[2J")
		// // Print the state of the board and piece
		// fmt.Println(board.ToStringWithPiece(piece))

		// Sleep so we can see the animation
		// time.Sleep(100 * time.Millisecond)
	}
	return board.Height(), pieceIndex
}

func solve(filename string) {
	// Read the input file into a string
	input, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	// Remove the newline at the end of the file.
	input = input[:len(input)-1]
	inputStr := string(input)

	fmt.Println("num moves:", len(inputStr))

	// // Simulate the pieces for 2022 moves.
	// totalHeight := simulatePieces(inputStr, 2022)
	// // Print the result.
	// fmt.Println(totalHeight)

	// Part 2: Find a loop.
	// simulatePieces(inputStr, 20*len(inputStr))

	height1, pieces1 := simulatePieces(inputStr, len(inputStr), math.MaxInt)
	height2, pieces2 := simulatePieces(inputStr, 2*len(inputStr), math.MaxInt)
	heightDiff := height2 - height1
	piecesDiff := pieces2 - pieces1
	fmt.Println("height diff", heightDiff)
	fmt.Println("pieces diff", piecesDiff)

	// Find the height after 1000000000000 pieces.
	piecesLeft := 1000000000000
	totalHeight := 0

	piecesLeft -= pieces1
	totalHeight += height1

	extraLoopsNeeded := piecesLeft / piecesDiff
	piecesLeft -= extraLoopsNeeded * piecesDiff
	totalHeight += extraLoopsNeeded * heightDiff

	// For the remaining pieces, simulate them.
	height, _ := simulatePieces(inputStr, math.MaxInt, pieces1+piecesLeft)
	piecesLeft = 0
	totalHeight += height - height1

	fmt.Println("total height", totalHeight)
}

func main() {
	solve("17/input.txt")
}
