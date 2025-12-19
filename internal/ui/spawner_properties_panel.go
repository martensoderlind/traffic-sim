package ui

import (
	"fmt"
	"image/color"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
)

type SpawnerPropertiesPanel struct {
	X, Y          float64
	Width, Height, shadowOffset float64
	Visible       bool
	
	bgColor     color.RGBA
	borderColor color.RGBA
	textColor   color.RGBA
	shadowColor color.RGBA

	labels      []*Label
	inputs []*NumberInput
	IntervalInput *NumberInput 
	MinSpeedInput *NumberInput 
	MaxSpeedInput *NumberInput
	MaxVehiclesInput *NumberInput
	EnabledInput *TextInput

	applyBtn    *Button
	closeBtn    *Button

	btnWidth, btnHeight float64
	
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
		shadowOffset: 3,
		Visible:     false,
		bgColor:     color.RGBA{40, 40, 50, 240},
		borderColor: color.RGBA{100, 100, 110, 255},
		textColor:   color.RGBA{220, 220, 220, 255},
		shadowColor: color.RGBA{0, 0, 0, 80},
		labels:      make([]*Label, 0),
		btnWidth: 120.0,
		btnHeight: 28.0,
	}
	
	panel.setupUI()
	return panel
}

func (p *SpawnerPropertiesPanel) setupUI() {
	titleLabel := NewLabel(p.X+15, p.Y+15, "Road Properties")
	titleLabel.Size = 16
	titleLabel.Color = color.RGBA{255, 255, 255, 255}
	p.labels = append(p.labels, titleLabel)
	
	yOffset := 50

	intervalLabel := NewLabel(p.X+15, p.Y+50, "Interval:")
	intervalLabel.Size = 14
	p.labels = append(p.labels, intervalLabel)
	
	p.IntervalInput = NewNumberInput(p.X+15, p.Y+70, 270, 35, 3.0)
	p.inputs=append(p.inputs, p.IntervalInput)

	yOffset += 100
	minSpeedLabel := NewLabel(p.X+15, p.Y+150, "Min Speed:")
	minSpeedLabel.Size = 14
	p.labels = append(p.labels, minSpeedLabel)
	
	p.MinSpeedInput = NewNumberInput(p.X+15, p.Y+170, 270, 35, 20)
	p.inputs=append(p.inputs, p.MinSpeedInput)
	
	yOffset += 100

	maxSpeedLabel := NewLabel(p.X+15, p.Y+250, "Max Speed:")
	maxSpeedLabel.Size = 14
	p.labels = append(p.labels, maxSpeedLabel)
	
	p.MaxSpeedInput = NewNumberInput(p.X+15, p.Y+270, 270, 35, 50)
	p.inputs=append(p.inputs, p.MaxSpeedInput)

	yOffset += 100

	maxVehiclesLabel := NewLabel(p.X+15, p.Y+350, "Max vehicles:")
	maxVehiclesLabel.Size = 14
	p.labels = append(p.labels, maxVehiclesLabel)
	
	p.MaxVehiclesInput = NewNumberInput(p.X+15, p.Y+350+20, 270, 35, 50)
	p.inputs=append(p.inputs, p.MaxVehiclesInput)

	yOffset += 100

	enabledLabel := NewLabel(p.X+15, p.Y+550.0, "Enabled:")
	enabledLabel.Size = 14
	p.labels = append(p.labels, enabledLabel)
	
	p.EnabledInput = NewTextInput(p.X+15, p.Y+250+20, 270, 35, "true")
	
	yOffset += 50

	p.applyBtn = NewButton(p.X+140, p.Y+500, p.btnWidth, p.btnHeight, "Apply ", nil)
	p.closeBtn = NewButton(p.X+225,p.Y+500,p.btnWidth, p.btnHeight, "Close ", func() {
	})
	
	p.calculateHeight()
}

func (p *SpawnerPropertiesPanel) Show(Interval,MinSpeed,MaxSpeed float64, MaxVehicles int,Enabled bool) {
	p.Visible = true
	p.IntervalInput.SetNumber( Interval)
	p.MinSpeedInput.SetNumber( MinSpeed)
	p.MaxSpeedInput.SetNumber( MaxSpeed)
	p.MaxVehiclesInput.SetNumber( float64(MaxVehicles))
	p.EnabledInput.SetText(fmt.Sprintf("%t", Enabled))
}

func (p *SpawnerPropertiesPanel) Hide() {
	p.Visible = false
}

func (p *SpawnerPropertiesPanel) SetOnApply(callback func(Interval,MinSpeed,MaxSpeed float64, MaxVehicles int,Enabled bool)) {
	p.onApply = callback
}

func (p *SpawnerPropertiesPanel) SetPosition(x, y float64) {
	p.X = x
	p.Y = y
	p.updateUIPositions()
}

func (p *SpawnerPropertiesPanel) updateUIPositions() {
	yOffset := 15.0
	
	for i, label := range p.labels {
		label.X = p.X + 15
		label.Y = p.Y + yOffset
		if i == 0 {
			yOffset += 35
		} else {
			yOffset += 70
		}
	}
	
	yOffset = 70.0
	for _, input := range p.inputs {
		input.X = p.X + 15
		input.Y = p.Y + yOffset
		input.incrementValueBtn.X = input.X + input.Width - 30
		input.incrementValueBtn.Y = input.Y + 5
		input.decrementValueBtn.X = input.X + input.Width - 60
		input.decrementValueBtn.Y = input.Y + 5
		yOffset += 70
	}
	p.EnabledInput.X = p.X + 15
	p.EnabledInput.Y = p.Y + yOffset

	p.applyBtn.X = p.X + 140
	p.applyBtn.Y = p.Y + 480
	
	p.closeBtn.X = p.X + 225
	p.closeBtn.Y = p.Y + 480
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
		Interval:=p.IntervalInput.GetNumber()
		MinSpeed:=p.MinSpeedInput.GetNumber()
		MaxSpeed:=p.MaxSpeedInput.GetNumber()
		MaxVehicles:= int(p.MaxVehiclesInput.GetNumber())
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
	for _, input := range p.inputs {
		p.Height += input.Height + 7
	}
	p.Height += p.EnabledInput.Height + 7
	p.Height += p.applyBtn.Height 
	fmt.Println("Calculated panel height:", p.Height)
}