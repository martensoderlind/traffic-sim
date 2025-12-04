package ui

import (
	"fmt"
	"image/color"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type SpawnerPropertiesPanel struct {
	X, Y          float64
	Width, Height float64
	Visible       bool
	
	bgColor     color.RGBA
	borderColor color.RGBA
	textColor   color.RGBA
	
	labels      []*Label
	inputs []*TextInput
	IntervalInput *TextInput 
	MinSpeedInput *TextInput 
	MaxSpeedInput *TextInput
	MaxVehiclesInput *TextInput
	EnabledInput *TextInput

	applyBtn    *Button
	closeBtn    *Button
	
	onApply func(Interval,MinSpeed,MaxSpeed float64, MaxVehicles int,Enabled bool)
}

func (p *SpawnerPropertiesPanel) Contains(x, y int) bool {
	if !p.Visible {
		return false
	}
	fx, fy := float64(x), float64(y)
	return fx >= p.X && fx <= p.X+p.Width && fy >= p.Y && fy <= p.Y+p.Height
}

func NewSpawnerPropertiesPanel(x, y float64) *SpawnerPropertiesPanel {
	panel := &SpawnerPropertiesPanel{
		X:           x,
		Y:           y,
		Width:       300,
		Height:      0,
		Visible:     false,
		bgColor:     color.RGBA{40, 40, 50, 240},
		borderColor: color.RGBA{100, 100, 110, 255},
		textColor:   color.RGBA{220, 220, 220, 255},
		labels:      make([]*Label, 0),
	}
	
	panel.setupUI()
	return panel
}

func (p *SpawnerPropertiesPanel) setupUI() {
	titleLabel := NewLabel(p.X+15, p.Y+15, "Road Properties")
	titleLabel.Size = 16
	titleLabel.Color = color.RGBA{255, 255, 255, 255}
	p.labels = append(p.labels, titleLabel)
	
	intervalLabel := NewLabel(p.X+15, p.Y+50, "Interval:")
	intervalLabel.Size = 14
	p.labels = append(p.labels, intervalLabel)
	
	p.IntervalInput = NewTextInput(p.X+15, p.Y+70, 270, 35, "3.0")
	p.inputs=append(p.inputs, p.IntervalInput)

	minSpeedLabel := NewLabel(p.X+15, p.Y+120, "Min Speed:")
	minSpeedLabel.Size = 14
	p.labels = append(p.labels, minSpeedLabel)
	
	p.MinSpeedInput = NewTextInput(p.X+15, p.Y+140, 270, 35, "20.0")
	p.inputs=append(p.inputs, p.MinSpeedInput)

	maxSpeedLabel := NewLabel(p.X+15, p.Y+190, "Max Speed:")
	maxSpeedLabel.Size = 14
	p.labels = append(p.labels, maxSpeedLabel)
	
	p.MaxSpeedInput = NewTextInput(p.X+15, p.Y+210, 270, 35, "50.0")
	p.inputs=append(p.inputs, p.MaxSpeedInput)

	maxVehiclesLabel := NewLabel(p.X+15, p.Y+260, "Max vehicles:")
	maxVehiclesLabel.Size = 14
	p.labels = append(p.labels, maxVehiclesLabel)
	
	p.MaxVehiclesInput = NewTextInput(p.X+15, p.Y+280, 270, 35, "50")
	p.inputs=append(p.inputs, p.MaxVehiclesInput)

	enabledLabel := NewLabel(p.X+15, p.Y+330, "Enabled:")
	enabledLabel.Size = 14
	p.labels = append(p.labels, enabledLabel)
	
	p.EnabledInput = NewTextInput(p.X+15, p.Y+350, 270, 35, "true")
	p.inputs=append(p.inputs, p.EnabledInput)

	p.applyBtn = NewButton(p.X+22, p.Y+420, 90, 30, "Apply", nil)
	p.closeBtn = NewButton(p.X+188,p.Y+420, 90, 30, "Close", nil)

	p.calculateHeight()
}

func (p *SpawnerPropertiesPanel) Show(Interval,MinSpeed,MaxSpeed float64, MaxVehicles int,Enabled bool) {
	p.Visible = true
	p.IntervalInput.SetText(fmt.Sprintf("%.1f", Interval))
	p.MinSpeedInput.SetText(fmt.Sprintf("%.1f", MinSpeed))
	p.MaxSpeedInput.SetText(fmt.Sprintf("%.1f", MaxSpeed))
	p.MaxVehiclesInput.SetText(fmt.Sprintf("%d", MaxVehicles))
	p.EnabledInput.SetText(fmt.Sprintf("%t", Enabled))
}

func (p *SpawnerPropertiesPanel) Hide() {
	p.Visible = false
}

func (p *SpawnerPropertiesPanel) SetOnApply(callback func(Interval,MinSpeed,MaxSpeed float64, MaxVehicles int,Enabled bool)) {
	p.onApply = callback
}

func (p *SpawnerPropertiesPanel) Update(mouseX, mouseY int, clicked bool) {
	if !p.Visible {
		return
	}

	p.IntervalInput.Update(mouseX, mouseY, clicked)
	p.MinSpeedInput.Update(mouseX, mouseY, clicked)
	p.MaxSpeedInput.Update(mouseX, mouseY, clicked)
	p.MaxVehiclesInput.Update(mouseX, mouseY, clicked)
	p.EnabledInput.Update(mouseX, mouseY, clicked)
	
	p.applyBtn.Update(mouseX, mouseY, clicked)
	if p.applyBtn.pressed && p.onApply != nil {
		Interval,_:=strconv.ParseFloat(p.IntervalInput.GetText(), 64)
		MinSpeed,_:=strconv.ParseFloat(p.MinSpeedInput.GetText(), 64)
		MaxSpeed,_:=strconv.ParseFloat(p.MaxSpeedInput.GetText(), 64)
		MaxVehicles, _:= strconv.Atoi(p.MaxVehiclesInput.GetText())
		Enabled, _:= strconv.ParseBool(p.EnabledInput.GetText())
	
		if MinSpeed <= 0 {
			MinSpeed = 20.0
		}
		if MaxSpeed <= 0 || MaxSpeed < MinSpeed ||MaxSpeed> 200 {
			MaxSpeed = 50.0
		}
		if Interval <= 0 {
			Interval = 3.0
		}
		if MaxVehicles <= 0 {
			MaxVehicles = 50
		}
		
		p.onApply(Interval,MinSpeed,MaxSpeed, MaxVehicles,Enabled)
	}
	
	p.closeBtn.Update(mouseX, mouseY, clicked)
	if p.closeBtn.pressed {
		p.Hide()
	}
}

func (p *SpawnerPropertiesPanel) Draw(screen *ebiten.Image) {
	if !p.Visible {
		return
	}
	vector.FillRect(screen, float32(p.X), float32(p.Y), float32(p.Width), float32(p.Height), p.bgColor, false)
	vector.StrokeRect(screen, float32(p.X), float32(p.Y), float32(p.Width), float32(p.Height), 2, p.borderColor, false)
	
	for _, label := range p.labels {
		label.Draw(screen)
	}
	p.IntervalInput.Draw(screen)
	p.MinSpeedInput.Draw(screen)
	p.MaxSpeedInput.Draw(screen)
	p.MaxVehiclesInput.Draw(screen)
	p.EnabledInput.Draw(screen)
	p.applyBtn.Draw(screen)
	p.closeBtn.Draw(screen)
}

func (p *SpawnerPropertiesPanel) calculateHeight() {
	for _, label := range p.labels {
		p.Height += label.calculateHeight()+7
	}
	fmt.Println("Calculated panel height:", p.Height)
	for _, input := range p.inputs {
		p.Height += input.Height + 7
	}
	p.Height += p.applyBtn.Height 
	fmt.Println("Calculated panel height:", p.Height)
}