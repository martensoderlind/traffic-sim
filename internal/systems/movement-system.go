package systems

import (
	"traffic-sim/internal/road"
	"traffic-sim/internal/vehicle"
	"traffic-sim/internal/world"
)

type MovementSystem struct {
	lookAheadDist float64
	deceleration  float64
	acceleration  float64
	minSpeed      float64
}

func NewMovementSystem() *MovementSystem {
	return &MovementSystem{
		lookAheadDist: 50.0,
		deceleration:  15.0,
		acceleration:  10.0,
		minSpeed:      5.0,
	}
}

func (ms *MovementSystem) Update(w *world.World, dt float64) {
	w.Mu.Lock()
	defer w.Mu.Unlock()

	for _, v := range w.Vehicles {
		targetSpeed := ms.calculateTargetSpeed(w, v)
		
		ms.adjustSpeed(v, targetSpeed, dt)
		
		newDist := v.Distance + v.Speed*dt
		
		if newDist > v.Road.Length {
			newDist = v.Road.Length
		}
		
		v.Distance = newDist

		x, y := v.Road.PosAt(v.Distance)
		v.Pos.X = x
		v.Pos.Y = y
	}
}

func (ms *MovementSystem) calculateTargetSpeed(w *world.World, v *vehicle.Vehicle) float64 {
	distToEnd := v.Road.Length - v.Distance
	
	if distToEnd > ms.lookAheadDist {
		return v.Road.MaxSpeed
	}
	
	if ms.hasDespawnPoint(w, v.Road) {
		return v.Road.MaxSpeed
	}
	
	intersection := w.IntersectionsByNode[v.Road.To.ID]
	if intersection == nil || len(intersection.Outgoing) == 0 {
		ratio := distToEnd / ms.lookAheadDist
		return ms.minSpeed + (v.Road.MaxSpeed-ms.minSpeed)*ratio
	}
	
	hasValidExit := false
	for _, r := range intersection.Outgoing {
		if !(r.From == v.Road.To && r.To == v.Road.From) {
			hasValidExit = true
			break
		}
	}
	
	if !hasValidExit {
		ratio := distToEnd / ms.lookAheadDist
		return ms.minSpeed + (v.Road.MaxSpeed-ms.minSpeed)*ratio
	}
	
	return v.Road.MaxSpeed
}

func (ms *MovementSystem) hasDespawnPoint(w *world.World, rd *road.Road) bool {
	for _, dp := range w.DespawnPoints {
		if dp.Enabled && dp.Road == rd {
			return true
		}
	}
	return false
}

func (ms *MovementSystem) adjustSpeed(v *vehicle.Vehicle, targetSpeed, dt float64) {
	if v.Speed < targetSpeed {
		v.Speed += ms.acceleration * dt
		if v.Speed > targetSpeed {
			v.Speed = targetSpeed
		}
	} else if v.Speed > targetSpeed {
		v.Speed -= ms.deceleration * dt
		if v.Speed < targetSpeed {
			v.Speed = targetSpeed
		}
		if v.Speed < ms.minSpeed {
			v.Speed = ms.minSpeed
		}
	}
}