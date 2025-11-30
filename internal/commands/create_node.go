package commands

import (
	"traffic-sim/internal/road"
	"traffic-sim/internal/world"
)

type CreateNodeCommand struct {
	X, Y   float64
	NodeID string
}

func (c *CreateNodeCommand) Execute(w *world.World) error {
	w.Mu.Lock()
	defer w.Mu.Unlock()

	newNode := &road.Node{
		ID: c.NodeID,
		X:  c.X,
		Y:  c.Y,
	}

	w.Nodes = append(w.Nodes, newNode)
	w.CreateIntersection(c.NodeID)

	return nil
}