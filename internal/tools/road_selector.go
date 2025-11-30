package tools

import (
	"traffic-sim/internal/road"
)

type RoadSelector struct {
	selectedNode *road.Node
	selectedRoad *road.Road
}

func NewRoadSelector() *RoadSelector {
	return &RoadSelector{}
}

func (rs *RoadSelector) GetSelectedNode() *road.Node {
	return rs.selectedNode
}

func (rs *RoadSelector) GetSelectedRoad() *road.Road {
	return rs.selectedRoad
}

func (rs *RoadSelector) SelectNode(node *road.Node) {
	rs.selectedNode = node
	rs.selectedRoad = nil
}

func (rs *RoadSelector) SelectRoad(road *road.Road) {
	rs.selectedRoad = road
}

func (rs *RoadSelector) Clear() {
	rs.selectedNode = nil
	rs.selectedRoad = nil
}

func (rs *RoadSelector) IsComplete() bool {
	return rs.selectedNode != nil && rs.selectedRoad != nil
}

func (rs *RoadSelector) HasNode() bool {
	return rs.selectedNode != nil
}

func (rs *RoadSelector) CycleRoad(roads []*road.Road) {
	if len(roads) == 0 {
		return
	}

	if rs.selectedRoad == nil {
		rs.selectedRoad = roads[0]
		return
	}

	for i, r := range roads {
		if r == rs.selectedRoad {
			rs.selectedRoad = roads[(i+1)%len(roads)]
			return
		}
	}

	rs.selectedRoad = roads[0]
}

func (rs *RoadSelector) AutoSelectFirstRoad(roads []*road.Road) {
	if len(roads) == 0 {
		return
	}

	if len(roads) == 1 {
		rs.selectedRoad = roads[0]
	} else if len(roads) > 1 {
		rs.selectedRoad = roads[1]
	}
}