package commands

import (
	"traffic-sim/internal/events"
	"traffic-sim/internal/road"
	"traffic-sim/internal/world"
)

type CreateTrafficLightCommand struct {
	LightID string
	Node    *road.Node
	Roads   []*road.Road
}

func (c *CreateTrafficLightCommand) ExecuteUnlocked(w *world.World) error {
	intersection := w.IntersectionsByNode[c.Node.ID]
	if intersection == nil {
		return nil
	}

	existingLights := 0
	for _, existingLight := range w.TrafficLights {
		if existingLight.Intersection.ID == intersection.ID {
			existingLights++
		}
	}

	startGreen := existingLights%2 == 0

	light := road.NewTrafficLight(c.LightID, intersection, startGreen)
	
	for _, rd := range c.Roads {
		light.AddControlledRoad(rd)
	}

	w.TrafficLights = append(w.TrafficLights, light)

	if w.Events != nil {
		w.Events.Emit(events.EventTrafficLightCreated, events.TrafficLightCreatedEvent{
			TrafficLight: light,
		})
	}

	return nil
}

func (c *CreateTrafficLightCommand) Execute(w *world.World) error {
    return nil
}