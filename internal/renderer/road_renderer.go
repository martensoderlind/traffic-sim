package renderer

import (
	"image/color"
	"math"
	"traffic-sim/internal/road"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type RoadRenderer struct{}

func NewRoadRenderer() *RoadRenderer {
	return &RoadRenderer{}
}

func (rr *RoadRenderer) RenderRoads(screen *ebiten.Image, roads []*road.Road) {
	for _, rd := range roads {
		rr.renderSingleRoad(screen, rd)
	}
}

func (rr *RoadRenderer) renderSingleRoad(screen *ebiten.Image, rd *road.Road) {
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

func (rr *RoadRenderer) RenderNodes(screen *ebiten.Image, nodes []*road.Node) {
	for _, node := range nodes {
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

func emptyImage(screen *ebiten.Image) *ebiten.Image {
	const size = 1
	img := ebiten.NewImage(size, size)
	img.Fill(color.White)
	return img
}