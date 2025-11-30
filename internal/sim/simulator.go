package sim

import (
	"time"

	"traffic-sim/internal/systems"
	"traffic-sim/internal/world"
)

type Simulator struct {
	world         *world.World
	tickRate      time.Duration
	systemManager *systems.SystemManager
}

func NewSimulator(w *world.World, tickRate time.Duration) *Simulator {
	sm := systems.NewSystemManager()
	sm.AddSystem(systems.NewSpawnSystem())
	sm.AddSystem(systems.NewCollisionSystem())
	sm.AddSystem(systems.NewPathfindingSystem())
	sm.AddSystem(systems.NewMovementSystem())

	return &Simulator{
		world:         w,
		tickRate:      tickRate,
		systemManager: sm,
	}
}

func (s *Simulator) Start() {
	ticker := time.NewTicker(s.tickRate)
	for range ticker.C {
		s.update()
	}
}

func (s *Simulator) UpdateOnce() {
	s.update()
}

func (s *Simulator) update() {
	dt := s.tickRate.Seconds()
	s.systemManager.Update(s.world, dt)
}