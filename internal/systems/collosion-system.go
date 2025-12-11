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
		emergencyBrake:   20.0,
		anticipationDist: 30.0,
	}
}

func (cs *CollisionSystem) Update(w *world.World, dt float64) {
	w.Mu.Lock()
	defer w.Mu.Unlock()

	roadVehicles := cs.groupVehiclesByRoad(w)

	for _, v := range w.Vehicles {
		if v.InTransition {
			cs.applyCollisionAvoidanceForTransition(v, w, roadVehicles, dt)
			continue
		}
		
		vehiclesAhead := roadVehicles[v.Road.ID]
		
		nearestAhead := cs.findNearestVehicleAhead(v, vehiclesAhead)
		
		var distToVehicle float64
		if nearestAhead != nil {
			distToVehicle = nearestAhead.Distance - v.Distance
		} else {
			distToVehicle = v.Road.Length
		}

		if v.NextRoad != nil {
			nextRoadVehicles := roadVehicles[v.NextRoad.ID]
			nearestOnNext := cs.findNearestVehicleFromStart(nextRoadVehicles)
			
			if nearestOnNext != nil {
				distToNextRoadVehicle := (v.Road.Length - v.Distance) + nearestOnNext.Distance
				if nearestAhead == nil || distToNextRoadVehicle < distToVehicle {
					distToVehicle = distToNextRoadVehicle
				}
			}
		}
		
		if distToVehicle < cs.safeDistance {
			v.Speed = 0
			continue
		}

		if distToVehicle < cs.anticipationDist {
			targetSpeed := cs.calculateSafeSpeed(distToVehicle, 0)
			if nearestAhead != nil && nearestAhead.Distance-v.Distance < cs.anticipationDist {
				targetSpeed = cs.calculateSafeSpeed(nearestAhead.Distance-v.Distance, nearestAhead.Speed)
			}
			cs.adjustSpeedForCollision(v, targetSpeed, dt)
		}
	}
}

func (cs *CollisionSystem) groupVehiclesByRoad(w *world.World) map[string][]*vehicle.Vehicle {
	roadVehicles := make(map[string][]*vehicle.Vehicle)
	
	for _, v := range w.Vehicles {
		if !v.InTransition {
			roadVehicles[v.Road.ID] = append(roadVehicles[v.Road.ID], v)
		}
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

func (cs *CollisionSystem) findNearestVehicleFromStart(vehicles []*vehicle.Vehicle) *vehicle.Vehicle {
	var nearest *vehicle.Vehicle
	minDist := 1000000.0

	for _, v := range vehicles {
		if v.Distance < minDist {
			minDist = v.Distance
			nearest = v
		}
	}

	return nearest
}

func (cs *CollisionSystem) findSafeStartDistanceOnRoad(vehicles []*vehicle.Vehicle, roadLength float64, desiredStartDist float64) float64 {
	sorted := make([]*vehicle.Vehicle, 0, len(vehicles))
	for _, v := range vehicles {
		sorted = append(sorted, v)
	}
	
	for i := 1; i < len(sorted); i++ {
		for j := i; j > 0 && sorted[j].Distance < sorted[j-1].Distance; j-- {
			sorted[j], sorted[j-1] = sorted[j-1], sorted[j]
		}
	}
	
	if len(sorted) == 0 {
		return desiredStartDist
	}
	
	lastVehicleDist := sorted[len(sorted)-1].Distance
	safeStartDist := lastVehicleDist + cs.safeDistance + 10.0
	
	if safeStartDist < desiredStartDist {
		safeStartDist = desiredStartDist
	}
	
	return safeStartDist
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

func (cs *CollisionSystem) applyCollisionAvoidanceForTransition(v *vehicle.Vehicle, w *world.World, roadVehicles map[string][]*vehicle.Vehicle, dt float64) {
	if v.NextRoad == nil {
		return
	}
	
	nextRoadVehicles := roadVehicles[v.NextRoad.ID]
	
	var nearestOnNext *vehicle.Vehicle
	var minDist float64 = 1000000.0
	
	for _, other := range nextRoadVehicles {
		if other.ID == v.ID {
			continue
		}
		dist := other.Distance
		if dist < minDist {
			minDist = dist
			nearestOnNext = other
		}
	}
	
	if nearestOnNext != nil {
		transitionEndDist := 25.0
		if v.NextRoad.Length < 50.0 {
			transitionEndDist = v.NextRoad.Length * 0.3
		}
		
		distToVehicle := nearestOnNext.Distance - transitionEndDist
		
		if distToVehicle < cs.safeDistance+15.0 {
			targetSpeed := cs.calculateSafeSpeed(distToVehicle, nearestOnNext.Speed)
			cs.adjustSpeedForCollision(v, targetSpeed, dt)
		}
	}
	
	currentRoadVehicles := roadVehicles[v.Road.ID]
	for _, other := range currentRoadVehicles {
		if other.ID == v.ID {
			continue
		}
		
		if other.Distance > v.Distance {
			distToVehicle := other.Distance - v.Distance
			if distToVehicle < cs.anticipationDist {
				targetSpeed := cs.calculateSafeSpeed(distToVehicle, other.Speed)
				cs.adjustSpeedForCollision(v, targetSpeed, dt)
			}
		}
	}
}