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
)

type InputHandler struct {
	mode           Mode
	roadTool       *tools.RoadBuildingTool
	moveTool       *tools.NodeMoveTool
	spawnTool      *tools.SpawnTool
	mouseX, mouseY int
}

func NewInputHandler(w *world.World) *InputHandler {
	executor := commands.NewCommandExecutor(w)
	query := query.NewWorldQuery(w)
	roadTool := tools.NewRoadBuildingTool(executor, query)
	moveTool := tools.NewNodeMoveTool(executor, query)
	spawnTool := tools.NewSpawnTool(executor, query)
	
	return &InputHandler{
		mode:      ModeNormal,
		roadTool:  roadTool,
		moveTool:  moveTool,
		spawnTool: spawnTool,
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
	
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		h.mode = ModeNormal
		h.roadTool.Cancel()
		h.moveTool.EndDrag()
		h.spawnTool.Cancel()
	}
	
	if inpututil.IsKeyJustPressed(ebiten.KeyB) {
		h.roadTool.SetBidirectional(!h.roadTool.IsBidirectional())
	}

	if h.mode == ModeSpawning && inpututil.IsKeyJustPressed(ebiten.KeyTab) {
		h.spawnTool.CycleRoad()
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