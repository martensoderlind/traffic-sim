package commands

import (
	"traffic-sim/internal/road"
	"traffic-sim/internal/world"
)

type UpdateSpawnPointPropertiesCommand struct {
	SpawnPoint     *road.SpawnPoint
	Interval      float64
	MinSpeed      float64
	MaxSpeed      float64
	MaxVehicles   int
	Enabled       bool
	VehicleCounter int
}
func (c *UpdateSpawnPointPropertiesCommand) Execute(w *world.World) error {
	w.Mu.Lock()
	defer w.Mu.Unlock()

	if c.MaxSpeed > 0 {
		c.SpawnPoint.MaxSpeed = c.MaxSpeed
	}
	if c.MinSpeed >= 0 {
		c.SpawnPoint.MinSpeed = c.MinSpeed
	}
	if c.Interval > 0 {
		c.SpawnPoint.Interval = c.Interval
	}
	if c.MaxVehicles > 0 {
		c.SpawnPoint.MaxVehicles = c.MaxVehicles
	}
	c.SpawnPoint.Enabled = c.Enabled

	return nil
}