package commands

import (
	"traffic-sim/internal/events"
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

	if w.Events != nil {
		w.Events.Emit(events.EventDespawnPointCreated, events.DespawnPointCreatedEvent{
			DespawnPoint: despawnPoint,
		})
	}

	return nil
}

func (c *CreateDespawnPointCommand) Execute(w *world.World) error {
    return nil
}