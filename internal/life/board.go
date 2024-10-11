package life

import (
	"errors"
)

type (
	Cells [][]bool

	Board struct {
		cells Cells
		yLen  uint
		xLen  uint
	}

	Position struct {
		y uint
		x uint
	}

	offset struct {
		y int
		x int
	}
)

var Directions = [8]offset{
	{-1, -1},
	{-1, 0},
	{-1, 1},
	{0, -1},
	{0, 1},
	{1, -1},
	{1, 0},
	{1, 1},
}

const OutOfBoundsState = false

func nextCellState(alive bool, numLiveNeighbors uint) bool {
	if alive {
		return numLiveNeighbors == 2 || numLiveNeighbors == 3
	}
	return numLiveNeighbors == 3
}

// Return nil if adding the offset to the position would be out of bounds, otherwise a pointer to the new axis position
func addAxisOffset(axisPos uint, axisOffset int, axisLen uint) *uint {
	newAxisPos := int(axisPos) + axisOffset
	if newAxisPos < 0 || newAxisPos >= int(axisLen) {
		return nil
	}

	validNewAxisPos := uint(newAxisPos)
	return &validNewAxisPos
}

func NewPosition(y uint, x uint) Position {
	return Position{y, x}
}

func NewBoard(yLen uint, xLen uint) *Board {
	cells := make([][]bool, yLen)
	for i := range cells {
		cells[i] = make([]bool, xLen)
	}

	return &Board{
		cells,
		yLen,
		xLen,
	}
}

func NewBoardFromCells(cells Cells) (*Board, error) {
	yLen := len(cells)
	if yLen == 0 {
		return nil, errors.New("Invalid cells, no rows")
	}

	xLen := len(cells[0])
	if xLen == 0 {
		return nil, errors.New("Invalid cells, empty row")
	}

	// Make sure all rows are same length
	for i := 1; i < yLen; i++ {
		if len(cells[i]) != xLen {
			return nil, errors.New("Invalid cells, rows do not have a consistent length")
		}
	}

	return &Board{
		cells,
		uint(yLen),
		uint(xLen),
	}, nil
}

func (b *Board) YLen() uint {
	return b.yLen
}

func (b *Board) XLen() uint {
	return b.xLen
}

func (b *Board) CellAlive(pos Position) bool {
	if !b.inBounds(pos) {
		panic("Tried to access an out of bounds cell state")
	}
	return b.cells[pos.y][pos.x]
}

func (b *Board) inBounds(pos Position) bool {
	return pos.y < b.yLen && pos.x < b.xLen
}

// Return nil if the neighbor is out of bounds, otherwise a pointer to the position
func (b *Board) neighbor(pos Position, offset offset) *Position {
	newY := addAxisOffset(pos.y, offset.y, b.yLen)
	newX := addAxisOffset(pos.x, offset.x, b.xLen)
	if newY == nil || newX == nil {
		return nil
	}

	return &Position{*newY, *newX}
}

// Return an array of position pointers, nil when the neighbor is out of bounds
func (b *Board) neighbors(pos Position) []*Position {
	var neighbors []*Position

	for _, direction := range Directions {
		neighborPos := b.neighbor(pos, direction)
		neighbors = append(neighbors, neighborPos)
	}

	return neighbors
}

func (b *Board) numLiveNeighbors(pos Position) uint {
	var numLive uint
	neighbors := b.neighbors(pos)

	for _, neighbor := range neighbors {
		neighborAlive := OutOfBoundsState
		if neighbor != nil {
			neighborAlive = b.CellAlive(*neighbor)
		}

		if neighborAlive {
			numLive++
		}

	}

	return numLive
}

func (b *Board) nextBoard() *Board {
	nextBoard := NewBoard(b.yLen, b.xLen)

	for y := range b.yLen {
		for x := range b.xLen {
			pos := Position{y, x}
			alive := b.CellAlive(pos)
			numLiveNeighbors := b.numLiveNeighbors(pos)
			nextBoard.cells[y][x] = nextCellState(alive, numLiveNeighbors)
		}
	}

	return nextBoard
}

// Add a pattern to a board, returning a new board to retain immutability
func (b *Board) Add(pattern *Board, pos Position) (*Board, error) {
	furthestPatternPos := Position{pos.y + pattern.yLen, pos.x + pattern.xLen}
	if !b.inBounds(furthestPatternPos) {
		return nil, errors.New("Pattern will not fit in board at that position (if any)")
	}

	// TODO: don't bother copying the cells which will be overwritten by the pattern
	newBoard := NewBoard(b.yLen, b.xLen)
	for y := range b.yLen {
		for x := range b.xLen {
			newBoard.cells[y][x] = b.cells[y][x]
		}
	}

	for y := range pattern.yLen {
		for x := range pattern.xLen {
			boardY := pos.y + y
			boardX := pos.x + x
			newBoard.cells[boardY][boardX] = pattern.cells[y][x]
		}
	}

	return newBoard, nil
}
