package commands

import (
	"traffic-sim/internal/road"
	"traffic-sim/internal/world"
)

type DeleteRoadCommand struct {
	Road *road.Road
}

func (c *DeleteRoadCommand) Execute(w *world.World) error {
	w.Mu.Lock()
	defer w.Mu.Unlock()

	if c.Road.ReverseRoad != nil {
		c.Road.ReverseRoad.ReverseRoad = nil
	}

	for i := len(w.Vehicles) - 1; i >= 0; i-- {
		if w.Vehicles[i].Road == c.Road {
			w.Vehicles = append(w.Vehicles[:i], w.Vehicles[i+1:]...)
		}
	}

	for i := len(w.SpawnPoints) - 1; i >= 0; i-- {
		if w.SpawnPoints[i].Road == c.Road {
			w.SpawnPoints = append(w.SpawnPoints[:i], w.SpawnPoints[i+1:]...)
		}
	}

	for i := len(w.DespawnPoints) - 1; i >= 0; i-- {
		if w.DespawnPoints[i].Road == c.Road {
			w.DespawnPoints = append(w.DespawnPoints[:i], w.DespawnPoints[i+1:]...)
		}
	}

	w.RemoveRoadFromIntersections(c.Road)

	for i, r := range w.Roads {
		if r == c.Road {
			w.Roads = append(w.Roads[:i], w.Roads[i+1:]...)
			break
		}
	}

	return nil
}