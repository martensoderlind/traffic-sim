package ui

import (
	"fmt"
	"image/color"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
)

type RoadPropertiesPanel struct {
	X, Y          float64
	Width, Height, shadowOffset float64
	Visible       bool
	
	bgColor     color.RGBA
	borderColor color.RGBA
	textColor   color.RGBA
	shadowColor color.RGBA
	
	labels      []*Label
	speedInput  *TextInput
	widthInput  *TextInput
	applyBtn    *Button
	closeBtn    *Button
	
	onApply func(maxSpeed, width float64)
}

func (p *RoadPropertiesPanel) Contains(x, y int) bool {
	if !p.Visible {
		return false
	}
	fx, fy := float64(x), float64(y)
	return fx >= p.X && fx <= p.X+p.Width && fy >= p.Y && fy <= p.Y+p.Height
}

func NewRoadPropertiesPanel(x, y float64) *RoadPropertiesPanel {
	panel := &RoadPropertiesPanel{
		X:           x,
		Y:           y,
		Width:       300,
		Height:      260,
		shadowOffset: 3,
		Visible:     false,
		bgColor:     color.RGBA{40, 40, 50, 240},
		borderColor: color.RGBA{100, 100, 110, 255},
		textColor:   color.RGBA{220, 220, 220, 255},
		shadowColor: color.RGBA{0, 0, 0, 80},
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
	
	p.applyBtn = NewButton(p.X+22, p.Y+200, 90, 30, "Apply", nil)
	p.closeBtn = NewButton(p.X+188,p.Y+200, 90, 30, "Close", nil)
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
	NewRect(
		float32(p.X+p.shadowOffset), float32(p.Y+p.shadowOffset),float32(p.Width),float32(p.Height), 13,p.shadowColor,
		).draw(screen)
	NewRect(
		float32(p.X),
		float32(p.Y),
		float32(p.Width),float32(p.Height),
		10,
		p.bgColor,
		).draw(screen)
	
	for _, label := range p.labels {
		label.Draw(screen)
	}
	
	p.speedInput.Draw(screen)
	p.widthInput.Draw(screen)
	p.applyBtn.Draw(screen)
	p.closeBtn.Draw(screen)
}