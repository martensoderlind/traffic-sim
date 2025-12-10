package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)
type Rect struct {
	X, Y          float32
	Width, Height float32
	Radius	   float32
	Text          string
	Active        bool

	bgColor       color.RGBA
}
func NewRect( x, y, w, h, r float32, bgCol color.Color) *Rect {
	return &Rect{
		X:x	,
		Y:y,
		Width:w,
		Height:h,
		Radius:r,
		bgColor: bgCol.(color.RGBA),
	}
}

func (r *Rect) draw(screen *ebiten.Image) {
	var path vector.Path

	path.MoveTo(r.X+r.Radius, r.Y)
	path.LineTo(r.X+r.Width-r.Radius, r.Y)
	path.ArcTo(r.X+r.Width, r.Y, r.X+r.Width, r.Y+r.Radius, r.Radius)
	path.LineTo(r.X+r.Width, r.Y+r.Height-r.Radius)
	path.ArcTo(r.X+r.Width, r.Y+r.Height, r.X+r.Width-r.Radius, r.Y+r.Height, r.Radius)
	path.LineTo(r.X+r.Radius, r.Y+r.Height)
	path.ArcTo(r.X, r.Y+r.Height, r.X, r.Y+r.Height-r.Radius, r.Radius) 
	path.LineTo(r.X, r.Y+r.Radius)
	path.ArcTo(r.X, r.Y, r.X+r.Radius, r.Y, r.Radius) 
	path.Close()

	fillOpts := &vector.FillOptions{} 

	drawOpts := &vector.DrawPathOptions{}
	drawOpts.ColorScale.ScaleWithColor(r.bgColor)

	vector.FillPath(screen, &path, fillOpts, drawOpts)
}

// drawRoundedRect(screen, 50, 50, 200, 100, 10, color.RGBA{255, 100, 50, 255})