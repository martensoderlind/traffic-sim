package tools

import (
	"traffic-sim/internal/commands"
	"traffic-sim/internal/query"
	"traffic-sim/internal/road"
)

type SpawnPointPropertiesTool struct {
	executor     *commands.CommandExecutor
	query        *query.WorldQuery
	selectedSpawnPoint *road.SpawnPoint
	maxSnapDist  float64
}

func NewSpawnPointPropertiesTool(executor *commands.CommandExecutor, query *query.WorldQuery) *SpawnPointPropertiesTool {
	return &SpawnPointPropertiesTool{
		executor:    executor,
		query:       query,
		maxSnapDist: 15.0,
	}
}

func (t *SpawnPointPropertiesTool) GetHoverSpawnPoint(mouseX, mouseY float64) *road.SpawnPoint {
	sp, _, _ := t.query.FindNearestSpawnPoint(mouseX, mouseY, t.maxSnapDist)
	return sp
}

func (t *SpawnPointPropertiesTool) GetSelectedSpawnPoint() *road.SpawnPoint {
	return t.selectedSpawnPoint
}

func (t *SpawnPointPropertiesTool) Click(mouseX, mouseY float64) error {
	hoverSpawnPoint := t.GetHoverSpawnPoint(mouseX, mouseY)
	
	if hoverSpawnPoint == nil {
		t.selectedSpawnPoint = nil
		return nil
	}

	t.selectedSpawnPoint = hoverSpawnPoint
	return nil
}

func (t *SpawnPointPropertiesTool) UpdateSpawnPointProperties(Interval,MinSpeed,MaxSpeed float64, MaxVehicles int,Enabled bool) error {
	if t.selectedSpawnPoint == nil {
		return nil
	}

	cmd := &commands.UpdateSpawnPointPropertiesCommand{
		SpawnPoint:     t.selectedSpawnPoint,
		Interval:	  	Interval,
		MinSpeed: 		MinSpeed,
		MaxSpeed:		MaxSpeed,
		MaxVehicles: 	MaxVehicles,
		Enabled:	   	Enabled,
	}

	return t.executor.Execute(cmd)
}

func (t *SpawnPointPropertiesTool) Cancel() {
	t.selectedSpawnPoint = nil
}