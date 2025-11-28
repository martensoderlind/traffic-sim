package commands

import (
	"fmt"
	"traffic-sim/internal/road"
	"traffic-sim/internal/world"
)

type Command interface {
	Execute(w *world.World) error
}

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

type CreateRoadCommand struct {
	From, To *road.Node
	MaxSpeed float64
}

func (c *CreateRoadCommand) Execute(w *world.World) error {
	w.Mu.Lock()
	defer w.Mu.Unlock()

	roadID := fmt.Sprintf("%s-%s", c.From.ID, c.To.ID)
	newRoad := road.NewRoad(roadID, c.From, c.To, c.MaxSpeed)

	w.Roads = append(w.Roads, newRoad)

	fromIntersection := w.IntersectionsByNode[c.From.ID]
	toIntersection := w.IntersectionsByNode[c.To.ID]

	if fromIntersection != nil {
		fromIntersection.AddOutgoing(newRoad)
	}

	if toIntersection != nil {
		toIntersection.AddIncoming(newRoad)
	}

	return nil
}

type CommandExecutor struct {
	world *world.World
}

func NewCommandExecutor(w *world.World) *CommandExecutor {
	return &CommandExecutor{world: w}
}

func (e *CommandExecutor) Execute(cmd Command) error {
	return cmd.Execute(e.world)
}