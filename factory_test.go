package geom

import (
	"testing"

	"vistarmedia.com/vistar/geom/geos-go"
	"vistarmedia.com/vistar/geom/geos-go/handle"
)

var fact = NewFactory(handle.NewPooledHandleProvider())

func TestNewEmptyPoint(t *testing.T) {
	p := fact.NewEmptyPoint()
	if p.Type() != POINT {
		t.Errorf("Unexpected geom type: %d", p.Type())
	}
	if isEmpty, err := p.IsEmpty(); err != nil {
		t.Error(err)
	} else if !isEmpty {
		t.Error("Point is not empty")
	}
}

func TestNewEmptyPolygon(t *testing.T) {
	p := fact.NewEmptyPolygon()
	if p.Type() != POLYGON {
		t.Errorf("Unexpected geom type: %d", p.Type())
	}
	if isEmpty, err := p.IsEmpty(); err != nil {
		t.Error(err)
	} else if !isEmpty {
		t.Error("Polygon is not empty")
	}
}

func TestNewPoint(t *testing.T) {
	p, err := fact.NewPoint(Coord{4, 2})
	if err != nil {
		t.Error(err)
	}
	if p.Type() != POINT {
		t.Errorf("Unexpected geom type: %d", p.Type())
	}
}

func TestNewLinearRing(t *testing.T) {
	coords := []Coord{
		{2, 2},
		{2, 4},
		{4, 4},
		{4, 2},
		{2, 2},
	}

	ring, err := fact.NewLinearRing(coords)
	if err != nil {
		t.Error(err)
	}
	if ring.Type() != LINEARRING {
		t.Errorf("Unexpected geom type: %d", ring.Type())
	}
	// LinearRings dont have an area
	if ring.Area() != 0 {
		t.Errorf("Expected area of 0, got %f", ring.Area())
	}
}

func TestNewLinearRingInvalid(t *testing.T) {
	// Not a closed ring
	coords := []Coord{{2, 2}, {2, 4}}
	_, err := fact.NewLinearRing(coords)
	if err != geos.ErrGeos {
		t.Errorf("Expected GEOS error, got %v", err)
	}
}

func TestNewPolygon(t *testing.T) {
	shell := []Coord{
		{2, 2},
		{2, 4},
		{4, 4},
		{4, 2},
		{2, 2},
	}

	poly, err := fact.NewPolygon(shell)
	if err != nil {
		t.Error(err)
	}
	if poly.Type() != POLYGON {
		t.Errorf("Unexpected geom type: %d", poly.Type())
	}
	if poly.Area() != 4 {
		t.Errorf("Expected area of 4, got %f", poly.Area())
	}
}

func TestNewPolygonWithHoles(t *testing.T) {
	//10x10 rect
	shell := []Coord{
		{0, 0},
		{10, 0},
		{10, 10},
		{0, 10},
		{0, 0},
	}

	//2x2 rect
	hole1 := []Coord{
		{2, 2},
		{2, 4},
		{4, 4},
		{4, 2},
		{2, 2},
	}

	//3x3 rect
	hole2 := []Coord{
		{6, 6},
		{6, 9},
		{9, 9},
		{9, 6},
		{6, 6},
	}

	poly, err := fact.NewPolygon(shell, hole1, hole2)
	if err != nil {
		t.Error(err)
	}
	if poly.Type() != POLYGON {
		t.Errorf("Unexpected geom type: %d", poly.Type())
	}
	// 10x10 - 3x3 - 2x2 = 87
	if poly.Area() != 87 {
		t.Errorf("Expected area of 87, got %f", poly.Area())
	}
}

func TestNewPolygonInvalid(t *testing.T) {
	// not a ring
	shell := []Coord{{2, 2}, {2, 4}}
	_, err := fact.NewPolygon(shell)
	if err != geos.ErrGeos {
		t.Errorf("Expected GEOS error, got %v", err)
	}
}
