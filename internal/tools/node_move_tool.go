package tools

import (
	"traffic-sim/internal/commands"
	"traffic-sim/internal/query"
	"traffic-sim/internal/road"
)

type NodeMoveTool struct {
	executor    *commands.CommandExecutor
	query       *query.WorldQuery
	draggedNode *road.Node
	maxSnapDist float64
}

func NewNodeMoveTool(executor *commands.CommandExecutor, query *query.WorldQuery) *NodeMoveTool {
	return &NodeMoveTool{
		executor:    executor,
		query:       query,
		maxSnapDist: 20.0,
	}
}

func (t *NodeMoveTool) GetHoverNode(mouseX, mouseY float64) *road.Node {
	return t.query.FindNearestNode(mouseX, mouseY, t.maxSnapDist)
}

func (t *NodeMoveTool) GetDraggedNode() *road.Node {
	return t.draggedNode
}

func (t *NodeMoveTool) StartDrag(mouseX, mouseY float64) {
	node := t.GetHoverNode(mouseX, mouseY)
	if node != nil {
		t.draggedNode = node
	}
}

func (t *NodeMoveTool) UpdateDrag(mouseX, mouseY float64) error {
	if t.draggedNode == nil {
		return nil
	}

	cmd := &commands.MoveNodeCommand{
		Node: t.draggedNode,
		NewX: mouseX,
		NewY: mouseY,
	}

	return t.executor.Execute(cmd)
}

func (t *NodeMoveTool) EndDrag() {
	t.draggedNode = nil
}

func (t *NodeMoveTool) Cancel() {
	t.EndDrag()
}

func (t *NodeMoveTool) Click(mouseX, mouseY float64) error {
	t.StartDrag(mouseX, mouseY)
	return nil
}

func (t *NodeMoveTool) IsDragging() bool {
	return t.draggedNode != nil
}