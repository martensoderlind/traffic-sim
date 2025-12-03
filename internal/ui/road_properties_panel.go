package ui

import (
	"fmt"
	"image/color"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type RoadPropertiesPanel struct {
	X, Y          float64
	Width, Height float64
	Visible       bool
	
	bgColor     color.RGBA
	borderColor color.RGBA
	textColor   color.RGBA
	
	labels      []*Label
	speedInput  *TextInput
	widthInput  *TextInput
	applyBtn    *Button
	closeBtn    *Button
	
	onApply func(maxSpeed, width float64)
}

func NewRoadPropertiesPanel(x, y float64) *RoadPropertiesPanel {
	panel := &RoadPropertiesPanel{
		X:           x,
		Y:           y,
		Width:       300,
		Height:      220,
		Visible:     false,
		bgColor:     color.RGBA{40, 40, 50, 240},
		borderColor: color.RGBA{100, 100, 110, 255},
		textColor:   color.RGBA{220, 220, 220, 255},
		labels:      make([]*Label, 0),
	}
	
	panel.setupUI()
	return panel
}

func (p *RoadPropertiesPanel) setupUI() {
	titleLabel := NewLabel(p.X+15, p.Y+15, "Road Properties")
	titleLabel.Size = 16
	titleLabel.Color = color.RGBA{255, 255, 255, 255}
	p.labels = append(p.labels, titleLabel)
	
	speedLabel := NewLabel(p.X+15, p.Y+50, "Max Speed:")
	speedLabel.Size = 14
	p.labels = append(p.labels, speedLabel)
	
	p.speedInput = NewTextInput(p.X+15, p.Y+70, 270, 35, "40.0")
	
	widthLabel := NewLabel(p.X+15, p.Y+115, "Width:")
	widthLabel.Size = 14
	p.labels = append(p.labels, widthLabel)
	
	p.widthInput = NewTextInput(p.X+15, p.Y+135, 270, 35, "8.0")
	
	p.applyBtn = NewButton(p.X+15, p.Y+180, 125, 30, "Apply", nil)
	p.closeBtn = NewButton(p.X+160, p.Y+180, 125, 30, "Close", nil)
}

func (p *RoadPropertiesPanel) Show(maxSpeed, width float64) {
	p.Visible = true
	p.speedInput.SetText(fmt.Sprintf("%.1f", maxSpeed))
	p.widthInput.SetText(fmt.Sprintf("%.1f", width))
}

func (p *RoadPropertiesPanel) Hide() {
	p.Visible = false
}

func (p *RoadPropertiesPanel) SetOnApply(callback func(maxSpeed, width float64)) {
	p.onApply = callback
}

func (p *RoadPropertiesPanel) Update(mouseX, mouseY int, clicked bool) {
	if !p.Visible {
		return
	}
	
	p.speedInput.Update(mouseX, mouseY, clicked)
	p.widthInput.Update(mouseX, mouseY, clicked)
	
	p.applyBtn.Update(mouseX, mouseY, clicked)
	if p.applyBtn.pressed && p.onApply != nil {
		maxSpeed, _ := strconv.ParseFloat(p.speedInput.GetText(), 64)
		width, _ := strconv.ParseFloat(p.widthInput.GetText(), 64)
		
		if maxSpeed <= 0 {
			maxSpeed = 40.0
		}
		if width <= 0 {
			width = 8.0
		}
		
		p.onApply(maxSpeed, width)
	}
	
	p.closeBtn.Update(mouseX, mouseY, clicked)
	if p.closeBtn.pressed {
		p.Hide()
	}
}

func (p *RoadPropertiesPanel) Draw(screen *ebiten.Image) {
	if !p.Visible {
		return
	}
	
	vector.FillRect(screen, float32(p.X), float32(p.Y), float32(p.Width), float32(p.Height), p.bgColor, false)
	vector.StrokeRect(screen, float32(p.X), float32(p.Y), float32(p.Width), float32(p.Height), 2, p.borderColor, false)
	
	for _, label := range p.labels {
		label.Draw(screen)
	}
	
	p.speedInput.Draw(screen)
	p.widthInput.Draw(screen)
	p.applyBtn.Draw(screen)
	p.closeBtn.Draw(screen)
}

type TextInput struct {
	X, Y          float64
	Width, Height float64
	Text          string
	Active        bool
	
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
		bgColor:       color.RGBA{30, 30, 40, 255},
		activeBgColor: color.RGBA{40, 40, 50, 255},
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
	
	if ebiten.IsKeyPressed(ebiten.KeyBackspace) {
		if len(ti.Text) > 0 {
			ti.Text = ti.Text[:len(ti.Text)-1]
		}
	}
}

func (ti *TextInput) Draw(screen *ebiten.Image) {
	bgColor := ti.bgColor
	if ti.Active {
		bgColor = ti.activeBgColor
	}
	
	vector.FillRect(screen, float32(ti.X), float32(ti.Y), float32(ti.Width), float32(ti.Height), bgColor, false)
	vector.StrokeRect(screen, float32(ti.X), float32(ti.Y), float32(ti.Width), float32(ti.Height), 2, ti.borderColor, false)
	
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