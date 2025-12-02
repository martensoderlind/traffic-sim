package renderer

import (
	"image/color"
	"math"
	"traffic-sim/internal/input"
	"traffic-sim/internal/road"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type OverlayRenderer struct{}

func NewOverlayRenderer() *OverlayRenderer {
	return &OverlayRenderer{}
}

func (or *OverlayRenderer) RenderToolOverlay(screen *ebiten.Image, inputHandler *input.InputHandler) {
	mode := inputHandler.Mode()

	switch mode {
	case input.ModeRoadBuilding:
		or.renderRoadBuildingOverlay(screen, inputHandler)
	case input.ModeNodeMoving:
		or.renderNodeMovingOverlay(screen, inputHandler)
	case input.ModeSpawning:
		or.renderSpawningOverlay(screen, inputHandler)
	case input.ModeDespawning:
		or.renderDespawningOverlay(screen, inputHandler)
	case input.ModeRoadDeleting:
		or.renderRoadDeletingOverlay(screen, inputHandler)
	case input.ModeNodeDeleting:
		or.renderNodeDeletingOverlay(screen, inputHandler)
	case input.ModeTrafficLight:
		or.renderTrafficLightOverlay(screen, inputHandler)
	}
}

func (or *OverlayRenderer) renderRoadBuildingOverlay(screen *ebiten.Image, inputHandler *input.InputHandler) {
	mouseX, mouseY := inputHandler.MousePos()
	mx := float64(mouseX)
	my := float64(mouseY)

	roadTool := inputHandler.RoadTool()
	hoverNode := roadTool.GetHoverNode(mx, my)

	if hoverNode != nil {
		vector.StrokeCircle(screen, float32(hoverNode.X), float32(hoverNode.Y), 10, 2, color.RGBA{100, 255, 100, 255}, false)
	} else {
		hoverRoad, snapX, snapY := roadTool.GetHoverRoad(mx, my)

		if hoverRoad != nil {
			vector.FillCircle(screen, float32(snapX), float32(snapY), 5, color.RGBA{100, 200, 255, 200}, false)
			vector.StrokeCircle(screen, float32(snapX), float32(snapY), 8, 2, color.RGBA{100, 200, 255, 255}, false)
		}
	}

	selectedNode := roadTool.GetSelectedNode()
	if selectedNode != nil {
		vector.StrokeCircle(screen, float32(selectedNode.X), float32(selectedNode.Y), 12, 2, color.RGBA{255, 255, 100, 255}, false)

		if hoverNode != nil && hoverNode != selectedNode {
			vector.StrokeLine(screen, float32(selectedNode.X), float32(selectedNode.Y), float32(hoverNode.X), float32(hoverNode.Y), 2, color.RGBA{255, 255, 100, 150}, false)
		} else {
			hoverRoad, snapX, snapY := roadTool.GetHoverRoad(mx, my)
			if hoverRoad != nil {
				vector.StrokeLine(screen, float32(selectedNode.X), float32(selectedNode.Y), float32(snapX), float32(snapY), 2, color.RGBA{255, 255, 100, 150}, false)
			} else {
				vector.StrokeLine(screen, float32(selectedNode.X), float32(selectedNode.Y), float32(mouseX), float32(mouseY), 2, color.RGBA{255, 255, 100, 100}, false)
			}
		}
	}
}

func (or *OverlayRenderer) renderNodeMovingOverlay(screen *ebiten.Image, inputHandler *input.InputHandler) {
	mouseX, mouseY := inputHandler.MousePos()
	mx := float64(mouseX)
	my := float64(mouseY)

	moveTool := inputHandler.MoveTool()

	if moveTool.IsDragging() {
		draggedNode := moveTool.GetDraggedNode()
		if draggedNode != nil {
			vector.StrokeCircle(screen, float32(draggedNode.X), float32(draggedNode.Y), 15, 3, color.RGBA{255, 100, 255, 255}, false)
		}
	} else {
		hoverNode := moveTool.GetHoverNode(mx, my)
		if hoverNode != nil {
			vector.StrokeCircle(screen, float32(hoverNode.X), float32(hoverNode.Y), 12, 2, color.RGBA{255, 150, 255, 200}, false)
		}
	}
}

func (or *OverlayRenderer) renderSpawningOverlay(screen *ebiten.Image, inputHandler *input.InputHandler) {
	mouseX, mouseY := inputHandler.MousePos()
	mx := float64(mouseX)
	my := float64(mouseY)

	spawnTool := inputHandler.SpawnTool()
	hoverNode := spawnTool.GetHoverNode(mx, my)

	if hoverNode != nil {
		vector.StrokeCircle(screen, float32(hoverNode.X), float32(hoverNode.Y), 12, 2, color.RGBA{100, 255, 100, 255}, false)
	}

	selectedNode := spawnTool.GetSelectedNode()
	if selectedNode != nil {
		vector.StrokeCircle(screen, float32(selectedNode.X), float32(selectedNode.Y), 15, 3, color.RGBA{255, 255, 100, 255}, false)

		selectedRoad := spawnTool.GetSelectedRoad()
		if selectedRoad != nil {
			dx := float32(selectedRoad.To.X - selectedRoad.From.X)
			dy := float32(selectedRoad.To.Y - selectedRoad.From.Y)
			length := float32(math.Sqrt(float64(dx*dx + dy*dy)))
			
			if length > 0 {
				dx /= length
				dy /= length

				x := float32(selectedNode.X)
				y := float32(selectedNode.Y)
				arrowLen := float32(25.0)
				x2 := x + dx*arrowLen
				y2 := y + dy*arrowLen

				vector.StrokeLine(screen, x, y, x2, y2, 4, color.RGBA{255, 255, 100, 255}, false)
			}
		}
	}
}

func (or *OverlayRenderer) renderDespawningOverlay(screen *ebiten.Image, inputHandler *input.InputHandler) {
	mouseX, mouseY := inputHandler.MousePos()
	mx := float64(mouseX)
	my := float64(mouseY)

	despawnTool := inputHandler.DespawnTool()
	hoverNode := despawnTool.GetHoverNode(mx, my)

	if hoverNode != nil {
		vector.StrokeCircle(screen, float32(hoverNode.X), float32(hoverNode.Y), 12, 2, color.RGBA{255, 100, 100, 255}, false)
	}

	selectedNode := despawnTool.GetSelectedNode()
	if selectedNode != nil {
		vector.StrokeCircle(screen, float32(selectedNode.X), float32(selectedNode.Y), 15, 3, color.RGBA{255, 255, 100, 255}, false)

		selectedRoad := despawnTool.GetSelectedRoad()
		if selectedRoad != nil {
			dx := float32(selectedRoad.To.X - selectedRoad.From.X)
			dy := float32(selectedRoad.To.Y - selectedRoad.From.Y)
			length := float32(math.Sqrt(float64(dx*dx + dy*dy)))
			
			if length > 0 {
				dx /= length
				dy /= length

				x := float32(selectedNode.X)
				y := float32(selectedNode.Y)
				arrowLen := float32(25.0)
				x1 := x - dx*arrowLen
				y1 := y - dy*arrowLen

				vector.StrokeLine(screen, x1, y1, x, y, 4, color.RGBA{255, 255, 100, 255}, false)
			}
		}
	}
}

func (or *OverlayRenderer) renderRoadDeletingOverlay(screen *ebiten.Image, inputHandler *input.InputHandler) {
	mouseX, mouseY := inputHandler.MousePos()
	mx := float64(mouseX)
	my := float64(mouseY)

	roadDeleteTool := inputHandler.RoadDeleteTool()
	hoverRoad := roadDeleteTool.GetHoverRoad(mx, my)

	if hoverRoad != nil {
		x1 := float32(hoverRoad.From.X)
		y1 := float32(hoverRoad.From.Y)
		x2 := float32(hoverRoad.To.X)
		y2 := float32(hoverRoad.To.Y)

		vector.StrokeLine(screen, x1, y1, x2, y2, float32(hoverRoad.Width+4), color.RGBA{255, 50, 50, 200}, false)
	}
}

func (or *OverlayRenderer) renderNodeDeletingOverlay(screen *ebiten.Image, inputHandler *input.InputHandler) {
	mouseX, mouseY := inputHandler.MousePos()
	mx := float64(mouseX)
	my := float64(mouseY)

	nodeDeleteTool := inputHandler.NodeDeleteTool()
	hoverNode := nodeDeleteTool.GetHoverNode(mx, my)

	if hoverNode != nil {
		vector.StrokeCircle(screen, float32(hoverNode.X), float32(hoverNode.Y), 15, 3, color.RGBA{255, 50, 50, 255}, false)
		vector.FillCircle(screen, float32(hoverNode.X), float32(hoverNode.Y), 12, color.RGBA{255, 50, 50, 100}, false)
	}
}

func (or *OverlayRenderer) renderTrafficLightOverlay(screen *ebiten.Image, inputHandler *input.InputHandler) {
	mouseX, mouseY := inputHandler.MousePos()
	mx := float64(mouseX)
	my := float64(mouseY)

	tlTool := inputHandler.TrafficLightTool()
	hoverNode := tlTool.GetHoverNode(mx, my)

	if hoverNode != nil {
		vector.StrokeCircle(screen, float32(hoverNode.X), float32(hoverNode.Y), 12, 2, color.RGBA{255, 255, 100, 255}, false)
	}

	selectedNode := tlTool.GetSelectedNode()
	if selectedNode != nil {
		vector.StrokeCircle(screen, float32(selectedNode.X), float32(selectedNode.Y), 15, 3, color.RGBA{255, 255, 100, 255}, false)

		availableRoads := tlTool.GetAvailableRoads()
		selectedRoads := tlTool.GetSelectedRoads()

		for _, rd := range availableRoads {
			isSelected := false
			for _, selRd := range selectedRoads {
				if selRd == rd {
					isSelected = true
					break
				}
			}

			roadColor := color.RGBA{100, 100, 120, 150}
			arrowColor := color.RGBA{150, 150, 170, 255}
			if isSelected {
				roadColor = color.RGBA{100, 255, 100, 200}
				arrowColor = color.RGBA{100, 255, 100, 255}
			}

			x1 := float32(rd.From.X)
			y1 := float32(rd.From.Y)
			x2 := float32(rd.To.X)
			y2 := float32(rd.To.Y)

			vector.StrokeLine(screen, x1, y1, x2, y2, float32(rd.Width+2), roadColor, false)

			dx := float32(rd.To.X - rd.From.X)
			dy := float32(rd.To.Y - rd.From.Y)
			length := float32(math.Sqrt(float64(dx*dx + dy*dy)))
			
			if length > 0 {
				dx /= length
				dy /= length

				x := float32(selectedNode.X)
				y := float32(selectedNode.Y)
				arrowLen := float32(25.0)
				x1 := x - dx*arrowLen
				y1 := y - dy*arrowLen

				vector.StrokeLine(screen, x1, y1, x, y, 4, arrowColor, false)
			}
		}
	}
}

func (or *OverlayRenderer) drawDirectionIndicator(screen *ebiten.Image, node *road.Node, rd *road.Road, clr color.RGBA) {
	dx := float32(rd.To.X - rd.From.X)
	dy := float32(rd.To.Y - rd.From.Y)
	length := float32(math.Sqrt(float64(dx*dx + dy*dy)))

	if length == 0 {
		return
	}

	dx /= length
	dy /= length

	x := float32(node.X)
	y := float32(node.Y)
	arrowLen := float32(25.0)
	x2 := x + dx*arrowLen
	y2 := y + dy*arrowLen

	vector.StrokeLine(screen, x, y, x2, y2, 4, clr, false)
}

func (or *OverlayRenderer) drawReverseDirectionIndicator(screen *ebiten.Image, node *road.Node, rd *road.Road, clr color.RGBA) {
	dx := float32(rd.To.X - rd.From.X)
	dy := float32(rd.To.Y - rd.From.Y)
	length := float32(math.Sqrt(float64(dx*dx + dy*dy)))

	if length == 0 {
		return
	}

	dx /= length
	dy /= length

	x := float32(node.X)
	y := float32(node.Y)
	arrowLen := float32(25.0)
	x1 := x - dx*arrowLen
	y1 := y - dy*arrowLen

	vector.StrokeLine(screen, x1, y1, x, y, 4, clr, false)
}