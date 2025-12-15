package commands

import (
	"traffic-sim/internal/events"
	"traffic-sim/internal/road"
	"traffic-sim/internal/world"
)

type UpdateRoadPropertiesCommand struct {
	Road     *road.Road
	MaxSpeed float64
	Width    float64
}

func (c *UpdateRoadPropertiesCommand) ExecuteUnlocked(w *world.World) error {
	if c.MaxSpeed > 0 {
		c.Road.MaxSpeed = c.MaxSpeed
	}

	if c.Width > 0 {
		c.Road.Width = c.Width
	}

	if w.Events != nil {
		w.Events.Emit(events.EventRoadPropertiesUpdated, events.RoadPropertiesUpdatedEvent{
			Road:     c.Road,
			MaxSpeed: c.MaxSpeed,
			Width:    c.Width,
		})
	}
	
	return nil
}

func (c *UpdateRoadPropertiesCommand) Execute(w *world.World) error {
    return nil
}