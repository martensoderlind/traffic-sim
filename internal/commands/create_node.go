package commands

import (
	"traffic-sim/internal/road"
	"traffic-sim/internal/world"
)

type CreateNodeCommand struct {
	X, Y   float64
	NodeID string
}

func (c *CreateNodeCommand) ExecuteUnlocked(w *world.World) error {
	newNode := &road.Node{
		ID: c.NodeID,
		X:  c.X,
		Y:  c.Y,
	}

	w.Nodes = append(w.Nodes, newNode)
	w.CreateIntersection(c.NodeID)

	return nil
}

// Execute satisfies Command interface (not called when ExecuteWithLocking is detected)
func (c *CreateNodeCommand) Execute(w *world.World) error {
	return nil
}