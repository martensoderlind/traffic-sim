package input

import (
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"traffic-sim/internal/commands"
	"traffic-sim/internal/query"
	"traffic-sim/internal/road"
	"traffic-sim/internal/sim"
	"traffic-sim/internal/tools"
	"traffic-sim/internal/world"
)

type Mode int

const (
	ModeNormal Mode = iota
	ModeRoadBuilding
	ModeNodeMoving
	ModeSpawning
	ModeDespawning
	ModeRoadDeleting
	ModeNodeDeleting
	ModeTrafficLight
	ModeRoadProperties
	ModeSpawnPointProperties
)

type InputHandler struct {
	mode             Mode
	roadTool         *tools.RoadBuildingTool
	moveTool         *tools.NodeMoveTool
	spawnTool        *tools.SpawnTool
	despawnTool      *tools.DespawnTool
	roadDeleteTool   *tools.RoadDeleteTool
	nodeDeleteTool   *tools.NodeDeleteTool
	trafficLightTool *tools.TrafficLightTool
	roadPropTool     *tools.RoadPropertiesTool
	spawnPointPropTool *tools.SpawnPointPropertiesTool
	mouseX, mouseY   int
	Simulator 	     *sim.Simulator
	roadPropertiesPanel interface{ Contains(x, y int) bool } 
	spawnPointPropertiesPanel interface{ Contains(x, y int) bool }
	world            *world.World
	executor         *commands.CommandExecutor
	onWorldReplaced  func(*world.World)
}

func NewInputHandler(w *world.World, s *sim.Simulator) *InputHandler {
	executor := commands.NewCommandExecutor(w)
	query := query.NewWorldQuery(w)
	roadTool := tools.NewRoadBuildingTool(executor, query)
	moveTool := tools.NewNodeMoveTool(executor, query)
	spawnTool := tools.NewSpawnTool(executor, query)
	despawnTool := tools.NewDespawnTool(executor, query)
	roadDeleteTool := tools.NewRoadDeleteTool(executor, query)
	nodeDeleteTool := tools.NewNodeDeleteTool(executor, query)
	trafficLightTool := tools.NewTrafficLightTool(executor, query)
	roadPropTool := tools.NewRoadPropertiesTool(executor, query)
	SpawnPointPropTool := tools.NewSpawnPointPropertiesTool(executor, query)
	simulator := s
	
	return &InputHandler{
		mode:             ModeNormal,
		roadTool:         roadTool,
		moveTool:         moveTool,
		spawnTool:        spawnTool,
		despawnTool:      despawnTool,
		roadDeleteTool:   roadDeleteTool,
		nodeDeleteTool:   nodeDeleteTool,
		trafficLightTool: trafficLightTool,
		roadPropTool:     roadPropTool,
		spawnPointPropTool: SpawnPointPropTool,
		Simulator:        simulator,
		world:            w,
		executor:         executor,
	}
}

func (h *InputHandler) SetOnWorldReplaced(callback func(*world.World)) {
	h.onWorldReplaced = callback
}

func (h *InputHandler) Mode() Mode {
	return h.mode
}

func (h *InputHandler) SetMode(mode Mode) {
	if h.mode == mode {
		return
	}
	
	h.roadTool.Cancel()
	h.moveTool.EndDrag()
	h.spawnTool.Cancel()
	h.despawnTool.Cancel()
	h.roadDeleteTool.Cancel()
	h.nodeDeleteTool.Cancel()
	h.trafficLightTool.Cancel()
	h.roadPropTool.Cancel()
	h.spawnPointPropTool.Cancel()
	h.mode = mode
}

func (h *InputHandler) ToggleBidirectional() {
	h.roadTool.SetBidirectional(!h.roadTool.IsBidirectional())
}

func (h *InputHandler) RoadTool() *tools.RoadBuildingTool {
	return h.roadTool
}

func (h *InputHandler) MoveTool() *tools.NodeMoveTool {
	return h.moveTool
}

func (h *InputHandler) SpawnTool() *tools.SpawnTool {
	return h.spawnTool
}

func (h *InputHandler) DespawnTool() *tools.DespawnTool {
	return h.despawnTool
}

func (h *InputHandler) RoadDeleteTool() *tools.RoadDeleteTool {
	return h.roadDeleteTool
}

func (h *InputHandler) NodeDeleteTool() *tools.NodeDeleteTool {
	return h.nodeDeleteTool
}

func (h *InputHandler) TrafficLightTool() *tools.TrafficLightTool {
	return h.trafficLightTool
}

func (h *InputHandler) MousePos() (int, int) {
	return h.mouseX, h.mouseY
}

func (h *InputHandler) RoadPropTool() *tools.RoadPropertiesTool {
    return h.roadPropTool
}

func (h *InputHandler) SpawnPointPropTool() *tools.SpawnPointPropertiesTool {
    return h.spawnPointPropTool
}

func (h *InputHandler) Update() {
	h.mouseX, h.mouseY = ebiten.CursorPosition()
	
	h.handleModeSwitch()
	h.handleToolInput()
	h.handleSaveLoad()
}

func (h *InputHandler) handleSaveLoad() {
	if ebiten.IsKeyPressed(ebiten.KeyControl) || ebiten.IsKeyPressed(ebiten.KeyMeta) {
		if inpututil.IsKeyJustPressed(ebiten.KeyS) {
			h.handleSave()
		}
		
		if inpututil.IsKeyJustPressed(ebiten.KeyO) {
			h.handleLoad()
		}
	}
}

func (h *InputHandler) handleSave() {
	cmd := &commands.SaveWorldCommand{}
	if err := h.executor.Execute(cmd); err != nil {
		log.Printf("Failed to save world: %v", err)
	}
}

func (h *InputHandler) handleLoad() {
	cmd := &commands.LoadWorldCommand{
		OnWorldLoaded: func(newWorld *world.World) {
			if h.onWorldReplaced != nil {
				h.onWorldReplaced(newWorld)
			}
		},
	}
	
	if err := cmd.Execute(h.world); err != nil {
		log.Printf("Failed to load world: %v", err)
	}
}

func (h *InputHandler) ReplaceWorld(newWorld *world.World) {
	h.world = newWorld
	h.executor = commands.NewCommandExecutor(newWorld)
	query := query.NewWorldQuery(newWorld)
	
	h.roadTool = tools.NewRoadBuildingTool(h.executor, query)
	h.moveTool = tools.NewNodeMoveTool(h.executor, query)
	h.spawnTool = tools.NewSpawnTool(h.executor, query)
	h.despawnTool = tools.NewDespawnTool(h.executor, query)
	h.roadDeleteTool = tools.NewRoadDeleteTool(h.executor, query)
	h.nodeDeleteTool = tools.NewNodeDeleteTool(h.executor, query)
	h.trafficLightTool = tools.NewTrafficLightTool(h.executor, query)
	h.roadPropTool = tools.NewRoadPropertiesTool(h.executor, query)
	h.spawnPointPropTool = tools.NewSpawnPointPropertiesTool(h.executor, query)
	
	h.SetMode(ModeNormal)
}

func (h *InputHandler) handleModeSwitch() {
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		if h.mode == ModeNormal {
			h.mode = ModeRoadBuilding
		} else {
			h.mode = ModeNormal
			h.roadTool.Cancel()
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyM) {
		if h.mode == ModeNormal {
			h.mode = ModeNodeMoving
		} else {
			h.mode = ModeNormal
			h.moveTool.EndDrag()
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyS) && !ebiten.IsKeyPressed(ebiten.KeyControl) && !ebiten.IsKeyPressed(ebiten.KeyMeta) {
		if h.mode == ModeNormal {
			h.mode = ModeSpawning
		} else {
			h.mode = ModeNormal
			h.spawnTool.Cancel()
		}
	}
	
	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		if h.mode == ModeNormal {
			h.mode = ModeSpawnPointProperties
		} else {
			h.mode = ModeNormal
			h.spawnPointPropTool.Cancel()
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		if h.mode == ModeNormal {
			h.mode = ModeDespawning
		} else {
			h.mode = ModeNormal
			h.despawnTool.Cancel()
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyX) {
		if h.mode == ModeNormal {
			h.mode = ModeRoadDeleting
		} else {
			h.mode = ModeNormal
			h.roadDeleteTool.Cancel()
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyDelete) {
		if h.mode == ModeNormal {
			h.mode = ModeNodeDeleting
		} else {
			h.mode = ModeNormal
			h.nodeDeleteTool.Cancel()
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyT) {
		if h.mode == ModeNormal {
			h.mode = ModeTrafficLight
		} else {
			h.mode = ModeNormal
			h.trafficLightTool.Cancel()
		}
	}
	
	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		if h.mode == ModeNormal {
			h.mode = ModeRoadProperties
		} else {
			h.mode = ModeNormal
			h.roadPropTool.Cancel()
		}
	}
	
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		h.mode = ModeNormal
		h.roadTool.Cancel()
		h.moveTool.EndDrag()
		h.spawnTool.Cancel()
		h.despawnTool.Cancel()
		h.roadDeleteTool.Cancel()
		h.nodeDeleteTool.Cancel()
		h.trafficLightTool.Cancel()
		h.roadPropTool.Cancel()
		h.spawnPointPropTool.Cancel()
	}
	
	if inpututil.IsKeyJustPressed(ebiten.KeyB) {
		h.roadTool.SetBidirectional(!h.roadTool.IsBidirectional())
	}

	if h.mode == ModeSpawning && inpututil.IsKeyJustPressed(ebiten.KeyTab) {
		h.spawnTool.CycleRoad()
	}

	if h.mode == ModeDespawning && inpututil.IsKeyJustPressed(ebiten.KeyTab) {
		h.despawnTool.CycleRoad()
	}

	if h.mode == ModeTrafficLight && inpututil.IsKeyJustPressed(ebiten.KeyTab) {
		h.trafficLightTool.CycleRoad()
	}

	if h.mode == ModeTrafficLight && inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		h.trafficLightTool.Click(float64(h.mouseX), float64(h.mouseY))
	}
	
	if h.mode == ModeNormal && inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		h.Simulator.TogglePause()
	}
}

func (h *InputHandler) handleToolInput() {
	switch h.mode {
	case ModeRoadBuilding:
		h.handleRoadBuildingInput()
	case ModeNodeMoving:
		h.handleNodeMovingInput()
	case ModeSpawning:
		h.handleSpawningInput()
	case ModeDespawning:
		h.handleDespawningInput()
	case ModeRoadDeleting:
		h.handleRoadDeletingInput()
	case ModeNodeDeleting:
		h.handleNodeDeletingInput()
	case ModeTrafficLight:
		h.handleTrafficLightInput()
	case ModeRoadProperties:
		h.handleRoadPropertiesInput()
	case ModeSpawnPointProperties:
		h.handleSpawnPointPropertiesInput()
	}
}

func (h *InputHandler) handleRoadBuildingInput() {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		h.roadTool.Click(float64(h.mouseX), float64(h.mouseY))
	}
	
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		h.roadTool.Cancel()
	}
}

func (h *InputHandler) handleNodeMovingInput() {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		h.moveTool.StartDrag(float64(h.mouseX), float64(h.mouseY))
	}

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && h.moveTool.IsDragging() {
		h.moveTool.UpdateDrag(float64(h.mouseX), float64(h.mouseY))
	}

	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		h.moveTool.EndDrag()
	}
}

func (h *InputHandler) handleSpawningInput() {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		h.spawnTool.Click(float64(h.mouseX), float64(h.mouseY))
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		h.spawnTool.Cancel()
	}
}

func (h *InputHandler) handleDespawningInput() {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		h.despawnTool.Click(float64(h.mouseX), float64(h.mouseY))
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		h.despawnTool.Cancel()
	}
}

func (h *InputHandler) handleRoadDeletingInput() {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		h.roadDeleteTool.Click(float64(h.mouseX), float64(h.mouseY))
	}
}

func (h *InputHandler) handleNodeDeletingInput() {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		h.nodeDeleteTool.Click(float64(h.mouseX), float64(h.mouseY))
	}
}

func (h *InputHandler) handleRoadPropertiesInput() {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if h.roadPropertiesPanel != nil && h.roadPropertiesPanel.Contains(h.mouseX, h.mouseY) {
			return 
		}
		h.roadPropTool.Click(float64(h.mouseX), float64(h.mouseY))
	}
	
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		h.roadPropTool.Cancel()
	}
}

func (h *InputHandler) handleSpawnPointPropertiesInput() {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if h.spawnPointPropertiesPanel != nil && h.spawnPointPropertiesPanel.Contains(h.mouseX, h.mouseY) {
			return 
		}
		h.spawnPointPropTool.Click(float64(h.mouseX), float64(h.mouseY))
	}
	
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		h.spawnPointPropTool.Cancel()
	}
}

func (h *InputHandler) SetRoadPropertiesPanel(panel interface{ Contains(x, y int) bool }) {
	h.roadPropertiesPanel = panel
}

func (h *InputHandler) SetSpawnPointPropertiesPanel(panel interface{ Contains(x, y int) bool }) {
	h.spawnPointPropertiesPanel = panel
}

func (h *InputHandler) handleTrafficLightInput() {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		mouseX := float64(h.mouseX)
		mouseY := float64(h.mouseY)
		
		node := h.trafficLightTool.GetSelectedNode()
		if node != nil {
			for _, road := range h.trafficLightTool.GetAvailableRoads() {
				if h.isMouseNearRoad(mouseX, mouseY, road) {
					h.trafficLightTool.ToggleRoad(road)
					return
				}
			}
		}
		
		h.trafficLightTool.Click(mouseX, mouseY)
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		h.trafficLightTool.Cancel()
	}
}

func (h *InputHandler) isMouseNearRoad(mouseX, mouseY float64, rd *road.Road) bool {
	x1, y1 := rd.From.X, rd.From.Y
	x2, y2 := rd.To.X, rd.To.Y
	
	dx := x2 - x1
	dy := y2 - y1
	
	if dx == 0 && dy == 0 {
		dist := math.Sqrt((mouseX-x1)*(mouseX-x1) + (mouseY-y1)*(mouseY-y1))
		return dist < 15.0
	}
	
	t := ((mouseX-x1)*dx + (mouseY-y1)*dy) / (dx*dx + dy*dy)
	
	if t < 0 {
		t = 0
	} else if t > 1 {
		t = 1
	}
	
	px := x1 + t*dx
	py := y1 + t*dy
	
	dist := math.Sqrt((mouseX-px)*(mouseX-px) + (mouseY-py)*(mouseY-py))
	
	return dist < 15.0
}