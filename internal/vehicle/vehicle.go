package vehicle

import "traffic-sim/internal/road"



type Vec2 struct {
    X float64
    Y float64
}

type Vehicle struct {
ID string
Road *road.Road
Distance float64
Speed float64
Pos Vec2
}


func (v *Vehicle) Position() Vec2 {
    return v.Pos
}
