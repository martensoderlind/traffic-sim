package systems

import (
	"math"
	"traffic-sim/internal/road"
	"traffic-sim/internal/vehicle"
	"traffic-sim/internal/world"
)

type MovementSystem struct {
	lookAheadDist float64
	deceleration  float64
	acceleration  float64
	minSpeed      float64
	jerkLimit     float64
}

func NewMovementSystem() *MovementSystem {
	return &MovementSystem{
		lookAheadDist: 50.0,
		deceleration:  15.0,
		acceleration:  10.0,
		minSpeed:      5.0,
		jerkLimit:     25.0,
	}
}

func (ms *MovementSystem) Update(w *world.World, dt float64) {
	w.Mu.Lock()
	defer w.Mu.Unlock()

	for _, v := range w.Vehicles {
		targetSpeed := ms.calculateTargetSpeed(w, v)
		
		ms.adjustSpeedSmooth(v, targetSpeed, dt)
		
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
	
	if ms.hasDespawnPoint(w, v.Road) {
		return v.Road.MaxSpeed
	}
	
	hasValidExit := ms.hasValidExit(w, v)
	
	if hasValidExit {
		return v.Road.MaxSpeed
	}
	
	effectiveLookAhead := math.Min(ms.lookAheadDist, v.Road.Length * 0.4)
	if effectiveLookAhead < 20.0 {
		effectiveLookAhead = 20.0
	}
	
	if distToEnd > effectiveLookAhead {
		return v.Road.MaxSpeed
	}
	
	return ms.calculateApproachSpeed(distToEnd, v.Road.MaxSpeed, effectiveLookAhead)
}

func (ms *MovementSystem) hasValidExit(w *world.World, v *vehicle.Vehicle) bool {
	intersection := w.IntersectionsByNode[v.Road.To.ID]
	if intersection == nil || len(intersection.Outgoing) == 0 {
		return false
	}
	
	for _, r := range intersection.Outgoing {
		if !(r.From == v.Road.To && r.To == v.Road.From) {
			return true
		}
	}
	
	return false
}

func (ms *MovementSystem) calculateApproachSpeed(distToEnd, maxSpeed, effectiveLookAhead float64) float64 {
	if distToEnd <= 0 {
		return ms.minSpeed
	}
	
	ratio := distToEnd / effectiveLookAhead
	
	smoothRatio := ms.smoothStep(ratio)
	
	return ms.minSpeed + (maxSpeed-ms.minSpeed)*smoothRatio
}

func (ms *MovementSystem) smoothStep(x float64) float64 {
	if x <= 0 {
		return 0
	}
	if x >= 1 {
		return 1
	}
	
	return x * x * (3 - 2*x)
}

func (ms *MovementSystem) hasDespawnPoint(w *world.World, rd *road.Road) bool {
	for _, dp := range w.DespawnPoints {
		if dp.Enabled && dp.Road == rd {
			return true
		}
	}
	return false
}

func (ms *MovementSystem) adjustSpeedSmooth(v *vehicle.Vehicle, targetSpeed, dt float64) {
	speedDiff := targetSpeed - v.Speed
	
	if math.Abs(speedDiff) < 0.1 {
		v.Speed = targetSpeed
		return
	}
	
	if speedDiff > 0 {
		maxAccel := ms.acceleration
		
		if v.Speed < ms.minSpeed {
			maxAccel *= 1.5
		}
		
		speedRatio := v.Speed / v.Road.MaxSpeed
		accelModifier := 1.0 - 0.3*speedRatio
		maxAccel *= accelModifier
		
		change := math.Min(speedDiff, maxAccel*dt)
		change = math.Min(change, ms.jerkLimit*dt)
		
		v.Speed += change
		
	} else {
		maxDecel := ms.deceleration
		
		urgency := math.Abs(speedDiff) / v.Road.MaxSpeed
		if urgency > 0.5 {
			maxDecel *= (1.0 + urgency)
		}
		
		change := math.Max(speedDiff, -maxDecel*dt)
		change = math.Max(change, -ms.jerkLimit*dt)
		
		v.Speed += change
		
		if v.Speed < ms.minSpeed && targetSpeed < ms.minSpeed {
			v.Speed = targetSpeed
		}
	}
	
	if v.Speed < 0 {
		v.Speed = 0
	}
	if v.Speed > v.Road.MaxSpeed {
		v.Speed = v.Road.MaxSpeed
	}
}