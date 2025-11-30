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