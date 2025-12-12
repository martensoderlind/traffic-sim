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
	waitingVehicles     map[string]float64
}

func NewRightOfWaySystem() *RightOfWaySystem {
	return &RightOfWaySystem{
		rules:               make(map[string]*road.RightOfWayRule),
		approachDistance:    60.0,
		yieldDistance:       30.0,
		stopDistance:        10.0,
		vehicleArrivalTimes: make(map[string]map[string]float64),
		waitingVehicles:     make(map[string]float64),
	}
}

// Reset clears internal state when world changes (e.g., load from file)
func (rows *RightOfWaySystem) Reset() {
	rows.rules = make(map[string]*road.RightOfWayRule)
	rows.vehicleArrivalTimes = make(map[string]map[string]float64)
	rows.waitingVehicles = make(map[string]float64)
}

func (rows *RightOfWaySystem) Update(w *world.World, dt float64) {
	w.Mu.Lock()
	defer w.Mu.Unlock()

	rows.updateRules(w)
	rows.updateVehicleArrivalTimes(w, dt)
	rows.applyRightOfWayRules(w, dt)
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

func (rows *RightOfWaySystem) updateVehicleArrivalTimes(w *world.World, dt float64) {
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

			// Lazy init: ensure the intersection map exists
			if rows.vehicleArrivalTimes[intersection.ID] == nil {
				rows.vehicleArrivalTimes[intersection.ID] = make(map[string]float64)
			}

			if _, exists := rows.vehicleArrivalTimes[intersection.ID][v.ID]; !exists {
				rows.vehicleArrivalTimes[intersection.ID][v.ID] = 0
			} else {
				rows.vehicleArrivalTimes[intersection.ID][v.ID] += dt
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

	for vehicleID := range rows.waitingVehicles {
		if !currentVehicles[vehicleID] {
			delete(rows.waitingVehicles, vehicleID)
		} else {
			rows.waitingVehicles[vehicleID] += dt
		}
	}
}

func (rows *RightOfWaySystem) applyRightOfWayRules(w *world.World, dt float64) {
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

		waitTime := rows.waitingVehicles[v.ID]
		if waitTime > 5.0 {
			rows.applyGracefulPassage(v, distToEnd)
			continue
		}

		shouldYield := rows.shouldVehicleYield(w, v, intersection, rule)

		if shouldYield {
			if v.Speed < 1.0 {
				rows.waitingVehicles[v.ID] = waitTime
			}
			rows.applyYieldBehavior(v, distToEnd)
		}
	}
}

func (rows *RightOfWaySystem) shouldVehicleYield(w *world.World, v *vehicle.Vehicle, intersection *road.Intersection, rule *road.RightOfWayRule) bool {
	conflictingVehicles := rows.findConflictingVehicles(w, v, intersection)

	if len(conflictingVehicles) == 0 {
		return false
	}

	for _, conflicting := range conflictingVehicles {
		vIsTurning := !road.IsMinorDirectionChange(v.Road, v.NextRoad)
		conflictingIsTurning := !road.IsMinorDirectionChange(conflicting.Road, conflicting.NextRoad)
		
		if vIsTurning && !conflictingIsTurning {
			return true
		}
		
		if !vIsTurning && conflictingIsTurning {
			continue
		}

		if rows.hasHigherPriority(v, conflicting, rule) {
			continue
		}

		if rows.hasHigherPriority(conflicting, v, rule) {
			return true
		}

		if rows.arrivedEarlier(intersection.ID, conflicting.ID, v.ID) {
			return true
		}

		if road.IsComingFromRight(v.Road, conflicting.Road) {
			distToEndOther := conflicting.Road.Length - conflicting.Distance
			if distToEndOther < rows.yieldDistance {
				return true
			}
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

		var otherIntersection *road.Intersection
		if other.InTransition {
			otherIntersection = w.IntersectionsByNode[other.NextRoad.To.ID]
		} else {
			otherIntersection = w.IntersectionsByNode[other.Road.To.ID]
		}
		
		if otherIntersection == nil || otherIntersection.ID != intersection.ID {
			continue
		}

		var distToEnd float64
		if other.InTransition {
			distToEnd = 5.0
		} else {
			distToEnd = other.Road.Length - other.Distance
		}
		
		if distToEnd > rows.approachDistance {
			continue
		}

		if rows.pathsConflict(v, other, intersection) {
			conflicting = append(conflicting, other)
		}
	}

	return conflicting
}

func (rows *RightOfWaySystem) pathsConflict(v1, v2 *vehicle.Vehicle, intersection *road.Intersection) bool {
	if v1.NextRoad == nil || v2.NextRoad == nil {
		return false
	}

	if v1.Road.ID == v2.Road.ID {
		return false
	}

	if v1.NextRoad == v2.NextRoad {
		return true
	}

	if road.IsLeftTurn(v1.Road, v1.NextRoad) {
		if road.IsStraight(v2.Road, v2.NextRoad) || road.IsLeftTurn(v2.Road, v2.NextRoad) {
			angle1 := road.CalculateRoadAngle(v1.Road)
			angle2 := road.CalculateRoadAngle(v2.Road)
			angleDiff := math.Abs(angle1 - angle2)
			
			if angleDiff > math.Pi/4 && angleDiff < 3*math.Pi/4 {
				return true
			}
		}
	}

	if road.IsStraight(v1.Road, v1.NextRoad) && road.IsLeftTurn(v2.Road, v2.NextRoad) {
		angle1 := road.CalculateRoadAngle(v1.Road)
		angle2 := road.CalculateRoadAngle(v2.Road)
		angleDiff := math.Abs(angle1 - angle2)
		
		if angleDiff > math.Pi/4 && angleDiff < 3*math.Pi/4 {
			return true
		}
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

	timeDiff := time1 - time2
	
	return timeDiff > 0.5
}

func (rows *RightOfWaySystem) applyYieldBehavior(v *vehicle.Vehicle, distToEnd float64) {
	if distToEnd < rows.stopDistance {
		targetSpeed := 0.0
		if v.Speed > targetSpeed {
			v.Speed = math.Max(0, v.Speed-15.0*0.016)
		}
		return
	}

	if distToEnd < rows.yieldDistance {
		ratio := (distToEnd - rows.stopDistance) / (rows.yieldDistance - rows.stopDistance)
		targetSpeed := v.Road.MaxSpeed * ratio * 0.4
		
		if v.Speed > targetSpeed {
			v.Speed = math.Max(targetSpeed, v.Speed-10.0*0.016)
		}
		return
	}

	slowdownRatio := (distToEnd - rows.yieldDistance) / (rows.approachDistance - rows.yieldDistance)
	targetSpeed := v.Road.MaxSpeed * (0.4 + slowdownRatio*0.6)
	
	if v.Speed > targetSpeed {
		v.Speed = math.Max(targetSpeed, v.Speed-8.0*0.016)
	}
}

func (rows *RightOfWaySystem) applyGracefulPassage(v *vehicle.Vehicle, distToEnd float64) {
	if distToEnd < rows.stopDistance*2 {
		targetSpeed := v.Road.MaxSpeed * 0.3
		if v.Speed < targetSpeed {
			v.Speed = math.Min(targetSpeed, v.Speed+5.0*0.016)
		}
	}
}