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

func (c *SplitRoadCommand) Execute(w *world.World) error {
	w.Mu.Lock()
	defer w.Mu.Unlock()

	splitNode := &road.Node{
		ID: c.NodeID,
		X:  c.X,
		Y:  c.Y,
	}

	w.Nodes = append(w.Nodes, splitNode)

	newIntersection := road.NewIntersection(c.NodeID)
	w.Intersections = append(w.Intersections, newIntersection)
	w.IntersectionsByNode[c.NodeID] = newIntersection

	road1ID := fmt.Sprintf("%s-%s", c.Road.From.ID, splitNode.ID)
	road2ID := fmt.Sprintf("%s-%s", splitNode.ID, c.Road.To.ID)

	newRoad1 := road.NewRoad(road1ID, c.Road.From, splitNode, c.Road.MaxSpeed)
	newRoad1.Width = c.Road.Width
	
	newRoad2 := road.NewRoad(road2ID, splitNode, c.Road.To, c.Road.MaxSpeed)
	newRoad2.Width = c.Road.Width

	if c.Road.ReverseRoad != nil {
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
		
		for i, r := range w.Roads {
			if r == reverseRoad {
				w.Roads = append(w.Roads[:i], w.Roads[i+1:]...)
				break
			}
		}
		
		w.Roads = append(w.Roads, reverseNewRoad1, reverseNewRoad2)
		
		fromIntersectionRev := w.IntersectionsByNode[reverseRoad.From.ID]
		if fromIntersectionRev != nil {
			for i, r := range fromIntersectionRev.Outgoing {
				if r == reverseRoad {
					fromIntersectionRev.Outgoing = append(fromIntersectionRev.Outgoing[:i], fromIntersectionRev.Outgoing[i+1:]...)
					break
				}
			}
			fromIntersectionRev.AddOutgoing(reverseNewRoad2)
		}
		
		toIntersectionRev := w.IntersectionsByNode[reverseRoad.To.ID]
		if toIntersectionRev != nil {
			for i, r := range toIntersectionRev.Incoming {
				if r == reverseRoad {
					toIntersectionRev.Incoming = append(toIntersectionRev.Incoming[:i], toIntersectionRev.Incoming[i+1:]...)
					break
				}
			}
			toIntersectionRev.AddIncoming(reverseNewRoad1)
		}
		
		newIntersection.AddIncoming(reverseNewRoad2)
		newIntersection.AddOutgoing(reverseNewRoad1)
		
		c.updateVehiclesOnRoad(w, reverseRoad, reverseNewRoad2, reverseNewRoad1)
	}

	for i, r := range w.Roads {
		if r == c.Road {
			w.Roads = append(w.Roads[:i], w.Roads[i+1:]...)
			break
		}
	}

	w.Roads = append(w.Roads, newRoad1, newRoad2)

	fromIntersection := w.IntersectionsByNode[c.Road.From.ID]
	if fromIntersection != nil {
		for i, r := range fromIntersection.Outgoing {
			if r == c.Road {
				fromIntersection.Outgoing = append(fromIntersection.Outgoing[:i], fromIntersection.Outgoing[i+1:]...)
				break
			}
		}
		fromIntersection.AddOutgoing(newRoad1)
	}

	toIntersection := w.IntersectionsByNode[c.Road.To.ID]
	if toIntersection != nil {
		for i, r := range toIntersection.Incoming {
			if r == c.Road {
				toIntersection.Incoming = append(toIntersection.Incoming[:i], toIntersection.Incoming[i+1:]...)
				break
			}
		}
		toIntersection.AddIncoming(newRoad2)
	}

	newIntersection.AddIncoming(newRoad1)
	newIntersection.AddOutgoing(newRoad2)

	c.updateVehiclesOnRoad(w, c.Road, newRoad1, newRoad2)

	return nil
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