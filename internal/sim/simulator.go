package sim

import (
	"log"
	"time"

	"traffic-sim/internal/config"
	"traffic-sim/internal/systems"
	"traffic-sim/internal/world"
)

type Simulator struct {
	world         *world.World
	tickRate      time.Duration
	systemManager *systems.SystemManager
	accumulator   float64
	paused		bool
}

func NewSimulator(w *world.World, tickRate time.Duration) *Simulator {

	cfg,err:= config.LoadConfig()
	if err != nil {
        log.Fatalf("Could not load config: %v", err)
    }
	sm := systems.NewSystemManager()
	sm.AddSystem(systems.NewSpawnSystem())
	sm.AddSystem(systems.NewCollisionSystem())
	sm.AddSystem(systems.NewTrafficLightSystem())
	if cfg.FeatureFlags.RightOfWaySystem && err	== nil {
		sm.AddSystem(systems.NewRightOfWaySystem())
	}
	sm.AddSystem(systems.NewPathfindingSystem())
	sm.AddSystem(systems.NewMovementSystem())
	sm.AddSystem(systems.NewDespawnSystem())

	return &Simulator{
		world:         w,
		tickRate:      tickRate,
		systemManager: sm,
		accumulator:   0.0,
		paused: false,
	}
}

// ResetSystems clears internal state in all systems (called when world changes)
func (s *Simulator) ResetSystems() {
	s.systemManager.ResetAll()
}

func (s *Simulator) Start() {
	ticker := time.NewTicker(s.tickRate)
	for range ticker.C {
		s.update()
	}
}

func (s *Simulator) UpdateOnce(deltaTime float64) {
	if s.paused {
		return
	}
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
func (s *Simulator) TogglePause() {
	s.paused = !s.paused
}
func (s *Simulator) IsPaused() bool {
	return s.paused
}