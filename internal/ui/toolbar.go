package ui

import (
	"fmt"
	"image/color"
	"traffic-sim/internal/input"
	"traffic-sim/internal/world"

	"github.com/hajimehoshi/ebiten/v2"
)

type Toolbar struct {
	uiManager     *UIManager
	inputHandler  *input.InputHandler
	modeIndicator *Label
	simulationState *Label
	statsPanel    *StatsPanel
	
	roadBuildBtn    *Button
	moveNodeBtn     *Button
	normalModeBtn   *Button
	spawnBtn        *Button
	despawnBtn      *Button
	roadDeleteBtn   *Button
	nodeDeleteBtn   *Button
	trafficLightBtn *Button
	bidirToggle     *Button
	roadPropBtn *Button
	spawnPointPropBtn *Button
	roadCurveBtn *Button
	saveBtn         *Button
	loadBtn         *Button
    roadPropertiesPanel *RoadPropertiesPanel
    spawnPointPropertiesPanel *SpawnerPropertiesPanel

	world *world.World
}

func NewToolbar(inputHandler *input.InputHandler,w *world.World) *Toolbar {
	tb := &Toolbar{
		uiManager:    NewUIManager(),
		inputHandler: inputHandler,
		world: 	  w,
	}
	
	tb.setupUI()
	return tb
}

func (tb *Toolbar) setupUI() {
	btnY := 10.0
	btnWidth := 120.0
	btnHeight := 35.0
	spacingX := 10.0
	spacingY:= 25.0
	currentX := 15.0
	
	tb.normalModeBtn = NewButton(currentX, btnY, btnWidth, btnHeight, "Normal (ESC)", func() {
		tb.inputHandler.SetMode(input.ModeNormal)
	})
	tb.uiManager.AddButton(tb.normalModeBtn)
	currentX += float64(tb.normalModeBtn.calculateWidth()) + spacingX
	
	tb.roadBuildBtn = NewButton(currentX, btnY, btnWidth, btnHeight, "Build Road (R)", func() {
		tb.inputHandler.SetMode(input.ModeRoadBuilding)
	})
	tb.uiManager.AddButton(tb.roadBuildBtn)
	currentX += float64(tb.roadBuildBtn.calculateWidth()) + spacingX
	
	tb.moveNodeBtn = NewButton(currentX, btnY, btnWidth, btnHeight, "Move Node (M)", func() {
		tb.inputHandler.SetMode(input.ModeNodeMoving)
	})
	tb.uiManager.AddButton(tb.moveNodeBtn)
	currentX += float64(tb.moveNodeBtn.calculateWidth()) + spacingX
	
	tb.spawnBtn = NewButton(currentX, btnY, btnWidth, btnHeight, "Add Spawn (S)", func() {
		tb.inputHandler.SetMode(input.ModeSpawning)
	})
	tb.uiManager.AddButton(tb.spawnBtn)
	currentX += float64(tb.spawnBtn.calculateWidth()) + spacingX
	
	tb.despawnBtn = NewButton(currentX, btnY, btnWidth, btnHeight, "Add Despawn (D)", func() {
		tb.inputHandler.SetMode(input.ModeDespawning)
	})
	tb.uiManager.AddButton(tb.despawnBtn)
	currentX += float64(tb.despawnBtn.calculateWidth()) + spacingX

	tb.roadDeleteBtn = NewButton(currentX, btnY, btnWidth, btnHeight, "Delete Road (X)", func() {
		tb.inputHandler.SetMode(input.ModeRoadDeleting)
	})
	tb.uiManager.AddButton(tb.roadDeleteBtn)
	currentX += float64(tb.roadDeleteBtn.calculateWidth()) + spacingX
	
	tb.nodeDeleteBtn = NewButton(currentX, btnY, btnWidth, btnHeight, "Delete Node (Del)", func() {
		tb.inputHandler.SetMode(input.ModeNodeDeleting)
	})
	tb.uiManager.AddButton(tb.nodeDeleteBtn)
	currentX += float64(tb.nodeDeleteBtn.calculateWidth()) + spacingX

	tb.trafficLightBtn = NewButton(currentX, btnY, btnWidth, btnHeight, "Traffic Light (T)", func() {
		tb.inputHandler.SetMode(input.ModeTrafficLight)
	})
	tb.uiManager.AddButton(tb.trafficLightBtn)
	currentX += float64(tb.trafficLightBtn.calculateWidth()) + spacingX
	
	tb.roadPropBtn = NewButton(currentX, btnY, btnWidth, btnHeight, "Road Props (P)", func() {
		tb.inputHandler.SetMode(input.ModeRoadProperties)
	})
	tb.uiManager.AddButton(tb.roadPropBtn)
	currentX += float64(tb.roadPropBtn.calculateWidth()) + spacingX
	
	tb.spawnPointPropBtn = NewButton(currentX, btnY, btnWidth, btnHeight, "Spawn point Props (Q)", func() {
		tb.inputHandler.SetMode(input.ModeSpawnPointProperties)
	})
	tb.uiManager.AddButton(tb.spawnPointPropBtn)
	currentX += float64(tb.spawnPointPropBtn.calculateWidth()) + spacingX
	
	tb.roadCurveBtn = NewButton(currentX, btnY, btnWidth, btnHeight, "Curve Road (C)", func() {
		tb.inputHandler.SetMode(input.ModeRoadCurving)
	})
	tb.uiManager.AddButton(tb.roadCurveBtn)
	
	currentX = 15.0
	btnY += btnHeight + spacingY
	
	tb.bidirToggle = NewButton(currentX, btnY, btnWidth, btnHeight, "Bidir: ON (B)", func() {
		tb.inputHandler.ToggleBidirectional()
	})
	tb.uiManager.AddButton(tb.bidirToggle)
	currentX += float64(tb.bidirToggle.calculateWidth()) + spacingX
	
	tb.saveBtn = NewButton(currentX, btnY, btnWidth, btnHeight, "Save (Ctrl+S)", nil)
	tb.uiManager.AddButton(tb.saveBtn)
	currentX += float64(tb.saveBtn.calculateWidth()) + spacingX
	
	tb.loadBtn = NewButton(currentX, btnY, btnWidth, btnHeight, "Load (Ctrl+O)", nil)
	tb.uiManager.AddButton(tb.loadBtn)
	btnY += btnHeight + spacingY
	tb.modeIndicator = NewLabel(18, btnY+3, "Mode: Normal")
	tb.modeIndicator.Size = 13
	tb.modeIndicator.Color = color.RGBA{240, 240, 245, 255}
	tb.modeIndicator.SetBackground(color.RGBA{45, 50, 65, 240})
	tb.uiManager.AddLabel(tb.modeIndicator)
	btnY += 27 + spacingY

	tb.simulationState = NewLabel(18, btnY, "Simulation: Running")
	tb.simulationState.Size = 13
	tb.simulationState.Color = color.RGBA{240, 240, 245, 255}
	tb.simulationState.SetBackground(color.RGBA{45, 55, 45, 240})
	tb.uiManager.AddLabel(tb.simulationState)
	
	tb.statsPanel = NewStatsPanel(15, btnY+40, tb.world)

	tb.roadPropertiesPanel = NewRoadPropertiesPanel(1600, 200)
	tb.roadPropertiesPanel.SetOnApply(func(maxSpeed, width float64) {
		if tb.inputHandler.RoadPropTool().GetSelectedRoad() != nil {
			tb.inputHandler.RoadPropTool().UpdateRoadProperties(maxSpeed, width)
			tb.roadPropertiesPanel.Hide()
		}
	})

	tb.spawnPointPropertiesPanel = NewSpawnerPropertiesPanel(1600, 200)
	tb.spawnPointPropertiesPanel.SetOnApply(func(Interval,MinSpeed,MaxSpeed float64, MaxVehicles int,Enabled bool) {
		if tb.inputHandler.SpawnPointPropTool().GetSelectedSpawnPoint() != nil {
			tb.inputHandler.SpawnPointPropTool().UpdateSpawnPointProperties(Interval,MinSpeed,MaxSpeed, MaxVehicles ,Enabled)
			tb.spawnPointPropertiesPanel.Hide()
		}
	})
	
	tb.inputHandler.SetRoadPropertiesPanel(tb.roadPropertiesPanel)
	tb.inputHandler.SetSpawnPointPropertiesPanel(tb.spawnPointPropertiesPanel)
}

func (tb *Toolbar) UpdatePanelPositions(screenWidth, screenHeight int) {
	panelMargin := 20.0
	panelX := float64(screenWidth) - tb.roadPropertiesPanel.Width - panelMargin
	panelY := 200.0
	
	tb.roadPropertiesPanel.SetPosition(panelX, panelY)
	tb.spawnPointPropertiesPanel.SetPosition(panelX, panelY)
}

func (tb *Toolbar) Update(mouseX, mouseY int, clicked bool) {
	tb.uiManager.Update(mouseX, mouseY, clicked)
	tb.updateModeIndicator()
	tb.updateSimulationStatus()
	tb.updateButtonStates()
	tb.statsPanel.Update()

	mode := tb.inputHandler.Mode()
	if mode == input.ModeRoadProperties {
		selectedRoad := tb.inputHandler.RoadPropTool().GetSelectedRoad()
		if selectedRoad != nil && !tb.roadPropertiesPanel.Visible {
			tb.roadPropertiesPanel.Show(selectedRoad.MaxSpeed, selectedRoad.Width)
		} else if selectedRoad == nil {
			tb.roadPropertiesPanel.Hide()
		}
	}else if mode == input.ModeSpawnPointProperties {
		selectedSpawnPoint := tb.inputHandler.SpawnPointPropTool().GetSelectedSpawnPoint()
		if selectedSpawnPoint != nil && !tb.spawnPointPropertiesPanel.Visible {
			tb.spawnPointPropertiesPanel.Show(selectedSpawnPoint.Interval,selectedSpawnPoint.MinSpeed,selectedSpawnPoint.MaxSpeed, selectedSpawnPoint.MaxVehicles,selectedSpawnPoint.Enabled)
		} else if selectedSpawnPoint == nil {
			tb.spawnPointPropertiesPanel.Hide()
		}
	} else {
		tb.roadPropertiesPanel.Hide()
	}
	
	tb.roadPropertiesPanel.Update(mouseX, mouseY, clicked)
	tb.spawnPointPropertiesPanel.Update(mouseX, mouseY, clicked)
}

func (tb *Toolbar) updateModeIndicator() {
	mode := tb.inputHandler.Mode()
	
	var modeText string
	var bgColor color.RGBA
	
	switch mode {
	case input.ModeNormal:
		modeText = "Mode: Normal"
		bgColor = color.RGBA{45, 50, 65, 240}
	case input.ModeRoadBuilding:
		modeText = "Mode: Build Road"
		bgColor = color.RGBA{65, 90, 50, 240}
		if tb.inputHandler.RoadTool().GetSelectedNode() != nil {
			modeText = "Mode: Build Road (Node Selected)"
		}
	case input.ModeNodeMoving:
		modeText = "Mode: Move Node"
		bgColor = color.RGBA{90, 60, 90, 240}
		if tb.inputHandler.MoveTool().IsDragging() {
			modeText = "Mode: Move Node (Dragging)"
		}
	case input.ModeSpawning:
		modeText = "Mode: Add Spawn Point"
		bgColor = color.RGBA{50, 90, 75, 240}
		if tb.inputHandler.SpawnTool().GetSelectedNode() != nil {
			if tb.inputHandler.SpawnTool().GetSelectedRoad() != nil {
				modeText = "Mode: Add Spawn (Road Selected - Click to Confirm)"
			} else {
				modeText = "Mode: Add Spawn (Node Selected - Tab to Cycle)"
			}
		}
	case input.ModeDespawning:
		modeText = "Mode: Add Despawn Point"
		bgColor = color.RGBA{90, 55, 55, 240}
		if tb.inputHandler.DespawnTool().GetSelectedNode() != nil {
			if tb.inputHandler.DespawnTool().GetSelectedRoad() != nil {
				modeText = "Mode: Add Despawn (Road Selected - Click to Confirm)"
			} else {
				modeText = "Mode: Add Despawn (Node Selected - Tab to Cycle)"
			}
		}
	case input.ModeRoadDeleting:
		modeText = "Mode: Delete Road (Click on road)"
		bgColor = color.RGBA{95, 45, 45, 240}
	case input.ModeNodeDeleting:
		modeText = "Mode: Delete Node (Click on node - deletes all connected roads)"
		bgColor = color.RGBA{95, 45, 45, 240}
	case input.ModeTrafficLight:
		modeText = "Mode: Traffic Light - Click node, then click roads (Space to confirm)"
		bgColor = color.RGBA{95, 95, 40, 240}
		if tb.inputHandler.TrafficLightTool().GetSelectedNode() != nil {
			selectedRoads := tb.inputHandler.TrafficLightTool().GetSelectedRoads()
			if len(selectedRoads) > 0 {
				modeText = fmt.Sprintf("Mode: Traffic Light (%d roads selected - Space to confirm)", len(selectedRoads))
			} else {
				modeText = "Mode: Traffic Light (Node Selected - Click roads to control)"
			}
		}
	case input.ModeRoadProperties:
		modeText = "Mode: Edit Road Properties (Click road)"
		bgColor = color.RGBA{50, 90, 95, 240}
		if tb.inputHandler.RoadPropTool().GetSelectedRoad() != nil {
			modeText = "Mode: Edit Road Properties (Selected - Edit in panel)"
		}
	case input.ModeSpawnPointProperties:
		modeText = "Mode: Edit Spawn Point Properties (Click Spawn Point)"
		bgColor = color.RGBA{50, 90, 95, 240}
		if tb.inputHandler.SpawnPointPropTool().GetSelectedSpawnPoint() != nil {
			modeText = "Mode: Edit Spawn point Properties (Selected - Edit in panel)"
		}
	case input.ModeRoadCurving:
		modeText = tb.inputHandler.RoadCurveTool().GetStatusMessage()
		bgColor = color.RGBA{75, 60, 90, 240}
	}
	
	tb.modeIndicator.Text = modeText
	tb.modeIndicator.SetBackground(bgColor)
}

func (tb *Toolbar) updateSimulationStatus() {
	mode := tb.inputHandler.Simulator.IsPaused()
	
	var modeText string
	var bgColor color.RGBA
	
	switch mode {
	case true:
		modeText = "Simulation: Paused"
		bgColor = color.RGBA{80, 50, 50, 240}
	case false:
		modeText = "Simulation: Running"
		bgColor = color.RGBA{45, 55, 45, 240}
	}
	
	tb.simulationState.Text = modeText
	tb.simulationState.SetBackground(bgColor)
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
	
	if mode == input.ModeSpawning {
		tb.spawnBtn.SetColors(activeColor, activeHover, activePress, textColor, borderColor)
	} else {
		tb.spawnBtn.SetColors(normalColor, normalHover, normalPress, textColor, borderColor)
	}
	
	if mode == input.ModeDespawning {
		tb.despawnBtn.SetColors(activeColor, activeHover, activePress, textColor, borderColor)
	} else {
		tb.despawnBtn.SetColors(normalColor, normalHover, normalPress, textColor, borderColor)
	}
	
	if mode == input.ModeRoadDeleting {
		tb.roadDeleteBtn.SetColors(activeColor, activeHover, activePress, textColor, borderColor)
	} else {
		tb.roadDeleteBtn.SetColors(normalColor, normalHover, normalPress, textColor, borderColor)
	}
	
	if mode == input.ModeNodeDeleting {
		tb.nodeDeleteBtn.SetColors(activeColor, activeHover, activePress, textColor, borderColor)
	} else {
		tb.nodeDeleteBtn.SetColors(normalColor, normalHover, normalPress, textColor, borderColor)
	}

	if mode == input.ModeTrafficLight {
		tb.trafficLightBtn.SetColors(activeColor, activeHover, activePress, textColor, borderColor)
	} else {
		tb.trafficLightBtn.SetColors(normalColor, normalHover, normalPress, textColor, borderColor)
	}
	
	if mode == input.ModeRoadProperties {
		tb.roadPropBtn.SetColors(activeColor, activeHover, activePress, textColor, borderColor)
	} else {
		tb.roadPropBtn.SetColors(normalColor, normalHover, normalPress, textColor, borderColor)
	}
	
	if mode == input.ModeSpawnPointProperties {
		tb.spawnPointPropBtn.SetColors(activeColor, activeHover, activePress, textColor, borderColor)
	} else {
		tb.spawnPointPropBtn.SetColors(normalColor, normalHover, normalPress, textColor, borderColor)
	}
	
	if mode == input.ModeRoadCurving {
		tb.roadCurveBtn.SetColors(activeColor, activeHover, activePress, textColor, borderColor)
	} else {
		tb.roadCurveBtn.SetColors(normalColor, normalHover, normalPress, textColor, borderColor)
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
	tb.statsPanel.Draw(screen)
	tb.roadPropertiesPanel.Draw(screen)
	tb.spawnPointPropertiesPanel.Draw(screen)
}

func (tb *Toolbar) GetUIManager() *UIManager {
	return tb.uiManager
}

func (tb *Toolbar) ReplaceWorld(newWorld *world.World) {
	tb.world = newWorld
	tb.statsPanel.ReplaceWorld(newWorld)
}

func (tb *Toolbar) Cleanup() {
	tb.statsPanel.Cleanup()
}