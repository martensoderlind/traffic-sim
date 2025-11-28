package ui

import (
	"image/color"
	"traffic-sim/internal/input"

	"github.com/hajimehoshi/ebiten/v2"
)

type Toolbar struct {
	uiManager     *UIManager
	inputHandler  *input.InputHandler
	modeIndicator *Label
	
	roadBuildBtn  *Button
	moveNodeBtn   *Button
	normalModeBtn *Button
	bidirToggle   *Button
}

func NewToolbar(inputHandler *input.InputHandler) *Toolbar {
	tb := &Toolbar{
		uiManager:    NewUIManager(),
		inputHandler: inputHandler,
	}
	
	tb.setupUI()
	return tb
}

func (tb *Toolbar) setupUI() {
	btnY := 10.0
	btnWidth := 120.0
	btnHeight := 35.0
	spacing := 10.0
	currentX := 10.0
	
	tb.normalModeBtn = NewButton(currentX, btnY, btnWidth, btnHeight, "Normal (ESC)", func() {
		tb.inputHandler.SetMode(input.ModeNormal)
	})
	tb.uiManager.AddButton(tb.normalModeBtn)
	currentX += btnWidth + spacing
	
	tb.roadBuildBtn = NewButton(currentX, btnY, btnWidth, btnHeight, "Build Road (R)", func() {
		tb.inputHandler.SetMode(input.ModeRoadBuilding)
	})
	tb.uiManager.AddButton(tb.roadBuildBtn)
	currentX += btnWidth + spacing
	
	tb.moveNodeBtn = NewButton(currentX, btnY, btnWidth, btnHeight, "Move Node (M)", func() {
		tb.inputHandler.SetMode(input.ModeNodeMoving)
	})
	tb.uiManager.AddButton(tb.moveNodeBtn)
	currentX += btnWidth + spacing
	
	tb.bidirToggle = NewButton(currentX, btnY, btnWidth, btnHeight, "Bidir: ON (B)", func() {
		tb.inputHandler.ToggleBidirectional()
	})
	tb.uiManager.AddButton(tb.bidirToggle)
	
	tb.modeIndicator = NewLabel(10, 60, "Mode: Normal")
	tb.modeIndicator.Size = 14
	bgColor := color.RGBA{40, 40, 50, 230}
	tb.modeIndicator.SetBackground(bgColor)
	tb.uiManager.AddLabel(tb.modeIndicator)
}

func (tb *Toolbar) Update(mouseX, mouseY int, clicked bool) {
	tb.uiManager.Update(mouseX, mouseY, clicked)
	tb.updateModeIndicator()
	tb.updateButtonStates()
}

func (tb *Toolbar) updateModeIndicator() {
	mode := tb.inputHandler.Mode()
	
	var modeText string
	var bgColor color.RGBA
	
	switch mode {
	case input.ModeNormal:
		modeText = "Mode: Normal"
		bgColor = color.RGBA{40, 40, 50, 230}
	case input.ModeRoadBuilding:
		modeText = "Mode: Build Road"
		bgColor = color.RGBA{60, 80, 40, 230}
		if tb.inputHandler.RoadTool().GetSelectedNode() != nil {
			modeText = "Mode: Build Road (Node Selected)"
		}
	case input.ModeNodeMoving:
		modeText = "Mode: Move Node"
		bgColor = color.RGBA{80, 40, 80, 230}
		if tb.inputHandler.MoveTool().IsDragging() {
			modeText = "Mode: Move Node (Dragging)"
		}
	}
	
	tb.modeIndicator.Text = modeText
	tb.modeIndicator.SetBackground(bgColor)
}

func (tb *Toolbar) updateButtonStates() {
	mode := tb.inputHandler.Mode()
	
	activeColor := color.RGBA{80, 120, 80, 255}
	activeHover := color.RGBA{100, 140, 100, 255}
	activePress := color.RGBA{60, 100, 60, 255}
	
	normalColor := color.RGBA{60, 60, 70, 255}
	normalHover := color.RGBA{80, 80, 90, 255}
	normalPress := color.RGBA{50, 50, 60, 255}
	
	textColor := color.RGBA{220, 220, 220, 255}
	borderColor := color.RGBA{100, 100, 110, 255}
	
	if mode == input.ModeNormal {
		tb.normalModeBtn.SetColors(activeColor, activeHover, activePress, textColor, borderColor)
	} else {
		tb.normalModeBtn.SetColors(normalColor, normalHover, normalPress, textColor, borderColor)
	}
	
	if mode == input.ModeRoadBuilding {
		tb.roadBuildBtn.SetColors(activeColor, activeHover, activePress, textColor, borderColor)
	} else {
		tb.roadBuildBtn.SetColors(normalColor, normalHover, normalPress, textColor, borderColor)
	}
	
	if mode == input.ModeNodeMoving {
		tb.moveNodeBtn.SetColors(activeColor, activeHover, activePress, textColor, borderColor)
	} else {
		tb.moveNodeBtn.SetColors(normalColor, normalHover, normalPress, textColor, borderColor)
	}
	
	if tb.inputHandler.RoadTool().IsBidirectional() {
		tb.bidirToggle.Text = "Bidir: ON (B)"
		tb.bidirToggle.SetColors(activeColor, activeHover, activePress, textColor, borderColor)
	} else {
		tb.bidirToggle.Text = "Bidir: OFF (B)"
		tb.bidirToggle.SetColors(normalColor, normalHover, normalPress, textColor, borderColor)
	}
}

func (tb *Toolbar) Draw(screen *ebiten.Image) {
	tb.uiManager.Draw(screen)
}

func (tb *Toolbar) GetUIManager() *UIManager {
	return tb.uiManager
}