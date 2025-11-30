package commands

import (
	"traffic-sim/internal/road"
	"traffic-sim/internal/world"
)

type CreateSpawnPointCommand struct {
	SpawnID string
	Node    *road.Node
	Road    *road.Road
}

func (c *CreateSpawnPointCommand) Execute(w *world.World) error {
	w.Mu.Lock()
	defer w.Mu.Unlock()

	spawnPoint := road.NewSpawnPoint(c.SpawnID, c.Node, c.Road)
	w.SpawnPoints = append(w.SpawnPoints, spawnPoint)

	return nil
}