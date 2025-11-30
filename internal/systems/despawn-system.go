package systems

import (
	"traffic-sim/internal/vehicle"
	"traffic-sim/internal/world"
)

type DespawnSystem struct{}

func NewDespawnSystem() *DespawnSystem {
	return &DespawnSystem{}
}

func (ds *DespawnSystem) Update(w *world.World, dt float64) {
	w.Mu.Lock()
	defer w.Mu.Unlock()

	despawnRoads := ds.buildDespawnRoadSet(w)

	toRemove := make(map[int]bool)

	for i, v := range w.Vehicles {
		if v.Distance >= v.Road.Length {
			if despawnRoads[v.Road.ID] {
				toRemove[i] = true
			}
		}
	}

	if len(toRemove) > 0 {
		newVehicles := make([]*vehicle.Vehicle, 0, len(w.Vehicles)-len(toRemove))
		for i, v := range w.Vehicles {
			if !toRemove[i] {
				newVehicles = append(newVehicles, v)
			}
		}
		w.Vehicles = newVehicles
	}
}

func (ds *DespawnSystem) buildDespawnRoadSet(w *world.World) map[string]bool {
	despawnRoads := make(map[string]bool)
	
	for _, dp := range w.DespawnPoints {
		if dp.Enabled {
			despawnRoads[dp.Road.ID] = true
		}
	}
	
	return despawnRoads
}