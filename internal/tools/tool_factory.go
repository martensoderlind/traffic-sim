package tools

import (
	"traffic-sim/internal/commands"
	"traffic-sim/internal/query"
)

type ToolSet struct {
	RoadBuilding       *RoadBuildingTool
	NodeMoving         *NodeMoveTool
	Spawning           *SpawnTool
	Despawning         *DespawnTool
	RoadDeleting       *RoadDeleteTool
	NodeDeleting       *NodeDeleteTool
	TrafficLight       *TrafficLightTool
	RoadProperties     *RoadPropertiesTool
	SpawnPointProperties *SpawnPointPropertiesTool
	RoadCurving        *RoadCurveTool
}

type ToolFactory struct {
	executor *commands.CommandExecutor
	query    *query.WorldQuery
}

func NewToolFactory(executor *commands.CommandExecutor, query *query.WorldQuery) *ToolFactory {
	return &ToolFactory{
		executor: executor,
		query:    query,
	}
}

func (tf *ToolFactory) CreateAll() *ToolSet {
	return &ToolSet{
		RoadBuilding:        NewRoadBuildingTool(tf.executor, tf.query),
		NodeMoving:          NewNodeMoveTool(tf.executor, tf.query),
		Spawning:            NewSpawnTool(tf.executor, tf.query),
		Despawning:          NewDespawnTool(tf.executor, tf.query),
		RoadDeleting:        NewRoadDeleteTool(tf.executor, tf.query),
		NodeDeleting:        NewNodeDeleteTool(tf.executor, tf.query),
		TrafficLight:        NewTrafficLightTool(tf.executor, tf.query),
		RoadProperties:      NewRoadPropertiesTool(tf.executor, tf.query),
		SpawnPointProperties: NewSpawnPointPropertiesTool(tf.executor, tf.query),
		RoadCurving:         NewRoadCurveTool(tf.executor, tf.query),
	}
}
