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

type SplitRoadCommand struct {
	Road   *road.Road
	X, Y   float64
	NodeID string
}

func (c *SplitRoadCommand) Execute(w *world.World) error {
	w.Mu.Lock()
	defer w.Mu.Unlock()

	splitNode := &road.Node{
		ID: c.NodeID,
		X:  c.X,
		Y:  c.Y,
	}

	w.Nodes = append(w.Nodes, splitNode)

	newIntersection := road.NewIntersection(c.NodeID)
	w.Intersections = append(w.Intersections, newIntersection)
	w.IntersectionsByNode[c.NodeID] = newIntersection

	road1ID := fmt.Sprintf("%s-%s", c.Road.From.ID, splitNode.ID)
	road2ID := fmt.Sprintf("%s-%s", splitNode.ID, c.Road.To.ID)

	newRoad1 := road.NewRoad(road1ID, c.Road.From, splitNode, c.Road.MaxSpeed)
	newRoad2 := road.NewRoad(road2ID, splitNode, c.Road.To, c.Road.MaxSpeed)

	for i, r := range w.Roads {
		if r == c.Road {
			w.Roads = append(w.Roads[:i], w.Roads[i+1:]...)
			break
		}
	}

	w.Roads = append(w.Roads, newRoad1, newRoad2)

	fromIntersection := w.IntersectionsByNode[c.Road.From.ID]
	if fromIntersection != nil {
		for i, r := range fromIntersection.Outgoing {
			if r == c.Road {
				fromIntersection.Outgoing = append(fromIntersection.Outgoing[:i], fromIntersection.Outgoing[i+1:]...)
				break
			}
		}
		fromIntersection.AddOutgoing(newRoad1)
	}

	toIntersection := w.IntersectionsByNode[c.Road.To.ID]
	if toIntersection != nil {
		for i, r := range toIntersection.Incoming {
			if r == c.Road {
				toIntersection.Incoming = append(toIntersection.Incoming[:i], toIntersection.Incoming[i+1:]...)
				break
			}
		}
		toIntersection.AddIncoming(newRoad2)
	}

	newIntersection.AddIncoming(newRoad1)
	newIntersection.AddOutgoing(newRoad2)

	c.updateVehiclesOnRoad(w, c.Road, newRoad1, newRoad2)

	return nil
}

func (c *SplitRoadCommand) updateVehiclesOnRoad(w *world.World, oldRoad, newRoad1, newRoad2 *road.Road) {
	for _, v := range w.Vehicles {
		if v.Road == oldRoad {
			if v.Distance <= newRoad1.Length {
				v.Road = newRoad1
			} else {
				v.Distance = v.Distance - newRoad1.Length
				v.Road = newRoad2
			}
		}
	}
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