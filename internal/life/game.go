package life

import (
	"errors"
)

type Game struct {
	boards  []*Board
	stepNum uint
}

func NewGame(yLen uint, xLen uint) *Game {
	return &Game{
		boards:  []*Board{NewBoard(yLen, xLen)},
		stepNum: 0,
	}
}

func NewGameFromBoard(board *Board) *Game {
	return &Game{
		boards:  []*Board{board},
		stepNum: 0,
	}

}

func (g *Game) CurrentStepNum() uint {
	return g.stepNum
}

func (g *Game) NumTotalSteps() uint {
	return uint(len(g.boards))
}

func (g *Game) CurrentBoard() *Board {
	return g.boards[g.stepNum]
}

func (g *Game) InPast() bool {
	return int(g.stepNum) < len(g.boards)-1
}

func (g *Game) Step() {
	if !g.InPast() {
		nextBoard := g.CurrentBoard().nextBoard()
		g.boards = append(g.boards, nextBoard)
	}

	g.stepNum++
}

func (g *Game) StepBack() error {
	if g.stepNum == 0 {
		return errors.New("Cannot step back, at step 0")
	}

	g.stepNum--
	return nil
}
