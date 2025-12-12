package commands

import (
	"fmt"
	"traffic-sim/internal/road"
	"traffic-sim/internal/world"
)

type CreateRoadCommand struct {
	From, To *road.Node
	MaxSpeed float64
	Width    float64
	StartOffset road.Point
	EndOffset   road.Point
}

func (c *CreateRoadCommand) ExecuteUnlocked(w *world.World) error {
	roadID := fmt.Sprintf("%s-%s", c.From.ID, c.To.ID)
	newRoad := road.NewRoad(roadID, c.From, c.To, c.MaxSpeed)
	
	if c.Width > 0 {
		newRoad.Width = c.Width
	}

	if c.StartOffset != (road.Point{}) || c.EndOffset != (road.Point{}) {
		newRoad.StartOffset = c.StartOffset
		newRoad.EndOffset = c.EndOffset
		newRoad.UpdateLength()
	} else if c.From == c.To {
		newRoad.StartOffset, newRoad.EndOffset = road.CalculateLoopRoadOffsets(0, 0, newRoad.Width)
		newRoad.UpdateLength()
	}

	for _, existingRoad := range w.Roads {
		if existingRoad.From == c.To && existingRoad.To == c.From {
			newRoad.ReverseRoad = existingRoad
			existingRoad.ReverseRoad = newRoad
			break
		}
	}

	w.Roads = append(w.Roads, newRoad)

	fromIntersection := w.IntersectionsByNode[c.From.ID]
	toIntersection := w.IntersectionsByNode[c.To.ID]

	if fromIntersection != nil {
		fromIntersection.AddOutgoing(newRoad)
	}

	if toIntersection != nil {
		toIntersection.AddIncoming(newRoad)
	}

	return nil
}

func (c *CreateRoadCommand) Execute(w *world.World) error {
	return nil
}