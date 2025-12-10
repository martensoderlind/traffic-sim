package renderer

import (
	"image/color"
	"math"
	"traffic-sim/internal/road"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var (
	colorRoadBase = color.RGBA{45, 45, 50, 255}
	// colorRoadEdge   = color.RGBA{234, 231, 228, 125}
	colorRoadShadow = color.RGBA{0, 0, 0, 80}
)

type RoadRenderer struct {
	shadowOffset float64
}

func NewRoadRenderer() *RoadRenderer {
	return &RoadRenderer{
		shadowOffset: 3.0,
	}
}

func (rr *RoadRenderer) RenderRoads(screen *ebiten.Image, roads []*road.Road, nodes []*road.Node) {
	for _, rd := range roads {
		rr.drawSingleRoadShadow(screen, rd)
	}

	for _, rd := range roads {
		rr.drawSingleRoadBase(screen, rd)
	}

	rr.drawIntersections(screen, nodes, roads)
}

func (rr *RoadRenderer) drawIntersections(screen *ebiten.Image, nodes []*road.Node, roads []*road.Road) {
	nodeRadiusMap := make(map[*road.Node]float64)

	for _, rd := range roads {
		radius := rd.Width
		if r, exists := nodeRadiusMap[rd.From]; !exists || radius > r {
			nodeRadiusMap[rd.From] = radius
		}
		if r, exists := nodeRadiusMap[rd.To]; !exists || radius > r {
			nodeRadiusMap[rd.To] = radius
		}
	}
	for node, radius := range nodeRadiusMap {
		x := float32(node.X)
		y := float32(node.Y)

		vector.FillCircle(screen, x, y, float32(radius), colorRoadBase, true)
	}
}

func getRoadGeometry(rd *road.Road) (x1, y1, x2, y2, perpX, perpY, width float64, valid bool) {
	x1, y1 = rd.From.X, rd.From.Y
	x2, y2 = rd.To.X, rd.To.Y
	width = rd.Width

	dx := x2 - x1
	dy := y2 - y1
	length := math.Sqrt(dx*dx + dy*dy)

	if length == 0 {
		return 0, 0, 0, 0, 0, 0, 0, false
	}

	perpX = -dy / length
	perpY = dx / length

	offset := 0.0
	if rd.ReverseRoad != nil {
		offset = rd.Width * 0.5
	}

	x1 += perpX * offset
	y1 += perpY * offset
	x2 += perpX * offset
	y2 += perpY * offset

	return x1, y1, x2, y2, perpX, perpY, width, true
}

func (rr *RoadRenderer) drawSingleRoadShadow(screen *ebiten.Image, rd *road.Road) {
	if rd.Curve != nil {
		rr.drawCurvedRoadShadow(screen, rd)
	} else {
		rr.drawStraightRoadShadow(screen, rd)
	}
}

func (rr *RoadRenderer) drawSingleRoadBase(screen *ebiten.Image, rd *road.Road) {
	if rd.Curve != nil {
		rr.drawCurvedRoadBase(screen, rd)
	} else {
		rr.drawStraightRoadBase(screen, rd)
	}
}

func (rr *RoadRenderer) drawStraightRoadShadow(screen *ebiten.Image, rd *road.Road) {
	x1, y1, x2, y2, perpX, perpY, width, ok := getRoadGeometry(rd)
	if !ok {
		return
	}

	rr.drawRoadRect(screen, x1, y1, x2, y2, width, perpX, perpY, rr.shadowOffset, rr.shadowOffset, colorRoadShadow)
}

func (rr *RoadRenderer) drawStraightRoadBase(screen *ebiten.Image, rd *road.Road) {
	x1, y1, x2, y2, perpX, perpY, width, ok := getRoadGeometry(rd)
	if !ok {
		return
	}

	rr.drawRoadRect(screen, x1, y1, x2, y2, width, perpX, perpY, 0, 0, colorRoadBase)
}

func (rr *RoadRenderer) drawCurvedRoadShadow(screen *ebiten.Image, rd *road.Road) {
	rr.drawCurvedRoad(screen, rd, rr.shadowOffset, rr.shadowOffset, colorRoadShadow)
}

func (rr *RoadRenderer) drawCurvedRoadBase(screen *ebiten.Image, rd *road.Road) {
	rr.drawCurvedRoad(screen, rd, 0, 0, colorRoadBase)
}

func (rr *RoadRenderer) drawCurvedRoad(screen *ebiten.Image, rd *road.Road, offX, offY float64, clr color.Color) {
	steps := 50
	width := rd.Width
	if rd.ReverseRoad != nil {
		width += 2.0
	}
	halfWidth := width / 2.0

	p0x, p0y := rd.PosAt(0)
	p3x, p3y := rd.PosAt(rd.Length)

	offX0 := p0x - rd.From.X
	offY0 := p0y - rd.From.Y

	p0 := road.Point{X: p0x, Y: p0y}
	p1 := road.Point{X: rd.Curve.ControlP1.X + offX0, Y: rd.Curve.ControlP1.Y + offY0}
	p2 := road.Point{X: rd.Curve.ControlP2.X + offX0, Y: rd.Curve.ControlP2.Y + offY0}
	p3 := road.Point{X: p3x, Y: p3y}

	r, g, b, a := clr.RGBA()
	cr := float32(r) / 65535.0
	cg := float32(g) / 65535.0
	cb := float32(b) / 65535.0
	ca := float32(a) / 65535.0

	for i := 0; i < steps; i++ {
		t1 := float64(i) / float64(steps)
		t2 := float64(i+1) / float64(steps)

		pt1 := rr.cubicBezierPoint(p0, p1, p2, p3, t1)
		pt2 := rr.cubicBezierPoint(p0, p1, p2, p3, t2)

		dx := pt2.X - pt1.X
		dy := pt2.Y - pt1.Y
		length := math.Sqrt(dx*dx + dy*dy)

		if length == 0 {
			continue
		}

		perpX := -dy / length
		perpY := dx / length

		p1x := float32(pt1.X + perpX*halfWidth + offX)
		p1y := float32(pt1.Y + perpY*halfWidth + offY)
		p2x := float32(pt1.X - perpX*halfWidth + offX)
		p2y := float32(pt1.Y - perpY*halfWidth + offY)

		dx2 := pt2.X - pt1.X
		dy2 := pt2.Y - pt1.Y
		length2 := math.Sqrt(dx2*dx2 + dy2*dy2)
		var perpX2, perpY2 float64
		if length2 > 0 {
			perpX2 = -dy2 / length2
			perpY2 = dx2 / length2
		}

		p3x := float32(pt2.X - perpX2*halfWidth + offX)
		p3y := float32(pt2.Y - perpY2*halfWidth + offY)
		p4x := float32(pt2.X + perpX2*halfWidth + offX)
		p4y := float32(pt2.Y + perpY2*halfWidth + offY)

		vertices := []ebiten.Vertex{
			{DstX: p1x, DstY: p1y, SrcX: 0, SrcY: 0, ColorR: cr, ColorG: cg, ColorB: cb, ColorA: ca},
			{DstX: p2x, DstY: p2y, SrcX: 0, SrcY: 0, ColorR: cr, ColorG: cg, ColorB: cb, ColorA: ca},
			{DstX: p3x, DstY: p3y, SrcX: 0, SrcY: 0, ColorR: cr, ColorG: cg, ColorB: cb, ColorA: ca},
			{DstX: p4x, DstY: p4y, SrcX: 0, SrcY: 0, ColorR: cr, ColorG: cg, ColorB: cb, ColorA: ca},
		}

		indices := []uint16{0, 1, 2, 2, 3, 0}
		screen.DrawTriangles(vertices, indices, createWhiteImage(), nil)
	}
}

func (rr *RoadRenderer) cubicBezierPoint(p0, p1, p2, p3 road.Point, t float64) road.Point {
	mt := 1 - t
	mt2 := mt * mt
	mt3 := mt2 * mt
	t2 := t * t
	t3 := t2 * t

	return road.Point{
		X: mt3*p0.X + 3*mt2*t*p1.X + 3*mt*t2*p2.X + t3*p3.X,
		Y: mt3*p0.Y + 3*mt2*t*p1.Y + 3*mt*t2*p2.Y + t3*p3.Y,
	}
}

func (rr *RoadRenderer) drawRoadRect(screen *ebiten.Image, x1, y1, x2, y2, width, perpX, perpY, offX, offY float64, clr color.Color) {
	halfWidth := width / 2.0

	r, g, b, a := clr.RGBA()
	cr := float32(r) / 65535.0
	cg := float32(g) / 65535.0
	cb := float32(b) / 65535.0
	ca := float32(a) / 65535.0

	p1x := float32(x1 + perpX*halfWidth + offX)
	p1y := float32(y1 + perpY*halfWidth + offY)
	p2x := float32(x1 - perpX*halfWidth + offX)
	p2y := float32(y1 - perpY*halfWidth + offY)
	p3x := float32(x2 - perpX*halfWidth + offX)
	p3y := float32(y2 - perpY*halfWidth + offY)
	p4x := float32(x2 + perpX*halfWidth + offX)
	p4y := float32(y2 + perpY*halfWidth + offY)

	vertices := []ebiten.Vertex{
		{DstX: p1x, DstY: p1y, SrcX: 0, SrcY: 0, ColorR: cr, ColorG: cg, ColorB: cb, ColorA: ca},
		{DstX: p2x, DstY: p2y, SrcX: 0, SrcY: 0, ColorR: cr, ColorG: cg, ColorB: cb, ColorA: ca},
		{DstX: p3x, DstY: p3y, SrcX: 0, SrcY: 0, ColorR: cr, ColorG: cg, ColorB: cb, ColorA: ca},
		{DstX: p4x, DstY: p4y, SrcX: 0, SrcY: 0, ColorR: cr, ColorG: cg, ColorB: cb, ColorA: ca},
	}

	indices := []uint16{0, 1, 2, 2, 3, 0}

	screen.DrawTriangles(vertices, indices, createWhiteImage(), nil)
}

// func (rr *RoadRenderer) drawSingleRoadEdges(screen *ebiten.Image, rd *road.Road) {
// 	x1, y1, x2, y2, perpX, perpY, width, ok := getRoadGeometry(rd)
// 	if !ok { return }

// 	halfWidth := width / 2.0
// 	edge1x1 := float32(x1 + perpX*halfWidth)
// 	edge1y1 := float32(y1 + perpY*halfWidth)
// 	edge1x2 := float32(x2 + perpX*halfWidth)
// 	edge1y2 := float32(y2 + perpY*halfWidth)

// 	edge2x1 := float32(x1 - perpX*halfWidth)
// 	edge2y1 := float32(y1 - perpY*halfWidth)
// 	edge2x2 := float32(x2 - perpX*halfWidth)
// 	edge2y2 := float32(y2 - perpY*halfWidth)

// 	vector.StrokeLine(screen, edge1x1, edge1y1, edge1x2, edge1y2, 1.5, colorRoadEdge, true)
// 	vector.StrokeLine(screen, edge2x1, edge2y1, edge2x2, edge2y2, 1.5, colorRoadEdge, true)
// }

func createWhiteImage() *ebiten.Image {
	img := ebiten.NewImage(1, 1)
	img.Fill(color.White)
	return img
}
