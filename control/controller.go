package control

import (
	"github.com/shanehowearth/bejeweled/state"
)

type GameController struct {
	Board    *state.Board
	Selected *state.Position
	OnUpdate func()
}

func NewGameController(board *state.Board, onUpdate func()) *GameController {
	return &GameController{
		Board:    board,
		OnUpdate: onUpdate,
	}
}

func (c *GameController) HandleGemClick(r, cPos int) {
	if c.Board.Moves <= 0 {
		return
	}

	clicked := state.Position{Row: r, Col: cPos}

	if c.Selected == nil {
		// First click: Select tile
		c.Selected = &clicked
		c.OnUpdate()
		return
	}

	first := *c.Selected
	c.Selected = nil

	// Double clicking same gem deselects
	if first == clicked {
		c.OnUpdate()
		return
	}

	// Must be adjacent to swap
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
	} else {
		// Valid move -> Process cascades
		c.Board.Moves--
		for len(matches[0]) > 0 {
			c.Board.ApplyCascade(matches[0])
			matches = c.Board.FindMatches()
		}
	}

	c.OnUpdate()
}

func (c *GameController) Reset() {
	c.Board.Score = 0
	c.Board.Moves = 30
	c.Board.FillBoardWithoutMatches()
	c.Selected = nil
	c.OnUpdate()
}
