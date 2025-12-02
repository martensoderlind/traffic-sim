package world

import (
	"sync"

	"traffic-sim/internal/road"
	"traffic-sim/internal/vehicle"
)

type World struct {
	Roads         []*road.Road
	Nodes         []*road.Node
	Vehicles      []*vehicle.Vehicle
	Intersections []*road.Intersection
	SpawnPoints   []*road.SpawnPoint
	DespawnPoints []*road.DespawnPoint
	TrafficLights []*road.TrafficLight 

	IntersectionsByNode map[string]*road.Intersection

	Mu sync.RWMutex
}

func New() *World {
	w := &World{
		Roads:               make([]*road.Road,0),      
		Nodes:               make([]*road.Node,0),    
		Vehicles:            make([]*vehicle.Vehicle,0),
		SpawnPoints:         make([]*road.SpawnPoint, 0),
		DespawnPoints:       make([]*road.DespawnPoint, 0),
		IntersectionsByNode: make(map[string]*road.Intersection),
	}

	// w.Intersections = road.BuildIntersections(roads, nodes)
	// for _, i := range w.Intersections {
	// 	w.IntersectionsByNode[i.ID] = i
	// }

	return w
}

func (w *World) GetIntersection(nodeID string) *road.Intersection {
	return w.IntersectionsByNode[nodeID]
}

func (w *World) CreateIntersection(nodeID string) *road.Intersection {
	intersection := road.NewIntersection(nodeID)
	w.Intersections = append(w.Intersections, intersection)
	w.IntersectionsByNode[nodeID] = intersection
	return intersection
}

func (w *World) DeleteIntersection(nodeID string) {
	delete(w.IntersectionsByNode, nodeID)
	
	for i, inter := range w.Intersections {
		if inter.ID == nodeID {
			w.Intersections = append(w.Intersections[:i], w.Intersections[i+1:]...)
			break
		}
	}
}

func (w *World) AddRoadToIntersections(rd *road.Road) {
	fromIntersection := w.GetIntersection(rd.From.ID)
	if fromIntersection != nil {
		fromIntersection.AddOutgoing(rd)
	}

	toIntersection := w.GetIntersection(rd.To.ID)
	if toIntersection != nil {
		toIntersection.AddIncoming(rd)
	}
}

func (w *World) RemoveRoadFromIntersections(rd *road.Road) {
	fromIntersection := w.GetIntersection(rd.From.ID)
	if fromIntersection != nil {
		for i, r := range fromIntersection.Outgoing {
			if r == rd {
				fromIntersection.Outgoing = append(fromIntersection.Outgoing[:i], fromIntersection.Outgoing[i+1:]...)
				break
			}
		}
	}

	toIntersection := w.GetIntersection(rd.To.ID)
	if toIntersection != nil {
		for i, r := range toIntersection.Incoming {
			if r == rd {
				toIntersection.Incoming = append(toIntersection.Incoming[:i], toIntersection.Incoming[i+1:]...)
				break
			}
		}
	}
}