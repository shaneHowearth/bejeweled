package screen

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/shanehowearth/bejeweled/control"
	"github.com/shanehowearth/bejeweled/state"
)

type BejeweledView struct {
	App        fyne.App
	Win        fyne.Window
	Board      *state.Board
	Ctrl       *control.GameController
	GridBox    *fyne.Container
	ScoreLabel *widget.Label
	MovesLabel *widget.Label
}

func StartApp() {
	myApp := app.New()
	win := myApp.NewWindow("Bejeweled")

	board := state.NewBoard(8, 8)
	view := &BejeweledView{
		App:   myApp,
		Win:   win,
		Board: board,
	}

	ctrl := control.NewGameController(board, view.RefreshUI)
	view.Ctrl = ctrl

	view.setupUI()
	win.ShowAndRun()
}

func (v *BejeweledView) setupUI() {
	v.ScoreLabel = widget.NewLabel("Score: 0")
	v.MovesLabel = widget.NewLabel("Moves: 30")

	resetBtn := widget.NewButton("New Game", func() {
		v.Ctrl.Reset()
	})

	header := container.NewHBox(v.ScoreLabel, widget.NewSeparator(), v.MovesLabel, resetBtn)

	v.GridBox = container.NewGridWithColumns(v.Board.Cols)
	v.rebuildGrid()

	mainLayout := container.NewVBox(
		header,
		widget.NewSeparator(),
		v.GridBox,
	)

	v.Win.SetContent(container.NewPadded(mainLayout))
}

func (v *BejeweledView) rebuildGrid() {
	v.GridBox.Objects = nil
	for r := 0; r < v.Board.Rows; r++ {
		for c := 0; c < v.Board.Cols; c++ {
			row, col := r, c
			gw := NewGemWidget(row, col, v.Board.Grid[row][col], func(rPos, cPos int) {
				v.Ctrl.HandleGemClick(rPos, cPos)
			})
			if v.Ctrl.Selected != nil && v.Ctrl.Selected.Row == row && v.Ctrl.Selected.Col == col {
				gw.Selected = true
			}
			v.GridBox.Add(gw)
		}
	}
}

func (v *BejeweledView) RefreshUI() {
	v.ScoreLabel.SetText(fmt.Sprintf("Score: %d", v.Board.Score))
	v.MovesLabel.SetText(fmt.Sprintf("Moves: %d", v.Board.Moves))
	v.rebuildGrid()
	v.GridBox.Refresh()
}
