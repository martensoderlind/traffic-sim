package commands

import (
	"traffic-sim/internal/road"
	"traffic-sim/internal/world"
)

type DeleteNodeCommand struct {
	Node *road.Node
}

func (c *DeleteNodeCommand) Execute(w *world.World) error {
	w.Mu.Lock()
	defer w.Mu.Unlock()

	intersection := w.IntersectionsByNode[c.Node.ID]
	if intersection != nil {
		roadsToDelete := make([]*road.Road, 0)
		roadsToDelete = append(roadsToDelete, intersection.Incoming...)
		roadsToDelete = append(roadsToDelete, intersection.Outgoing...)

		for _, r := range roadsToDelete {
			deleteRoadCmd := &DeleteRoadCommand{Road: r}
			w.Mu.Unlock()
			deleteRoadCmd.Execute(w)
			w.Mu.Lock()
		}

		delete(w.IntersectionsByNode, c.Node.ID)

		for i, inter := range w.Intersections {
			if inter.ID == c.Node.ID {
				w.Intersections = append(w.Intersections[:i], w.Intersections[i+1:]...)
				break
			}
		}
	}

	for i, n := range w.Nodes {
		if n == c.Node {
			w.Nodes = append(w.Nodes[:i], w.Nodes[i+1:]...)
			break
		}
	}

	return nil
}