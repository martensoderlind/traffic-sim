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
	accumulator   float64
}

func NewSimulator(w *world.World, tickRate time.Duration) *Simulator {
	sm := systems.NewSystemManager()
	sm.AddSystem(systems.NewSpawnSystem())
	sm.AddSystem(systems.NewCollisionSystem())
	sm.AddSystem(systems.NewTrafficLightSystem())
	sm.AddSystem(systems.NewPathfindingSystem())
	sm.AddSystem(systems.NewMovementSystem())
	sm.AddSystem(systems.NewDespawnSystem())

	return &Simulator{
		world:         w,
		tickRate:      tickRate,
		systemManager: sm,
		accumulator:   0.0,
	}
}

func (s *Simulator) Start() {
	ticker := time.NewTicker(s.tickRate)
	for range ticker.C {
		s.update()
	}
}

func (s *Simulator) UpdateOnce(deltaTime float64) {
	s.accumulator += deltaTime
	fixedDt := s.tickRate.Seconds()
	
	for s.accumulator >= fixedDt {
		s.systemManager.Update(s.world, fixedDt)
		s.accumulator -= fixedDt
	}
}

func (s *Simulator) update() {
	dt := s.tickRate.Seconds()
	s.systemManager.Update(s.world, dt)
}