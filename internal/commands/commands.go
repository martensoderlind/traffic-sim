package commands

import (
	"fmt"
	"traffic-sim/internal/road"
	"traffic-sim/internal/sim"
)

type Command interface {
	Execute(world *sim.World) error
}

type CreateNodeCommand struct {
	X, Y   float64
	NodeID string
}

type CreateRoadCommand struct {
	From, To *road.Node
	MaxSpeed float64
}

func (c *CreateNodeCommand) Execute(world *sim.World) error {
	world.Mu.Lock()
	defer world.Mu.Unlock()

	newNode := &road.Node{
		ID: c.NodeID,
		X:  c.X,
		Y:  c.Y,
	}

	world.Nodes = append(world.Nodes, newNode)

	newIntersection := road.NewIntersection(c.NodeID)
	world.Intersections = append(world.Intersections, newIntersection)
	world.IntersectionsByNode[c.NodeID] = newIntersection

	return nil
}

func (c *CreateRoadCommand) Execute(world *sim.World) error {
	world.Mu.Lock()
	defer world.Mu.Unlock()

	roadID := fmt.Sprintf("%s-%s", c.From.ID, c.To.ID)
	newRoad := road.NewRoad(roadID, c.From, c.To, c.MaxSpeed)

	world.Roads = append(world.Roads, newRoad)

	fromIntersection := world.IntersectionsByNode[c.From.ID]
	toIntersection := world.IntersectionsByNode[c.To.ID]

	if fromIntersection != nil {
		fromIntersection.AddOutgoing(newRoad)
	}

	if toIntersection != nil {
		toIntersection.AddIncoming(newRoad)
	}

	return nil
}

type CommandExecutor struct {
	world *sim.World
}

func NewCommandExecutor(world *sim.World) *CommandExecutor {
	return &CommandExecutor{world: world}
}

func (e *CommandExecutor) Execute(cmd Command) error {
	return cmd.Execute(e.world)
}