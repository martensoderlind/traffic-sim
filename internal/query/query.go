package query

import (
	"math"
	"traffic-sim/internal/road"
	"traffic-sim/internal/world"
)

type WorldQuery struct {
	world *world.World
}

func NewWorldQuery(w *world.World) *WorldQuery {
	return &WorldQuery{world: w}
}

func (q *WorldQuery) FindNearestNode(x, y, maxDistance float64) *road.Node {
	q.world.Mu.RLock()
	defer q.world.Mu.RUnlock()

	var nearest *road.Node
	minDist := maxDistance

	for _, node := range q.world.Nodes {
		dx := node.X - x
		dy := node.Y - y
		dist := math.Sqrt(dx*dx + dy*dy)

		if dist < minDist {
			minDist = dist
			nearest = node
		}
	}

	return nearest
}

func (q *WorldQuery) CanPlaceNodeAt(x, y, minDistance float64) bool {
	q.world.Mu.RLock()
	defer q.world.Mu.RUnlock()

	for _, node := range q.world.Nodes {
		dx := node.X - x
		dy := node.Y - y
		dist := math.Sqrt(dx*dx + dy*dy)

		if dist < minDistance {
			return false
		}
	}

	return true
}

func (q *WorldQuery) GetNodes() []*road.Node {
	q.world.Mu.RLock()
	defer q.world.Mu.RUnlock()
	
	nodes := make([]*road.Node, len(q.world.Nodes))
	copy(nodes, q.world.Nodes)
	return nodes
}

func (q *WorldQuery) GetRoads() []*road.Road {
	q.world.Mu.RLock()
	defer q.world.Mu.RUnlock()
	
	roads := make([]*road.Road, len(q.world.Roads))
	copy(roads, q.world.Roads)
	return roads
}