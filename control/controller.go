package control

import (
	"time"

	"fyne.io/fyne/v2"
	"github.com/shanehowearth/bejeweled/state"
)

type GameController struct {
	Board       *state.Board
	Selected    *state.Position
	Highlighted []state.Position
	IsAnimating bool
	OnUpdate    func()
}

func NewGameController(board *state.Board, onUpdate func()) *GameController {
	return &GameController{
		Board:    board,
		OnUpdate: onUpdate,
	}
}

func (c *GameController) HandleGemClick(r, cPos int) {
	// Block clicks while cascade animation is in progress
	if c.Board.GameOver || c.IsAnimating {
		return
	}

	clicked := state.Position{Row: r, Col: cPos}

	if c.Selected == nil {
		c.Selected = &clicked
		c.OnUpdate()
		return
	}

	first := *c.Selected
	c.Selected = nil

	if first == clicked {
		c.OnUpdate()
		return
	}

	if !state.IsAdjacent(first, clicked) {
		c.Selected = &clicked
		c.OnUpdate()
		return
	}

	// Try Swap
	c.Board.Swap(first, clicked)
	matches := c.Board.FindMatches()

	if len(matches[0]) == 0 {
		// Invalid move -> Swap back
		c.Board.Swap(first, clicked)
		c.OnUpdate()
	} else {
		// Valid move -> Launch Paced Cascade sequence in background
		go c.processCascadesPaced()
	}
}

func (c *GameController) processCascadesPaced() {
	c.IsAnimating = true

	for {
		matches := c.Board.FindMatches()
		if len(matches[0]) == 0 {
			break
		}

		// Phase 1: Highlight matched gems briefly on the main Fyne thread
		c.Highlighted = matches[0]
		fyne.DoAndWait(func() {
			c.OnUpdate()
		})
		time.Sleep(350 * time.Millisecond)

		// Phase 2: Clear highlights & apply cascade drop, updating UI on main thread
		c.Highlighted = nil
		c.Board.ApplyCascade(matches[0])
		fyne.DoAndWait(func() {
			c.OnUpdate()
		})

		// Phase 3: Short delay to let player observe dropped gems
		time.Sleep(250 * time.Millisecond)
	}

	// Auto-reshuffle if no valid moves left
	for !c.Board.HasValidMoves() {
		c.Board.FillBoardWithoutMatches()
	}

	c.IsAnimating = false
	fyne.Do(func() {
		c.OnUpdate()
	})
}

func (c *GameController) Reset() {
	c.Board.Score = 0
	c.Board.GameOver = false
	c.Board.FillBoardWithoutMatches()
	c.Selected = nil
	c.Highlighted = nil
	c.IsAnimating = false
	c.OnUpdate()
}
