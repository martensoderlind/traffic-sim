package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type Button struct {
	X, Y          float64
	Width, Height float64
	Text          string
	OnClick       func()
	Padding float64
	size	 float64

	hovered bool
	pressed bool
	
	bgColor      color.RGBA
	hoverColor   color.RGBA
	pressColor   color.RGBA
	textColor    color.RGBA
	borderColor  color.RGBA

	Icon       *ebiten.Image
	IconWidth  float64
	IconHeight float64

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
		size:	 14,
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
	return fx >= b.X-b.Padding && fx <= b.X+b.Width+b.Padding && fy >= b.Y-b.Padding && fy <= b.Y+b.Height+b.Padding
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

func (b *Button) SetIcon(img *ebiten.Image, w, h float64) {
	b.Icon = img
	b.IconWidth = w
	b.IconHeight = h
}


func (b *Button) Draw(screen *ebiten.Image) {
	bgColor := b.bgColor
	buttonHeight := b.calculateHeight()
	buttonWidth:= b.calculateWidth()
	if b.pressed {
		bgColor = b.pressColor
	} else if b.hovered {
		bgColor = b.hoverColor
	}
	
	NewRect(
		float32(b.X-b.Padding),
		float32(b.Y-b.Padding),
		buttonWidth,
		buttonHeight,
		6,
		bgColor,
		).draw(screen)

	cursorX := b.X + b.Padding
	cursorY := b.Y + b.Padding

	if b.Icon != nil {
		op := &ebiten.DrawImageOptions{}

		innerHeight := b.Height - b.Padding*2
		dy := b.Y + b.Padding + (innerHeight - b.IconHeight) / 2

		op.GeoM.Translate(cursorX, dy)

		w, h := b.Icon.Bounds().Dx(), b.Icon.Bounds().Dy()
		op.GeoM.Scale(b.IconWidth/float64(w), b.IconHeight/float64(h))

		screen.DrawImage(b.Icon, op)

		cursorX += b.IconWidth + 6
	}
	textOp := &text.DrawOptions{}
	textOp.GeoM.Translate(cursorX, cursorY)
	textOp.ColorScale.ScaleWithColor(b.textColor)

	text.Draw(screen, b.Text, &text.GoTextFace{
		Source: getDefaultFontSource(),
		Size:   b.size,
	}, textOp)
}

func (b *Button) calculateHeight()float32 {
	return float32(b.Height + b.Padding*2)
}
func (b * Button) calculateWidth() float32 {
	textWidth := float64(len(b.Text)) * (b.size * 0.6)
	return float32(textWidth+ b.Padding*2 )
}