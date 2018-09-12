package wkt

import (
	"testing"

	"github.com/vistarmedia/geom"
	"github.com/vistarmedia/geom/geos-go/handle"
)

func TestWKT(t *testing.T) {
	hp := handle.NewPooledHandleProvider()
	fact := geom.NewFactory(hp)
	encoder := NewEncoder(hp)
	decoder := NewDecoder(hp, fact)

	point, err := fact.NewPoint(geom.Coord{2, 4})
	if err != nil {
		t.Error(err)
	}
	wkt := encoder.Encode(point)
	if wkt != "POINT (2.0000000000000000 4.0000000000000000)" {
		t.Errorf("Incorrect WKT: %s", wkt)
	}
	geometry, err := decoder.Decode(wkt)
	if err != nil {
		t.Error(err)
	}
	if geometry.Type() != geom.POINT {
		t.Errorf("Unexpected type: %d", geometry.Type())
	}
}
