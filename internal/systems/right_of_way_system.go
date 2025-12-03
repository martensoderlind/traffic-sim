package systems

import (
	"math"
	"traffic-sim/internal/road"
	"traffic-sim/internal/vehicle"
	"traffic-sim/internal/world"
)

type RightOfWaySystem struct {
	rules               map[string]*road.RightOfWayRule
	approachDistance    float64
	yieldDistance       float64
	stopDistance        float64
	vehicleArrivalTimes map[string]map[string]float64
}

func NewRightOfWaySystem() *RightOfWaySystem {
	return &RightOfWaySystem{
		rules:               make(map[string]*road.RightOfWayRule),
		approachDistance:    60.0,
		yieldDistance:       30.0,
		stopDistance:        15.0,
		vehicleArrivalTimes: make(map[string]map[string]float64),
	}
}

func (rows *RightOfWaySystem) Update(w *world.World, dt float64) {
	w.Mu.Lock()
	defer w.Mu.Unlock()

	rows.updateRules(w)
	rows.updateVehicleArrivalTimes(w)
	rows.applyRightOfWayRules(w)
}

func (rows *RightOfWaySystem) updateRules(w *world.World) {
	for _, intersection := range w.Intersections {
		if rows.hasTrafficLight(w, intersection) {
			continue
		}

		if _, exists := rows.rules[intersection.ID]; !exists {
			rule := road.NewRightOfWayRule(intersection.ID)
			rule.Type = road.AnalyzeIntersection(intersection)
			rows.assignPriorities(rule, intersection)
			rows.rules[intersection.ID] = rule
		}
	}
}

func (rows *RightOfWaySystem) hasTrafficLight(w *world.World, intersection *road.Intersection) bool {
	for _, light := range w.TrafficLights {
		if light.Intersection.ID == intersection.ID {
			return true
		}
	}
	return false
}

func (rows *RightOfWaySystem) assignPriorities(rule *road.RightOfWayRule, intersection *road.Intersection) {
	allRoads := append([]*road.Road{}, intersection.Incoming...)
	allRoads = append(allRoads, intersection.Outgoing...)

	roadAngles := make(map[string]float64)
	for _, rd := range allRoads {
		roadAngles[rd.ID] = road.CalculateRoadAngle(rd)
	}

	for _, rd := range allRoads {
		priority := road.PriorityNormal
		
		if rd.MaxSpeed > 50.0 {
			priority = road.PriorityHigh
		} else if rd.MaxSpeed < 30.0 {
			priority = road.PriorityLow
		}

		rule.SetRoadPriority(rd.ID, priority)
	}
}

func (rows *RightOfWaySystem) updateVehicleArrivalTimes(w *world.World) {
	for intersectionID := range rows.rules {
		if rows.vehicleArrivalTimes[intersectionID] == nil {
			rows.vehicleArrivalTimes[intersectionID] = make(map[string]float64)
		}
	}

	currentVehicles := make(map[string]bool)

	for _, v := range w.Vehicles {
		if v.NextRoad == nil {
			continue
		}

		intersection := w.IntersectionsByNode[v.Road.To.ID]
		if intersection == nil {
			continue
		}

		distToEnd := v.Road.Length - v.Distance

		if distToEnd < rows.approachDistance {
			currentVehicles[v.ID] = true

			if _, exists := rows.vehicleArrivalTimes[intersection.ID][v.ID]; !exists {
				rows.vehicleArrivalTimes[intersection.ID][v.ID] = 0
			} else {
				rows.vehicleArrivalTimes[intersection.ID][v.ID] += 0.016
			}
		}
	}

	for intersectionID := range rows.vehicleArrivalTimes {
		for vehicleID := range rows.vehicleArrivalTimes[intersectionID] {
			if !currentVehicles[vehicleID] {
				delete(rows.vehicleArrivalTimes[intersectionID], vehicleID)
			}
		}
	}
}

func (rows *RightOfWaySystem) applyRightOfWayRules(w *world.World) {
	for _, v := range w.Vehicles {
		if v.NextRoad == nil {
			continue
		}

		intersection := w.IntersectionsByNode[v.Road.To.ID]
		if intersection == nil {
			continue
		}

		if rows.hasTrafficLight(w, intersection) {
			continue
		}

		distToEnd := v.Road.Length - v.Distance

		if distToEnd > rows.approachDistance {
			continue
		}

		rule := rows.rules[intersection.ID]
		if rule == nil {
			continue
		}

		shouldYield := rows.shouldVehicleYield(w, v, intersection, rule)

		if shouldYield {
			rows.applyYieldBehavior(v, distToEnd)
		}
	}
}

func (rows *RightOfWaySystem) shouldVehicleYield(w *world.World, v *vehicle.Vehicle, intersection *road.Intersection, rule *road.RightOfWayRule) bool {
	conflictingVehicles := rows.findConflictingVehicles(w, v, intersection)

	for _, conflicting := range conflictingVehicles {
		if rows.hasHigherPriority(v, conflicting, rule) {
			continue
		}

		if rows.arrivedEarlier(intersection.ID, conflicting.ID, v.ID) {
			return true
		}

		if road.IsComingFromRight(v.Road, conflicting.Road) {
			return true
		}
	}

	return false
}

func (rows *RightOfWaySystem) findConflictingVehicles(w *world.World, v *vehicle.Vehicle, intersection *road.Intersection) []*vehicle.Vehicle {
	conflicting := make([]*vehicle.Vehicle, 0)

	for _, other := range w.Vehicles {
		if other.ID == v.ID {
			continue
		}

		if other.NextRoad == nil {
			continue
		}

		otherIntersection := w.IntersectionsByNode[other.Road.To.ID]
		if otherIntersection == nil || otherIntersection.ID != intersection.ID {
			continue
		}

		distToEnd := other.Road.Length - other.Distance
		if distToEnd > rows.approachDistance {
			continue
		}

		if rows.pathsConflict(v, other) {
			conflicting = append(conflicting, other)
		}
	}

	return conflicting
}

func (rows *RightOfWaySystem) pathsConflict(v1, v2 *vehicle.Vehicle) bool {
	if v1.NextRoad == nil || v2.NextRoad == nil {
		return false
	}

	if v1.NextRoad == v2.NextRoad {
		return true
	}

	if road.IsLeftTurn(v1.Road, v1.NextRoad) && road.IsStraight(v2.Road, v2.NextRoad) {
		return true
	}

	if road.IsStraight(v1.Road, v1.NextRoad) && road.IsLeftTurn(v2.Road, v2.NextRoad) {
		return true
	}

	if road.IsLeftTurn(v1.Road, v1.NextRoad) && road.IsLeftTurn(v2.Road, v2.NextRoad) {
		return true
	}

	return false
}

func (rows *RightOfWaySystem) hasHigherPriority(v, conflicting *vehicle.Vehicle, rule *road.RightOfWayRule) bool {
	myPriority := rule.GetRoadPriority(v.Road.ID)
	theirPriority := rule.GetRoadPriority(conflicting.Road.ID)

	return myPriority > theirPriority
}

func (rows *RightOfWaySystem) arrivedEarlier(intersectionID, vehicleID1, vehicleID2 string) bool {
	times := rows.vehicleArrivalTimes[intersectionID]
	if times == nil {
		return false
	}

	time1, exists1 := times[vehicleID1]
	time2, exists2 := times[vehicleID2]

	if !exists1 || !exists2 {
		return false
	}

	return time1 > time2
}

func (rows *RightOfWaySystem) applyYieldBehavior(v *vehicle.Vehicle, distToEnd float64) {
	if distToEnd < rows.stopDistance {
		v.Speed = 0
		return
	}

	if distToEnd < rows.yieldDistance {
		ratio := distToEnd / rows.yieldDistance
		targetSpeed := v.Road.MaxSpeed * ratio * 0.3
		
		if v.Speed > targetSpeed {
			v.Speed = math.Max(targetSpeed, 0)
		}
		return
	}

	slowdownRatio := (distToEnd - rows.yieldDistance) / (rows.approachDistance - rows.yieldDistance)
	targetSpeed := v.Road.MaxSpeed * (0.3 + slowdownRatio*0.7)
	
	if v.Speed > targetSpeed {
		v.Speed = targetSpeed
	}
}