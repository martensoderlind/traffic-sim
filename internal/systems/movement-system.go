package systems

import "traffic-sim/internal/world"

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

		x, y := v.Road.PosAt(v.Distance)
		v.Pos.X = x
		v.Pos.Y = y
	}
}
