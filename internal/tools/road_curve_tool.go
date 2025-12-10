package tools

import (
	"traffic-sim/internal/commands"
	"traffic-sim/internal/query"
	"traffic-sim/internal/road"
)

type RoadCurveTool struct {
	executor       *commands.CommandExecutor
	query          *query.WorldQuery
	maxSnapDist    float64
	selectedRoad   *road.Road
	incomingRoad   *road.Road
	outgoingRoad   *road.Road
	stage          int 
}

func NewRoadCurveTool(executor *commands.CommandExecutor, query *query.WorldQuery) *RoadCurveTool {
	return &RoadCurveTool{
		executor:    executor,
		query:       query,
		maxSnapDist: 15.0,
		stage:       0,
	}
}

func (t *RoadCurveTool) GetHoverRoad(mouseX, mouseY float64) *road.Road {
	rd, _, _ := t.query.FindNearestRoad(mouseX, mouseY, t.maxSnapDist)
	return rd
}

func (t *RoadCurveTool) Click(mouseX, mouseY float64) error {
	hoverRoad := t.GetHoverRoad(mouseX, mouseY)

	if t.stage == 0 {
		if hoverRoad == nil {
			return nil
		}
		t.selectedRoad = hoverRoad
		t.stage = 1
		return nil
	}

	if t.stage == 1 {
		if hoverRoad == nil {
			t.stage = 0
			t.selectedRoad = nil
			return nil
		}

		if (hoverRoad.To != t.selectedRoad.From) && (hoverRoad.From != t.selectedRoad.From) {
			return nil
		}

		t.incomingRoad = hoverRoad
		t.stage = 2
		return nil
	}

	if t.stage == 2 {
		if hoverRoad == nil {
			t.stage = 1
			t.incomingRoad = nil
			return nil
		}

		if (hoverRoad.From != t.selectedRoad.To) && (hoverRoad.To != t.selectedRoad.To) {
			return nil
		}

		t.outgoingRoad = hoverRoad
		
		cmd := &commands.CurveRoadCommand{
			Road:         t.selectedRoad,
			IncomingRoad: t.incomingRoad,
			OutgoingRoad: t.outgoingRoad,
		}

		if err := t.executor.Execute(cmd); err != nil {
			return err
		}

		t.Cancel()
		return nil
	}

	return nil
}

func (t *RoadCurveTool) GetSelectedRoad() *road.Road {
	return t.selectedRoad
}

func (t *RoadCurveTool) GetIncomingRoad() *road.Road {
	return t.incomingRoad
}

func (t *RoadCurveTool) GetOutgoingRoad() *road.Road {
	return t.outgoingRoad
}

func (t *RoadCurveTool) GetStage() int {
	return t.stage
}

func (t *RoadCurveTool) Cancel() {
	t.selectedRoad = nil
	t.incomingRoad = nil
	t.outgoingRoad = nil
	t.stage = 0
}

func (t *RoadCurveTool) GetStatusMessage() string {
	switch t.stage {
	case 0:
		return "Click on a road to curve"
	case 1:
		return "Click on the incoming road (connected to the start node)"
	case 2:
		return "Click on the outgoing road (connected to the end node)"
	}
	return ""
}
