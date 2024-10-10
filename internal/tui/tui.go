package tui

import (
	"fmt"
	"strings"
	"time"

	b "conway-life-go/internal/board"
	"conway-life-go/internal/patterns"
	u "conway-life-go/internal/util"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
	"golang.org/x/term"
)

var (
	ColsPerCell   = 2
	LiveCellChar  = "█"
	DeadCellChar  = " "
	LiveCellStr   = strings.Repeat(LiveCellChar, ColsPerCell)
	DeadCellStr   = strings.Repeat(DeadCellChar, ColsPerCell)
	InitFps       = uint(10)
	MinFps        = uint(1)
	MaxFps        = uint(120)
	FpsChangeFrac = 0.2
	UiLines       = uint(1) // Number of lines to leave available in the viewport for ui elements other than the board
	bold          = color.New(color.Bold).SprintFunc()
)

type (
	model struct {
		termHeight uint
		termWidth  uint
		board      *b.Board
		fps        uint
		stepNum    uint
	}

	tickMsg struct{}
)

func (m model) Init() tea.Cmd {
	return tick(m.fps)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		m.board = m.board.NextBoard()
		m.stepNum++
		return m, tick(m.fps)
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
	return m.renderBoard() + m.renderToolbar()
}

func Run() error {
	termHeight, termWidth := getTermDims()

	m := model{
		termHeight: termHeight,
		termWidth:  termWidth,
		board:      setupBoard(),
		fps:        InitFps,
		stepNum:    0,
	}

	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())
	_, err := p.Run()
	u.FailIf(err, "running TUI")

	return nil

}

func getTermDims() (uint, uint) {
	termWidthInt, termHeightInt, err := term.GetSize(0)
	u.FailIf(err, "getting terminal dimensions")
	if termWidthInt < 0 || termHeightInt < 0 {
		u.Fail("getting terminal dimensions: Got a negative width nand height")
	}
	return uint(termHeightInt), uint(termWidthInt)
}

func setupBoard() *b.Board {
	termWidth, termHeight, err := term.GetSize(0)
	u.FailIf(err, "getting terminal dimensions")
	if termWidth < 0 || termHeight < 0 {
		u.Fail("getting terminal dimensions: Got a negative width nand height")
	}
	height := uint(termHeight) - UiLines
	width := uint(termWidth / ColsPerCell)
	board := b.NewBoard(height, width)

	patternPos := b.NewPosition(30, 50)
	pattern := patterns.AcornPattern()
	err = board.Add(pattern, patternPos)
	u.FailIf(err, "setting initial board state")

	return board
}

func tick(fps uint) tea.Cmd {
	return tea.Tick(time.Second/time.Duration(fps), func(_ time.Time) tea.Msg {
		return tickMsg{}
	})
}

func (m model) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit
	case "+", "=":
		m.changeFps(true)
		return m, nil
	case "-":
		m.changeFps(false)
		return m, nil
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

func (m model) renderBoard() string {
	boardStr := ""

	for _, row := range m.board.Cells {
		rowStr := ""

		for _, cell := range row {
			cellStr := DeadCellStr
			if cell {
				cellStr = LiveCellStr
			}
			rowStr += cellStr
		}

		rowStr += "\n"
		boardStr += rowStr
	}

	return boardStr
}

func (m model) renderToolbar() string {
	left := fmt.Sprintf("fps: %d, stepNum: %d", m.fps, m.stepNum)

	keyInfo := []string{
		bold("+") + "/" + bold("-") + " change speed",
		bold("q") + " quit",
	}
	right := strings.Join(keyInfo, " • ")

	// Calculate the padding to make the left side left-aligned, and the right side right-aligned.
	// This is a bit tricky, as Sprintf will treat ANSI escape code as actual text to be displayed, so will not pad
	// enough. We deal with this by adding the number of escape chars to the padding.
	// We also add an extra padding space to the right side if the terminal width is odd, to make sure it's right-aligned.
	paddingLeft := m.termWidth / 2
	paddingRight := paddingLeft
	if m.termWidth%2 == 1 {
		paddingRight++
	}
	paddingLeft += u.NumEscapeChars(left)
	paddingRight += u.NumEscapeChars(right)

	fmtStr := fmt.Sprintf("%%-%ds%%%ds", paddingLeft, paddingRight)
	return fmt.Sprintf(fmtStr, left, right)
}

func (m *model) changeFps(increase bool) {
	change := int(float64(m.fps) * FpsChangeFrac)
	// Make sure we don't get stuck not being able to change the fps
	if change == 0 {
		change = 1
	}

	if !increase {
		change = -change
	}

	newFps := int(m.fps) + change
	// Make sure we stay within the bounds of allowed fps
	if newFps < int(MinFps) {
		newFps = int(MinFps)
	}
	if newFps > int(MaxFps) {
		newFps = int(MaxFps)
	}

	m.fps = uint(newFps)
}
