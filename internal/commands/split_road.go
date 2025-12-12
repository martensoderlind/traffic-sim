package commands

import (
	"fmt"
	"traffic-sim/internal/road"
	"traffic-sim/internal/world"
)

type SplitRoadCommand struct {
	Road   *road.Road
	X, Y   float64
	NodeID string
}

func (c *SplitRoadCommand) ExecuteUnlocked(w *world.World) error {
	splitNode := &road.Node{
		ID: c.NodeID,
		X:  c.X,
		Y:  c.Y,
	}

	w.Nodes = append(w.Nodes, splitNode)
	newIntersection := w.CreateIntersection(c.NodeID)

	road1ID := fmt.Sprintf("%s-%s", c.Road.From.ID, splitNode.ID)
	road2ID := fmt.Sprintf("%s-%s", splitNode.ID, c.Road.To.ID)

	newRoad1 := road.NewRoad(road1ID, c.Road.From, splitNode, c.Road.MaxSpeed)
	newRoad1.Width = c.Road.Width
	
	newRoad2 := road.NewRoad(road2ID, splitNode, c.Road.To, c.Road.MaxSpeed)
	newRoad2.Width = c.Road.Width

	if c.Road.ReverseRoad != nil {
		c.handleReverseRoad(w, splitNode, newRoad1, newRoad2, newIntersection)
	}

	w.RemoveRoadFromIntersections(c.Road)

	for i, r := range w.Roads {
		if r == c.Road {
			w.Roads = append(w.Roads[:i], w.Roads[i+1:]...)
			break
		}
	}

	w.Roads = append(w.Roads, newRoad1, newRoad2)

	fromIntersection := w.GetIntersection(c.Road.From.ID)
	if fromIntersection != nil {
		fromIntersection.AddOutgoing(newRoad1)
	}

	toIntersection := w.GetIntersection(c.Road.To.ID)
	if toIntersection != nil {
		toIntersection.AddIncoming(newRoad2)
	}

	newIntersection.AddIncoming(newRoad1)
	newIntersection.AddOutgoing(newRoad2)

	c.updateVehiclesOnRoad(w, c.Road, newRoad1, newRoad2)

	return nil
}

func (c *SplitRoadCommand) handleReverseRoad(w *world.World, splitNode *road.Node, newRoad1, newRoad2 *road.Road, newIntersection *road.Intersection) {
	reverseRoad := c.Road.ReverseRoad
	
	reverseRoad1ID := fmt.Sprintf("%s-%s", splitNode.ID, c.Road.From.ID)
	reverseRoad2ID := fmt.Sprintf("%s-%s", c.Road.To.ID, splitNode.ID)
	
	reverseNewRoad1 := road.NewRoad(reverseRoad1ID, splitNode, c.Road.From, reverseRoad.MaxSpeed)
	reverseNewRoad1.Width = reverseRoad.Width
	
	reverseNewRoad2 := road.NewRoad(reverseRoad2ID, c.Road.To, splitNode, reverseRoad.MaxSpeed)
	reverseNewRoad2.Width = reverseRoad.Width
	
	newRoad1.ReverseRoad = reverseNewRoad1
	reverseNewRoad1.ReverseRoad = newRoad1
	
	newRoad2.ReverseRoad = reverseNewRoad2
	reverseNewRoad2.ReverseRoad = newRoad2
	
	w.RemoveRoadFromIntersections(reverseRoad)
	
	for i, r := range w.Roads {
		if r == reverseRoad {
			w.Roads = append(w.Roads[:i], w.Roads[i+1:]...)
			break
		}
	}
	
	w.Roads = append(w.Roads, reverseNewRoad1, reverseNewRoad2)
	
	fromIntersectionRev := w.GetIntersection(reverseRoad.From.ID)
	if fromIntersectionRev != nil {
		fromIntersectionRev.AddOutgoing(reverseNewRoad2)
	}
	
	toIntersectionRev := w.GetIntersection(reverseRoad.To.ID)
	if toIntersectionRev != nil {
		toIntersectionRev.AddIncoming(reverseNewRoad1)
	}
	
	newIntersection.AddIncoming(reverseNewRoad2)
	newIntersection.AddOutgoing(reverseNewRoad1)
	
	c.updateVehiclesOnRoad(w, reverseRoad, reverseNewRoad2, reverseNewRoad1)
}

func (c *SplitRoadCommand) updateVehiclesOnRoad(w *world.World, oldRoad, newRoad1, newRoad2 *road.Road) {
	for _, v := range w.Vehicles {
		if v.Road == oldRoad {
			if v.Distance <= newRoad1.Length {
				v.Road = newRoad1
			} else {
				v.Distance = v.Distance - newRoad1.Length
				v.Road = newRoad2
			}
		}
	}
}

func (c *SplitRoadCommand) Execute(w *world.World) error {
    return nil
}