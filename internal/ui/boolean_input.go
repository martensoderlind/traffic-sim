package ui

import (
	"image/color"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
)
type BoolInput struct {
	X, Y          	float64
	Width, Height 	float64

	text          	string
	value		 	bool
	
	Step 			float64

	TrueValueLabel		*Label
	FalseValueLabel		*Label
	TrueValueBtn    *Button
	FalseValueBtn     *Button

	Active        	bool

	bgColor       	color.RGBA
	activeBgColor 	color.RGBA
	borderColor   	color.RGBA
	textColor     	color.RGBA
}

func NewBoolInput(x, y, width, height float64, initial bool) *BoolInput {
	bi := &BoolInput{
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

	bi.text = strconv.FormatBool(initial)
	bi.setupUI()
	return bi
}

func (bi *BoolInput) setupUI() {
	bi.TrueValueLabel = NewLabel(bi.X+5, bi.Y+5, "True")
	bi.TrueValueBtn = NewButton(bi.X+30, bi.Y+5,bi.Height-10, bi.Height-10, "X", func() {
		bi.ToggleValueTrue()
	})
	bi.TrueValueBtn.SizeMode = ButtonFixedSize
	bi.TrueValueBtn.size = 10
	bi.TrueValueBtn.Padding = 9
	
	bi.FalseValueLabel = NewLabel(bi.X+35, bi.Y+5, "False")
	bi.FalseValueBtn = NewButton(bi.X+60, bi.Y+5,bi.Height-10, bi.Height-10, " ", func(){bi.ToggleValueFalse()})
	bi.FalseValueBtn.SizeMode = ButtonFixedSize
	bi.FalseValueBtn.size = 10
	bi.FalseValueBtn.Padding = 9
}

func (bi *BoolInput) Contains(x, y int) bool {
	fx, fy := float64(x), float64(y)
	return fx >= bi.X && fx <= bi.X+bi.Width && fy >= bi.Y && fy <= bi.Y+bi.Height
}

func (bi *BoolInput) Update(mouseX, mouseY int, clicked bool) {
	if bi.TrueValueBtn != nil {
		bi.TrueValueBtn.Update(mouseX, mouseY, clicked)
	}
	if bi.FalseValueBtn != nil {
		bi.FalseValueBtn.Update(mouseX, mouseY, clicked)
	}

	if clicked {
		if bi.TrueValueBtn != nil && bi.TrueValueBtn.Contains(mouseX, mouseY) {
			return
		}
		if bi.FalseValueBtn != nil && bi.FalseValueBtn.Contains(mouseX, mouseY) {
			return
		}
	}

	if clicked {
		bi.Active = bi.Contains(mouseX, mouseY)
	}

}

func (bi *BoolInput) Draw(screen *ebiten.Image) {
	

	// display := bi.text
	// if bi.Active {
	// 	display += "|"
	// }

	// op := &text.DrawOptions{}
	// op.GeoM.Translate(bi.X+10, bi.Y+8)
	// op.ColorScale.ScaleWithColor(bi.textColor)

	// text.Draw(
	// 	screen,
	// 	display,
	// 	&text.GoTextFace{
	// 		Source: getDefaultFontSource(),
	// 		Size:   14,
	// 	},
	// 	op,
	// )
	bi.TrueValueLabel.Draw(screen)
	bi.FalseValueLabel.Draw(screen)
	bi.TrueValueBtn.Draw(screen)
	bi.FalseValueBtn.Draw(screen)
}

func (bi *BoolInput) GetValue() bool {
	return bi.value
}

func (bi *BoolInput) SetValue(v bool) {
	bi.value = v
	bi.text = strconv.FormatBool(v)
}

func (bi *BoolInput) ToggleValueTrue() {
	bi.value = true
	bi.text = strconv.FormatBool(bi.value)
	bi.TrueValueBtn.Text = "X"
	bi.FalseValueBtn.Text = " "
}

func (bi *BoolInput) ToggleValueFalse() {
	bi.value = false
	bi.text = strconv.FormatBool(bi.value)
	bi.TrueValueBtn.Text = " "
	bi.FalseValueBtn.Text = "X"
}
