package renderer

import (
	"image/color"
	"math"
	"traffic-sim/internal/input"
	"traffic-sim/internal/road"
	"traffic-sim/internal/ui"
	"traffic-sim/internal/world"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Renderer struct {
	World        *world.World
	InputHandler *input.InputHandler
	Toolbar      *ui.Toolbar
	screenWidth  int
	screenHeight int
}

func NewRenderer(w *world.World, inputHandler *input.InputHandler) *Renderer {
	return &Renderer{
		World:        w,
		InputHandler: inputHandler,
		Toolbar:      ui.NewToolbar(inputHandler),
		screenWidth:  1920,
		screenHeight: 1080,
	}
}

func (r *Renderer) Update() error {
	return nil
}

func (r *Renderer) renderVehicles(screen *ebiten.Image){
	for _, v := range r.World.Vehicles {
		pos := v.Position()
		angle := v.GetAngle()
		
		width := float32(5.0)
		height := float32(10.0)
		
		cx := float32(pos.X)
		cy := float32(pos.Y)
		
		cos := float32(math.Cos(angle))
		sin := float32(math.Sin(angle))
		
		hw := width / 2
		hh := height / 2
		
		corners := [][2]float32{
			{-hw, -hh},
			{hw, -hh},
			{hw, hh},
			{-hw, hh},
		}
		
		rotated := make([][2]float32, 4)
		for i, corner := range corners {
			rotated[i][0] = cx + corner[0]*cos - corner[1]*sin
			rotated[i][1] = cy + corner[0]*sin + corner[1]*cos
		}
		
		vertices := []ebiten.Vertex{
			{DstX: rotated[0][0], DstY: rotated[0][1], SrcX: 0, SrcY: 0, ColorR: 1.0, ColorG: 0.2, ColorB: 0.2, ColorA: 1},
			{DstX: rotated[1][0], DstY: rotated[1][1], SrcX: 0, SrcY: 0, ColorR: 1.0, ColorG: 0.2, ColorB: 0.2, ColorA: 1},
			{DstX: rotated[2][0], DstY: rotated[2][1], SrcX: 0, SrcY: 0, ColorR: 1.0, ColorG: 0.2, ColorB: 0.2, ColorA: 1},
			{DstX: rotated[3][0], DstY: rotated[3][1], SrcX: 0, SrcY: 0, ColorR: 1.0, ColorG: 0.2, ColorB: 0.2, ColorA: 1},
		}
		
		indices := []uint16{0, 1, 2, 0, 2, 3}
		
		screen.DrawTriangles(vertices, indices, emptyImage(screen), nil)
		
		vector.StrokeLine(screen, rotated[0][0], rotated[0][1], rotated[1][0], rotated[1][1], 1, color.RGBA{255, 150, 150, 255}, false)
		vector.StrokeLine(screen, rotated[1][0], rotated[1][1], rotated[2][0], rotated[2][1], 1, color.RGBA{255, 150, 150, 255}, false)
		vector.StrokeLine(screen, rotated[2][0], rotated[2][1], rotated[3][0], rotated[3][1], 1, color.RGBA{255, 150, 150, 255}, false)
		vector.StrokeLine(screen, rotated[3][0], rotated[3][1], rotated[0][0], rotated[0][1], 1, color.RGBA{255, 150, 150, 255}, false)
	}
}

func (r *Renderer) renderRoads(screen *ebiten.Image){
	for _, rd := range r.World.Roads {
		r.renderSingleRoad(screen, rd)
	}
}

func (r *Renderer) renderSingleRoad(screen *ebiten.Image, rd *road.Road) {
	x1 := rd.From.X
	y1 := rd.From.Y
	x2 := rd.To.X
	y2 := rd.To.Y
	
	dx := x2 - x1
	dy := y2 - y1
	length := math.Sqrt(dx*dx + dy*dy)
	
	if length == 0 {
		return
	}
	
	perpX := -dy / length
	perpY := dx / length
	
	offset := 0.0
	if rd.ReverseRoad != nil {
		offset = rd.Width * 0.5
	}
	
	halfWidth := rd.Width / 2.0
	
	x1 += perpX * offset
	y1 += perpY * offset
	x2 += perpX * offset
	y2 += perpY * offset
	
	p1x := float32(x1 + perpX*halfWidth)
	p1y := float32(y1 + perpY*halfWidth)
	
	p2x := float32(x1 - perpX*halfWidth)
	p2y := float32(y1 - perpY*halfWidth)
	
	p3x := float32(x2 - perpX*halfWidth)
	p3y := float32(y2 - perpY*halfWidth)
	
	p4x := float32(x2 + perpX*halfWidth)
	p4y := float32(y2 + perpY*halfWidth)
	
	vertices := []ebiten.Vertex{
		{DstX: p1x, DstY: p1y, SrcX: 0, SrcY: 0, ColorR: 0.31, ColorG: 0.31, ColorB: 0.35, ColorA: 1},
		{DstX: p2x, DstY: p2y, SrcX: 0, SrcY: 0, ColorR: 0.31, ColorG: 0.31, ColorB: 0.35, ColorA: 1},
		{DstX: p3x, DstY: p3y, SrcX: 0, SrcY: 0, ColorR: 0.31, ColorG: 0.31, ColorB: 0.35, ColorA: 1},
		{DstX: p4x, DstY: p4y, SrcX: 0, SrcY: 0, ColorR: 0.31, ColorG: 0.31, ColorB: 0.35, ColorA: 1},
	}
	
	indices := []uint16{0, 1, 2, 0, 2, 3}
	
	screen.DrawTriangles(vertices, indices, emptyImage(screen), nil)
}

func emptyImage(screen *ebiten.Image) *ebiten.Image {
	const size = 1
	img := ebiten.NewImage(size, size)
	img.Fill(color.White)
	return img
}

func (r *Renderer) renderNodes(screen *ebiten.Image){
	for _, node := range r.World.Nodes {
		x := float32(node.X)
		y := float32(node.Y)
		
		vector.FillCircle(
			screen,
			x, y,
			8,
			color.RGBA{100, 100, 120, 255},
			false,
		)
	}
}

func (r *Renderer) renderSpawnPoints(screen *ebiten.Image) {
	for _, sp := range r.World.SpawnPoints {
		if !sp.Enabled {
			continue
		}

		x := float32(sp.Node.X)
		y := float32(sp.Node.Y)

		vector.FillCircle(
			screen,
			x, y,
			10,
			color.RGBA{50, 255, 50, 200},
			false,
		)

		vector.StrokeCircle(
			screen,
			x, y,
			10,
			2,
			color.RGBA{100, 255, 100, 255},
			false,
		)

		dx := float32(sp.Road.To.X - sp.Road.From.X)
		dy := float32(sp.Road.To.Y - sp.Road.From.Y)
		length := float32(math.Sqrt(float64(dx*dx + dy*dy)))
		
		if length > 0 {
			dx /= length
			dy /= length

			arrowLen := float32(20.0)
			x2 := x + dx*arrowLen
			y2 := y + dy*arrowLen

			vector.StrokeLine(
				screen,
				x, y,
				x2, y2,
				3,
				color.RGBA{50, 255, 50, 255},
				false,
			)

			arrowSize := float32(8.0)
			angle := float32(0.5)
			
			leftX := x2 - (dx*arrowSize*float32(math.Cos(float64(angle))) - dy*arrowSize*float32(math.Sin(float64(angle))))
			leftY := y2 - (dy*arrowSize*float32(math.Cos(float64(angle))) + dx*arrowSize*float32(math.Sin(float64(angle))))
			
			rightX := x2 - (dx*arrowSize*float32(math.Cos(float64(-angle))) - dy*arrowSize*float32(math.Sin(float64(-angle))))
			rightY := y2 - (dy*arrowSize*float32(math.Cos(float64(-angle))) + dx*arrowSize*float32(math.Sin(float64(-angle))))

			vector.StrokeLine(screen, x2, y2, leftX, leftY, 3, color.RGBA{50, 255, 50, 255}, false)
			vector.StrokeLine(screen, x2, y2, rightX, rightY, 3, color.RGBA{50, 255, 50, 255}, false)
		}
	}
}

func (r *Renderer) renderDespawnPoints(screen *ebiten.Image) {
	for _, dp := range r.World.DespawnPoints {
		if !dp.Enabled {
			continue
		}

		x := float32(dp.Node.X)
		y := float32(dp.Node.Y)

		vector.FillCircle(
			screen,
			x, y,
			10,
			color.RGBA{255, 50, 50, 200},
			false,
		)

		vector.StrokeCircle(
			screen,
			x, y,
			10,
			2,
			color.RGBA{255, 100, 100, 255},
			false,
		)

		dx := float32(dp.Road.To.X - dp.Road.From.X)
		dy := float32(dp.Road.To.Y - dp.Road.From.Y)
		length := float32(math.Sqrt(float64(dx*dx + dy*dy)))
		
		if length > 0 {
			dx /= length
			dy /= length

			arrowLen := float32(20.0)
			x1 := x - dx*arrowLen
			y1 := y - dy*arrowLen

			vector.StrokeLine(
				screen,
				x1, y1,
				x, y,
				3,
				color.RGBA{255, 50, 50, 255},
				false,
			)

			arrowSize := float32(8.0)
			angle := float32(0.5)
			
			leftX := x - (dx*arrowSize*float32(math.Cos(float64(angle))) + dy*arrowSize*float32(math.Sin(float64(angle))))
			leftY := y - (dy*arrowSize*float32(math.Cos(float64(angle))) - dx*arrowSize*float32(math.Sin(float64(angle))))
			
			rightX := x - (dx*arrowSize*float32(math.Cos(float64(-angle))) + dy*arrowSize*float32(math.Sin(float64(-angle))))
			rightY := y - (dy*arrowSize*float32(math.Cos(float64(-angle))) - dx*arrowSize*float32(math.Sin(float64(-angle))))

			vector.StrokeLine(screen, x, y, leftX, leftY, 3, color.RGBA{255, 50, 50, 255}, false)
			vector.StrokeLine(screen, x, y, rightX, rightY, 3, color.RGBA{255, 50, 50, 255}, false)
		}
	}
}

func (r *Renderer) renderToolOverlay(screen *ebiten.Image) {
	mode := r.InputHandler.Mode()
	
	if mode == input.ModeRoadBuilding {
		r.renderRoadBuildingOverlay(screen)
	} else if mode == input.ModeNodeMoving {
		r.renderNodeMovingOverlay(screen)
	} else if mode == input.ModeSpawning {
		r.renderSpawningOverlay(screen)
	} else if mode == input.ModeDespawning {
		r.renderDespawningOverlay(screen)
	} else if mode == input.ModeRoadDeleting {
		r.renderRoadDeletingOverlay(screen)
	} else if mode == input.ModeNodeDeleting {
		r.renderNodeDeletingOverlay(screen)
	} else if mode == input.ModeTrafficLight {
    r.renderTrafficLightOverlay(screen)
	}
}

func (r *Renderer) renderSpawningOverlay(screen *ebiten.Image) {
	mouseX, mouseY := r.InputHandler.MousePos()
	mx := float64(mouseX)
	my := float64(mouseY)

	spawnTool := r.InputHandler.SpawnTool()
	hoverNode := spawnTool.GetHoverNode(mx, my)

	if hoverNode != nil {
		vector.StrokeCircle(
			screen,
			float32(hoverNode.X),
			float32(hoverNode.Y),
			12,
			2,
			color.RGBA{100, 255, 100, 255},
			false,
		)
	}

	selectedNode := spawnTool.GetSelectedNode()
	if selectedNode != nil {
		vector.StrokeCircle(
			screen,
			float32(selectedNode.X),
			float32(selectedNode.Y),
			15,
			3,
			color.RGBA{255, 255, 100, 255},
			false,
		)

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

				vector.StrokeLine(
					screen,
					x, y,
					x2, y2,
					4,
					color.RGBA{255, 255, 100, 255},
					false,
				)
			}
		}
	}
}

func (r *Renderer) renderDespawningOverlay(screen *ebiten.Image) {
	mouseX, mouseY := r.InputHandler.MousePos()
	mx := float64(mouseX)
	my := float64(mouseY)

	despawnTool := r.InputHandler.DespawnTool()
	hoverNode := despawnTool.GetHoverNode(mx, my)

	if hoverNode != nil {
		vector.StrokeCircle(
			screen,
			float32(hoverNode.X),
			float32(hoverNode.Y),
			12,
			2,
			color.RGBA{255, 100, 100, 255},
			false,
		)
	}

	selectedNode := despawnTool.GetSelectedNode()
	if selectedNode != nil {
		vector.StrokeCircle(
			screen,
			float32(selectedNode.X),
			float32(selectedNode.Y),
			15,
			3,
			color.RGBA{255, 255, 100, 255},
			false,
		)

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

				vector.StrokeLine(
					screen,
					x1, y1,
					x, y,
					4,
					color.RGBA{255, 255, 100, 255},
					false,
				)
			}
		}
	}
}

func (r *Renderer) renderRoadDeletingOverlay(screen *ebiten.Image) {
	mouseX, mouseY := r.InputHandler.MousePos()
	mx := float64(mouseX)
	my := float64(mouseY)

	roadDeleteTool := r.InputHandler.RoadDeleteTool()
	hoverRoad := roadDeleteTool.GetHoverRoad(mx, my)

	if hoverRoad != nil {
		x1 := float32(hoverRoad.From.X)
		y1 := float32(hoverRoad.From.Y)
		x2 := float32(hoverRoad.To.X)
		y2 := float32(hoverRoad.To.Y)

		vector.StrokeLine(
			screen,
			x1, y1,
			x2, y2,
			float32(hoverRoad.Width+4),
			color.RGBA{255, 50, 50, 200},
			false,
		)
	}
}

func (r *Renderer) renderNodeDeletingOverlay(screen *ebiten.Image) {
	mouseX, mouseY := r.InputHandler.MousePos()
	mx := float64(mouseX)
	my := float64(mouseY)

	nodeDeleteTool := r.InputHandler.NodeDeleteTool()
	hoverNode := nodeDeleteTool.GetHoverNode(mx, my)

	if hoverNode != nil {
		vector.StrokeCircle(
			screen,
			float32(hoverNode.X),
			float32(hoverNode.Y),
			15,
			3,
			color.RGBA{255, 50, 50, 255},
			false,
		)
		
		vector.FillCircle(
			screen,
			float32(hoverNode.X),
			float32(hoverNode.Y),
			12,
			color.RGBA{255, 50, 50, 100},
			false,
		)
	}
}

func (r *Renderer) renderNodeMovingOverlay(screen *ebiten.Image) {
	mouseX, mouseY := r.InputHandler.MousePos()
	mx := float64(mouseX)
	my := float64(mouseY)

	moveTool := r.InputHandler.MoveTool()
	
	if moveTool.IsDragging() {
		draggedNode := moveTool.GetDraggedNode()
		if draggedNode != nil {
			vector.StrokeCircle(
				screen,
				float32(draggedNode.X),
				float32(draggedNode.Y),
				15,
				3,
				color.RGBA{255, 100, 255, 255},
				false,
			)
		}
	} else {
		hoverNode := moveTool.GetHoverNode(mx, my)
		if hoverNode != nil {
			vector.StrokeCircle(
				screen,
				float32(hoverNode.X),
				float32(hoverNode.Y),
				12,
				2,
				color.RGBA{255, 150, 255, 200},
				false)
		}
	}
}

func (r *Renderer) renderRoadBuildingOverlay(screen *ebiten.Image) {
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

func (r *Renderer) renderTrafficLights(screen *ebiten.Image) {
	for _, light := range r.World.TrafficLights {
		if !light.Enabled {
			continue
		}

		// x := float32(light.Intersection.ID)
		// y := float32(light.Intersection.ID)
		
		node := r.findNodeByID(light.Intersection.ID)
		if node == nil {
			continue
		}
		
		x := float32(node.X)
		y := float32(node.Y)

		offset := float32(15.0)
		
		var lightColor color.RGBA
		switch light.State {
		case road.LightRed:
			lightColor = color.RGBA{255, 50, 50, 255}
		case road.LightYellow:
			lightColor = color.RGBA{255, 255, 50, 255}
		case road.LightGreen:
			lightColor = color.RGBA{50, 255, 50, 255}
		}

		vector.FillCircle(
			screen,
			x+offset, y-offset,
			8,
			lightColor,
			false,
		)

		vector.StrokeCircle(
			screen,
			x+offset, y-offset,
			8,
			2,
			color.RGBA{40, 40, 40, 255},
			false,
		)
	}
}

func (r *Renderer) findNodeByID(id string) *road.Node {
	for _, node := range r.World.Nodes {
		if node.ID == id {
			return node
		}
	}
	return nil
}

func (r *Renderer) renderTrafficLightOverlay(screen *ebiten.Image) {
	mouseX, mouseY := r.InputHandler.MousePos()
	mx := float64(mouseX)
	my := float64(mouseY)

	tlTool := r.InputHandler.TrafficLightTool()
	hoverNode := tlTool.GetHoverNode(mx, my)

	if hoverNode != nil {
		vector.StrokeCircle(
			screen,
			float32(hoverNode.X),
			float32(hoverNode.Y),
			12,
			2,
			color.RGBA{255, 255, 100, 255},
			false,
		)
	}

	selectedNode := tlTool.GetSelectedNode()
	if selectedNode != nil {
		vector.StrokeCircle(
			screen,
			float32(selectedNode.X),
			float32(selectedNode.Y),
			15,
			3,
			color.RGBA{255, 255, 100, 255},
			false,
		)

		selectedRoad := tlTool.GetSelectedRoad()
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

				vector.StrokeLine(
					screen,
					x1, y1,
					x, y,
					4,
					color.RGBA{255, 255, 100, 255},
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
	r.renderSpawnPoints(screen)
	r.renderDespawnPoints(screen)
	r.renderVehicles(screen)
	r.renderToolOverlay(screen)
	r.renderTrafficLights(screen)
	r.Toolbar.Draw(screen)
}

func (r *Renderer) Layout(w, h int) (int, int) {
	r.screenWidth = w
	r.screenHeight = h
	return w,h
}