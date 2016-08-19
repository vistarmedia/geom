package geom

import (
	"testing"

	"vistarmedia.com/vistar/geom/geos-go/handle"
)

func TestIntersection(t *testing.T) {
	fact := NewFactory(handle.NewPooledHandleProvider())
	square1, err := fact.NewPolygon([]Coord{
		{0, 0},
		{10, 0},
		{10, 10},
		{0, 10},
		{0, 0},
	})
	if err != nil {
		t.Error(err)
	}

	// half the square overlaps with square1
	square2, err := fact.NewPolygon([]Coord{
		{5, 0},
		{15, 0},
		{15, 10},
		{5, 10},
		{5, 0},
	})
	if err != nil {
		t.Error(err)
	}

	if square1.Area() != 100 || square2.Area() != 100 {
		t.Errorf("Unexpected area. s1: %f, s2: %f", square1.Area(), square2.Area())
	}
	intersection, err := square1.Intersection(square2)
	if err != nil {
		t.Error(err)
	}
	if intersection.Area() != 50 {
		t.Errorf("Expected area of 50, got %f", intersection.Area())
	}
}

func TestPreparedCovers(t *testing.T) {
	fact := NewFactory(handle.NewPooledHandleProvider())
	square1, err := fact.NewPolygon([]Coord{
		{0, 0},
		{10, 0},
		{10, 10},
		{0, 10},
		{0, 0},
	})
	if err != nil {
		t.Error(err)
	}
	prep := square1.Prepared()
	insidePoint, err := fact.NewPoint(Coord{5, 5})
	if err != nil {
		t.Error(err)
	}
	outsidePoint, err := fact.NewPoint(Coord{14, 5})
	if err != nil {
		t.Error(err)
	}
	insideCovers, err := prep.Covers(insidePoint)
	if err != nil {
		t.Error(err)
	}
	outsideCovers, err := prep.Covers(outsidePoint)
	if err != nil {
		t.Error(err)
	}
	if !insideCovers || outsideCovers {
		t.Errorf("Inside point: %t, outside point: %t", insideCovers, outsideCovers)
	}
}
