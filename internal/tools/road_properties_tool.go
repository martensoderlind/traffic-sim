package tools

import (
	"traffic-sim/internal/commands"
	"traffic-sim/internal/query"
	"traffic-sim/internal/road"
)

type RoadPropertiesTool struct {
	executor     *commands.CommandExecutor
	query        *query.WorldQuery
	selectedRoad *road.Road
	maxSnapDist  float64
}

func NewRoadPropertiesTool(executor *commands.CommandExecutor, query *query.WorldQuery) *RoadPropertiesTool {
	return &RoadPropertiesTool{
		executor:    executor,
		query:       query,
		maxSnapDist: 15.0,
	}
}

func (t *RoadPropertiesTool) GetHoverRoad(mouseX, mouseY float64) *road.Road {
	rd, _, _ := t.query.FindNearestRoad(mouseX, mouseY, t.maxSnapDist)
	return rd
}

func (t *RoadPropertiesTool) GetSelectedRoad() *road.Road {
	return t.selectedRoad
}

func (t *RoadPropertiesTool) Click(mouseX, mouseY float64) error {
	hoverRoad := t.GetHoverRoad(mouseX, mouseY)
	
	if hoverRoad == nil {
		t.selectedRoad = nil
		return nil
	}

	t.selectedRoad = hoverRoad
	return nil
}

func (t *RoadPropertiesTool) UpdateRoadProperties(maxSpeed, width float64) error {
	if t.selectedRoad == nil {
		return nil
	}

	cmd := &commands.UpdateRoadPropertiesCommand{
		Road:     t.selectedRoad,
		MaxSpeed: maxSpeed,
		Width:    width,
	}

	return t.executor.Execute(cmd)
}

func (t *RoadPropertiesTool) Cancel() {
	t.selectedRoad = nil
}