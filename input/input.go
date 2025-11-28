package input

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"traffic-sim/internal/road"
	"traffic-sim/internal/sim"
)

type Mode int

const (
	ModeNormal Mode = iota
	ModePlacingRoad
)

type InputHandler struct {
	world *sim.World
	mode  Mode

	// State for road placement
	selectedNode     *road.Node
	hoverNode        *road.Node
	mouseX, mouseY   int
	maxSnapDistance  float64
}

func NewInputHandler(world *sim.World) *InputHandler {
	return &InputHandler{
		world:           world,
		mode:            ModeNormal,
		maxSnapDistance: 20.0, 
	}
}

func (h *InputHandler) Mode() Mode {
	return h.mode
}

func (h *InputHandler) SelectedNode() *road.Node {
	return h.selectedNode
}

func (h *InputHandler) HoverNode() *road.Node {
	return h.hoverNode
}

func (h *InputHandler) MousePos() (int, int) {
	return h.mouseX, h.mouseY
}

func (h *InputHandler) Update() {
	// Update mouse position
	h.mouseX, h.mouseY = ebiten.CursorPosition()

	// Toggle mode with 'R' key (Road mode)
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		if h.mode == ModeNormal {
			h.mode = ModePlacingRoad
			h.selectedNode = nil
		} else {
			h.mode = ModeNormal
			h.selectedNode = nil
		}
	}

	// Cancel current action with Escape
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		h.mode = ModeNormal
		h.selectedNode = nil
	}

	// Update hover node
	h.world.Mu.RLock()
	h.hoverNode = h.findNearestNode(float64(h.mouseX), float64(h.mouseY))
	h.world.Mu.RUnlock()

	// Handle mode-specific input
	switch h.mode {
	case ModePlacingRoad:
		h.updatePlacingRoad()
	}
}

func (h *InputHandler) updatePlacingRoad() {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if h.hoverNode == nil {
			return
		}

		// First click - select start node
		if h.selectedNode == nil {
			h.selectedNode = h.hoverNode
			return
		}

		// Second click - create road
		if h.selectedNode != h.hoverNode {
			h.createRoad(h.selectedNode, h.hoverNode)
			h.selectedNode = nil // Reset for next road
		}
	}

	// Right click to cancel selection
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		h.selectedNode = nil
	}
}

func (h *InputHandler) createRoad(from, to *road.Node) {
	h.world.Mu.Lock()
	defer h.world.Mu.Unlock()

	// Create road in both directions
	newRoad := road.NewRoad(
		from.ID+"-"+to.ID,
		from,
		to,
		40.0, // Default speed
	)

	h.world.Roads = append(h.world.Roads, newRoad)

	// Update intersections
	fromIntersection := h.world.IntersectionsByNode[from.ID]
	toIntersection := h.world.IntersectionsByNode[to.ID]

	if fromIntersection != nil {
		fromIntersection.AddOutgoing(newRoad)
	}

	if toIntersection != nil {
		toIntersection.AddIncoming(newRoad)
	}
}

func (h *InputHandler) findNearestNode(x, y float64) *road.Node {
	var nearest *road.Node
	minDist := h.maxSnapDistance

	for _, node := range h.world.Nodes {
		dx := node.X - x
		dy := node.Y - y
		dist := math.Sqrt(dx*dx + dy*dy)

		if dist < minDist {
			minDist = dist
			nearest = node
		}
	}

	return nearest
}