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

func BuildIntersections(roads []*road.Road, nodes []*road.Node) []*road.Intersection {
	m := make(map[string]*road.Intersection)

	for _, n := range nodes {
		m[n.ID] = road.NewIntersection(n.ID)
	}

	for _, r := range roads {
		in := m[r.From.ID]
		out := m[r.To.ID]
		in.AddOutgoing(r)
		out.AddIncoming(r)
	}

	intersections := make([]*road.Intersection, 0, len(m))
	for _, i := range m {
		intersections = append(intersections, i)
	}
	return intersections
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

	w.Intersections = BuildIntersections(roads, nodes)
	for _, i := range w.Intersections {
		w.IntersectionsByNode[i.ID] = i
	}

	return w
}