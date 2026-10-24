package screen

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
	"github.com/shanehowearth/bejeweled/state"
)

type GemWidget struct {
	widget.BaseWidget
	Row      int
	Col      int
	GemType  state.GemType
	Selected bool
	OnTapped func(row, col int)
}

func NewGemWidget(row, col int, gType state.GemType, onTapped func(row, col int)) *GemWidget {
	gw := &GemWidget{
		Row:      row,
		Col:      col,
		GemType:  gType,
		OnTapped: onTapped,
	}
	gw.ExtendBaseWidget(gw)
	return gw
}

func (gw *GemWidget) Tapped(_ *fyne.PointEvent) {
	if gw.OnTapped != nil {
		gw.OnTapped(gw.Row, gw.Col)
	}
}

func (gw *GemWidget) CreateRenderer() fyne.WidgetRenderer {
	bg := canvas.NewRectangle(color.RGBA{R: 35, G: 35, B: 45, A: 255})

	// Draw colorful Gem shape representation
	gemShape := canvas.NewCircle(getGemColor(gw.GemType))

	border := canvas.NewRectangle(color.Transparent)
	border.StrokeWidth = 3
	if gw.Selected {
		border.StrokeColor = color.RGBA{R: 255, G: 255, B: 0, A: 255} // Highlight yellow on selection
	} else {
		border.StrokeColor = color.RGBA{R: 60, G: 60, B: 70, A: 255}
	}

	return &gemWidgetRenderer{
		gw:      gw,
		bg:      bg,
		gem:     gemShape,
		border:  border,
		objects: []fyne.CanvasObject{bg, gemShape, border},
	}
}

func getGemColor(g state.GemType) color.Color {
	switch g {
	case state.GemRed:
		return color.RGBA{R: 230, G: 40, B: 40, A: 255}
	case state.GemBlue:
		return color.RGBA{R: 40, G: 120, B: 240, A: 255}
	case state.GemGreen:
		return color.RGBA{R: 40, G: 200, B: 80, A: 255}
	case state.GemYellow:
		return color.RGBA{R: 240, G: 210, B: 40, A: 255}
	case state.GemPurple:
		return color.RGBA{R: 160, G: 50, B: 220, A: 255}
	case state.GemOrange:
		return color.RGBA{R: 240, G: 130, B: 30, A: 255}
	case state.GemWhite:
		return color.RGBA{R: 230, G: 230, B: 240, A: 255}
	default:
		return color.Transparent
	}
}

type gemWidgetRenderer struct {
	gw      *GemWidget
	bg      *canvas.Rectangle
	gem     *canvas.Circle
	border  *canvas.Rectangle
	objects []fyne.CanvasObject
}

func (r *gemWidgetRenderer) Layout(size fyne.Size) {
	r.bg.Resize(size)
	r.border.Resize(size)

	// Inset gem inside tile
	padding := float32(6)
	r.gem.Move(fyne.NewPos(padding, padding))
	r.gem.Resize(fyne.NewSize(size.Width-padding*2, size.Height-padding*2))
}

func (r *gemWidgetRenderer) MinSize() fyne.Size {
	return fyne.NewSize(48, 48)
}

func (r *gemWidgetRenderer) Refresh() {
	r.gem.FillColor = getGemColor(r.gw.GemType)
	if r.gw.Selected {
		r.border.StrokeColor = color.RGBA{R: 255, G: 255, B: 0, A: 255}
	} else {
		r.border.StrokeColor = color.RGBA{R: 60, G: 60, B: 70, A: 255}
	}
	r.bg.Refresh()
	r.gem.Refresh()
	r.border.Refresh()
}

func (r *gemWidgetRenderer) Objects() []fyne.CanvasObject { return r.objects }
func (r *gemWidgetRenderer) Destroy()                     {}
