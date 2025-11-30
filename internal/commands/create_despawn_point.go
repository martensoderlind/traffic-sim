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

func (c *CreateDespawnPointCommand) Execute(w *world.World) error {
	w.Mu.Lock()
	defer w.Mu.Unlock()

	despawnPoint := road.NewDespawnPoint(c.DespawnID, c.Node, c.Road)
	w.DespawnPoints = append(w.DespawnPoints, despawnPoint)

	return nil
}