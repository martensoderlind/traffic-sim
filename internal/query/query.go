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

func (q *WorldQuery) FindNodeByID(id string) *road.Node {
	q.world.Mu.RLock()
	defer q.world.Mu.RUnlock()

	for _, node := range q.world.Nodes {
		if node.ID == id {
			return node
		}
	}

	return nil
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

func (q *WorldQuery) FindNearestRoad(x, y, maxDistance float64) (*road.Road, float64, float64) {
	q.world.Mu.RLock()
	defer q.world.Mu.RUnlock()

	var nearestRoad *road.Road
	var nearestX, nearestY float64
	minDist := maxDistance

	for _, rd := range q.world.Roads {
		px, py, dist := q.closestPointOnRoad(rd, x, y)
		
		if dist < minDist {
			minDist = dist
			nearestRoad = rd
			nearestX = px
			nearestY = py
		}
	}

	return nearestRoad, nearestX, nearestY
}

func (q *WorldQuery) closestPointOnRoad(rd *road.Road, x, y float64) (float64, float64, float64) {
	x1, y1 := rd.From.X, rd.From.Y
	x2, y2 := rd.To.X, rd.To.Y
	
	dx := x2 - x1
	dy := y2 - y1
	
	if dx == 0 && dy == 0 {
		return x1, y1, math.Sqrt((x-x1)*(x-x1) + (y-y1)*(y-y1))
	}
	
	t := ((x-x1)*dx + (y-y1)*dy) / (dx*dx + dy*dy)
	
	if t < 0 {
		t = 0
	} else if t > 1 {
		t = 1
	}
	
	px := x1 + t*dx
	py := y1 + t*dy
	
	dist := math.Sqrt((x-px)*(x-px) + (y-py)*(y-py))
	
	return px, py, dist
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

func (q *WorldQuery) GetOutgoingRoads(node *road.Node) []*road.Road {
	q.world.Mu.RLock()
	defer q.world.Mu.RUnlock()

	intersection := q.world.IntersectionsByNode[node.ID]
	if intersection == nil {
		return []*road.Road{}
	}

	outgoing := make([]*road.Road, len(intersection.Outgoing))
	copy(outgoing, intersection.Outgoing)
	return outgoing
}

func (q *WorldQuery) GetIncomingRoads(node *road.Node) []*road.Road {
	q.world.Mu.RLock()
	defer q.world.Mu.RUnlock()

	intersection := q.world.IntersectionsByNode[node.ID]
	if intersection == nil {
		return []*road.Road{}
	}

	incoming := make([]*road.Road, len(intersection.Incoming))
	copy(incoming, intersection.Incoming)
	return incoming
}