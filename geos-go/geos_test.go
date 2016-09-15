package geos

import (
	"strings"
	"testing"
)

func TestCoordSeq(t *testing.T) {
	h := NewHandle()
	defer h.Destroy()

	cs := NewCoordSeq(h, 1, 2)
	defer cs.Destroy(h)
	if cs == nil {
		t.Error("cs should not be nil")
	}
}

func TestCoordSeqSize(t *testing.T) {
	h := NewHandle()
	defer h.Destroy()

	cs := NewCoordSeq(h, 123, 2)
	defer cs.Destroy(h)
	if cs.Size(h) != 123 {
		t.Errorf("Expected cs size to be 123, was %d", cs.Size(h))
	}
}

func TestCoordSeqX(t *testing.T) {
	h := NewHandle()
	defer h.Destroy()

	cs := NewCoordSeq(h, 1, 2)
	defer cs.Destroy(h)
	err := cs.SetX(h, 0, 12.3)
	if err != nil {
		t.Error(err)
	}

	if cs.X(h, 0) != 12.3 {
		t.Errorf("Expected 12.3, got %f", cs.X(h, 0))
	}

	if err = cs.SetX(h, 1, 12.3); err != ErrIndexOutOfBounds {
		t.Errorf("Expected %v, got %v", ErrIndexOutOfBounds, err)
	}
}

func TestCoordSeqY(t *testing.T) {
	h := NewHandle()
	defer h.Destroy()

	cs := NewCoordSeq(h, 1, 2)
	defer cs.Destroy(h)

	err := cs.SetY(h, 0, 12.3)
	if err != nil {
		t.Error(err)
	}

	if cs.Y(h, 0) != 12.3 {
		t.Errorf("Expected 12.3, got %f", cs.Y(h, 0))
	}

	if err = cs.SetY(h, 1, 12.3); err != ErrIndexOutOfBounds {
		t.Errorf("Expected %v, got %v", ErrIndexOutOfBounds, err)
	}
}

func TestCoordSeqZ(t *testing.T) {
	h := NewHandle()
	defer h.Destroy()

	cs := NewCoordSeq(h, 1, 3)
	defer cs.Destroy(h)

	err := cs.SetZ(h, 0, 12.3)
	if err != nil {
		t.Error(err)
	}

	if cs.Z(h, 0) != 12.3 {
		t.Errorf("Expected 12.3, got %f", cs.Z(h, 0))
	}

	if err = cs.SetX(h, 1, 12.3); err != ErrIndexOutOfBounds {
		t.Errorf("Expected %v, got %v", ErrIndexOutOfBounds, err)
	}
}

func makeTenByTenSquare(h *Handle) *CoordSeq {
	cs := NewCoordSeq(h, 5, 2)
	cs.SetX(h, 0, 10)
	cs.SetY(h, 0, 10)
	cs.SetX(h, 1, 20)
	cs.SetY(h, 1, 10)
	cs.SetX(h, 2, 20)
	cs.SetY(h, 2, 20)
	cs.SetX(h, 3, 10)
	cs.SetY(h, 3, 20)
	cs.SetX(h, 4, 10)
	cs.SetY(h, 4, 10)
	return cs
}

func TestNewPolygon(t *testing.T) {
	h := NewHandle()
	defer h.Destroy()

	cs := makeTenByTenSquare(h)
	shell, err := cs.LinearRing(h)
	if err != nil {
		t.Error(err)
	}
	poly, err := NewPolygon(h, shell, nil)
	if err != nil {
		t.Error(err)
	}
	defer poly.Destroy(h)
	area := poly.Area(h)
	if area != 100 {
		t.Errorf("Expected 100, got %f", area)
	}
}

func TestNewPolygonWithHole(t *testing.T) {
	h := NewHandle()
	defer h.Destroy()

	cs := makeTenByTenSquare(h)
	shell, err := cs.LinearRing(h)
	if err != nil {
		t.Error(err)
	}

	holeCs := NewCoordSeq(h, 5, 2)
	// A 2x2 square
	holeCs.SetX(h, 0, 12)
	holeCs.SetY(h, 0, 12)
	holeCs.SetX(h, 1, 14)
	holeCs.SetY(h, 1, 12)
	holeCs.SetX(h, 2, 14)
	holeCs.SetY(h, 2, 14)
	holeCs.SetX(h, 3, 12)
	holeCs.SetY(h, 3, 14)
	holeCs.SetX(h, 4, 12)
	holeCs.SetY(h, 4, 12)
	hole, err := holeCs.LinearRing(h)
	if err != nil {
		t.Error(err)
	}
	poly, err := NewPolygon(h, shell, []*Geometry{hole})
	if err != nil {
		t.Error(err)
	}
	defer poly.Destroy(h)
	area := poly.Area(h)
	// (10x10) - (2x2) = 96
	if area != 96 {
		t.Errorf("Expected 96, got %f", area)
	}
}

func TestContains(t *testing.T) {
	h := NewHandle()
	defer h.Destroy()

	cs := makeTenByTenSquare(h)
	shell, err := cs.LinearRing(h)
	if err != nil {
		t.Error(err)
	}
	poly, err := NewPolygon(h, shell, nil)
	if err != nil {
		t.Error(err)
	}
	defer poly.Destroy(h)

	inside := NewCoordSeq(h, 1, 2)
	inside.SetX(h, 0, 12)
	inside.SetY(h, 0, 13)
	insidePoint, err := inside.Point(h)
	if err != nil {
		t.Error(err)
	}
	defer insidePoint.Destroy(h)
	contains1, err := poly.Contains(h, insidePoint)
	if err != nil {
		t.Error(err)
	}
	if !contains1 {
		t.Error("Expected point 12,13 to be contained")
	}

	outside := NewCoordSeq(h, 1, 2)
	outside.SetX(h, 0, 22)
	outside.SetY(h, 0, 13)
	outsidePoint, err := outside.Point(h)
	if err != nil {
		t.Error(err)
	}
	defer outsidePoint.Destroy(h)
	contains2, err := poly.Contains(h, outsidePoint)
	if err != nil {
		t.Error(err)
	}
	if contains2 {
		t.Error("Expected point 22,13 to not be contained")
	}
}

func TestWKT(t *testing.T) {
	h := NewHandle()
	defer h.Destroy()

	reader := NewWKTReader(h)
	defer reader.Destroy(h)

	writer := NewWKTWriter(h)
	defer writer.Destroy(h)

	geom, err := reader.Read(h, "POINT(10 20)")
	if err != nil {
		t.Error(err)
	}
	wkt := writer.Write(h, geom)
	if wkt != "POINT (10.0000000000000000 20.0000000000000000)" {
		t.Errorf("Unexpected WKT: %s", wkt)
	}
}

func TestWKB(t *testing.T) {
	h := NewHandle()
	defer h.Destroy()

	reader := NewWKBReader(h)
	defer reader.Destroy(h)

	writer := NewWKBWriter(h)
	defer writer.Destroy(h)

	cs := makeTenByTenSquare(h)
	shell, err := cs.LinearRing(h)
	if err != nil {
		t.Error(err)
	}
	poly, err := NewPolygon(h, shell, nil)
	if err != nil {
		t.Error(err)
	}
	defer poly.Destroy(h)
	wkb := writer.Write(h, poly)
	newPoly, err := reader.Read(h, wkb)
	if err != nil {
		t.Error(err)
	}
	if newPoly.Area(h) != 100 {
		t.Error("Error reading WKB")
	}
}

func TestNewEmptyMultipolygon(t *testing.T) {
	h := NewHandle()
	defer h.Destroy()

	writer := NewWKTWriter(h)
	defer writer.Destroy(h)

	wkt := writer.Write(h, NewEmptyGeometryCollection(h, MULTIPOLYGON))
	exp := "MULTIPOLYGON EMPTY"
	if wkt != exp {
		t.Fatalf("Expected '%s', got '%s'", exp, wkt)
	}
}

func TestNewMultipolygon(t *testing.T) {
	h := NewHandle()
	defer h.Destroy()

	writer := NewWKTWriter(h)
	defer writer.Destroy(h)

	cs0 := NewCoordSeq(h, 4, 2)
	cs0.SetX(h, 0, 0)
	cs0.SetY(h, 0, 0)

	cs0.SetX(h, 1, 1)
	cs0.SetY(h, 1, 1)

	cs0.SetX(h, 2, 1)
	cs0.SetY(h, 2, 0)

	cs0.SetX(h, 3, 0)
	cs0.SetY(h, 3, 0)

	lr0, _ := cs0.LinearRing(h)
	p0, _ := NewPolygon(h, lr0, nil)

	cs1 := NewCoordSeq(h, 4, 2)
	cs1.SetX(h, 0, 2)
	cs1.SetY(h, 0, 2)

	cs1.SetX(h, 1, 3)
	cs1.SetY(h, 1, 3)

	cs1.SetX(h, 2, 3)
	cs1.SetY(h, 2, 2)

	cs1.SetX(h, 3, 2)
	cs1.SetY(h, 3, 2)

	lr1, _ := cs1.LinearRing(h)
	p1, _ := NewPolygon(h, lr1, nil)

	mp, err := NewGeometryCollection(h, MULTIPOLYGON, []*Geometry{p0, p1})
	if err != nil {
		t.Fatal(err)
	}

	wkt := writer.Write(h, mp)
	if !strings.HasPrefix(wkt, "MULTIPOLYGON (") {
		t.Fatalf("Unexpected WKT: %s", wkt)
	}

	// It's brutal to do a big string compare here. Instead, we'll count the
	// braces in the MULTIPOLYGON
	// Each poly:        2 (
	// Each poly's ring: 2 (
	// MP def:           1 (
	// TOTAL             5
	nOpens := strings.Count(wkt, "(")
	if nOpens != 5 {
		t.Fatalf("expected 5 ('s, found %d", nOpens)
	}

	nCloses := strings.Count(wkt, ")")
	if nCloses != 5 {
		t.Fatalf("expected 5 ('s, found %d", nCloses)
	}
}
