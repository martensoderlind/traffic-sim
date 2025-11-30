package tools

import (
	"fmt"
	"traffic-sim/internal/commands"
	"traffic-sim/internal/query"
	"traffic-sim/internal/road"
)

type RoadBuildingTool struct {
	executor      *commands.CommandExecutor
	query         *query.WorldQuery
	nodeCounter   int
	bidirectional bool
	
	selectedNode *road.Node
	maxSnapDist  float64
	minNodeDist  float64
	roadSnapDist float64
}

func NewRoadBuildingTool(executor *commands.CommandExecutor, query *query.WorldQuery) *RoadBuildingTool {
	return &RoadBuildingTool{
		executor:      executor,
		query:         query,
		nodeCounter:   1000,
		bidirectional: true,
		maxSnapDist:   20.0,
		minNodeDist:   30.0,
		roadSnapDist:  15.0,
	}
}

func (t *RoadBuildingTool) SetBidirectional(enabled bool) {
	t.bidirectional = enabled
}

func (t *RoadBuildingTool) IsBidirectional() bool {
	return t.bidirectional
}

func (t *RoadBuildingTool) GetSelectedNode() *road.Node {
	return t.selectedNode
}

func (t *RoadBuildingTool) GetHoverNode(mouseX, mouseY float64) *road.Node {
	return t.query.FindNearestNode(mouseX, mouseY, t.maxSnapDist)
}

func (t *RoadBuildingTool) GetHoverRoad(mouseX, mouseY float64) (*road.Road, float64, float64) {
	return t.query.FindNearestRoad(mouseX, mouseY, t.roadSnapDist)
}

func (t *RoadBuildingTool) Cancel() {
	t.selectedNode = nil
}

func (t *RoadBuildingTool) Click(mouseX, mouseY float64) error {
	hoverNode := t.GetHoverNode(mouseX, mouseY)
	
	var clickedNode *road.Node
	
	if hoverNode != nil {
		clickedNode = hoverNode
	} else {
		hoverRoad, snapX, snapY := t.GetHoverRoad(mouseX, mouseY)
		
		if hoverRoad != nil {
			t.nodeCounter++
			nodeID := fmt.Sprintf("n%d", t.nodeCounter)
			
			cmd := &commands.SplitRoadCommand{
				Road:   hoverRoad,
				X:      snapX,
				Y:      snapY,
				NodeID: nodeID,
			}
			
			if err := t.executor.Execute(cmd); err != nil {
				return err
			}
			
			clickedNode = t.query.FindNodeByID(nodeID)
			if clickedNode == nil {
				return fmt.Errorf("failed to find created node")
			}
		} else {
			if !t.query.CanPlaceNodeAt(mouseX, mouseY, t.minNodeDist) {
				return nil
			}
			
			t.nodeCounter++
			nodeID := fmt.Sprintf("n%d", t.nodeCounter)
			
			cmd := &commands.CreateNodeCommand{
				X:      mouseX,
				Y:      mouseY,
				NodeID: nodeID,
			}
			
			if err := t.executor.Execute(cmd); err != nil {
				return err
			}
			
			clickedNode = t.query.FindNodeByID(nodeID)
			if clickedNode == nil {
				return fmt.Errorf("failed to find created node")
			}
		}
	}
	
	if t.selectedNode == nil {
		t.selectedNode = clickedNode
		return nil
	}
	
	if t.selectedNode != clickedNode {
		if err := t.createRoadBetween(t.selectedNode, clickedNode); err != nil {
			return err
		}
		
		if t.bidirectional {
			if err := t.createRoadBetween(clickedNode, t.selectedNode); err != nil {
				return err
			}
		}
		
		t.selectedNode = nil
	}
	
	return nil
}

func (t *RoadBuildingTool) createRoadBetween(from, to *road.Node) error {
	cmd := &commands.CreateRoadCommand{
		From:     from,
		To:       to,
		MaxSpeed: 40.0,
		Width:    8.0,
	}
	return t.executor.Execute(cmd)
}