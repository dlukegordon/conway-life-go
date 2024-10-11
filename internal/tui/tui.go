package tui

import (
	"fmt"
	"strings"
	"time"

	"conway-life-go/internal/life"
	"conway-life-go/internal/patterns"
	u "conway-life-go/internal/util"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
	"golang.org/x/term"
)

var (
	ColsPerCell           = 2
	LiveCellChar          = "█"
	DeadCellChar          = " "
	LiveCellStr           = strings.Repeat(LiveCellChar, ColsPerCell)
	DeadCellStr           = strings.Repeat(DeadCellChar, ColsPerCell)
	InitStepsPerSec       = uint(10)
	MaxStepsPerSec        = uint(120)
	StepsPerSecChangeFrac = 0.2
	UiLines               = uint(1) // Number of lines to leave available in the viewport for elements other than the board
	bold                  = color.New(color.Bold).SprintFunc()
)

type (
	model struct {
		termHeight  uint
		termWidth   uint
		game        *life.Game
		stepsPerSec uint
		stepBack    bool
	}

	tickMsg struct{}
)

func (m model) Init() tea.Cmd {
	return m.tick()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		return m.handleTick()
	case tea.KeyMsg:
		return m.handleKeyMsg(msg)
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
		termHeight:  termHeight,
		termWidth:   termWidth,
		game:        setupGame(),
		stepsPerSec: InitStepsPerSec,
		stepBack:    false,
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

func setupGame() *life.Game {
	termWidth, termHeight, err := term.GetSize(0)
	u.FailIf(err, "getting terminal dimensions")
	if termWidth < 0 || termHeight < 0 {
		u.Fail("getting terminal dimensions: Got a negative width nand height")
	}
	height := uint(termHeight) - UiLines
	width := uint(termWidth / ColsPerCell)
	board := life.NewBoard(height, width)

	patternPos := life.NewPosition(30, 50)
	pattern := patterns.AcornPattern()
	board, err = board.Add(pattern, patternPos)
	u.FailIf(err, "setting initial board state")

	return life.NewGameFromBoard(board)
}

func (m model) tick() tea.Cmd {
	// Even if we're paused, we still want to tick every second so we can resume later
	waitTime := time.Second
	if m.stepsPerSec > 0 {
		waitTime = time.Second / time.Duration(m.stepsPerSec)
	}

	return tea.Tick(waitTime, func(_ time.Time) tea.Msg {
		return tickMsg{}
	})
}

func (m model) handleTick() (tea.Model, tea.Cmd) {
	if m.stepsPerSec == 0 {
		// If we're paused, do nothing
		return m, m.tick()
	}

	if m.stepBack {
		err := m.game.StepBack()
		if err != nil {
			// We can't go back anymore
			m.stepsPerSec = 0
		}
	} else {
		m.game.Step()
	}

	return m, m.tick()
}

func (m model) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit
	case "+", "=":
		m.changeSpeed(true)
		return m, nil
	case "-":
		m.changeSpeed(false)
		return m, nil
	default:
		return m, nil
	}
}

func (m model) handleMouseMsg(_msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	// TODO:
	return m, nil
}

func (m model) renderBoard() string {
	boardStr := ""
	board := m.game.CurrentBoard()

	for y := range board.YLen() {
		rowStr := ""

		for x := range board.XLen() {
			pos := life.NewPosition(uint(y), uint(x))

			cellStr := DeadCellStr
			if board.CellAlive(pos) {
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
	timeDirectionStr := ""
	if m.stepBack && m.stepsPerSec > 0 {
		timeDirectionStr = "-"
	}
	left := fmt.Sprintf("stepsPerSec: %s%d, stepNum: %d", timeDirectionStr, m.stepsPerSec, m.game.CurrentStepNum())

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

func (m *model) changeSpeed(increase bool) {
	// If we're going back in time, then decreasing speed means increasing stepsPerSec, and vice versa
	if m.stepBack {
		increase = !increase
	}

	change := int(float64(m.stepsPerSec) * StepsPerSecChangeFrac)
	// Make sure we don't get stuck not being able to change the stepsPerSec
	if change == 0 {
		change = 1
	}

	if !increase {
		change = -change
	}

	newStepsPerSec := int(m.stepsPerSec) + change
	if newStepsPerSec < 0 {
		// We're changing time direction
		newStepsPerSec = 1
		m.stepBack = !m.stepBack
	}
	if newStepsPerSec > int(MaxStepsPerSec) {
		// Make sure we stay within the bounds of allowed stepsPerSec
		newStepsPerSec = int(MaxStepsPerSec)
	}

	m.stepsPerSec = uint(newStepsPerSec)
}
