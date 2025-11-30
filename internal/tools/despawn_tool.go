package tools

import (
	"fmt"
	"traffic-sim/internal/commands"
	"traffic-sim/internal/query"
	"traffic-sim/internal/road"
)

type DespawnTool struct {
	executor       *commands.CommandExecutor
	query          *query.WorldQuery
	maxSnapDist    float64
	despawnCounter int
	roadSelector   *RoadSelector
}

func NewDespawnTool(executor *commands.CommandExecutor, query *query.WorldQuery) *DespawnTool {
	return &DespawnTool{
		executor:       executor,
		query:          query,
		maxSnapDist:    20.0,
		despawnCounter: 0,
		roadSelector:   NewRoadSelector(),
	}
}

func (t *DespawnTool) GetHoverNode(mouseX, mouseY float64) *road.Node {
	return t.query.FindNearestNode(mouseX, mouseY, t.maxSnapDist)
}

func (t *DespawnTool) GetSelectedNode() *road.Node {
	return t.roadSelector.GetSelectedNode()
}

func (t *DespawnTool) GetSelectedRoad() *road.Road {
	return t.roadSelector.GetSelectedRoad()
}

func (t *DespawnTool) GetIncomingRoads(node *road.Node) []*road.Road {
	return t.query.GetIncomingRoads(node)
}

func (t *DespawnTool) Cancel() {
	t.roadSelector.Clear()
}

func (t *DespawnTool) Click(mouseX, mouseY float64) error {
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
		t.despawnCounter++
		despawnID := fmt.Sprintf("dp%d", t.despawnCounter)
		
		cmd := &commands.CreateDespawnPointCommand{
			DespawnID: despawnID,
			Node:      t.roadSelector.GetSelectedNode(),
			Road:      t.roadSelector.GetSelectedRoad(),
		}
		
		if err := t.executor.Execute(cmd); err != nil {
			return err
		}
		
		t.Cancel()
	}

	return nil
}

func (t *DespawnTool) CycleRoad() {
	if !t.roadSelector.HasNode() {
		return
	}

	incoming := t.GetIncomingRoads(t.roadSelector.GetSelectedNode())
	t.roadSelector.CycleRoad(incoming)
}