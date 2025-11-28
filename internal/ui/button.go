package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Button struct {
	X, Y          float64
	Width, Height float64
	Text          string
	OnClick       func()
	Padding float64
	
	hovered bool
	pressed bool
	
	bgColor      color.RGBA
	hoverColor   color.RGBA
	pressColor   color.RGBA
	textColor    color.RGBA
	borderColor  color.RGBA
}

func NewButton(x, y, width, height float64, text string, onClick func()) *Button {
	return &Button{
		X:           x,
		Y:           y,
		Width:       width,
		Height:      height,
		Text:        text,
		OnClick:     onClick,
		Padding: 8,
		bgColor:     color.RGBA{60, 60, 70, 255},
		hoverColor:  color.RGBA{80, 80, 90, 255},
		pressColor:  color.RGBA{50, 50, 60, 255},
		textColor:   color.RGBA{220, 220, 220, 255},
		borderColor: color.RGBA{100, 100, 110, 255},
	}
}

func (b *Button) SetColors(bg, hover, press, text, border color.RGBA) {
	b.bgColor = bg
	b.hoverColor = hover
	b.pressColor = press
	b.textColor = text
	b.borderColor = border
}

func (b *Button) Contains(x, y int) bool {
	fx, fy := float64(x), float64(y)
	return fx >= b.X && fx <= b.X+b.Width && fy >= b.Y && fy <= b.Y+b.Height
}

func (b *Button) Update(mouseX, mouseY int, clicked bool) {
	b.hovered = b.Contains(mouseX, mouseY)
	
	if b.hovered && clicked {
		b.pressed = true
		if b.OnClick != nil {
			b.OnClick()
		}
	} else {
		b.pressed = false
	}
}

func (b *Button) Draw(screen *ebiten.Image) {
	bgColor := b.bgColor
	if b.pressed {
		bgColor = b.pressColor
	} else if b.hovered {
		bgColor = b.hoverColor
	}
	
	vector.FillRect(
		screen,
			float32(b.X-b.Padding),
			float32(b.Y-b.Padding),
			float32(b.Width+b.Padding*2),
			float32(b.Height+b.Padding*2),
		bgColor,
		false,
	)
	
	vector.StrokeRect(
		screen,
		float32(b.X-b.Padding),
		float32(b.Y-b.Padding),
		float32(b.Width+b.Padding*2),
		float32(b.Height+b.Padding*2),
		2,
		b.borderColor,
		false,
	)
	
	op := &text.DrawOptions{}
	op.GeoM.Translate(b.X+b.Padding, b.Y+b.Padding)
	op.ColorScale.ScaleWithColor(b.textColor)
	text.Draw(screen, b.Text, &text.GoTextFace{
		Source: getDefaultFontSource(),
		Size:   14,
	}, op)
}