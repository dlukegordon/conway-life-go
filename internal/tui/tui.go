package tui

import (
	b "conway-life-go/internal/board"
	"conway-life-go/internal/patterns"
	u "conway-life-go/internal/util"

	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/term"
)

type model struct {
	board *b.Board
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyMsg(msg)
	case tea.WindowSizeMsg:
		return m.handleWindowSizeMsg(msg)
	case tea.MouseMsg:
		return m.handleMouseMsg(msg)
	default:
		return m, nil
	}
}

func (m model) View() string {
	return renderBoard(m.board)
}

func Run() error {
	board := setupBoard()
	m := model{board}

	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())
	_, err := p.Run()
	u.FailIf(err, "running TUI")

	return nil

}

func setupBoard() *b.Board {
	width, height, err := term.GetSize(0)
	u.FailIf(err, "getting terminal dimensions")
	if width < 0 || height < 0 {
		u.Fail("getting terminal dimensions: Got a negative width nand height")
	}
	board := b.NewBoard(uint(height), uint(width))

	patternPos := b.NewPosition(30, 60)
	pattern := patterns.AcornPattern()
	err = board.Add(pattern, patternPos)
	u.FailIf(err, "setting initial board state")

	return board
}

func (m model) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit
	default:
		return m, nil
	}
}

func (m model) handleWindowSizeMsg(_msg tea.WindowSizeMsg) (tea.Model, tea.Cmd) {
	// TODO:
	return m, nil
}

func (m model) handleMouseMsg(_msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	// TODO:
	return m, nil
}

func renderBoard(board *b.Board) string {
	boardStr := ""

	for _, row := range board.Cells {
		rowStr := ""

		for _, cell := range row {
			cellStr := "  "
			if cell {
				cellStr = "██"
			}
			rowStr += cellStr
		}

		rowStr += "\n"
		boardStr += rowStr
	}

	return boardStr
}
