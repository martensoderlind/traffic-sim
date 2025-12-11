package road

import "math"

type RoadPriority int

const (
	PriorityLow RoadPriority = iota
	PriorityNormal
	PriorityHigh
)

type IntersectionType int

const (
	IntersectionUncontrolled IntersectionType = iota
	IntersectionTWay
	IntersectionThreeWay
	IntersectionFourWay
	IntersectionRoundabout
)

type RightOfWayRule struct {
	IntersectionID string
	Type           IntersectionType
	RoadPriorities map[string]RoadPriority
}

func NewRightOfWayRule(intersectionID string) *RightOfWayRule {
	return &RightOfWayRule{
		IntersectionID: intersectionID,
		RoadPriorities: make(map[string]RoadPriority),
	}
}

func (r *RightOfWayRule) SetRoadPriority(roadID string, priority RoadPriority) {
	r.RoadPriorities[roadID] = priority
}

func (r *RightOfWayRule) GetRoadPriority(roadID string) RoadPriority {
	if priority, exists := r.RoadPriorities[roadID]; exists {
		return priority
	}
	return PriorityNormal
}

func (r *RightOfWayRule) HasPriority(approachingRoadID string, conflictingRoadID string) bool {
	approachingPriority := r.GetRoadPriority(approachingRoadID)
	conflictingPriority := r.GetRoadPriority(conflictingRoadID)
	
	if approachingPriority > conflictingPriority {
		return true
	}
	
	if approachingPriority < conflictingPriority {
		return false
	}
	
	return false
}

func AnalyzeIntersection(intersection *Intersection) IntersectionType {
	
	totalRoads := countUniqueRoads(intersection.Incoming, intersection.Outgoing)
	
	if totalRoads <= 2 {
		return IntersectionTWay
	} else if totalRoads == 3 {
		return IntersectionThreeWay
	} else if totalRoads >= 4 {
		return IntersectionFourWay
	}
	
	return IntersectionUncontrolled
}

func countUniqueRoads(incoming, outgoing []*Road) int {
	seen := make(map[string]bool)
	
	for _, r := range incoming {
		if r.ReverseRoad != nil {
			roadPair := getRoadPairID(r, r.ReverseRoad)
			seen[roadPair] = true
		} else {
			seen[r.ID] = true
		}
	}
	
	for _, r := range outgoing {
		if r.ReverseRoad != nil {
			roadPair := getRoadPairID(r, r.ReverseRoad)
			seen[roadPair] = true
		} else {
			seen[r.ID] = true
		}
	}
	
	return len(seen)
}

func getRoadPairID(r1, r2 *Road) string {
	if r1.ID < r2.ID {
		return r1.ID + "-" + r2.ID
	}
	return r2.ID + "-" + r1.ID
}

func CalculateRoadAngle(road *Road) float64 {
	dx := road.To.X - road.From.X
	dy := road.To.Y - road.From.Y
	return math.Atan2(dy, dx)
}

func IsRightTurn(fromRoad, toRoad *Road) bool {
	angleFrom := CalculateRoadAngle(fromRoad)
	angleTo := CalculateRoadAngle(toRoad)
	
	angleDiff := angleTo - angleFrom
	
	for angleDiff < -math.Pi {
		angleDiff += 2 * math.Pi
	}
	for angleDiff > math.Pi {
		angleDiff -= 2 * math.Pi
	}
	
	return angleDiff < 0 && angleDiff > -math.Pi
}

func IsLeftTurn(fromRoad, toRoad *Road) bool {
	angleFrom := CalculateRoadAngle(fromRoad)
	angleTo := CalculateRoadAngle(toRoad)
	
	angleDiff := angleTo - angleFrom
	
	for angleDiff < -math.Pi {
		angleDiff += 2 * math.Pi
	}
	for angleDiff > math.Pi {
		angleDiff -= 2 * math.Pi
	}
	
	return angleDiff > 0 && angleDiff < math.Pi
}

func IsMinorDirectionChange(fromRoad, toRoad *Road) bool {
	angleFrom := CalculateRoadAngle(fromRoad)
	angleTo := CalculateRoadAngle(toRoad)
	
	angleDiff := angleTo - angleFrom
	
	for angleDiff < -math.Pi {
		angleDiff += 2 * math.Pi
	}
	for angleDiff > math.Pi {
		angleDiff -= 2 * math.Pi
	}
	
	threshold := 25.0 * math.Pi / 180.0
	return math.Abs(angleDiff) < threshold
}

func IsStraight(fromRoad, toRoad *Road) bool {
	angleFrom := CalculateRoadAngle(fromRoad)
	angleTo := CalculateRoadAngle(toRoad)
	
	angleDiff := math.Abs(angleTo - angleFrom)
	
	threshold := math.Pi / 6
	
	return angleDiff < threshold || math.Abs(angleDiff-math.Pi) < threshold
}

func GetRelativeAngle(observerRoad, targetRoad *Road) float64 {
	observerAngle := CalculateRoadAngle(observerRoad)
	targetAngle := CalculateRoadAngle(targetRoad)
	
	angleDiff := targetAngle - observerAngle
	
	for angleDiff < -math.Pi {
		angleDiff += 2 * math.Pi
	}
	for angleDiff > math.Pi {
		angleDiff -= 2 * math.Pi
	}
	
	return angleDiff
}

func IsComingFromRight(observerRoad, targetRoad *Road) bool {
	angle := GetRelativeAngle(observerRoad, targetRoad)
	return angle > -math.Pi/2 && angle < math.Pi/2
}