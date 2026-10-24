package state

import (
	"math/rand"
	"time"
)

type GemType int

const (
	GemEmpty  GemType = iota
	GemRed            // Ruby
	GemBlue           // Sapphire
	GemGreen          // Emerald
	GemYellow         // Topaz
	GemPurple         // Amethyst
	GemOrange         // Amber
	GemWhite          // Diamond
)

type Position struct {
	Row int
	Col int
}

type Board struct {
	Rows  int
	Cols  int
	Grid  [][]GemType
	Score int
	Moves int
}

func NewBoard(rows, cols int) *Board {
	rand.Seed(time.Now().UnixNano())
	b := &Board{
		Rows:  rows,
		Cols:  cols,
		Grid:  make([][]GemType, rows),
		Score: 0,
		Moves: 30, // Default 30-move game
	}
	for r := 0; r < rows; r++ {
		b.Grid[r] = make([]GemType, cols)
	}
	b.FillBoardWithoutMatches()
	return b
}

func (b *Board) FillBoardWithoutMatches() {
	for r := 0; r < b.Rows; r++ {
		for c := 0; c < b.Cols; c++ {
			var g GemType
			for {
				g = GemType(rand.Intn(7) + 1)
				// Prevent horizontal match on fill
				if c >= 2 && b.Grid[r][c-1] == g && b.Grid[r][c-2] == g {
					continue
				}
				// Prevent vertical match on fill
				if r >= 2 && b.Grid[r-1][c] == g && b.Grid[r-2][c] == g {
					continue
				}
				break
			}
			b.Grid[r][c] = g
		}
	}
}

// Swap exchanges two gems on the board
func (b *Board) Swap(p1, p2 Position) {
	b.Grid[p1.Row][p1.Col], b.Grid[p2.Row][p2.Col] = b.Grid[p2.Row][p2.Col], b.Grid[p1.Row][p1.Col]
}

// IsAdjacent checks if two positions are next to each other
func IsAdjacent(p1, p2 Position) bool {
	dr := p1.Row - p2.Row
	dc := p1.Col - p2.Col
	if dr < 0 {
		dr = -dr
	}
	if dc < 0 {
		dc = -dc
	}
	return (dr == 1 && dc == 0) || (dr == 0 && dc == 1)
}

// FindMatches locates all horizontal and vertical line matches (length >= 3)
func (b *Board) FindMatches() [][]Position {
	matched := make(map[Position]bool)

	// Horizontal check
	for r := 0; r < b.Rows; r++ {
		matchLen := 1
		for c := 0; c < b.Cols; c++ {
			if c < b.Cols-1 && b.Grid[r][c] != GemEmpty && b.Grid[r][c] == b.Grid[r][c+1] {
				matchLen++
			} else {
				if matchLen >= 3 {
					for k := 0; k < matchLen; k++ {
						matched[Position{Row: r, Col: c - k}] = true
					}
				}
				matchLen = 1
			}
		}
	}

	// Vertical check
	for c := 0; c < b.Cols; c++ {
		matchLen := 1
		for r := 0; r < b.Rows; r++ {
			if r < b.Rows-1 && b.Grid[r][c] != GemEmpty && b.Grid[r][c] == b.Grid[r+1][c] {
				matchLen++
			} else {
				if matchLen >= 3 {
					for k := 0; k < matchLen; k++ {
						matched[Position{Row: r - k, Col: c}] = true
					}
				}
				matchLen = 1
			}
		}
	}

	var result []Position
	for pos := range matched {
		result = append(result, pos)
	}
	return [][]Position{result}
}

// ApplyCascade clears matched gems, drops existing gems down, and refills the top
func (b *Board) ApplyCascade(matches []Position) int {
	if len(matches) == 0 {
		return 0
	}

	// 1. Mark matched positions empty
	for _, p := range matches {
		b.Grid[p.Row][p.Col] = GemEmpty
	}

	// 2. Drop gems down column by column
	for c := 0; c < b.Cols; c++ {
		writeRow := b.Rows - 1
		for r := b.Rows - 1; r >= 0; r-- {
			if b.Grid[r][c] != GemEmpty {
				b.Grid[writeRow][c] = b.Grid[r][c]
				if writeRow != r {
					b.Grid[r][c] = GemEmpty
				}
				writeRow--
			}
		}
		// Refill top with new random gems
		for r := writeRow; r >= 0; r-- {
			b.Grid[r][c] = GemType(rand.Intn(7) + 1)
		}
	}

	points := len(matches) * 50
	b.Score += points
	return points
}
