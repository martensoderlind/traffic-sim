package input

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"traffic-sim/internal/commands"
	"traffic-sim/internal/query"
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
)

type InputHandler struct {
	mode           Mode
	roadTool       *tools.RoadBuildingTool
	moveTool       *tools.NodeMoveTool
	spawnTool      *tools.SpawnTool
	despawnTool    *tools.DespawnTool
	roadDeleteTool *tools.RoadDeleteTool
	nodeDeleteTool *tools.NodeDeleteTool
	mouseX, mouseY int
}

func NewInputHandler(w *world.World) *InputHandler {
	executor := commands.NewCommandExecutor(w)
	query := query.NewWorldQuery(w)
	roadTool := tools.NewRoadBuildingTool(executor, query)
	moveTool := tools.NewNodeMoveTool(executor, query)
	spawnTool := tools.NewSpawnTool(executor, query)
	despawnTool := tools.NewDespawnTool(executor, query)
	roadDeleteTool := tools.NewRoadDeleteTool(executor, query)
	nodeDeleteTool := tools.NewNodeDeleteTool(executor, query)
	
	return &InputHandler{
		mode:           ModeNormal,
		roadTool:       roadTool,
		moveTool:       moveTool,
		spawnTool:      spawnTool,
		despawnTool:    despawnTool,
		roadDeleteTool: roadDeleteTool,
		nodeDeleteTool: nodeDeleteTool,
	}
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

func (h *InputHandler) MousePos() (int, int) {
	return h.mouseX, h.mouseY
}

func (h *InputHandler) Update() {
	h.mouseX, h.mouseY = ebiten.CursorPosition()
	
	h.handleModeSwitch()
	h.handleToolInput()
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

	if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		if h.mode == ModeNormal {
			h.mode = ModeSpawning
		} else {
			h.mode = ModeNormal
			h.spawnTool.Cancel()
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
	
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		h.mode = ModeNormal
		h.roadTool.Cancel()
		h.moveTool.EndDrag()
		h.spawnTool.Cancel()
		h.despawnTool.Cancel()
		h.roadDeleteTool.Cancel()
		h.nodeDeleteTool.Cancel()
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