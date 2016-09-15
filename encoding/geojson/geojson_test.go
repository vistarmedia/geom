package geojson

import (
	"strings"
	"testing"

	"vistarmedia.com/vistar/geom"
	geomcontext "vistarmedia.com/vistar/geom/context"
)

var (
	ctx     = geomcontext.NewContext()
	geoFact = ctx.Factory()
	dec     = NewDecoder(geoFact)
	wkt     = ctx.WKTEncoder()
)

func decode(s string) (*geom.Geometry, error) {
	return dec.Decode([]byte(s))
}

func TestDecodeInvalidType(t *testing.T) {
	_, err := decode(`{"type":"party"}`)

	if err == nil {
		t.Fatalf("Expcted error")
	}

	if err.Error() != "geojson: Unsupported type 'party'" {
		t.Fatalf("Unexpected error: %s", err.Error())
	}
}

func TestDecodePointZero(t *testing.T) {
	g, err := decode(`{"type":"Point","coordinates":[0,0]}`)
	if err != nil {
		t.Fatal(err)
	}

	enc := wkt.Encode(g)
	exp := "POINT (0.0000000000000000 0.0000000000000000)"
	if enc != exp {
		t.Fatalf("expected '%s', got '%s'", exp, enc)
	}
}

func TestDecodePointXY(t *testing.T) {
	g, _ := decode(`{"type":"Point","coordinates":[1.2,4.0]}`)
	enc := wkt.Encode(g)
	exp := "POINT (1.2000000000000000 4.0000000000000000)"
	if enc != exp {
		t.Fatalf("expected '%s', got '%s'", exp, enc)
	}
}

// When deserialzing to a [2]float64, if 3 values are encountered, it will
// silently drop the third point
// TODO: This should error
func TestDecodePointXYZ(t *testing.T) {
	g, _ := decode(`{"type":"Point","coordinates":[1,2,3]}`)
	enc := wkt.Encode(g)
	exp := "POINT (1.0000000000000000 2.0000000000000000)"
	if enc != exp {
		t.Fatalf("expected '%s', got '%s'", exp, enc)
	}
}

// When deserialzing to a [2]float64, if only 1 value is encountered, it will
// silently drop leave the second value at 0.
// TODO: This should error
func TestDecodePointX(t *testing.T) {
	g, _ := decode(`{"type":"Point","coordinates":[1]}`)
	enc := wkt.Encode(g)
	exp := "POINT (1.0000000000000000 0.0000000000000000)"
	if enc != exp {
		t.Fatalf("expected '%s', got '%s'", exp, enc)
	}
}

func TestLineStringUnsupported(t *testing.T) {
	_, err := decode(`{"type":"LineString","coordinates":[]}`)
	if err == nil {
		t.Fatalf("expected non-nil error")
	}
	if err.Error() != "geojson: Unsupported type 'LineString'" {
		t.Fatalf("Bad error message")
	}
}

func TestDecodeEmptyPolygon(t *testing.T) {
	g, err := decode(`{"type":"Polygon","coordinates":[]}`)
	if err != nil {
		t.Fatal(err)
	}
	enc := wkt.Encode(g)
	exp := "POLYGON EMPTY"
	if enc != exp {
		t.Fatalf("expected '%s', got '%s'", exp, enc)
	}
}

func TestDecodePolygonNoHoles(t *testing.T) {
	g, err := decode(`{"type":"Polygon","coordinates":[[[1,2],[3,4],[5,6],[1,2]]]}`)
	if err != nil {
		t.Fatal(err)
	}
	enc := wkt.Encode(g)
	exp := `POLYGON ((1.0000000000000000 2.0000000000000000, 3.0000000000000000 4.0000000000000000, 5.0000000000000000 6.0000000000000000, 1.0000000000000000 2.0000000000000000))`
	if enc != exp {
		t.Fatalf("expected '%s', got '%s'", exp, enc)
	}
}

func TestDecodeMultipolygon(t *testing.T) {
	g, err := decode(`{"type":"MultiPolygon","coordinates":` +
		`[[[[1,2,3],[4,5,6],[7,8,9],[1,2,3]],` +
		`[[-1,-2,-3],[-4,-5,-6],[-7,-8,-9],[-1,-2,-3]]]]}`)

	if err != nil {
		t.Fatal(err)
	}

	enc := wkt.Encode(g)
	if !strings.HasPrefix(enc, "MULTIPOLYGON (") {
		t.Fatalf("Unexpected WKT: %s", enc)
	}
}
