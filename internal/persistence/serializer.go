package persistence

import (
	"time"
	"traffic-sim/internal/world"
)

func SerializeWorld(w *world.World) *SaveFormat {
	w.Mu.RLock()
	defer w.Mu.RUnlock()

	saveData := &SaveFormat{
		Version:       CurrentVersion,
		Timestamp:     time.Now().Format(time.RFC3339),
		Nodes:         make([]NodeData, 0, len(w.Nodes)),
		Roads:         make([]RoadData, 0, len(w.Roads)),
		SpawnPoints:   make([]SpawnPointData, 0, len(w.SpawnPoints)),
		DespawnPoints: make([]DespawnPointData, 0, len(w.DespawnPoints)),
		TrafficLights: make([]TrafficLightData, 0, len(w.TrafficLights)),
	}

	for _, node := range w.Nodes {
		saveData.Nodes = append(saveData.Nodes, NodeData{
			ID: node.ID,
			X:  node.X,
			Y:  node.Y,
		})
	}

	for _, rd := range w.Roads {
		roadData := RoadData{
			ID:           rd.ID,
			FromNodeID:   rd.From.ID,
			ToNodeID:     rd.To.ID,
			MaxSpeed:     rd.MaxSpeed,
			Width:        rd.Width,
			StartOffsetX: rd.StartOffset.X,
			StartOffsetY: rd.StartOffset.Y,
			EndOffsetX:   rd.EndOffset.X,
			EndOffsetY:   rd.EndOffset.Y,
		}
		
		if rd.ReverseRoad != nil {
			roadData.ReverseRoadID = rd.ReverseRoad.ID
		}
		
		saveData.Roads = append(saveData.Roads, roadData)
	}

	for _, sp := range w.SpawnPoints {
		saveData.SpawnPoints = append(saveData.SpawnPoints, SpawnPointData{
			ID:             sp.ID,
			NodeID:         sp.Node.ID,
			RoadID:         sp.Road.ID,
			Interval:       sp.Interval,
			MinSpeed:       sp.MinSpeed,
			MaxSpeed:       sp.MaxSpeed,
			MaxVehicles:    sp.MaxVehicles,
			Enabled:        sp.Enabled,
			VehicleCounter: sp.VehicleCounter,
		})
	}

	for _, dp := range w.DespawnPoints {
		saveData.DespawnPoints = append(saveData.DespawnPoints, DespawnPointData{
			ID:      dp.ID,
			NodeID:  dp.Node.ID,
			RoadID:  dp.Road.ID,
			Enabled: dp.Enabled,
		})
	}

	for _, light := range w.TrafficLights {
		controlledIDs := make([]string, len(light.ControlledRoads))
		for i, rd := range light.ControlledRoads {
			controlledIDs[i] = rd.ID
		}

		saveData.TrafficLights = append(saveData.TrafficLights, TrafficLightData{
			ID:                light.ID,
			IntersectionID:    light.Intersection.ID,
			ControlledRoadIDs: controlledIDs,
			State:             int(light.State),
			GreenTime:         light.GreenTime,
			YellowTime:        light.YellowTime,
			RedTime:           light.RedTime,
			Enabled:           light.Enabled,
		})
	}

	return saveData
}