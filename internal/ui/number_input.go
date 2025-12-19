package ui

import (
	"image/color"
	"strconv"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)
type NumberInput struct {
	X, Y          	float64
	Width, Height 	float64

	text          	string
	value		 	float64
	
	Step 			float64

	incrementValueBtn    *Button
	decrementValueBtn     *Button

	Active        	bool
	backspaceCooldown int

	bgColor       	color.RGBA
	activeBgColor 	color.RGBA
	borderColor   	color.RGBA
	textColor     	color.RGBA
}

func NewNumberInput(x, y, width, height float64, initial float64) *NumberInput {
	ni := &NumberInput{
		X:      x,
		Y:      y,
		Width:  width,
		Height: height,

		value: initial,
		Step:  1.0,

		bgColor:       color.RGBA{30, 30, 40, 255},
		activeBgColor: color.RGBA{35, 35, 45, 255},
		borderColor:   color.RGBA{100, 100, 110, 255},
		textColor:     color.RGBA{220, 220, 220, 255},
	}

	ni.text = strconv.FormatFloat(initial, 'f', -1, 64)
	ni.setupUI()
	return ni
}

func (ni *NumberInput) setupUI() {
	ni.incrementValueBtn = NewButton(ni.X+ni.Width+5, ni.Y, ni.Height/2, ni.Height/2, "+", func() {
		ni.IncrementNumber()
	})
	ni.incrementValueBtn.size = 8
	ni.incrementValueBtn.Padding = 4
	ni.decrementValueBtn = NewButton(ni.X+ni.Width+5, ni.Y+ni.Height/2, ni.Height/2, ni.Height/2, "-", func(){ni.decrementNumber()})
	ni.decrementValueBtn.size = 8
	ni.decrementValueBtn.Padding = 4
}

func (ni *NumberInput) Contains(x, y int) bool {
	fx, fy := float64(x), float64(y)
	return fx >= ni.X && fx <= ni.X+ni.Width && fy >= ni.Y && fy <= ni.Y+ni.Height
}

func (ni *NumberInput) Update(mouseX, mouseY int, clicked bool) {
	if ni.incrementValueBtn != nil {
		ni.incrementValueBtn.Update(mouseX, mouseY, clicked)
	}
	if ni.decrementValueBtn != nil {
		ni.decrementValueBtn.Update(mouseX, mouseY, clicked)
	}

	if clicked {
		if ni.incrementValueBtn != nil && ni.incrementValueBtn.Contains(mouseX, mouseY) {
			return
		}
		if ni.decrementValueBtn != nil && ni.decrementValueBtn.Contains(mouseX, mouseY) {
			return
		}
	}

	if clicked {
		ni.Active = ni.Contains(mouseX, mouseY)
	}

	if ni.Active {
		ni.handleInput()
	}
}


func (ni *NumberInput) handleInput() {
	inputChars := ebiten.AppendInputChars(nil)
	for _, ch := range inputChars {

		if ch >= '0' && ch <= '9' {
			ni.text += string(ch)
			continue
		}

		if ch == '.' && !strings.Contains(ni.text, ".") {
			ni.text += "."
		}
	}

	// backspace
	if ni.backspaceCooldown > 0 {
		ni.backspaceCooldown--
	}

	if ebiten.IsKeyPressed(ebiten.KeyBackspace) && ni.backspaceCooldown == 0 {
		if len(ni.text) > 0 {
			ni.text = ni.text[:len(ni.text)-1]
		}
		ni.backspaceCooldown = 6
	}
	ni.parseValue()
}

func (ni *NumberInput) parseValue() {
	if ni.text == "" || ni.text == "." {
		ni.value = 0
		return
	}

	if v, err := strconv.ParseFloat(ni.text, 64); err == nil {
		ni.value = v
	}
}

func (ni *NumberInput) Draw(screen *ebiten.Image) {
	bgColor := ni.bgColor
	if ni.Active {
		bgColor = ni.activeBgColor
	}

	vector.FillRect(
		screen,
		float32(ni.X), float32(ni.Y),
		float32(ni.Width), float32(ni.Height),
		bgColor, false,
	)

	display := ni.text
	if ni.Active {
		display += "|"
	}

	op := &text.DrawOptions{}
	op.GeoM.Translate(ni.X+10, ni.Y+8)
	op.ColorScale.ScaleWithColor(ni.textColor)

	text.Draw(
		screen,
		display,
		&text.GoTextFace{
			Source: getDefaultFontSource(),
			Size:   14,
		},
		op,
	)
	ni.incrementValueBtn.Draw(screen)
	ni.decrementValueBtn.Draw(screen)
}

func (ni *NumberInput) GetNumber() float64 {
	return ni.value
}

func (ni *NumberInput) SetNumber(v float64) {
	ni.value = v
	ni.text = strconv.FormatFloat(v, 'f', -1, 64)
}

func (ni *NumberInput) IncrementNumber() {
	ni.value += ni.Step
	ni.text = strconv.FormatFloat(ni.value, 'f', -1, 64)
}

func (ni *NumberInput) decrementNumber() {
	ni.value -= ni.Step
	ni.text = strconv.FormatFloat(ni.value, 'f', -1, 64)
}
