package systems

import (
	"math/rand"
	"traffic-sim/internal/road"
	"traffic-sim/internal/vehicle"
	"traffic-sim/internal/world"
)

type System interface {
	Update(w *world.World, dt float64)
}

type SystemManager struct {
	systems []System
}

func NewSystemManager() *SystemManager {
	return &SystemManager{
		systems: make([]System, 0),
	}
}

func (sm *SystemManager) AddSystem(s System) {
	sm.systems = append(sm.systems, s)
}

func (sm *SystemManager) Update(w *world.World, dt float64) {
	for _, system := range sm.systems {
		system.Update(w, dt)
	}
}

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
		if r.ID != v.Road.ID {
			available = append(available, r)
		}
	}

	if len(available) == 0 {
		return nil
	}

	return available[rand.Intn(len(available))]
}