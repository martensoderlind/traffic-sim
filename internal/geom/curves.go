package geom

import "math"

type BezierCurve struct {
	P0, P1, P2, P3 Point
	Length         float64
}

func NewQuadraticBezier(p0, p1, p2 Point) *BezierCurve {
	return &BezierCurve{
		P0:     p0,
		P1:     p1,
		P2:     p2,
		Length: estimateBezierLength(p0, p1, p2, Point{}),
	}
}

func NewCubicBezier(p0, p1, p2, p3 Point) *BezierCurve {
	return &BezierCurve{
		P0:     p0,
		P1:     p1,
		P2:     p2,
		P3:     p3,
		Length: estimateBezierLength(p0, p1, p2, p3),
	}
}

func (b *BezierCurve) PointAt(t float64) Point {
	if b.P3.X == 0 && b.P3.Y == 0 {
		return quadraticBezierPoint(b.P0, b.P1, b.P2, t)
	}
	return cubicBezierPoint(b.P0, b.P1, b.P2, b.P3, t)
}

func (b *BezierCurve) TangentAt(t float64) Point {
	if b.P3.X == 0 && b.P3.Y == 0 {
		return quadraticBezierTangent(b.P0, b.P1, b.P2, t)
	}
	return cubicBezierTangent(b.P0, b.P1, b.P2, b.P3, t)
}

func quadraticBezierPoint(p0, p1, p2 Point, t float64) Point {
	mt := 1 - t
	mt2 := mt * mt
	t2 := t * t
	
	return Point{
		X: mt2*p0.X + 2*mt*t*p1.X + t2*p2.X,
		Y: mt2*p0.Y + 2*mt*t*p1.Y + t2*p2.Y,
	}
}

func cubicBezierPoint(p0, p1, p2, p3 Point, t float64) Point {
	mt := 1 - t
	mt2 := mt * mt
	mt3 := mt2 * mt
	t2 := t * t
	t3 := t2 * t
	
	return Point{
		X: mt3*p0.X + 3*mt2*t*p1.X + 3*mt*t2*p2.X + t3*p3.X,
		Y: mt3*p0.Y + 3*mt2*t*p1.Y + 3*mt*t2*p2.Y + t3*p3.Y,
	}
}

func quadraticBezierTangent(p0, p1, p2 Point, t float64) Point {
	mt := 1 - t
	
	dx := 2*mt*(p1.X-p0.X) + 2*t*(p2.X-p1.X)
	dy := 2*mt*(p1.Y-p0.Y) + 2*t*(p2.Y-p1.Y)
	
	length := math.Sqrt(dx*dx + dy*dy)
	if length > 0 {
		dx /= length
		dy /= length
	}
	
	return Point{X: dx, Y: dy}
}

func cubicBezierTangent(p0, p1, p2, p3 Point, t float64) Point {
	mt := 1 - t
	mt2 := mt * mt
	t2 := t * t
	
	dx := 3*mt2*(p1.X-p0.X) + 6*mt*t*(p2.X-p1.X) + 3*t2*(p3.X-p2.X)
	dy := 3*mt2*(p1.Y-p0.Y) + 6*mt*t*(p2.Y-p1.Y) + 3*t2*(p3.Y-p2.Y)
	
	length := math.Sqrt(dx*dx + dy*dy)
	if length > 0 {
		dx /= length
		dy /= length
	}
	
	return Point{X: dx, Y: dy}
}

func estimateBezierLength(p0, p1, p2, p3 Point) float64 {
	steps := 20
	length := 0.0
	
	isCubic := !(p3.X == 0 && p3.Y == 0)
	
	var prevPoint Point
	for i := 0; i <= steps; i++ {
		t := float64(i) / float64(steps)
		
		var point Point
		if isCubic {
			point = cubicBezierPoint(p0, p1, p2, p3, t)
		} else {
			point = quadraticBezierPoint(p0, p1, p2, t)
		}
		
		if i > 0 {
			length += Distance(prevPoint, point)
		}
		prevPoint = point
	}
	
	return length
}

func CalculateControlPoint(from, intersection, to Point, curvature float64) Point {
	dirIn := Point{
		X: intersection.X - from.X,
		Y: intersection.Y - from.Y,
	}
	lenIn := math.Sqrt(dirIn.X*dirIn.X + dirIn.Y*dirIn.Y)
	if lenIn > 0 {
		dirIn.X /= lenIn
		dirIn.Y /= lenIn
	}
	
	dirOut := Point{
		X: to.X - intersection.X,
		Y: to.Y - intersection.Y,
	}
	lenOut := math.Sqrt(dirOut.X*dirOut.X + dirOut.Y*dirOut.Y)
	if lenOut > 0 {
		dirOut.X /= lenOut
		dirOut.Y /= lenOut
	}
	
	avgDir := Point{
		X: (dirIn.X + dirOut.X) / 2,
		Y: (dirIn.Y + dirOut.Y) / 2,
	}
	lenAvg := math.Sqrt(avgDir.X*avgDir.X + avgDir.Y*avgDir.Y)
	if lenAvg > 0 {
		avgDir.X /= lenAvg
		avgDir.Y /= lenAvg
	}
	
	distance := math.Min(lenIn, lenOut) * curvature
	
	return Point{
		X: intersection.X + avgDir.X*distance,
		Y: intersection.Y + avgDir.Y*distance,
	}
}