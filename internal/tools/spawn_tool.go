package tools

import (
	"fmt"
	"traffic-sim/internal/commands"
	"traffic-sim/internal/query"
	"traffic-sim/internal/road"
)

type SpawnTool struct {
	executor       *commands.CommandExecutor
	query          *query.WorldQuery
	maxSnapDist    float64
	spawnCounter   int
	selectedNode   *road.Node
	selectedRoad   *road.Road
}

func NewSpawnTool(executor *commands.CommandExecutor, query *query.WorldQuery) *SpawnTool {
	return &SpawnTool{
		executor:     executor,
		query:        query,
		maxSnapDist:  20.0,
		spawnCounter: 0,
	}
}

func (t *SpawnTool) GetHoverNode(mouseX, mouseY float64) *road.Node {
	return t.query.FindNearestNode(mouseX, mouseY, t.maxSnapDist)
}

func (t *SpawnTool) GetSelectedNode() *road.Node {
	return t.selectedNode
}

func (t *SpawnTool) GetSelectedRoad() *road.Road {
	return t.selectedRoad
}

func (t *SpawnTool) GetOutgoingRoads(node *road.Node) []*road.Road {
	return t.query.GetOutgoingRoads(node)
}

func (t *SpawnTool) Cancel() {
	t.selectedNode = nil
	t.selectedRoad = nil
}

func (t *SpawnTool) Click(mouseX, mouseY float64) error {
	hoverNode := t.GetHoverNode(mouseX, mouseY)
	
	if hoverNode == nil {
		return nil
	}

	if t.selectedNode == nil {
		t.selectedNode = hoverNode
		return nil
	}

	if t.selectedNode == hoverNode && t.selectedRoad == nil {
		outgoing := t.GetOutgoingRoads(hoverNode)
		if len(outgoing) == 0 {
			t.Cancel()
			return nil
		}

		t.selectedRoad = outgoing[0]
		
		for i, r := range outgoing {
			if i > 0 {
				t.selectedRoad = r
				break
			}
		}
	}

	if t.selectedNode != nil && t.selectedRoad != nil {
		t.spawnCounter++
		spawnID := fmt.Sprintf("sp%d", t.spawnCounter)
		
		cmd := &commands.CreateSpawnPointCommand{
			SpawnID: spawnID,
			Node:    t.selectedNode,
			Road:    t.selectedRoad,
		}
		
		if err := t.executor.Execute(cmd); err != nil {
			return err
		}
		
		t.Cancel()
	}

	return nil
}

func (t *SpawnTool) CycleRoad() {
	if t.selectedNode == nil {
		return
	}

	outgoing := t.GetOutgoingRoads(t.selectedNode)
	if len(outgoing) == 0 {
		return
	}

	if t.selectedRoad == nil {
		t.selectedRoad = outgoing[0]
		return
	}

	for i, r := range outgoing {
		if r == t.selectedRoad {
			t.selectedRoad = outgoing[(i+1)%len(outgoing)]
			return
		}
	}

	t.selectedRoad = outgoing[0]
}