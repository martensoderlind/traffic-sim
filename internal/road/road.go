package road

import "math"

type Node struct{
	ID string
	X float64
	Y float64
}

type RoadCurve struct {
	ControlP1 Point 
	ControlP2 Point 
}

type Point struct {
	X, Y float64
}

type Road struct {
	ID string
	From *Node
	To *Node
	MaxSpeed float64
	Length float64
	Width float64
	ReverseRoad *Road
	Curve *RoadCurve
}

func NewRoad(id string, from, to *Node, maxSpeed float64) *Road{
	dx := from.X-to.X
	dy := from.Y-to.Y
	length := math.Hypot(dx, dy) 

	return &Road{
		ID:id,
		From: from,
		To: to,
		MaxSpeed: maxSpeed,
		Length: length,
		Width: 12.0,
	}
}

func (r *Road) UpdateLength() {
	dx := r.From.X - r.To.X
	dy := r.From.Y - r.To.Y
	r.Length = math.Hypot(dx, dy)
}

func (r *Road) PosAt(dist float64) (float64, float64) {
    if r.Length == 0 {
        return r.From.X, r.From.Y
    }
    t := dist / r.Length
    if t > 1 {
        t = 1
    }
    if t < 0 {
        t = 0
    }
    
    var x, y float64
    
    if r.Curve != nil {
        p0 := Point{X: r.From.X, Y: r.From.Y}
        p1 := r.Curve.ControlP1
        p2 := r.Curve.ControlP2
        p3 := Point{X: r.To.X, Y: r.To.Y}
        
        pt := cubicBezierPoint(p0, p1, p2, p3, t)
        x = pt.X
        y = pt.Y
    } else {
        x = r.From.X + t*(r.To.X-r.From.X)
        y = r.From.Y + t*(r.To.Y-r.From.Y)
    }
    
    if r.ReverseRoad != nil {
        dx := r.To.X - r.From.X
        dy := r.To.Y - r.From.Y
        length := r.Length
        
        if length > 0 {
            perpX := -dy / length
            perpY := dx / length
            
            offset := r.Width * 0.5
            
            x += perpX * offset
            y += perpY * offset
        }
    }
    
    return x, y
}

func cubicBezierPoint(p0, p1, p2, p3 Point, t float64) Point {
	mt := 1 - t
	mt2 := mt * mt
	mt3 := mt2 * mt
	t2 := t * t
	t3 := t2 * t

	return Point{
		X: mt3*p0.X + 3*mt2*t*p1.X + 3*mt*t2*p2.X + t3*p3.X,
		Y: mt3*p0.Y + 3*mt2*t*p1.Y + 3*mt*t2*p2.Y + t3*p3.Y,
	}
}