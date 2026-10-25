package state

import (
	"math/rand"
	"time"
)

type GemType int

const (
	GemEmpty GemType = iota
	GemRed
	GemBlue
	GemGreen
	GemYellow
	GemPurple
	GemOrange
	GemWhite
)

type Position struct {
	Row int
	Col int
}

type Board struct {
	Rows     int
	Cols     int
	Grid     [][]GemType
	Score    int
	GameOver bool
}

func NewBoard(rows, cols int) *Board {
	rand.Seed(time.Now().UnixNano())
	b := &Board{
		Rows: rows,
		Cols: cols,
		Grid: make([][]GemType, rows),
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
				if c >= 2 && b.Grid[r][c-1] == g && b.Grid[r][c-2] == g {
					continue
				}
				if r >= 2 && b.Grid[r-1][c] == g && b.Grid[r-2][c] == g {
					continue
				}
				break
			}
			b.Grid[r][c] = g
		}
	}
	// Ensure there is at least one valid move available from the start
	if !b.HasValidMoves() {
		b.FillBoardWithoutMatches()
	}
}

func (b *Board) Swap(p1, p2 Position) {
	b.Grid[p1.Row][p1.Col], b.Grid[p2.Row][p2.Col] = b.Grid[p2.Row][p2.Col], b.Grid[p1.Row][p1.Col]
}

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

func (b *Board) FindMatches() [][]Position {
	matched := make(map[Position]bool)

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

// HasValidMoves checks if any adjacent swap can create a match of 3+
func (b *Board) HasValidMoves() bool {
	for r := 0; r < b.Rows; r++ {
		for c := 0; c < b.Cols; c++ {
			p1 := Position{Row: r, Col: c}
			// Check right neighbor
			if c+1 < b.Cols {
				p2 := Position{Row: r, Col: c + 1}
				b.Swap(p1, p2)
				matches := b.FindMatches()
				b.Swap(p1, p2) // Swap back
				if len(matches[0]) > 0 {
					return true
				}
			}
			// Check down neighbor
			if r+1 < b.Rows {
				p2 := Position{Row: r + 1, Col: c}
				b.Swap(p1, p2)
				matches := b.FindMatches()
				b.Swap(p1, p2) // Swap back
				if len(matches[0]) > 0 {
					return true
				}
			}
		}
	}
	return false
}

func (b *Board) ApplyCascade(matches []Position) int {
	if len(matches) == 0 {
		return 0
	}

	for _, p := range matches {
		b.Grid[p.Row][p.Col] = GemEmpty
	}

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
		for r := writeRow; r >= 0; r-- {
			b.Grid[r][c] = GemType(rand.Intn(7) + 1)
		}
	}

	points := len(matches) * 50
	b.Score += points
	return points
}
