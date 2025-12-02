package commands

import (
	"traffic-sim/internal/road"
	"traffic-sim/internal/world"
)

type CreateTrafficLightCommand struct {
	LightID string
	Node    *road.Node
	Road    *road.Road
}

func (c *CreateTrafficLightCommand) Execute(w *world.World) error {
	w.Mu.Lock()
	defer w.Mu.Unlock()

	intersection := w.IntersectionsByNode[c.Node.ID]
	if intersection == nil {
		return nil
	}

	light := road.NewTrafficLight(c.LightID, intersection)
	light.AddControlledRoad(c.Road)

	w.TrafficLights = append(w.TrafficLights, light)

	return nil
}