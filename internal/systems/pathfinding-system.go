package systems

import (
	"math/rand"
	"traffic-sim/internal/road"
	"traffic-sim/internal/vehicle"
	"traffic-sim/internal/world"
)

type PathfindingSystem struct{}

func NewPathfindingSystem() *PathfindingSystem {
	return &PathfindingSystem{}
}

func (ps *PathfindingSystem) Update(w *world.World, dt float64) {
	w.Mu.Lock()
	defer w.Mu.Unlock()

	for _, v := range w.Vehicles {
		if v.Distance >= v.Road.Length && v.Speed > 0 {
			nextRoad := ps.findNextRoad(w, v)
			
			if nextRoad != nil {
				v.Road = nextRoad
				v.Distance = 0
			} else {
				v.Speed = 0
			}
		}
	}
}

func (ps *PathfindingSystem) findNextRoad(w *world.World, v *vehicle.Vehicle) *road.Road {
	intersection := w.IntersectionsByNode[v.Road.To.ID]
	if intersection == nil || len(intersection.Outgoing) == 0 {
		return nil
	}

	available := make([]*road.Road, 0, len(intersection.Outgoing))
	for _, r := range intersection.Outgoing {
		if notSameRoad(r,v.Road){
			available = append(available, r)
		}
	}

	if len(available) == 0 {
		return nil
	}

	return available[rand.Intn(len(available))]
}

func notSameRoad(r1, r2 *road.Road) bool {
    return !(r1.From == r2.To && r1.To == r2.From)
}