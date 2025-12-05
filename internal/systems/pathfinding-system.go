package systems

import (
	"math/rand"
	"traffic-sim/internal/geom"
	"traffic-sim/internal/road"
	"traffic-sim/internal/vehicle"
	"traffic-sim/internal/world"
)

type PathfindingSystem struct{}

func NewPathfindingSystem() *PathfindingSystem {
	return &PathfindingSystem{}
}

func (ps *PathfindingSystem) Update(w *world.World, dt float64) {
	w.Mu.Lock()
	defer w.Mu.Unlock()

	for _, v := range w.Vehicles {
		if v.InTransition {
			ps.updateTransition(v, dt)
			continue
		}
		
		if v.NextRoad == nil {
			threshold := v.Road.Length * 0.5
			if v.Road.Length < 60.0 {
				threshold = v.Road.Length * 0.3
			}
			
			if v.Distance > threshold {
				v.NextRoad = ps.findNextRoad(w, v)
			}
		}

		if v.Distance >= v.Road.Length && v.Speed > 0 {
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


func (ps *PathfindingSystem) startTransition(v *vehicle.Vehicle) {
	fromRoad := v.Road
	toRoad := v.NextRoad
	
	intersectionNode := fromRoad.To
	
	p0 := geom.Point{X: fromRoad.To.X, Y: fromRoad.To.Y}
	p3 := geom.Point{X: toRoad.From.X, Y: toRoad.From.Y}
	
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
	
	controlDist := 15.0
	
	p1 := geom.Point{
		X: intersectionNode.X + dirIn.X*controlDist,
		Y: intersectionNode.Y + dirIn.Y*controlDist,
	}
	
	p2 := geom.Point{
		X: intersectionNode.X + dirOut.X*controlDist,
		Y: intersectionNode.Y + dirOut.Y*controlDist,
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
		v.Distance = 0
		v.TransitionCurve = nil
		
		x, y := v.Road.PosAt(0)
		v.Pos.X = x
		v.Pos.Y = y
	} else {
		point := v.TransitionCurve.PointAt(v.TransitionT)
		v.Pos.X = point.X
		v.Pos.Y = point.Y
	}
}

func (ps *PathfindingSystem) findNextRoad(w *world.World, v *vehicle.Vehicle) *road.Road {
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

func notSameRoad(r1, r2 *road.Road) bool {
	return !(r1.From == r2.To && r1.To == r2.From)
}