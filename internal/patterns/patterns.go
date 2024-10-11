package patterns

import (
	"conway-life-go/internal/life"
)

const (
	AliveChar = 'X'
	DeadChar  = ' '
)

func stringsToBoard(strs []string) *life.Board {
	var cells [][]bool

	for _, strRow := range strs {
		var row []bool

		for _, char := range strRow {
			alive := false
			if char == AliveChar {
				alive = true
			} else if char != DeadChar {
				panic("Invalid string representation of cells")
			}
			row = append(row, alive)
		}

		cells = append(cells, row)
	}

	board, err := life.NewBoardFromCells(cells)
	if err != nil {
		panic(err)
	}
	return board
}

func AcornPattern() *life.Board {
	strs := []string{
		"         ",
		"  X      ",
		"    X    ",
		" XX  XXX ",
		"         ",
	}

	return stringsToBoard(strs)
}

func BlinkerPattern() *life.Board {
	strs := []string{
		"     ",
		" XXX ",
		"     ",
	}

	return stringsToBoard(strs)
}
