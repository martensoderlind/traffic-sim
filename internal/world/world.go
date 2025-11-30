package world

import (
	"sync"

	"traffic-sim/internal/road"
	"traffic-sim/internal/vehicle"
)

type World struct {
	Roads         []*road.Road
	Nodes         []*road.Node
	Vehicles      []*vehicle.Vehicle
	Intersections []*road.Intersection
	SpawnPoints   []*road.SpawnPoint
	DespawnPoints []*road.DespawnPoint

	IntersectionsByNode map[string]*road.Intersection

	Mu sync.RWMutex
}

func New(roads []*road.Road, nodes []*road.Node, vehicles []*vehicle.Vehicle) *World {
	w := &World{
		Roads:               roads,
		Nodes:               nodes,
		Vehicles:            vehicles,
		SpawnPoints:         make([]*road.SpawnPoint, 0),
		DespawnPoints:       make([]*road.DespawnPoint, 0),
		IntersectionsByNode: make(map[string]*road.Intersection),
	}

	w.Intersections = road.BuildIntersections(roads, nodes)
	for _, i := range w.Intersections {
		w.IntersectionsByNode[i.ID] = i
	}

	return w
}