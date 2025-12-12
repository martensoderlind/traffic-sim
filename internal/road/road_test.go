package road

import (
	"math"
	"testing"
)

const epsilon = 1e-9

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) < epsilon
}

func TestUpdateLengthNoOffsets(t *testing.T) {
	n1 := &Node{ID: "n1", X: 0, Y: 0}
	n2 := &Node{ID: "n2", X: 3, Y: 4}

	r := NewRoad("r1", n1, n2, 40.0)
	r.UpdateLength()

	expected := 5.0
	if !almostEqual(r.Length, expected) {
		t.Errorf("Expected length %.2f, got %.2f", expected, r.Length)
	}
}

func TestUpdateLengthWithOffsets(t *testing.T) {
	n1 := &Node{ID: "n1", X: 0, Y: 0}
	n2 := &Node{ID: "n2", X: 10, Y: 0}

	r := NewRoad("r1", n1, n2, 40.0)
	r.StartOffset = Point{X: 1, Y: 0}
	r.EndOffset = Point{X: 1, Y: 0}
	r.UpdateLength()

	expected := 10.0
	if !almostEqual(r.Length, expected) {
		t.Errorf("Expected length %.2f, got %.2f", expected, r.Length)
	}
}

func TestUpdateLengthLoopWithOpposingOffsets(t *testing.T) {
	n := &Node{ID: "n1", X: 5, Y: 5}

	r := NewRoad("loop", n, n, 40.0)
	r.StartOffset = Point{X: 6, Y: 0}
	r.EndOffset = Point{X: -6, Y: 0}
	r.UpdateLength()

	expected := 12.0
	if !almostEqual(r.Length, expected) {
		t.Errorf("Expected length %.2f, got %.2f", expected, r.Length)
	}
}

func TestPosAtNoOffsetsStart(t *testing.T) {
	n1 := &Node{ID: "n1", X: 0, Y: 0}
	n2 := &Node{ID: "n2", X: 10, Y: 0}

	r := NewRoad("r1", n1, n2, 40.0)
	r.UpdateLength()

	x, y := r.PosAt(0)
	if !almostEqual(x, 0) || !almostEqual(y, 0) {
		t.Errorf("Expected (0, 0), got (%.2f, %.2f)", x, y)
	}
}

func TestPosAtNoOffsetsEnd(t *testing.T) {
	n1 := &Node{ID: "n1", X: 0, Y: 0}
	n2 := &Node{ID: "n2", X: 10, Y: 0}

	r := NewRoad("r1", n1, n2, 40.0)
	r.UpdateLength()

	x, y := r.PosAt(r.Length)
	if !almostEqual(x, 10) || !almostEqual(y, 0) {
		t.Errorf("Expected (10, 0), got (%.2f, %.2f)", x, y)
	}
}

func TestPosAtNoOffsetsMidpoint(t *testing.T) {
	n1 := &Node{ID: "n1", X: 0, Y: 0}
	n2 := &Node{ID: "n2", X: 10, Y: 0}

	r := NewRoad("r1", n1, n2, 40.0)
	r.UpdateLength()

	x, y := r.PosAt(r.Length / 2)
	if !almostEqual(x, 5) || !almostEqual(y, 0) {
		t.Errorf("Expected (5, 0), got (%.2f, %.2f)", x, y)
	}
}

func TestPosAtWithOffsetsStart(t *testing.T) {
	n1 := &Node{ID: "n1", X: 0, Y: 0}
	n2 := &Node{ID: "n2", X: 10, Y: 0}

	r := NewRoad("r1", n1, n2, 40.0)
	r.StartOffset = Point{X: 2, Y: 0}
	r.EndOffset = Point{X: 1, Y: 0}
	r.UpdateLength()

	x, y := r.PosAt(0)
	if !almostEqual(x, 2) || !almostEqual(y, 0) {
		t.Errorf("Expected (2, 0), got (%.2f, %.2f)", x, y)
	}
}

func TestPosAtWithOffsetsEnd(t *testing.T) {
	n1 := &Node{ID: "n1", X: 0, Y: 0}
	n2 := &Node{ID: "n2", X: 10, Y: 0}

	r := NewRoad("r1", n1, n2, 40.0)
	r.StartOffset = Point{X: 2, Y: 0}
	r.EndOffset = Point{X: 1, Y: 0}
	r.UpdateLength()

	x, y := r.PosAt(r.Length)
	if !almostEqual(x, 11) || !almostEqual(y, 0) {
		t.Errorf("Expected (11, 0), got (%.2f, %.2f)", x, y)
	}
}

func TestPosAtWithOffsetsLateral(t *testing.T) {
	n1 := &Node{ID: "n1", X: 0, Y: 0}
	n2 := &Node{ID: "n2", X: 10, Y: 0}

	r := NewRoad("r1", n1, n2, 40.0)
	r.StartOffset = Point{X: 0, Y: 2}
	r.EndOffset = Point{X: 0, Y: 3}
	r.UpdateLength()

	x, y := r.PosAt(0)
	if !almostEqual(x, 0) || !almostEqual(y, 2) {
		t.Errorf("Expected (0, 2), got (%.2f, %.2f)", x, y)
	}
}

func TestPosAtOutOfBoundsNegative(t *testing.T) {
	n1 := &Node{ID: "n1", X: 0, Y: 0}
	n2 := &Node{ID: "n2", X: 10, Y: 0}

	r := NewRoad("r1", n1, n2, 40.0)
	r.UpdateLength()

	x, y := r.PosAt(-5)
	if !almostEqual(x, 0) || !almostEqual(y, 0) {
		t.Errorf("Expected (0, 0), got (%.2f, %.2f)", x, y)
	}
}

func TestPosAtOutOfBoundsPositive(t *testing.T) {
	n1 := &Node{ID: "n1", X: 0, Y: 0}
	n2 := &Node{ID: "n2", X: 10, Y: 0}

	r := NewRoad("r1", n1, n2, 40.0)
	r.UpdateLength()

	x, y := r.PosAt(r.Length + 10)
	if !almostEqual(x, 10) || !almostEqual(y, 0) {
		t.Errorf("Expected (10, 0), got (%.2f, %.2f)", x, y)
	}
}

func TestPosAtZeroLengthRoad(t *testing.T) {
	n := &Node{ID: "n1", X: 5, Y: 3}

	r := NewRoad("zero", n, n, 40.0)
	r.UpdateLength()

	x, y := r.PosAt(0)
	if !almostEqual(x, 5) || !almostEqual(y, 3) {
		t.Errorf("Expected (5, 3), got (%.2f, %.2f)", x, y)
	}
}

func TestPosAtLoopRoad(t *testing.T) {
	n := &Node{ID: "n1", X: 0, Y: 0}

	r := NewRoad("loop", n, n, 40.0)
	r.StartOffset = Point{X: 10, Y: 0}
	r.EndOffset = Point{X: -10, Y: 0}
	r.UpdateLength()

	x, y := r.PosAt(r.Length / 2)
	if !almostEqual(x, 0) || !almostEqual(y, 0) {
		t.Errorf("Expected (0, 0), got (%.2f, %.2f)", x, y)
	}
}

func TestCubicBezierPointAtStart(t *testing.T) {
	p0 := Point{X: 0, Y: 0}
	p1 := Point{X: 1, Y: 1}
	p2 := Point{X: 2, Y: 1}
	p3 := Point{X: 3, Y: 0}

	result := cubicBezierPoint(p0, p1, p2, p3, 0)
	if !almostEqual(result.X, 0) || !almostEqual(result.Y, 0) {
		t.Errorf("Expected (0, 0), got (%.2f, %.2f)", result.X, result.Y)
	}
}

func TestCubicBezierPointAtEnd(t *testing.T) {
	p0 := Point{X: 0, Y: 0}
	p1 := Point{X: 1, Y: 1}
	p2 := Point{X: 2, Y: 1}
	p3 := Point{X: 3, Y: 0}

	result := cubicBezierPoint(p0, p1, p2, p3, 1)
	if !almostEqual(result.X, 3) || !almostEqual(result.Y, 0) {
		t.Errorf("Expected (3, 0), got (%.2f, %.2f)", result.X, result.Y)
	}
}

func TestCubicBezierPointAtMidpoint(t *testing.T) {
	p0 := Point{X: 0, Y: 0}
	p1 := Point{X: 2, Y: 2}
	p2 := Point{X: 4, Y: 2}
	p3 := Point{X: 6, Y: 0}

	result := cubicBezierPoint(p0, p1, p2, p3, 0.5)
	if !almostEqual(result.X, 3) || !almostEqual(result.Y, 1.5) {
		t.Errorf("Expected (3, 1.5), got (%.2f, %.2f)", result.X, result.Y)
	}
}

func TestNewRoadInitialLength(t *testing.T) {
	n1 := &Node{ID: "n1", X: 0, Y: 0}
	n2 := &Node{ID: "n2", X: 6, Y: 8}

	r := NewRoad("r1", n1, n2, 40.0)

	expected := 10.0
	if !almostEqual(r.Length, expected) {
		t.Errorf("Expected length %.2f, got %.2f", expected, r.Length)
	}
}

func TestNewRoadDefaultWidth(t *testing.T) {
	n1 := &Node{ID: "n1", X: 0, Y: 0}
	n2 := &Node{ID: "n2", X: 10, Y: 0}

	r := NewRoad("r1", n1, n2, 40.0)

	expected := 12.0
	if !almostEqual(r.Width, expected) {
		t.Errorf("Expected width %.2f, got %.2f", expected, r.Width)
	}
}

func TestNewRoadMaxSpeed(t *testing.T) {
	n1 := &Node{ID: "n1", X: 0, Y: 0}
	n2 := &Node{ID: "n2", X: 10, Y: 0}
	maxSpeed := 60.0

	r := NewRoad("r1", n1, n2, maxSpeed)

	if !almostEqual(r.MaxSpeed, maxSpeed) {
		t.Errorf("Expected maxSpeed %.2f, got %.2f", maxSpeed, r.MaxSpeed)
	}
}

func TestOffsetsCombined(t *testing.T) {
	n1 := &Node{ID: "n1", X: 0, Y: 0}
	n2 := &Node{ID: "n2", X: 10, Y: 0}

	r := NewRoad("r1", n1, n2, 40.0)
	r.StartOffset = Point{X: 1, Y: 1}
	r.EndOffset = Point{X: 2, Y: -1}
	r.UpdateLength()

	expected := math.Sqrt(125)
	if !almostEqual(r.Length, expected) {
		t.Errorf("Expected length %.4f, got %.4f", expected, r.Length)
	}

	x, y := r.PosAt(0)
	if !almostEqual(x, 1) || !almostEqual(y, 1) {
		t.Errorf("At t=0: Expected (1, 1), got (%.2f, %.2f)", x, y)
	}

	x, y = r.PosAt(r.Length)
	if !almostEqual(x, 12) || !almostEqual(y, -1) {
		t.Errorf("At t=end: Expected (12, -1), got (%.2f, %.2f)", x, y)
	}
}
