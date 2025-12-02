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
	selectedNode   *road.Node
	selectedRoads  []*road.Road
	availableRoads []*road.Road
	roadSelector   *RoadSelector
}

func NewTrafficLightTool(executor *commands.CommandExecutor, query *query.WorldQuery) *TrafficLightTool {
	return &TrafficLightTool{
		executor:       executor,
		query:          query,
		maxSnapDist:    20.0,
		lightCounter:   0,
		selectedRoads:  make([]*road.Road, 0),
		availableRoads: make([]*road.Road, 0),
		roadSelector:   NewRoadSelector(),
	}
}

func (t *TrafficLightTool) GetHoverNode(mouseX, mouseY float64) *road.Node {
	return t.query.FindNearestNode(mouseX, mouseY, t.maxSnapDist)
}

func (t *TrafficLightTool) GetSelectedNode() *road.Node {
	return t.selectedNode
}

func (t *TrafficLightTool) GetSelectedRoad() *road.Road {
	return t.roadSelector.GetSelectedRoad()
}

func (t *TrafficLightTool) GetSelectedRoads() []*road.Road {
	return t.selectedRoads
}

func (t *TrafficLightTool) GetAvailableRoads() []*road.Road {
	return t.availableRoads
}

func (t *TrafficLightTool) GetIncomingRoads(node *road.Node) []*road.Road {
	return t.query.GetIncomingRoads(node)
}

func (t *TrafficLightTool) Cancel() {
	t.selectedNode = nil
	t.selectedRoads = make([]*road.Road, 0)
	t.availableRoads = make([]*road.Road, 0)
}

func (t *TrafficLightTool) Click(mouseX, mouseY float64) error {
	hoverNode := t.GetHoverNode(mouseX, mouseY)

	if hoverNode == nil {
		return nil
	}

	if t.selectedNode == nil {
		t.selectedNode = hoverNode
		t.availableRoads = t.GetIncomingRoads(hoverNode)
		if len(t.availableRoads) == 0 {
			t.Cancel()
			return nil
		}
		return nil
	}

	if t.selectedNode == hoverNode {
		if len(t.selectedRoads) == 0 {
			t.Cancel()
			return nil
		}

		t.lightCounter++
		lightID := fmt.Sprintf("tl%d", t.lightCounter)

		cmd := &commands.CreateTrafficLightCommand{
			LightID: lightID,
			Node:    t.selectedNode,
			Roads:   t.selectedRoads,
		}

		if err := t.executor.Execute(cmd); err != nil {
			return err
		}

		t.Cancel()
	}

	return nil
}

func (t *TrafficLightTool) ToggleRoad(road *road.Road) {
	for i, r := range t.selectedRoads {
		if r == road {
			t.selectedRoads = append(t.selectedRoads[:i], t.selectedRoads[i+1:]...)
			return
		}
	}
	t.selectedRoads = append(t.selectedRoads, road)
}

func (t *TrafficLightTool) IsRoadSelected(road *road.Road) bool {
	for _, r := range t.selectedRoads {
		if r == road {
			return true
		}
	}
	return false
}

func (t *TrafficLightTool) CycleRoad() {
	if t.selectedNode == nil || len(t.availableRoads) == 0 {
		return
	}

	if len(t.selectedRoads) == 0 {
		t.selectedRoads = append(t.selectedRoads, t.availableRoads[0])
		return
	}

	lastRoad := t.selectedRoads[len(t.selectedRoads)-1]
	
	for i, r := range t.availableRoads {
		if r == lastRoad {
			nextIdx := (i + 1) % len(t.availableRoads)
			nextRoad := t.availableRoads[nextIdx]
			if !t.IsRoadSelected(nextRoad) {
				t.selectedRoads = append(t.selectedRoads, nextRoad)
			}
			return
		}
	}
	
	for _, r := range t.availableRoads {
		if !t.IsRoadSelected(r) {
			t.selectedRoads = append(t.selectedRoads, r)
			return
		}
	}
}