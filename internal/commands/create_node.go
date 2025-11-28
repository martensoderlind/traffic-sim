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

	newIntersection := road.NewIntersection(c.NodeID)
	w.Intersections = append(w.Intersections, newIntersection)
	w.IntersectionsByNode[c.NodeID] = newIntersection

	return nil
}