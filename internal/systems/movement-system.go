package systems

import (
	"traffic-sim/internal/vehicle"
	"traffic-sim/internal/world"
)

type MovementSystem struct{}

func NewMovementSystem() *MovementSystem {
	return &MovementSystem{}
}

func (ms *MovementSystem) Update(w *world.World, dt float64) {
	w.Mu.Lock()
	defer w.Mu.Unlock()

	for _, v := range w.Vehicles {
		newDist := v.Distance + v.Speed*dt

		if newDist >= v.Road.Length {
			v.Distance = v.Road.Length
		} else {
			v.Distance = newDist
		}

		x, y, t := v.Road.PosAt(v.Distance)
		v.Pos.X = x
		v.Pos.Y = y
		ms.handleSpeed(t,v)
	}
}

func (ms *MovementSystem) handleSpeed(t float64,v *vehicle.Vehicle){
	if t > 0.9 {
		ms.slowDown(v)
		return
	}
	if t<0.9 && v.Speed<v.Road.MaxSpeed{
		ms.accelerate(v)
		return
	} 
}

func (ms *MovementSystem) slowDown(v *vehicle.Vehicle){
	if v.Speed >10 {
		v.Speed = v.Speed*0.90
	}
}
func (ms *MovementSystem) accelerate(vehicle *vehicle.Vehicle ){
	if vehicle.Speed < vehicle.Road.MaxSpeed {
		vehicle.Speed = vehicle.Speed*1.1
	}
}