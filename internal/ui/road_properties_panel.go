package ui

import (
	"image/color"

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
	speedInput  *NumberInput
	widthInput  *NumberInput
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
	
	p.speedInput = NewNumberInput(p.X+15, p.Y+70, 270, 35, 40.0)
	
	widthLabel := NewLabel(p.X+15, p.Y+115, "Width:")
	widthLabel.Size = 14
	p.labels = append(p.labels, widthLabel)
	
	p.widthInput = NewNumberInput(p.X+15, p.Y+135, 270, 35, 8.0)
	
	p.applyBtn = NewButton(p.X+22, p.Y+200, 90, 30, "Apply", nil)
	p.closeBtn = NewButton(p.X+188,p.Y+200, 90, 30, "Close", nil)
}

func (p *RoadPropertiesPanel) Show(maxSpeed, width float64) {
	p.Visible = true
	p.speedInput.SetNumber( maxSpeed)
	p.widthInput.SetNumber( width)
}

func (p *RoadPropertiesPanel) Hide() {
	p.Visible = false
}

func (p *RoadPropertiesPanel) SetOnApply(callback func(maxSpeed, width float64)) {
	p.onApply = callback
}

func (p *RoadPropertiesPanel) SetPosition(x, y float64) {
	p.X = x
	p.Y = y
	p.updateUIPositions()
}

func (p *RoadPropertiesPanel) updateUIPositions() {
	for i, label := range p.labels {
		if i == 0 {
			label.X = p.X + 15
			label.Y = p.Y + 15
		} else if i == 1 {
			label.X = p.X + 15
			label.Y = p.Y + 50
		} else if i == 2 {
			label.X = p.X + 15
			label.Y = p.Y + 115
		}
	}
	
	p.speedInput.X = p.X + 15
	p.speedInput.Y = p.Y + 70
	
	p.widthInput.X = p.X + 15
	p.widthInput.Y = p.Y + 135
	
	
	p.applyBtn.X = p.X + 22
	p.applyBtn.Y = p.Y + 200
	
	p.closeBtn.X = p.X + 188
	p.closeBtn.Y = p.Y + 200
}

func (p *RoadPropertiesPanel) Update(mouseX, mouseY int, clicked bool) {
	if !p.Visible {
		return
	}
	
	p.speedInput.Update(mouseX, mouseY, clicked)
	p.widthInput.Update(mouseX, mouseY, clicked)
	
	p.applyBtn.Update(mouseX, mouseY, clicked)
	if p.applyBtn.pressed && p.onApply != nil {
		maxSpeed:=p.speedInput.GetNumber()
		
		width:= p.widthInput.GetNumber()
		
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