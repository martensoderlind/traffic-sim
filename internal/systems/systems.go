package systems

import (
	"traffic-sim/internal/world"
)

type System interface {
	Update(w *world.World, dt float64)
	Reset() // Called when world is replaced
}

type SystemManager struct {
	systems []System
}

func NewSystemManager() *SystemManager {
	return &SystemManager{
		systems: make([]System, 0),
	}
}

func (sm *SystemManager) AddSystem(s System) {
	sm.systems = append(sm.systems, s)
}

func (sm *SystemManager) Update(w *world.World, dt float64) {
	for _, system := range sm.systems {
		system.Update(w, dt)
	}
}

// ResetAll calls Reset on all systems (typically called when world changes)
func (sm *SystemManager) ResetAll() {
	for _, system := range sm.systems {
		system.Reset()
	}
}

