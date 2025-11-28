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
)

type InputHandler struct {
	mode           Mode
	roadTool       *tools.RoadBuildingTool
	mouseX, mouseY int
}

func NewInputHandler(w *world.World) *InputHandler {
	executor := commands.NewCommandExecutor(w)
	query := query.NewWorldQuery(w)
	roadTool := tools.NewRoadBuildingTool(executor, query)
	
	return &InputHandler{
		mode:     ModeNormal,
		roadTool: roadTool,
	}
}

func (h *InputHandler) Mode() Mode {
	return h.mode
}

func (h *InputHandler) RoadTool() *tools.RoadBuildingTool {
	return h.roadTool
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
	
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		h.mode = ModeNormal
		h.roadTool.Cancel()
	}
	
	if inpututil.IsKeyJustPressed(ebiten.KeyB) {
		h.roadTool.SetBidirectional(!h.roadTool.IsBidirectional())
	}
}

func (h *InputHandler) handleToolInput() {
	if h.mode != ModeRoadBuilding {
		return
	}
	
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		h.roadTool.Click(float64(h.mouseX), float64(h.mouseY))
	}
	
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		h.roadTool.Cancel()
	}
}