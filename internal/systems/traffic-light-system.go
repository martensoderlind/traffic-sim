package systems

import (
	"traffic-sim/internal/road"
	"traffic-sim/internal/world"
)

type TrafficLightSystem struct{}

func NewTrafficLightSystem() *TrafficLightSystem {
	return &TrafficLightSystem{}
}

func (tls *TrafficLightSystem) Update(w *world.World, dt float64) {
	w.Mu.Lock()
	defer w.Mu.Unlock()

	for _, light := range w.TrafficLights {
		light.Update(dt)
	}

	tls.enforceTrafficRules(w)
}

func (tls *TrafficLightSystem) enforceTrafficRules(w *world.World) {
	lightsByRoad := tls.buildLightsByRoadMap(w)
	for _, v := range w.Vehicles {
		if v.NextRoad == nil {
			continue
		}
		
		light := lightsByRoad[v.Road.ID]
		if light == nil {
			continue
		}

		distToEnd := v.Road.Length - v.Distance

		stopDist := 30.0

		if light.IsRed() && distToEnd < stopDist {
			if distToEnd > 5.0 {
				slowdownRatio := distToEnd / stopDist
				v.Speed *= slowdownRatio
				} else {
					v.Speed = 0
				}
				} else if light.ShouldSlow() && distToEnd < stopDist {
			v.Speed *= 0.5
		}
	}
}

func (tls *TrafficLightSystem) buildLightsByRoadMap(w *world.World) map[string]*road.TrafficLight {
	lightsByRoad := make(map[string]*road.TrafficLight)

	for _, light := range w.TrafficLights {
		for _, rd := range light.ControlledRoads {
			lightsByRoad[rd.ID] = light
		}
	}

	return lightsByRoad
}