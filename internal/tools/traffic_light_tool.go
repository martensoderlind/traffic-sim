package tools

import (
	"fmt"
	"traffic-sim/internal/commands"
	"traffic-sim/internal/query"
	"traffic-sim/internal/road"
)

type TrafficLightTool struct {
	executor       *commands.CommandExecutor
	query          *query.WorldQuery
	maxSnapDist    float64
	lightCounter   int
	roadSelector   *RoadSelector
}

func NewTrafficLightTool(executor *commands.CommandExecutor, query *query.WorldQuery) *TrafficLightTool {
	return &TrafficLightTool{
		executor:     executor,
		query:        query,
		maxSnapDist:  20.0,
		lightCounter: 0,
		roadSelector: NewRoadSelector(),
	}
}

func (t *TrafficLightTool) GetHoverNode(mouseX, mouseY float64) *road.Node {
	return t.query.FindNearestNode(mouseX, mouseY, t.maxSnapDist)
}

func (t *TrafficLightTool) GetSelectedNode() *road.Node {
	return t.roadSelector.GetSelectedNode()
}

func (t *TrafficLightTool) GetSelectedRoad() *road.Road {
	return t.roadSelector.GetSelectedRoad()
}

func (t *TrafficLightTool) GetIncomingRoads(node *road.Node) []*road.Road {
	return t.query.GetIncomingRoads(node)
}

func (t *TrafficLightTool) Cancel() {
	t.roadSelector.Clear()
}

func (t *TrafficLightTool) Click(mouseX, mouseY float64) error {
	hoverNode := t.GetHoverNode(mouseX, mouseY)

	if hoverNode == nil {
		return nil
	}

	if !t.roadSelector.HasNode() {
		t.roadSelector.SelectNode(hoverNode)
		return nil
	}

	if t.roadSelector.GetSelectedNode() == hoverNode && t.roadSelector.GetSelectedRoad() == nil {
		incoming := t.GetIncomingRoads(hoverNode)
		if len(incoming) == 0 {
			t.Cancel()
			return nil
		}

		t.roadSelector.AutoSelectFirstRoad(incoming)
	}

	if t.roadSelector.IsComplete() {
		t.lightCounter++
		lightID := fmt.Sprintf("tl%d", t.lightCounter)

		cmd := &commands.CreateTrafficLightCommand{
			LightID: lightID,
			Node:    t.roadSelector.GetSelectedNode(),
			Road:    t.roadSelector.GetSelectedRoad(),
		}

		if err := t.executor.Execute(cmd); err != nil {
			return err
		}

		t.Cancel()
	}

	return nil
}

func (t *TrafficLightTool) CycleRoad() {
	if !t.roadSelector.HasNode() {
		return
	}

	incoming := t.GetIncomingRoads(t.roadSelector.GetSelectedNode())
	t.roadSelector.CycleRoad(incoming)
}