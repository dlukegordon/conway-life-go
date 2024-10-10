package board

import (
	"errors"
)

type (
	Cells [][]bool

	Board struct {
		Cells Cells
		YLen  uint
		XLen  uint
	}

	Position struct {
		y uint
		x uint
	}

	Offset struct {
		y int
		x int
	}
)

var Directions = [8]Offset{
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

func (b Board) inBounds(pos Position) bool {
	return pos.y < b.YLen && pos.x < b.XLen
}

func (b Board) CellAlive(pos Position) bool {
	if !b.inBounds(pos) {
		panic("Tried to access an out of bounds cell state")
	}
	return b.Cells[pos.y][pos.x]
}

// Return nil if the neighbor is out of bounds, otherwise a pointer to the position
func (b Board) neighbor(pos Position, offset Offset) *Position {
	newY := addAxisOffset(pos.y, offset.y, b.YLen)
	newX := addAxisOffset(pos.x, offset.x, b.XLen)
	if newY == nil || newX == nil {
		return nil
	}

	return &Position{*newY, *newX}
}

// Return an array of position pointers, nil when the neighbor is out of bounds
func (b Board) neighbors(pos Position) []*Position {
	var neighbors []*Position

	for _, direction := range Directions {
		neighborPos := b.neighbor(pos, direction)
		neighbors = append(neighbors, neighborPos)
	}

	return neighbors
}

func (b Board) numLiveNeighbors(pos Position) uint {
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

func (b Board) NextBoard() *Board {
	nextBoard := NewBoard(b.YLen, b.XLen)

	for y := range b.YLen {
		for x := range b.XLen {
			pos := Position{y, x}
			alive := b.CellAlive(pos)
			numLiveNeighbors := b.numLiveNeighbors(pos)
			nextBoard.Cells[y][x] = nextCellState(alive, numLiveNeighbors)
		}
	}

	return nextBoard
}

func (b Board) Add(pattern *Board, pos Position) error {
	furthestPatternPos := Position{pos.y + pattern.YLen, pos.x + pattern.XLen}
	if !b.inBounds(furthestPatternPos) {
		return errors.New("Pattern will not fit in board at that position (if any)")
	}

	for y := range pattern.YLen {
		for x := range pattern.XLen {
			boardY := pos.y + y
			boardX := pos.x + x
			b.Cells[boardY][boardX] = pattern.Cells[y][x]
		}
	}

	return nil
}
