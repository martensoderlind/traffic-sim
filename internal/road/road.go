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
	}
}

func (r *Road) UpdateLength() {
	dx := r.From.X - r.To.X
	dy := r.From.Y - r.To.Y
	r.Length = math.Hypot(dx, dy)
}

func (r *Road) PosAt(dist float64) (float64, float64) {
    t := dist / r.Length
    x := r.From.X + t*(r.To.X-r.From.X)
    y := r.From.Y + t*(r.To.Y-r.From.Y)
    return x, y
}