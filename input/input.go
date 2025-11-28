package input

import (
	"fmt"
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

	selectedNode     *road.Node
	hoverNode        *road.Node
	mouseX, mouseY   int
	maxSnapDistance  float64
	minNodeDistance  float64
	nodeCounter      int
	bidirectional    bool
}

func NewInputHandler(world *sim.World) *InputHandler {
	return &InputHandler{
		world:           world,
		mode:            ModeNormal,
		maxSnapDistance: 20.0, 
		minNodeDistance: 30.0, 
		nodeCounter:     1000, 
		bidirectional:   true, 
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

func (h *InputHandler) Bidirectional() bool {
	return h.bidirectional
}

func (h *InputHandler) MousePos() (int, int) {
	return h.mouseX, h.mouseY
}

func (h *InputHandler) Update() {
	h.mouseX, h.mouseY = ebiten.CursorPosition()

	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		if h.mode == ModeNormal {
			h.mode = ModePlacingRoad
			h.selectedNode = nil
		} else {
			h.mode = ModeNormal
			h.selectedNode = nil
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		h.mode = ModeNormal
		h.selectedNode = nil
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyB) {
		h.bidirectional = !h.bidirectional
	}

	h.world.Mu.RLock()
	h.hoverNode = h.findNearestNode(float64(h.mouseX), float64(h.mouseY))
	h.world.Mu.RUnlock()

	switch h.mode {
	case ModePlacingRoad:
		h.updatePlacingRoad()
	}
}

func (h *InputHandler) updatePlacingRoad() {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		var clickedNode *road.Node

		if h.hoverNode != nil {
			clickedNode = h.hoverNode
		} else {
			if h.canPlaceNodeAt(float64(h.mouseX), float64(h.mouseY)) {
				clickedNode = h.createNode(float64(h.mouseX), float64(h.mouseY))
			} else {
				return
			}
		}

		if h.selectedNode == nil {
			h.selectedNode = clickedNode
			return
		}

		if h.selectedNode != clickedNode {
			h.createRoad(h.selectedNode, clickedNode)
			if h.bidirectional {
				h.createRoad(clickedNode, h.selectedNode)
			}
			h.selectedNode = nil 
		}
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		h.selectedNode = nil
	}
}

func (h *InputHandler) canPlaceNodeAt(x, y float64) bool {
	h.world.Mu.RLock()
	defer h.world.Mu.RUnlock()

	for _, node := range h.world.Nodes {
		dx := node.X - x
		dy := node.Y - y
		dist := math.Sqrt(dx*dx + dy*dy)

		if dist < h.minNodeDistance {
			return false
		}
	}

	return true
}

func (h *InputHandler) createNode(x, y float64) *road.Node {
	h.world.Mu.Lock()
	defer h.world.Mu.Unlock()

	h.nodeCounter++
	nodeID := fmt.Sprintf("n%d", h.nodeCounter)

	newNode := &road.Node{
		ID: nodeID,
		X:  x,
		Y:  y,
	}

	h.world.Nodes = append(h.world.Nodes, newNode)

	newIntersection := road.NewIntersection(nodeID)
	h.world.Intersections = append(h.world.Intersections, newIntersection)
	h.world.IntersectionsByNode[nodeID] = newIntersection

	return newNode
}

func (h *InputHandler) createRoad(from, to *road.Node) {
	h.world.Mu.Lock()
	defer h.world.Mu.Unlock()

	newRoad := road.NewRoad(
		from.ID+"-"+to.ID,
		from,
		to,
		40.0,
	)

	h.world.Roads = append(h.world.Roads, newRoad)

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