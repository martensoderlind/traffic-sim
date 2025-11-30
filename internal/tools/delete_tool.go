package tools

import (
	"traffic-sim/internal/commands"
	"traffic-sim/internal/query"
	"traffic-sim/internal/road"
)

type NodeDeleteTool struct {
	executor    *commands.CommandExecutor
	query       *query.WorldQuery
	maxSnapDist float64
}

func NewNodeDeleteTool(executor *commands.CommandExecutor, query *query.WorldQuery) *NodeDeleteTool {
	return &NodeDeleteTool{
		executor:    executor,
		query:       query,
		maxSnapDist: 20.0,
	}
}

func (t *NodeDeleteTool) GetHoverNode(mouseX, mouseY float64) *road.Node {
	return t.query.FindNearestNode(mouseX, mouseY, t.maxSnapDist)
}

func (t *NodeDeleteTool) Click(mouseX, mouseY float64) error {
	hoverNode := t.GetHoverNode(mouseX, mouseY)
	
	if hoverNode == nil {
		return nil
	}

	cmd := &commands.DeleteNodeCommand{
		Node: hoverNode,
	}

	return t.executor.Execute(cmd)
}

func (t *NodeDeleteTool) Cancel() {
}