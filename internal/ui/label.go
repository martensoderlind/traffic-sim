package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Label struct {
	X, Y      float64
	height   float64
	Text      string
	Color     color.RGBA
	Size      float64
	BgColor   *color.RGBA
	Padding   float64
}

func NewLabel(x, y float64, text string) *Label {
	return &Label{
		X:       x,
		Y:       y,
		Text:    text,
		Color:   color.RGBA{220, 220, 220, 255},
		Size:    14,
		Padding: 8,
		height: 0,
	}
}

func (l *Label) SetBackground(bg color.RGBA) {
	l.BgColor = &bg
}

func (l *Label) Draw(screen *ebiten.Image) {
	if l.BgColor != nil {
		textWidth := float64(len(l.Text)) * (l.Size * 0.6)
		l.height = l.calculateHeight()
		
		vector.FillRect(
			screen,
			float32(l.X-l.Padding),
			float32(l.Y-l.Padding),
			float32(textWidth+l.Padding*2),
			float32(l.height),
			*l.BgColor,
			false,
		)
		
		vector.StrokeRect(
			screen,
			float32(l.X-l.Padding),
			float32(l.Y-l.Padding),
			float32(textWidth+l.Padding*2),
			float32(l.height),
			1,
			color.RGBA{100, 100, 110, 255},
			false,
		)
	}
	
	op := &text.DrawOptions{}
	op.GeoM.Translate(l.X, l.Y)
	op.ColorScale.ScaleWithColor(l.Color)
	text.Draw(screen, l.Text, &text.GoTextFace{
		Source: getDefaultFontSource(),
		Size:   l.Size,
	}, op)
}

var defaultFontSource *text.GoTextFaceSource

func (l *Label) calculateHeight() float64 {
	return l.Size + 4 + l.Padding*2
}