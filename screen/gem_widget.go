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
	Gem      state.GemType
	Selected bool
	Matched  bool
	OnClick  func(r, c int)
}

func NewGemWidget(r, c int, g state.GemType, onClick func(r, c int)) *GemWidget {
	gw := &GemWidget{
		Row:     r,
		Col:     c,
		Gem:     g,
		OnClick: onClick,
	}
	gw.ExtendBaseWidget(gw)
	return gw
}

func (gw *GemWidget) Tapped(*fyne.PointEvent) {
	if gw.OnClick != nil {
		gw.OnClick(gw.Row, gw.Col)
	}
}

func (gw *GemWidget) CreateRenderer() fyne.WidgetRenderer {
	circle := canvas.NewCircle(getGemColor(gw.Gem))
	border := canvas.NewRectangle(color.Transparent)
	border.StrokeWidth = 2
	border.StrokeColor = color.Transparent

	return &gemWidgetRenderer{
		gw:      gw,
		circle:  circle,
		border:  border,
		objects: []fyne.CanvasObject{circle, border},
	}
}

func getGemColor(g state.GemType) color.Color {
	switch g {
	case state.GemRed:
		return color.RGBA{R: 220, G: 50, B: 50, A: 255}
	case state.GemBlue:
		return color.RGBA{R: 50, G: 100, B: 220, A: 255}
	case state.GemGreen:
		return color.RGBA{R: 50, G: 200, B: 80, A: 255}
	case state.GemYellow:
		return color.RGBA{R: 240, G: 220, B: 50, A: 255}
	case state.GemPurple:
		return color.RGBA{R: 160, G: 50, B: 200, A: 255}
	case state.GemOrange:
		return color.RGBA{R: 240, G: 140, B: 40, A: 255}
	case state.GemWhite:
		return color.RGBA{R: 240, G: 240, B: 240, A: 255}
	default:
		return color.Transparent
	}
}

type gemWidgetRenderer struct {
	gw      *GemWidget
	circle  *canvas.Circle
	border  *canvas.Rectangle
	objects []fyne.CanvasObject
}

func (r *gemWidgetRenderer) Layout(size fyne.Size) {
	r.circle.Resize(size)
	r.border.Resize(size)
}

func (r *gemWidgetRenderer) MinSize() fyne.Size {
	return fyne.NewSize(40, 40)
}

func (r *gemWidgetRenderer) Refresh() {
	r.circle.FillColor = getGemColor(r.gw.Gem)

	// Highlight styling
	if r.gw.Matched {
		// Bright gold outline + white glow inside for matched tiles
		r.border.StrokeWidth = 4
		r.border.StrokeColor = color.RGBA{R: 255, G: 255, B: 255, A: 255}
		r.circle.FillColor = color.RGBA{R: 255, G: 255, B: 200, A: 255}
	} else if r.gw.Selected {
		r.border.StrokeWidth = 3
		r.border.StrokeColor = color.RGBA{R: 255, G: 215, B: 0, A: 255} // Gold outline
	} else {
		r.border.StrokeWidth = 1
		r.border.StrokeColor = color.RGBA{R: 40, G: 40, B: 50, A: 255}
	}

	r.circle.Refresh()
	r.border.Refresh()
}

func (r *gemWidgetRenderer) Objects() []fyne.CanvasObject { return r.objects }
func (r *gemWidgetRenderer) Destroy()                     {}
