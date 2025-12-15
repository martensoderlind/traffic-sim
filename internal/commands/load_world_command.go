package commands

import (
	"fmt"
	"traffic-sim/internal/events"
	"traffic-sim/internal/persistence"
	"traffic-sim/internal/world"
)

type LoadWorldCommand struct {
	OnWorldLoaded func(*world.World)
}

func (c *LoadWorldCommand) Execute(w *world.World) error {
	saveData, err := persistence.LoadFromFile()
	if err != nil {
		return fmt.Errorf("failed to load world: %w", err)
	}

	if saveData == nil {
		return nil
	}

	newWorld, err := persistence.DeserializeWorld(saveData)
	if err != nil {
		return fmt.Errorf("failed to deserialize world: %w", err)
	}

	if w != nil && w.Events != nil {
		w.Events.Emit(events.EventWorldLoaded, events.WorldLoadedEvent{World: newWorld})
	}

	if c.OnWorldLoaded != nil {
		c.OnWorldLoaded(newWorld)
	}

	return nil
}