package patterns

import (
	b "conway-life-go/internal/board"
)

const (
	AliveChar = 'X'
	DeadChar  = ' '
)

func stringsToBoard(strs []string) *b.Board {
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

	board, err := b.NewBoardFromCells(cells)
	if err != nil {
		panic(err)
	}
	return board
}

func AcornPattern() *b.Board {
	strs := []string{
		"         ",
		"  X      ",
		"    X    ",
		" XX  XXX ",
		"         ",
	}

	return stringsToBoard(strs)
}

func BlinkerPattern() *b.Board {
	strs := []string{
		"     ",
		" XXX ",
		"     ",
	}

	return stringsToBoard(strs)
}
