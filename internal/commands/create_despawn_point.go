package commands

import (
	"traffic-sim/internal/road"
	"traffic-sim/internal/world"
)

type CreateDespawnPointCommand struct {
	DespawnID string
	Node      *road.Node
	Road      *road.Road
}

func (c *CreateDespawnPointCommand) ExecuteUnlocked(w *world.World) error {
	despawnPoint := road.NewDespawnPoint(c.DespawnID, c.Node, c.Road)
	w.DespawnPoints = append(w.DespawnPoints, despawnPoint)

	return nil
}

func (c *CreateDespawnPointCommand) Execute(w *world.World) error {
    return nil
}