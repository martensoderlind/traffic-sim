package tools

import (
	"traffic-sim/internal/commands"
	"traffic-sim/internal/query"
	"traffic-sim/internal/road"
)

type RoadDeleteTool struct {
	executor    *commands.CommandExecutor
	query       *query.WorldQuery
	maxSnapDist float64
}

func NewRoadDeleteTool(executor *commands.CommandExecutor, query *query.WorldQuery) *RoadDeleteTool {
	return &RoadDeleteTool{
		executor:    executor,
		query:       query,
		maxSnapDist: 15.0,
	}
}

func (t *RoadDeleteTool) GetHoverRoad(mouseX, mouseY float64) *road.Road {
	rd, _, _ := t.query.FindNearestRoad(mouseX, mouseY, t.maxSnapDist)
	return rd
}

func (t *RoadDeleteTool) Click(mouseX, mouseY float64) error {
	hoverRoad := t.GetHoverRoad(mouseX, mouseY)
	
	if hoverRoad == nil {
		return nil
	}

	cmd := &commands.DeleteRoadCommand{
		Road: hoverRoad,
	}

	return t.executor.Execute(cmd)
}

func (t *RoadDeleteTool) Cancel() {
}