package commands

import (
	"traffic-sim/internal/road"
	"traffic-sim/internal/world"
)

type MoveNodeCommand struct {
	Node *road.Node
	NewX float64
	NewY float64
}

func (c *MoveNodeCommand) Execute(w *world.World) error {
	w.Mu.Lock()
	defer w.Mu.Unlock()

	c.Node.X = c.NewX
	c.Node.Y = c.NewY

	for _, rd := range w.Roads {
		if rd.From == c.Node || rd.To == c.Node {
			rd.UpdateLength()
		}
	}
	
	for _, v := range w.Vehicles {
		if v.Road.From == c.Node || v.Road.To == c.Node {
			x, y := v.Road.PosAt(v.Distance)
			v.Pos.X = x
			v.Pos.Y = y
		}
	}

	return nil
}