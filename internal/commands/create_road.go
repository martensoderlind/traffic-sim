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
}

func (c *CreateRoadCommand) Execute(w *world.World) error {
	w.Mu.Lock()
	defer w.Mu.Unlock()

	roadID := fmt.Sprintf("%s-%s", c.From.ID, c.To.ID)
	newRoad := road.NewRoad(roadID, c.From, c.To, c.MaxSpeed)
	
	if c.Width > 0 {
		newRoad.Width = c.Width
	}

	if c.From == c.To {
		offsetLen := newRoad.Width * 0.75
		if offsetLen < 6.0 {
			offsetLen = 6.0
		}

		newRoad.StartOffset = road.Point{X: offsetLen, Y: 0}
		newRoad.EndOffset = road.Point{X: -offsetLen, Y: 0}
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