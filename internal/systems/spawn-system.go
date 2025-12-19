package systems

import (
	"fmt"
	"math/rand"
	"traffic-sim/internal/road"
	"traffic-sim/internal/vehicle"
	"traffic-sim/internal/world"
)

type SpawnSystem struct{}

func NewSpawnSystem() *SpawnSystem {
	return &SpawnSystem{}
}

func (ss *SpawnSystem) Reset() {
}

func (ss *SpawnSystem) Update(w *world.World, dt float64) {
	w.Mu.Lock()
	defer w.Mu.Unlock()

	for _, sp := range w.SpawnPoints {
		if !sp.Enabled {
			continue
		}

		if len(w.Vehicles) >= sp.MaxVehicles {
			continue
		}

		sp.Timer += dt

		randomInterval := sp.Interval * (0.5 + rand.Float64())
		
		if sp.Timer >= randomInterval {
			sp.Timer = 0.0

			speed := sp.MinSpeed + rand.Float64()*(sp.MaxSpeed-sp.MinSpeed)

			sp.VehicleCounter++
			vehicleID := fmt.Sprintf("%s-v%d", sp.ID, sp.VehicleCounter)

			ss.spawnVehicle(w, vehicleID, sp.Road, speed)
		}
	}
}

func (ss *SpawnSystem) spawnVehicle(w *world.World, id string, rd *road.Road, speed float64) {
	newVehicle := &vehicle.Vehicle{
		ID:       id,
		Road:     rd,
		Distance: 0,
		Speed:    speed,
		Pos:      vehicle.Vec2{X: rd.From.X, Y: rd.From.Y},
	}
	
	ss.assignTargetDespawn(newVehicle, w)

	w.Vehicles = append(w.Vehicles, newVehicle)
}

func (ss *SpawnSystem) assignTargetDespawn(v *vehicle.Vehicle, w *world.World) {
	activeDespawns := make([]*road.DespawnPoint, 0)
	for _, dp := range w.DespawnPoints {
		if dp.Enabled {
			activeDespawns = append(activeDespawns, dp)
		}
	}
	
	if len(activeDespawns) == 0 {
		return
	}
	
	v.TargetDespawn = activeDespawns[rand.Intn(len(activeDespawns))]
}