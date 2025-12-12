package persistence

import (
	"fmt"
	"traffic-sim/internal/road"
	"traffic-sim/internal/world"
)

func DeserializeWorld(saveData *SaveFormat) (*world.World, error) {
	if saveData.Version != CurrentVersion {
		return nil, fmt.Errorf("incompatible save version: %s (expected %s)", saveData.Version, CurrentVersion)
	}

	w := world.New()

	nodeMap := make(map[string]*road.Node)
	for _, nodeData := range saveData.Nodes {
		node := &road.Node{
			ID: nodeData.ID,
			X:  nodeData.X,
			Y:  nodeData.Y,
		}
		w.Nodes = append(w.Nodes, node)
		nodeMap[node.ID] = node
		w.CreateIntersection(node.ID)
	}

	roadMap := make(map[string]*road.Road)
	for _, roadData := range saveData.Roads {
		fromNode, fromExists := nodeMap[roadData.FromNodeID]
		toNode, toExists := nodeMap[roadData.ToNodeID]

		if !fromExists || !toExists {
			return nil, fmt.Errorf("road %s references non-existent nodes", roadData.ID)
		}

		rd := road.NewRoad(roadData.ID, fromNode, toNode, roadData.MaxSpeed)
		rd.Width = roadData.Width
		rd.StartOffset = road.Point{X: roadData.StartOffsetX, Y: roadData.StartOffsetY}
		rd.EndOffset = road.Point{X: roadData.EndOffsetX, Y: roadData.EndOffsetY}
		rd.UpdateLength()

		w.Roads = append(w.Roads, rd)
		roadMap[rd.ID] = rd

		w.AddRoadToIntersections(rd)
	}

	for _, roadData := range saveData.Roads {
		if roadData.ReverseRoadID != "" {
			rd := roadMap[roadData.ID]
			reverseRd := roadMap[roadData.ReverseRoadID]
			
			if rd != nil && reverseRd != nil {
				rd.ReverseRoad = reverseRd
			}
		}
	}

	for _, spData := range saveData.SpawnPoints {
		node, nodeExists := nodeMap[spData.NodeID]
		rd, roadExists := roadMap[spData.RoadID]

		if !nodeExists || !roadExists {
			return nil, fmt.Errorf("spawn point %s references non-existent node or road", spData.ID)
		}

		sp := &road.SpawnPoint{
			ID:             spData.ID,
			Node:           node,
			Road:           rd,
			Interval:       spData.Interval,
			Timer:          0.0,
			MinSpeed:       spData.MinSpeed,
			MaxSpeed:       spData.MaxSpeed,
			MaxVehicles:    spData.MaxVehicles,
			Enabled:        spData.Enabled,
			VehicleCounter: spData.VehicleCounter,
		}

		w.SpawnPoints = append(w.SpawnPoints, sp)
	}

	for _, dpData := range saveData.DespawnPoints {
		node, nodeExists := nodeMap[dpData.NodeID]
		rd, roadExists := roadMap[dpData.RoadID]

		if !nodeExists || !roadExists {
			return nil, fmt.Errorf("despawn point %s references non-existent node or road", dpData.ID)
		}

		dp := &road.DespawnPoint{
			ID:      dpData.ID,
			Node:    node,
			Road:    rd,
			Enabled: dpData.Enabled,
		}

		w.DespawnPoints = append(w.DespawnPoints, dp)
	}

	for _, lightData := range saveData.TrafficLights {
		intersection := w.IntersectionsByNode[lightData.IntersectionID]
		if intersection == nil {
			return nil, fmt.Errorf("traffic light %s references non-existent intersection", lightData.ID)
		}

		light := &road.TrafficLight{
			ID:              lightData.ID,
			Intersection:    intersection,
			ControlledRoads: make([]*road.Road, 0, len(lightData.ControlledRoadIDs)),
			State:           road.LightState(lightData.State),
			Timer:           0.0,
			GreenTime:       lightData.GreenTime,
			YellowTime:      lightData.YellowTime,
			RedTime:         lightData.RedTime,
			Enabled:         lightData.Enabled,
		}

		for _, roadID := range lightData.ControlledRoadIDs {
			rd, exists := roadMap[roadID]
			if !exists {
				return nil, fmt.Errorf("traffic light %s references non-existent road %s", lightData.ID, roadID)
			}
			light.ControlledRoads = append(light.ControlledRoads, rd)
		}

		w.TrafficLights = append(w.TrafficLights, light)
	}

	return w, nil
}