package vehicle

import "traffic-sim/internal/road"


type Vehicle struct {
ID string
Road *road.Road
Position float64
Speed float64
MaxSpeed float64
}