package vehicle

import (
	"math"
	"traffic-sim/internal/road"
)

type Vec2 struct {
	X float64
	Y float64
}

type Vehicle struct {
	ID       string
	Road     *road.Road
	NextRoad *road.Road
	Distance float64
	Speed    float64
	Pos      Vec2
}

func (v *Vehicle) Position() Vec2 {
	return v.Pos
}

func (v *Vehicle) GetAngle() float64 {
	if v.Road == nil {
		return 0
	}

	dx := v.Road.To.X - v.Road.From.X
	dy := v.Road.To.Y - v.Road.From.Y

	return math.Atan2(dy, dx)+math.Pi/2
}