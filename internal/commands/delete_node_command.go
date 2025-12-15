package commands

import (
	"traffic-sim/internal/events"
	"traffic-sim/internal/road"
	"traffic-sim/internal/world"
)

type DeleteNodeCommand struct {
	Node *road.Node
}

func (c *DeleteNodeCommand) ExecuteUnlocked(w *world.World) error {
	nodeID := c.Node.ID
	intersection := w.IntersectionsByNode[c.Node.ID]
	if intersection != nil {
		roadsToDelete := make([]*road.Road, 0)
		roadsToDelete = append(roadsToDelete, intersection.Incoming...)
		roadsToDelete = append(roadsToDelete, intersection.Outgoing...)

		for _, r := range roadsToDelete {
			deleteRoadCmd := &DeleteRoadCommand{Road: r}
			deleteRoadCmd.ExecuteUnlocked(w)
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

	if w.Events != nil {
		w.Events.Emit(events.EventNodeDeleted, events.NodeDeletedEvent{NodeID: nodeID})
	}

	return nil
}

func (c *DeleteNodeCommand) Execute(w *world.World) error {
    return nil
}