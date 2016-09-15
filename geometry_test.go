package geom

import (
	"fmt"
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

func TestCoordFromPoint(t *testing.T) {
	fact := NewFactory(handle.NewPooledHandleProvider())
	point, err := fact.NewPoint(Coord{5, 10})
	if err != nil {
		t.Error(err)
	}
	coord, err := point.Coord()
	if err != nil {
		t.Error(err)
	}
	if coord.X != 5 || coord.Y != 10 {
		t.Errorf("Something aint right, X: %f, Y: %f", coord.X, coord.Y)
	}
}

func TestCoordsFromLinearRing(t *testing.T) {
	coords := []Coord{
		{0, 0},
		{10, 0},
		{10, 10},
		{0, 10},
		{0, 0},
	}
	fact := NewFactory(handle.NewPooledHandleProvider())
	lr, err := fact.NewLinearRing(coords)
	if err != nil {
		t.Error(err)
	}
	outCoords, err := lr.Coords()
	if err != nil {
		t.Error(err)
	}
	if !compareCoordSlice(coords, outCoords) {
		t.Errorf("Coordinates dont match, in: %v, out %v", coords, outCoords)
	}
}

func TestCoordsFromPolygon(t *testing.T) {
	shell := []Coord{
		{0, 0},
		{10, 0},
		{10, 10},
		{0, 10},
		{0, 0},
	}

	hole1 := []Coord{
		{2, 2},
		{2, 4},
		{4, 4},
		{4, 2},
		{2, 2},
	}

	hole2 := []Coord{
		{6, 6},
		{6, 9},
		{9, 9},
		{9, 6},
		{6, 6},
	}

	fact := NewFactory(handle.NewPooledHandleProvider())
	poly, err := fact.NewPolygon(shell, hole1, hole2)
	if err != nil {
		t.Error(err)
	}
	outShellCoords, err := poly.Shell()
	if err != nil {
		t.Error(err)
	}
	if !compareCoordSlice(shell, outShellCoords) {
		t.Errorf("Shell coords dont match, in: %v, out %v", shell, outShellCoords)
	}
	holes, err := poly.Holes()
	if err != nil {
		t.Error(err)
	}
	if len(holes) != 2 {
		t.Errorf("Wrong number of holes, %d", len(holes))
	}
	if !compareCoordSlice(hole1, holes[0]) {
		t.Errorf("Hole coords dont match, in: %v, out %v", hole1, holes[0])
	}
	if !compareCoordSlice(hole2, holes[1]) {
		t.Errorf("Hole coords dont match, in: %v, out %v", hole2, holes[1])
	}
}

func TestGeometriesOfGeom(t *testing.T) {
	coords := []Coord{
		{0, 0},
		{10, 0},
		{10, 10},
		{0, 0},
	}
	fact := NewFactory(handle.NewPooledHandleProvider())
	lr, err := fact.NewLinearRing(coords)
	if err != nil {
		t.Fatal(err)
	}

	n, err := lr.NumGeometries()
	if err != nil {
		t.Fatal(err)
	}
	if n != 1 {
		t.Error("expected 1 geometry")
	}

	gs, err := lr.Geometries()
	if err != nil {
		t.Fatal(err)
	}

	if len(gs) != 1 {
		t.Errorf("Expected 1 geometry, got %d", len(gs))
	}

	// From this package, we can't reach inside the geos structs to check pointer
	// equality without reflection. The fmt code will print out the addresses,
	// though. It's hacky, but we can prove they're not the same instance
	act := fmt.Sprintf("%v", gs[0].g)
	exp := fmt.Sprintf("%v", lr.g)
	if act == exp {
		t.Errorf("Expected to get a cloned geometry")
	}
}

func TestBounds(t *testing.T) {
	coords := []Coord{
		{30, 10},
		{40, 40},
		{20, 40},
		{10, 20},
		{30, 10},
	}
	fact := NewFactory(handle.NewPooledHandleProvider())
	lr, err := fact.NewLinearRing(coords)
	if err != nil {
		t.Fatal(err)
	}
	c0, c1, err := lr.Bounds()

	if err != nil {
		t.Fatal(err)
	}

	if c0.X != 10 {
		t.Errorf("Expected x0 to be 10, got %f", c0.X)
	}
	if c0.Y != 10 {
		t.Errorf("Expected y0 to be 10, got %f", c0.Y)
	}

	if c1.X != 40 {
		t.Errorf("Expected x1 to be 40, got %f", c1.X)
	}
	if c1.Y != 40 {
		t.Errorf("Expected y1 to be 40, got %f", c1.Y)
	}
}

func compareCoordSlice(a, b []Coord) bool {
	if a == nil || b == nil {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
