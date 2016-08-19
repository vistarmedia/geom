package wkb

import (
	"encoding/hex"
	"testing"

	"vistarmedia.com/vistar/geom"
	"vistarmedia.com/vistar/geom/geos-go/handle"
)

func TestWKB(t *testing.T) {
	hp := handle.NewPooledHandleProvider()
	fact := geom.NewFactory(hp)
	encoder := NewEncoder(hp)
	decoder := NewDecoder(hp, fact)

	point, err := fact.NewPoint(geom.Coord{2, 4})
	if err != nil {
		t.Error(err)
	}
	wkb := encoder.Encode(point)
	hex := hex.EncodeToString(wkb)
	// Example from (converted to little endian):
	// https://en.wikipedia.org/wiki/Well-known_text#Well-known_binary
	if hex != "010100000000000000000000400000000000001040" {
		t.Errorf("Incorrect WKB: %s", hex)
	}
	geometry, err := decoder.Decode(wkb)
	if err != nil {
		t.Error(err)
	}
	if geometry.Type() != geom.POINT {
		t.Errorf("Unexpected type: %d", geometry.Type())
	}
}
