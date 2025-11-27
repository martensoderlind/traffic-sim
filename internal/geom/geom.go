package geom

import "math"

type Point struct {
	X,Y float64
}

func Distance(a,b Point) float64{
	dx:= a.X-b.X
	dy:= a.Y-b.Y
	return math.Hypot(dx, dy) 
}
