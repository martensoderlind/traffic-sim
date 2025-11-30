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
	selectedNode   *road.Node
	selectedRoad   *road.Road
}

func NewDespawnTool(executor *commands.CommandExecutor, query *query.WorldQuery) *DespawnTool {
	return &DespawnTool{
		executor:       executor,
		query:          query,
		maxSnapDist:    20.0,
		despawnCounter: 0,
	}
}

func (t *DespawnTool) GetHoverNode(mouseX, mouseY float64) *road.Node {
	return t.query.FindNearestNode(mouseX, mouseY, t.maxSnapDist)
}

func (t *DespawnTool) GetSelectedNode() *road.Node {
	return t.selectedNode
}

func (t *DespawnTool) GetSelectedRoad() *road.Road {
	return t.selectedRoad
}

func (t *DespawnTool) GetIncomingRoads(node *road.Node) []*road.Road {
	return t.query.GetIncomingRoads(node)
}

func (t *DespawnTool) Cancel() {
	t.selectedNode = nil
	t.selectedRoad = nil
}

func (t *DespawnTool) Click(mouseX, mouseY float64) error {
	hoverNode := t.GetHoverNode(mouseX, mouseY)
	
	if hoverNode == nil {
		return nil
	}

	if t.selectedNode == nil {
		t.selectedNode = hoverNode
		return nil
	}

	if t.selectedNode == hoverNode && t.selectedRoad == nil {
		incoming := t.GetIncomingRoads(hoverNode)
		if len(incoming) == 0 {
			t.Cancel()
			return nil
		}

		t.selectedRoad = incoming[0]
		
		for i, r := range incoming {
			if i > 0 {
				t.selectedRoad = r
				break
			}
		}
	}

	if t.selectedNode != nil && t.selectedRoad != nil {
		t.despawnCounter++
		despawnID := fmt.Sprintf("dp%d", t.despawnCounter)
		
		cmd := &commands.CreateDespawnPointCommand{
			DespawnID: despawnID,
			Node:      t.selectedNode,
			Road:      t.selectedRoad,
		}
		
		if err := t.executor.Execute(cmd); err != nil {
			return err
		}
		
		t.Cancel()
	}

	return nil
}

func (t *DespawnTool) CycleRoad() {
	if t.selectedNode == nil {
		return
	}

	incoming := t.GetIncomingRoads(t.selectedNode)
	if len(incoming) == 0 {
		return
	}

	if t.selectedRoad == nil {
		t.selectedRoad = incoming[0]
		return
	}

	for i, r := range incoming {
		if r == t.selectedRoad {
			t.selectedRoad = incoming[(i+1)%len(incoming)]
			return
		}
	}

	t.selectedRoad = incoming[0]
}