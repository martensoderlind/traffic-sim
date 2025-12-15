package commands

import (
	"traffic-sim/internal/events"
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

	if w.Events != nil {
		w.Events.Emit(events.EventNodeCreated, events.NodeCreatedEvent{Node: newNode})
	}

	return nil
}

func (c *CreateNodeCommand) Execute(w *world.World) error {
	return nil
}