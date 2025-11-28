package systems

import (
	"traffic-sim/internal/vehicle"
	"traffic-sim/internal/world"
)

type CollisionSystem struct {
	safeDistance     float64
	emergencyBrake   float64
	anticipationDist float64
}

func NewCollisionSystem() *CollisionSystem {
	return &CollisionSystem{
		safeDistance:     15.0,
		emergencyBrake:   30.0,
		anticipationDist: 40.0,
	}
}

func (cs *CollisionSystem) Update(w *world.World, dt float64) {
	w.Mu.Lock()
	defer w.Mu.Unlock()

	roadVehicles := cs.groupVehiclesByRoad(w)

	for _, v := range w.Vehicles {
		vehiclesAhead := roadVehicles[v.Road.ID]
		
		nearestAhead := cs.findNearestVehicleAhead(v, vehiclesAhead)
		if nearestAhead == nil {
			continue
		}

		distToVehicle := nearestAhead.Distance - v.Distance
		
		if distToVehicle < cs.safeDistance {
			v.Speed = 0
			continue
		}

		if distToVehicle < cs.anticipationDist {
			targetSpeed := cs.calculateSafeSpeed(distToVehicle, nearestAhead.Speed)
			cs.adjustSpeedForCollision(v, targetSpeed, dt)
		}
	}
}

func (cs *CollisionSystem) groupVehiclesByRoad(w *world.World) map[string][]*vehicle.Vehicle {
	roadVehicles := make(map[string][]*vehicle.Vehicle)
	
	for _, v := range w.Vehicles {
		roadVehicles[v.Road.ID] = append(roadVehicles[v.Road.ID], v)
	}
	
	return roadVehicles
}

func (cs *CollisionSystem) findNearestVehicleAhead(current *vehicle.Vehicle, vehicles []*vehicle.Vehicle) *vehicle.Vehicle {
	var nearest *vehicle.Vehicle
	minDist := current.Road.Length

	for _, v := range vehicles {
		if v == current {
			continue
		}

		if v.Distance > current.Distance {
			dist := v.Distance - current.Distance
			if dist < minDist {
				minDist = dist
				nearest = v
			}
		}
	}

	return nearest
}

func (cs *CollisionSystem) calculateSafeSpeed(distToVehicle, vehicleAheadSpeed float64) float64 {
	ratio := distToVehicle / cs.anticipationDist
	
	if ratio < 0.5 {
		return vehicleAheadSpeed * 0.8
	}
	
	return vehicleAheadSpeed * (0.8 + ratio*0.4)
}

func (cs *CollisionSystem) adjustSpeedForCollision(v *vehicle.Vehicle, targetSpeed, dt float64) {
	if v.Speed > targetSpeed {
		v.Speed -= cs.emergencyBrake * dt
		if v.Speed < targetSpeed {
			v.Speed = targetSpeed
		}
		if v.Speed < 0 {
			v.Speed = 0
		}
	}
}