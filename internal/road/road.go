package road

import "math"

type Node struct{
	ID string
	X float64
	Y float64
}

type Road struct {
	ID string
	From *Node
	To *Node
	MaxSpeed float64
	Length float64
	Width float64
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
		Width: 8.0,
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
    x := r.From.X + t*(r.To.X-r.From.X)
    y := r.From.Y + t*(r.To.Y-r.From.Y)
    return x, y
}