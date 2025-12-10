package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)
type TextInput struct {
	X, Y          float64
	Width, Height float64
	Text          string
	Active        bool
	backspaceCooldown int

	bgColor       color.RGBA
	activeBgColor color.RGBA
	borderColor   color.RGBA
	textColor     color.RGBA
}

func NewTextInput(x, y, width, height float64, initialText string) *TextInput {
	return &TextInput{
		X:             x,
		Y:             y,
		Width:         width,
		Height:        height,
		Text:          initialText,
		Active:        false,
		backspaceCooldown: 0,
		bgColor:       color.RGBA{30, 30, 40, 255},
		activeBgColor: color.RGBA{35, 35, 45, 255},
		borderColor:   color.RGBA{100, 100, 110, 255},
		textColor:     color.RGBA{220, 220, 220, 255},
	}
}

func (ti *TextInput) Contains(x, y int) bool {
	fx, fy := float64(x), float64(y)
	return fx >= ti.X && fx <= ti.X+ti.Width && fy >= ti.Y && fy <= ti.Y+ti.Height
}

func (ti *TextInput) Update(mouseX, mouseY int, clicked bool) {
	if clicked {
		ti.Active = ti.Contains(mouseX, mouseY)
	}
	
	if ti.Active {
		ti.handleInput()
	}
}


func (ti *TextInput) handleInput() {
	inputChars := ebiten.AppendInputChars(nil)
	for _, ch := range inputChars {
		if (ch >= '0' && ch <= '9') || ch == '.' {
			ti.Text += string(ch)
		}
	}
	
	if ti.backspaceCooldown > 0 {
        ti.backspaceCooldown--
    }

    if ebiten.IsKeyPressed(ebiten.KeyBackspace) && ti.backspaceCooldown == 0 {
        if len(ti.Text) > 0 {
            ti.Text = ti.Text[:len(ti.Text)-1]
        }
        ti.backspaceCooldown = 6
    }
}

func (ti *TextInput) Draw(screen *ebiten.Image) {
	bgColor := ti.bgColor
	if ti.Active {
		bgColor = ti.activeBgColor
	}
	
	vector.FillRect(screen, float32(ti.X), float32(ti.Y), float32(ti.Width), float32(ti.Height), bgColor, false)
	
	displayText := ti.Text
	if ti.Active && len(displayText) < 20 {
		displayText += "|"
	}
	
	op := &text.DrawOptions{}
	op.GeoM.Translate(ti.X+10, ti.Y+8)
	op.ColorScale.ScaleWithColor(ti.textColor)
	text.Draw(screen, displayText, &text.GoTextFace{
		Source: getDefaultFontSource(),
		Size:   14,
	}, op)
}

func (ti *TextInput) SetText(text string) {
	ti.Text = text
}

func (ti *TextInput) GetText() string {
	return ti.Text
}