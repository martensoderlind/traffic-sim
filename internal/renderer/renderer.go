package renderer

import (
	"image/color"
	"traffic-sim/internal/input"
	"traffic-sim/internal/world"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Renderer struct {
	World        *world.World
	InputHandler *input.InputHandler
}

func (r *Renderer) Update() error {
	return nil
}

func (r *Renderer) renderVehicles(screen *ebiten.Image){
	for _, v := range r.World.Vehicles {
		pos := v.Position()
		px := float32(pos.X)
		py := float32(pos.Y)

		vector.FillRect(
			screen,
			px-3, py-3,
			6, 6,
			color.RGBA{255, 50, 50, 255},
			false,
		)

		vector.StrokeRect(
			screen,
			px-3, py-3,
			6, 6,
			1,
			color.RGBA{255, 150, 150, 255},
			false,
		)
	}
}

func (r *Renderer) renderRoads(screen *ebiten.Image){
	for _, rd := range r.World.Roads {
		x1 := float32(rd.From.X)
		y1 := float32(rd.From.Y)
		x2 := float32(rd.To.X)
		y2 := float32(rd.To.Y)

		vector.StrokeLine(
			screen,
			x1, y1,
			x2, y2,
			3,
			color.RGBA{80, 80, 90, 255},
			false,
		)

		vector.StrokeLine(
			screen,
			x1, y1,
			x2, y2,
			1,
			color.RGBA{150, 150, 150, 120},
			false,
		)
	}
}

func (r *Renderer) renderNodes(screen *ebiten.Image){
	for _, node := range r.World.Nodes {
		x := float32(node.X)
		y := float32(node.Y)
		
		vector.FillCircle(
			screen,
			x, y,
			6,
			color.RGBA{100, 100, 120, 255},
			false,
		)
	}
}

func (r *Renderer) renderToolOverlay(screen *ebiten.Image) {
	if r.InputHandler.Mode() != input.ModeRoadBuilding {
		return
	}

	mouseX, mouseY := r.InputHandler.MousePos()
	mx := float64(mouseX)
	my := float64(mouseY)

	roadTool := r.InputHandler.RoadTool()
	hoverNode := roadTool.GetHoverNode(mx, my)

	if hoverNode != nil {
		vector.StrokeCircle(
			screen,
			float32(hoverNode.X),
			float32(hoverNode.Y),
			10,
			2,
			color.RGBA{100, 255, 100, 255},
			false,
		)
	} else {
		hoverRoad, snapX, snapY := roadTool.GetHoverRoad(mx, my)
		
		if hoverRoad != nil {
			vector.FillCircle(
				screen,
				float32(snapX),
				float32(snapY),
				5,
				color.RGBA{100, 200, 255, 200},
				false,
			)
			
			vector.StrokeCircle(
				screen,
				float32(snapX),
				float32(snapY),
				8,
				2,
				color.RGBA{100, 200, 255, 255},
				false,
			)
		}
	}

	selectedNode := roadTool.GetSelectedNode()
	if selectedNode != nil {
		vector.StrokeCircle(
			screen,
			float32(selectedNode.X),
			float32(selectedNode.Y),
			12,
			2,
			color.RGBA{255, 255, 100, 255},
			false,
		)

		if hoverNode != nil && hoverNode != selectedNode {
			vector.StrokeLine(
				screen,
				float32(selectedNode.X),
				float32(selectedNode.Y),
				float32(hoverNode.X),
				float32(hoverNode.Y),
				2,
				color.RGBA{255, 255, 100, 150},
				false,
			)
		} else {
			hoverRoad, snapX, snapY := roadTool.GetHoverRoad(mx, my)
			if hoverRoad != nil {
				vector.StrokeLine(
					screen,
					float32(selectedNode.X),
					float32(selectedNode.Y),
					float32(snapX),
					float32(snapY),
					2,
					color.RGBA{255, 255, 100, 150},
					false,
				)
			} else {
				vector.StrokeLine(
					screen,
					float32(selectedNode.X),
					float32(selectedNode.Y),
					float32(mouseX),
					float32(mouseY),
					2,
					color.RGBA{255, 255, 100, 100},
					false,
				)
			}
		}
	}
}

func (r *Renderer) Draw(screen *ebiten.Image) {
	r.World.Mu.RLock()
	defer r.World.Mu.RUnlock()

	screen.Fill(color.RGBA{20, 20, 30, 255})
	
	r.renderNodes(screen)
	r.renderRoads(screen)
	r.renderVehicles(screen)
	r.renderToolOverlay(screen)
}

func (r *Renderer) Layout(w, h int) (int, int) {
	return 800, 600
}