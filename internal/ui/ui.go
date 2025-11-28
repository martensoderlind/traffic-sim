package ui

import (
	"bytes"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type UIManager struct {
	buttons []*Button
	labels  []*Label
}

func NewUIManager() *UIManager {
	return &UIManager{
		buttons: make([]*Button, 0),
		labels:  make([]*Label, 0),
	}
}

func (ui *UIManager) AddButton(btn *Button) {
	ui.buttons = append(ui.buttons, btn)
}

func (ui *UIManager) AddLabel(lbl *Label) {
	ui.labels = append(ui.labels, lbl)
}

func (ui *UIManager) Update(mouseX, mouseY int, clicked bool) {
	for _, btn := range ui.buttons {
		btn.Update(mouseX, mouseY, clicked)
	}
}

func (ui *UIManager) Draw(screen *ebiten.Image) {
	for _, btn := range ui.buttons {
		btn.Draw(screen)
	}
	
	for _, lbl := range ui.labels {
		lbl.Draw(screen)
	}
}

func (ui *UIManager) Clear() {
	ui.buttons = make([]*Button, 0)
	ui.labels = make([]*Label, 0)
}

type Button struct {
	X, Y          float64
	Width, Height float64
	Text          string
	OnClick       func()
	
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
		float32(b.X),
		float32(b.Y),
		float32(b.Width),
		float32(b.Height),
		bgColor,
		false,
	)
	
	vector.StrokeRect(
		screen,
		float32(b.X),
		float32(b.Y),
		float32(b.Width),
		float32(b.Height),
		2,
		b.borderColor,
		false,
	)
	
	op := &text.DrawOptions{}
	op.GeoM.Translate(b.X+8, b.Y+b.Height/2-6)
	op.ColorScale.ScaleWithColor(b.textColor)
	text.Draw(screen, b.Text, &text.GoTextFace{
		Source: getDefaultFontSource(),
		Size:   14,
	}, op)
}

type Label struct {
	X, Y      float64
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
	}
}

func (l *Label) SetBackground(bg color.RGBA) {
	l.BgColor = &bg
}

func (l *Label) Draw(screen *ebiten.Image) {
	if l.BgColor != nil {
		textWidth := float64(len(l.Text)) * (l.Size * 0.6)
		textHeight := l.Size + 4
		
		vector.FillRect(
			screen,
			float32(l.X-l.Padding),
			float32(l.Y-l.Padding),
			float32(textWidth+l.Padding*2),
			float32(textHeight+l.Padding*2),
			*l.BgColor,
			false,
		)
		
		vector.StrokeRect(
			screen,
			float32(l.X-l.Padding),
			float32(l.Y-l.Padding),
			float32(textWidth+l.Padding*2),
			float32(textHeight+l.Padding*2),
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

func getDefaultFontSource() *text.GoTextFaceSource {
	if defaultFontSource == nil {
		source, err := text.NewGoTextFaceSource(bytes.NewReader(defaultFontData))
		if err != nil {
			panic(err)
		}
		defaultFontSource = source
	}
	return defaultFontSource
}