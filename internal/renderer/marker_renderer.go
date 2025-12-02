package renderer

import (
	"image/color"
	"math"
	"traffic-sim/internal/road"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type MarkerRenderer struct{}

func NewMarkerRenderer() *MarkerRenderer {
	return &MarkerRenderer{}
}

func (mr *MarkerRenderer) RenderSpawnPoints(screen *ebiten.Image, spawnPoints []*road.SpawnPoint) {
	for _, sp := range spawnPoints {
		if !sp.Enabled {
			continue
		}

		x := float32(sp.Node.X)
		y := float32(sp.Node.Y)

		vector.FillCircle(screen, x, y, 10, color.RGBA{50, 255, 50, 200}, false)
		vector.StrokeCircle(screen, x, y, 10, 2, color.RGBA{100, 255, 100, 255}, false)

		mr.drawArrowFromNode(screen, sp.Node, sp.Road, color.RGBA{50, 255, 50, 255})
	}
}

func (mr *MarkerRenderer) RenderDespawnPoints(screen *ebiten.Image, despawnPoints []*road.DespawnPoint) {
	for _, dp := range despawnPoints {
		if !dp.Enabled {
			continue
		}

		x := float32(dp.Node.X)
		y := float32(dp.Node.Y)

		vector.FillCircle(screen, x, y, 10, color.RGBA{255, 50, 50, 200}, false)
		vector.StrokeCircle(screen, x, y, 10, 2, color.RGBA{255, 100, 100, 255}, false)

		mr.drawArrowToNode(screen, dp.Node, dp.Road, color.RGBA{255, 50, 50, 255})
	}
}

func (mr *MarkerRenderer) RenderTrafficLights(screen *ebiten.Image, trafficLights []*road.TrafficLight, nodes []*road.Node) {
	for _, light := range trafficLights {
		if !light.Enabled {
			continue
		}

		node := mr.findNodeByID(light.Intersection.ID, nodes)
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

		vector.FillCircle(screen, x+offset, y-offset, 8, lightColor, false)
		vector.StrokeCircle(screen, x+offset, y-offset, 8, 2, color.RGBA{40, 40, 40, 255}, false)
	}
}

func (mr *MarkerRenderer) drawArrowFromNode(screen *ebiten.Image, node *road.Node, rd *road.Road, clr color.RGBA) {
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
	arrowLen := float32(20.0)
	x2 := x + dx*arrowLen
	y2 := y + dy*arrowLen

	vector.StrokeLine(screen, x, y, x2, y2, 3, clr, false)

	arrowSize := float32(8.0)
	angle := float32(0.5)

	leftX := x2 - (dx*arrowSize*float32(math.Cos(float64(angle))) - dy*arrowSize*float32(math.Sin(float64(angle))))
	leftY := y2 - (dy*arrowSize*float32(math.Cos(float64(angle))) + dx*arrowSize*float32(math.Sin(float64(angle))))

	rightX := x2 - (dx*arrowSize*float32(math.Cos(float64(-angle))) - dy*arrowSize*float32(math.Sin(float64(-angle))))
	rightY := y2 - (dy*arrowSize*float32(math.Cos(float64(-angle))) + dx*arrowSize*float32(math.Sin(float64(-angle))))

	vector.StrokeLine(screen, x2, y2, leftX, leftY, 3, clr, false)
	vector.StrokeLine(screen, x2, y2, rightX, rightY, 3, clr, false)
}

func (mr *MarkerRenderer) drawArrowToNode(screen *ebiten.Image, node *road.Node, rd *road.Road, clr color.RGBA) {
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
	arrowLen := float32(20.0)
	x1 := x - dx*arrowLen
	y1 := y - dy*arrowLen

	vector.StrokeLine(screen, x1, y1, x, y, 3, clr, false)

	arrowSize := float32(8.0)
	angle := float32(0.5)

	leftX := x - (dx*arrowSize*float32(math.Cos(float64(angle))) + dy*arrowSize*float32(math.Sin(float64(angle))))
	leftY := y - (dy*arrowSize*float32(math.Cos(float64(angle))) - dx*arrowSize*float32(math.Sin(float64(angle))))

	rightX := x - (dx*arrowSize*float32(math.Cos(float64(-angle))) + dy*arrowSize*float32(math.Sin(float64(-angle))))
	rightY := y - (dy*arrowSize*float32(math.Cos(float64(-angle))) - dx*arrowSize*float32(math.Sin(float64(-angle))))

	vector.StrokeLine(screen, x, y, leftX, leftY, 3, clr, false)
	vector.StrokeLine(screen, x, y, rightX, rightY, 3, clr, false)
}

func (mr *MarkerRenderer) findNodeByID(id string, nodes []*road.Node) *road.Node {
	for _, node := range nodes {
		if node.ID == id {
			return node
		}
	}
	return nil
}