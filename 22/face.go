package main

// For input.txt
const FaceSize = 50
const FacePerRow = 3

// // For demo.txt
// const FaceSize = 4
// const FacePerRow = 4

type FaceIndex int

type FaceEdge struct {
	Index FaceIndex
	Edge  Direction
}

func (f FaceIndex) XIndex() int {
	return int(f) % FacePerRow
}

func (f FaceIndex) YIndex() int {
	return int(f) / FacePerRow
}

func (f FaceIndex) X() int {
	return f.XIndex() * FaceSize
}

func (f FaceIndex) Y() int {
	return f.YIndex() * FaceSize
}

func GetFaceIndex(x, y int) FaceIndex {
	return FaceIndex((y/FaceSize)*FacePerRow + (x / FaceSize))
}

// For input.txt
var FaceConnections = map[FaceEdge]FaceEdge{
	{1, West}: {6, West},
	{6, West}: {1, West},

	{1, North}: {9, West},
	{9, West}:  {1, North},

	{2, North}: {9, South},
	{9, South}: {2, North},

	{2, East}: {7, East},
	{7, East}: {2, East},

	{2, South}: {4, East},
	{4, East}:  {2, South},

	{4, West}:  {6, North},
	{6, North}: {4, West},

	{7, South}: {9, East},
	{9, East}:  {7, South},
}

// // For demo.txt
// var FaceConnections = map[FaceEdge]FaceEdge{
// 	{2, West}:  {5, North},
// 	{5, North}: {2, West},

// 	{2, North}: {4, North},
// 	{4, North}: {2, North},

// 	{2, East}:  {11, East},
// 	{11, East}: {2, East},

// 	{4, West}:   {11, South},
// 	{11, South}: {4, West},

// 	{4, South}:  {10, South},
// 	{10, South}: {4, South},

// 	{5, South}: {10, East},
// 	{10, East}: {5, South},

// 	{6, East}:   {11, North},
// 	{11, North}: {6, East},
// }
