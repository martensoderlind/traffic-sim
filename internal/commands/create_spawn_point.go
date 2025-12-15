package commands

import (
	"traffic-sim/internal/events"
	"traffic-sim/internal/road"
	"traffic-sim/internal/world"
)

type CreateSpawnPointCommand struct {
	SpawnID string
	Node    *road.Node
	Road    *road.Road
}

func (c *CreateSpawnPointCommand) ExecuteUnlocked(w *world.World) error {
	spawnPoint := road.NewSpawnPoint(c.SpawnID, c.Node, c.Road)
	w.SpawnPoints = append(w.SpawnPoints, spawnPoint)

	if w.Events != nil {
		w.Events.Emit(events.EventSpawnPointCreated, events.SpawnPointCreatedEvent{SpawnPoint: spawnPoint})
	}

	return nil
}

func (c *CreateSpawnPointCommand) Execute(w *world.World) error {
    return nil
}