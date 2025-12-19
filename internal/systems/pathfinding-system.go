package systems

import (
	"container/heap"
	"math"
	"math/rand"
	"traffic-sim/internal/geom"
	"traffic-sim/internal/road"
	"traffic-sim/internal/vehicle"
	"traffic-sim/internal/world"
)

const(
	curveRadius = 7.0
	startPointOffset = 12
)

type PathNode struct {
	nodeID   string
	distance float64
	index    int
}

type PriorityQueue []*PathNode

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].distance < pq[j].distance
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x any) {
	n := len(*pq)
	item := x.(*PathNode)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.index = -1
	*pq = old[0 : n-1]
	return item
}

type PathfindingSystem struct {
	roadGraph map[string][]*road.Road
}

func NewPathfindingSystem() *PathfindingSystem {
	return &PathfindingSystem{
		roadGraph: make(map[string][]*road.Road),
	}
}

func (ps *PathfindingSystem) Reset() {
	ps.roadGraph = make(map[string][]*road.Road)
}

func (ps *PathfindingSystem) Update(w *world.World, dt float64) {
	w.Mu.Lock()
	defer w.Mu.Unlock()

	ps.ensureRoadGraph(w)

	for _, v := range w.Vehicles {
		if v.InTransition {
			ps.updateTransition(v, dt)
			continue
		}
		
		if v.TargetDespawn == nil {
			ps.assignTarget(v, w)
		}
		
		if v.NextRoad == nil {
			threshold := v.Road.Length * 0.5
			if v.Road.Length < 60.0 {
				threshold = v.Road.Length * 0.3
			}
			
			if v.Distance > threshold {
				v.NextRoad = ps.findNextRoadToTarget(v, w)
			}
		}

		if v.Distance >= v.Road.Length-startPointOffset && v.Speed > 0 {
			if v.NextRoad != nil {
				ps.startTransition(v)
			} else {
				v.Speed = 0
				v.Distance = v.Road.Length
				x, y := v.Road.PosAt(v.Distance)
				v.Pos.X = x
				v.Pos.Y = y
			}
		}
	}
}

func (ps *PathfindingSystem) ensureRoadGraph(w *world.World) {
	if len(ps.roadGraph) == 0 {
		ps.buildRoadGraph(w)
	}
}

func (ps *PathfindingSystem) buildRoadGraph(w *world.World) {
	ps.roadGraph = make(map[string][]*road.Road)
	
	for _, intersection := range w.Intersections {
		ps.roadGraph[intersection.ID] = intersection.Outgoing
	}
}

func (ps *PathfindingSystem) assignTarget(v *vehicle.Vehicle, w *world.World) {
	activeDespawns := make([]*road.DespawnPoint, 0)
	for _, dp := range w.DespawnPoints {
		if dp.Enabled && dp.Node.ID != v.Road.From.ID {
			activeDespawns = append(activeDespawns, dp)
		}
	}
	
	if len(activeDespawns) == 0 {
		return
	}
	
	v.TargetDespawn = activeDespawns[rand.Intn(len(activeDespawns))]
}

func (ps *PathfindingSystem) findNextRoadToTarget(v *vehicle.Vehicle, w *world.World) *road.Road {
	if v.TargetDespawn == nil {
		return ps.findNextRoadRandom(v, w)
	}
	
	targetNodeID := v.TargetDespawn.Node.ID
	currentNodeID := v.Road.To.ID
	
	if currentNodeID == targetNodeID {
		for _, rd := range w.IntersectionsByNode[currentNodeID].Outgoing {
			if rd.ID == v.TargetDespawn.Road.ID {
				return v.TargetDespawn.Road
			}
		}
	}
	
	path := ps.findShortestPathDijkstra(currentNodeID, targetNodeID, w)
	
	if len(path) < 2 {
		return ps.findNextRoadRandom(v, w)
	}
	
	nextNodeID := path[1]
	
	intersection := w.IntersectionsByNode[currentNodeID]
	if intersection == nil {
		return ps.findNextRoadRandom(v, w)
	}
	
	for _, rd := range intersection.Outgoing {
		if notSameRoad(rd, v.Road) && rd.To.ID == nextNodeID {
			return rd
		}
	}
	
	return ps.findNextRoadRandom(v, w)
}

func (ps *PathfindingSystem) findShortestPathDijkstra(startNodeID, targetNodeID string, w *world.World) []string {
	if startNodeID == targetNodeID {
		return []string{startNodeID}
	}
	
	distances := make(map[string]float64)
	parent := make(map[string]string)
	visited := make(map[string]bool)
	
	for _, intersection := range w.Intersections {
		distances[intersection.ID] = math.Inf(1)
	}
	distances[startNodeID] = 0
	
	pq := make(PriorityQueue, 0)
	heap.Init(&pq)
	heap.Push(&pq, &PathNode{nodeID: startNodeID, distance: 0})
	
	for pq.Len() > 0 {
		current := heap.Pop(&pq).(*PathNode)
		
		if visited[current.nodeID] {
			continue
		}
		visited[current.nodeID] = true
		
		if current.nodeID == targetNodeID {
			return ps.reconstructPath(parent, startNodeID, targetNodeID)
		}
		
		intersection := w.IntersectionsByNode[current.nodeID]
		if intersection == nil {
			continue
		}
		
		for _, rd := range intersection.Outgoing {
			neighbor := rd.To.ID
			
			if visited[neighbor] {
				continue
			}
			
			newDistance := distances[current.nodeID] + rd.Length
			
			if newDistance < distances[neighbor] {
				distances[neighbor] = newDistance
				parent[neighbor] = current.nodeID
				heap.Push(&pq, &PathNode{nodeID: neighbor, distance: newDistance})
			}
		}
	}
	
	return nil
}

func (ps *PathfindingSystem) reconstructPath(parent map[string]string, start, end string) []string {
	path := []string{end}
	current := end
	
	for current != start {
		current = parent[current]
		path = append([]string{current}, path...)
	}
	
	return path
}

func (ps *PathfindingSystem) findNextRoadRandom(v *vehicle.Vehicle, w *world.World) *road.Road {
	intersection := w.IntersectionsByNode[v.Road.To.ID]
	if intersection == nil || len(intersection.Outgoing) == 0 {
		return nil
	}

	available := make([]*road.Road, 0, len(intersection.Outgoing))
	for _, r := range intersection.Outgoing {
		if notSameRoad(r, v.Road) {
			available = append(available, r)
		}
	}

	if len(available) == 0 {
		return nil
	}

	return available[rand.Intn(len(available))]
}

func (ps *PathfindingSystem) startTransition(v *vehicle.Vehicle) {
	fromRoad := v.Road
	toRoad := v.NextRoad
	
	startDist := 20.0
	if toRoad.Length < 40.0 {
		startDist = toRoad.Length * 0.3
	}
	x0, y0 := fromRoad.PosAt(fromRoad.Length-startPointOffset)
	x3, y3 := toRoad.PosAt(startDist)
	
	p0 := geom.Point{X: x0, Y: y0}
	p3 := geom.Point{X: x3, Y: y3}
	
	dirIn := geom.Point{
		X: fromRoad.To.X - fromRoad.From.X,
		Y: fromRoad.To.Y - fromRoad.From.Y,
	}
	lenIn := geom.Distance(geom.Point{}, dirIn)
	if lenIn > 0 {
		dirIn.X /= lenIn
		dirIn.Y /= lenIn
	}
	
	dirOut := geom.Point{
		X: toRoad.To.X - toRoad.From.X,
		Y: toRoad.To.Y - toRoad.From.Y,
	}
	lenOut := geom.Distance(geom.Point{}, dirOut)
	if lenOut > 0 {
		dirOut.X /= lenOut
		dirOut.Y /= lenOut
	}
	
	controlDist := curveRadius
	
	p1 := geom.Point{
		X: p0.X + dirIn.X*controlDist,
		Y: p0.Y + dirIn.Y*controlDist,
	}
	
	p2 := geom.Point{
		X: p3.X - dirOut.X*controlDist,
		Y: p3.Y - dirOut.Y*controlDist,
	}
	
	v.TransitionCurve = geom.NewCubicBezier(p0, p1, p2, p3)
	v.InTransition = true
	v.TransitionT = 0
	v.TransitionSpeed = v.Speed
}

func (ps *PathfindingSystem) updateTransition(v *vehicle.Vehicle, dt float64) {
	if v.TransitionCurve == nil {
		v.InTransition = false
		return
	}
	
	distanceToTravel := v.Speed * dt
	tStep := distanceToTravel / v.TransitionCurve.Length
	
	v.TransitionT += tStep
	
	if v.TransitionT >= 1.0 {
		v.TransitionT = 1.0
		v.InTransition = false
		v.Road = v.NextRoad
		v.NextRoad = nil
		
		startDist := 20.0
		if v.Road.Length < 40.0 {
			startDist = v.Road.Length * 0.3
		}
		
		v.Distance = startDist
		
		v.TransitionCurve = nil
		
		x, y := v.Road.PosAt(v.Distance)
		v.Pos.X = x
		v.Pos.Y = y
	} else {
		point := v.TransitionCurve.PointAt(v.TransitionT)
		v.Pos.X = point.X
		v.Pos.Y = point.Y
	}
}

func notSameRoad(r1, r2 *road.Road) bool {
	return !(r1.From == r2.To && r1.To == r2.From)
}