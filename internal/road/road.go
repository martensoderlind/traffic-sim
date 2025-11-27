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