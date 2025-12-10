package commands

import (
	"math"
	"traffic-sim/internal/road"
	"traffic-sim/internal/world"
)

type CurveRoadCommand struct {
	Road         *road.Road
	IncomingRoad *road.Road
	OutgoingRoad *road.Road
}

func (c *CurveRoadCommand) Execute(w *world.World) error {
	if c.Road == nil {
		return nil
	}

	fromDir := c.Point(c.getOutgoingDirection(c.Road, c.Road.From))
	var fromTangent road.Point

	if c.IncomingRoad != nil {
		incomingDir := c.Point(c.getIncomingDirection(c.IncomingRoad, c.Road.From))
		fromTangent = road.Point{
			X: incomingDir.X,
			Y: incomingDir.Y,
		}
		fromTangent = c.Point(fromTangent.X, fromTangent.Y)
	} else {
		fromTangent = fromDir
	}

	toDir := c.Point(c.getIncomingDirection(c.Road, c.Road.To))
	var toTangent road.Point

	if c.OutgoingRoad != nil {
		outgoingDir := c.Point(c.getOutgoingDirection(c.OutgoingRoad, c.Road.To))
		toTangent = road.Point{
			X: outgoingDir.X,
			Y: outgoingDir.Y,
		}
		toTangent = c.Point(toTangent.X, toTangent.Y)
	} else {
		toTangent = toDir
	}

	const baseMultiplier = 0.8
	dot := fromTangent.X*toTangent.X + fromTangent.Y*toTangent.Y
	if dot > 1 {
		dot = 1
	}
	if dot < -1 {
		dot = -1
	}
	angle := math.Acos(dot)
	angleFactor := 0.25 + 0.75*(angle/math.Pi)
	distance := c.Road.Length * baseMultiplier * angleFactor
	if distance < 60 {
		distance = 60
	}
	if distance > 800 {
		distance = 800
	}

	controlP1 := road.Point{
		X: c.Road.From.X + fromTangent.X*distance,
		Y: c.Road.From.Y + fromTangent.Y*distance,
	}

	controlP2 := road.Point{
		X: c.Road.To.X - toTangent.X*distance,
		Y: c.Road.To.Y - toTangent.Y*distance,
	}

	c.Road.Curve = &road.RoadCurve{
		ControlP1: controlP1,
		ControlP2: controlP2,
	}

	if c.Road.ReverseRoad != nil {
		offsetDist := c.Road.Width * 0.5
		dx := c.Road.To.X - c.Road.From.X
		dy := c.Road.To.Y - c.Road.From.Y
		length := math.Sqrt(dx*dx + dy*dy)

		var offsetX, offsetY float64
		if length > 0 {
			offsetX = -dy / length * offsetDist
			offsetY = dx / length * offsetDist
		}

		c.Road.ReverseRoad.Curve = &road.RoadCurve{
			ControlP1: road.Point{
				X: controlP2.X + offsetX,
				Y: controlP2.Y + offsetY,
			},
			ControlP2: road.Point{
				X: controlP1.X + offsetX,
				Y: controlP1.Y + offsetY,
			},
		}
	}

	return nil
}

func (c *CurveRoadCommand) Point(dx, dy float64) road.Point {
	length := math.Sqrt(dx*dx + dy*dy)
	if length > 0 {
		dx /= length
		dy /= length
	}
	return road.Point{X: dx, Y: dy}
}

func (c *CurveRoadCommand) getIncomingDirection(r *road.Road, node *road.Node) (float64, float64) {
	var dx, dy float64
	if r.To == node {
		dx = r.To.X - r.From.X
		dy = r.To.Y - r.From.Y
	} else {
		dx = r.From.X - r.To.X
		dy = r.From.Y - r.To.Y
	}

	return dx, dy
}

func (c *CurveRoadCommand) getOutgoingDirection(r *road.Road, node *road.Node) (float64, float64) {
	var dx, dy float64
	if r.From == node {
		dx = r.To.X - r.From.X
		dy = r.To.Y - r.From.Y
	} else {
		dx = r.From.X - r.To.X
		dy = r.From.Y - r.To.Y
	}

	return dx, dy
}
