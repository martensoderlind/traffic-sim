package commands

import (
	"traffic-sim/internal/road"
	"traffic-sim/internal/world"
)

type UpdateRoadPropertiesCommand struct {
	Road     *road.Road
	MaxSpeed float64
	Width    float64
}

func (c *UpdateRoadPropertiesCommand) Execute(w *world.World) error {
	w.Mu.Lock()
	defer w.Mu.Unlock()

	if c.MaxSpeed > 0 {
		c.Road.MaxSpeed = c.MaxSpeed
	}

	if c.Width > 0 {
		c.Road.Width = c.Width
	}

	return nil
}