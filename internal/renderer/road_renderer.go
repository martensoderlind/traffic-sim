package renderer

import (
	"image/color"
	"math"
	"traffic-sim/internal/road"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type RoadRenderer struct{
	shadowOffset float64
	nodeRadius   float64
}

func NewRoadRenderer() *RoadRenderer {
	return &RoadRenderer{
		shadowOffset:  3.0,
		nodeRadius:   8.0,
	}
}

func (rr *RoadRenderer) RenderRoads(screen *ebiten.Image, roads []*road.Road) {
	for _, rd := range roads {
		rr.renderSingleRoad(screen, rd)
	}
}

func (rr *RoadRenderer) renderSingleRoad(screen *ebiten.Image, rd *road.Road) {
	x1, y1 := rd.From.X, rd.From.Y
	x2, y2 := rd.To.X, rd.To.Y

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

	baseX1 := x1 + perpX * offset
	baseY1 := y1 + perpY * offset
	baseX2 := x2 + perpX * offset
	baseY2 := y2 + perpY * offset

	rr.drawRoadShadow(screen, baseX1, baseY1, baseX2, baseY2, rd.Width, perpX, perpY)
	rr.drawRoadBase(screen, baseX1, baseY1, baseX2, baseY2, rd.Width, perpX, perpY)
	rr.drawRoadEdges(screen, baseX1, baseY1, baseX2, baseY2, rd.Width, perpX, perpY)
}

func (rr *RoadRenderer) RenderNodes(screen *ebiten.Image, nodes []*road.Node) {
	for _, node := range nodes {
		x := float32(node.X)
		y := float32(node.Y)
		
		vector.FillCircle(screen, x+1, y+1, float32(rr.nodeRadius), color.RGBA{0, 0, 0, 80}, false)
		
		vector.FillCircle(screen, x, y, float32(rr.nodeRadius), color.RGBA{80, 80, 90, 255}, false)
		
		vector.FillCircle(screen, x, y, float32(rr.nodeRadius), color.RGBA{45, 45, 50, 255}, false)
		
		vector.StrokeCircle(screen, x, y, float32(rr.nodeRadius), 1.5, color.RGBA{234, 231, 228, 255}, false)
	}
}


func (rr *RoadRenderer) drawRoadShadow(screen *ebiten.Image, x1, y1, x2, y2, width, perpX, perpY float64) {
	halfWidth := width / 2.0
	
	shadowX := float32(rr.shadowOffset * perpX)
	shadowY := float32(rr.shadowOffset * perpY)
	
	p1x := float32(x1 + perpX*halfWidth) + shadowX
	p1y := float32(y1 + perpY*halfWidth) + shadowY
	p2x := float32(x1 - perpX*halfWidth) + shadowX
	p2y := float32(y1 - perpY*halfWidth) + shadowY
	p3x := float32(x2 - perpX*halfWidth) + shadowX
	p3y := float32(y2 - perpY*halfWidth) + shadowY
	p4x := float32(x2 + perpX*halfWidth) + shadowX
	p4y := float32(y2 + perpY*halfWidth) + shadowY

	vertices := []ebiten.Vertex{
		{DstX: p1x, DstY: p1y, SrcX: 0, SrcY: 0, ColorR: 0, ColorG: 0, ColorB: 0, ColorA: 0.3},
		{DstX: p2x, DstY: p2y, SrcX: 0, SrcY: 0, ColorR: 0, ColorG: 0, ColorB: 0, ColorA: 0.3},
		{DstX: p3x, DstY: p3y, SrcX: 0, SrcY: 0, ColorR: 0, ColorG: 0, ColorB: 0, ColorA: 0.3},
		{DstX: p4x, DstY: p4y, SrcX: 0, SrcY: 0, ColorR: 0, ColorG: 0, ColorB: 0, ColorA: 0.3},
	}

	indices := []uint16{0, 1, 2, 0, 2, 3}
	screen.DrawTriangles(vertices, indices, createWhiteImage(), nil)
}

func (rr *RoadRenderer) drawRoadBase(screen *ebiten.Image, x1, y1, x2, y2, width, perpX, perpY float64) {
	halfWidth := width / 2.0
	
	p1x := float32(x1 + perpX*halfWidth)
	p1y := float32(y1 + perpY*halfWidth)
	p2x := float32(x1 - perpX*halfWidth)
	p2y := float32(y1 - perpY*halfWidth)
	p3x := float32(x2 - perpX*halfWidth)
	p3y := float32(y2 - perpY*halfWidth)
	p4x := float32(x2 + perpX*halfWidth)
	p4y := float32(y2 + perpY*halfWidth)

	baseColor := color.RGBA{45, 45, 50, 255}
	
	vertices := []ebiten.Vertex{
		{DstX: p1x, DstY: p1y, SrcX: 0, SrcY: 0, ColorR: 0.18, ColorG: 0.18, ColorB: 0.2, ColorA: 1},
		{DstX: p2x, DstY: p2y, SrcX: 0, SrcY: 0, ColorR: 0.16, ColorG: 0.16, ColorB: 0.18, ColorA: 1},
		{DstX: p3x, DstY: p3y, SrcX: 0, SrcY: 0, ColorR: 0.16, ColorG: 0.16, ColorB: 0.18, ColorA: 1},
		{DstX: p4x, DstY: p4y, SrcX: 0, SrcY: 0, ColorR: 0.18, ColorG: 0.18, ColorB: 0.2, ColorA: 1},
	}

	indices := []uint16{0, 1, 2, 0, 2, 3}
	screen.DrawTriangles(vertices, indices, createWhiteImage(), nil)
	
	_ = baseColor
}

// func (rr *RoadRenderer) drawRoadMarkings(screen *ebiten.Image, x1, y1, x2, y2, width, length, perpX, perpY float64) {
// 	dx := x2 - x1
// 	dy := y2 - y1
	
// 	dashLength := 15.0
// 	gapLength := 10.0
// 	segmentLength := dashLength + gapLength
// 	numSegments := int(length / segmentLength)
	
// 	for i := 0; i < numSegments; i++ {
// 		t1 := float64(i) * segmentLength / length
// 		t2 := (float64(i)*segmentLength + dashLength) / length
		
// 		if t2 > 1.0 {
// 			t2 = 1.0
// 		}
		
// 		mx1 := float32(x1 + dx*t1)
// 		my1 := float32(y1 + dy*t1)
// 		mx2 := float32(x1 + dx*t2)
// 		my2 := float32(y1 + dy*t2)
		
// 		vector.StrokeLine(screen, mx1, my1, mx2, my2, 1.5,
// 			color.RGBA{200, 200, 180, 255}, false)
// 	}
// }

func (rr *RoadRenderer) drawRoadEdges(screen *ebiten.Image, x1, y1, x2, y2, width, perpX, perpY float64) {
	halfWidth := width / 2.0
	
	edge1x1 := float32(x1 + perpX*halfWidth)
	edge1y1 := float32(y1 + perpY*halfWidth)
	edge1x2 := float32(x2 + perpX*halfWidth)
	edge1y2 := float32(y2 + perpY*halfWidth)
	
	edge2x1 := float32(x1 - perpX*halfWidth)
	edge2y1 := float32(y1 - perpY*halfWidth)
	edge2x2 := float32(x2 - perpX*halfWidth)
	edge2y2 := float32(y2 - perpY*halfWidth)
	
	edgeColor := color.RGBA{234, 231, 228, 125}
	
	vector.StrokeLine(screen, edge1x1, edge1y1, edge1x2, edge1y2, 1.5, edgeColor, false)
	vector.StrokeLine(screen, edge2x1, edge2y1, edge2x2, edge2y2, 1.5, edgeColor, false)
}


func createWhiteImage() *ebiten.Image {
	img := ebiten.NewImage(1, 1)
	img.Fill(color.White)
	return img
}