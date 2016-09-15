// Package geojson implements basic GeoJSON decoding.
package geojson

import (
	"encoding/json"
	"fmt"

	"vistarmedia.com/vistar/geom"
)

type ErrUnsupportedType string

func (e ErrUnsupportedType) Error() string {
	return fmt.Sprintf("geojson: Unsupported type '%s'", string(e))
}

// -----------------------------------------------------------------------------
// Decoder
type Decoder struct {
	geoFact geom.Factory
}

func NewDecoder(f geom.Factory) Decoder {
	return Decoder{f}
}

func (d Decoder) Decode(b []byte) (*geom.Geometry, error) {
	m := struct {
		Type        string          `json:"type"`
		Coordinates json.RawMessage `json:"coordinates"`
		Geometry    json.RawMessage `json:"geometry"`
	}{}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	switch m.Type {
	case "Point":
		point, err := d.decodePoint(m.Coordinates)
		if err != nil {
			return nil, err
		}
		return point.Geometry, nil

	case "Polygon":
		poly, err := d.decodePolygon(m.Coordinates)
		if err != nil {
			return nil, err
		}
		return poly.Geometry, nil

	case "MultiPolygon":
		mpoly, err := d.decodeMultipolygon(m.Coordinates)
		if err != nil {
			return nil, err
		}
		return mpoly.Geometry, nil
	}

	return nil, ErrUnsupportedType(m.Type)
}

func (d Decoder) decodePoint(coords json.RawMessage) (p geom.Point, err error) {
	cs := [2]float64{}
	if err = json.Unmarshal(coords, &cs); err != nil {
		return
	}
	coord := geom.Coord{X: cs[0], Y: cs[1]}
	return d.geoFact.NewPoint(coord)
}

func (d Decoder) decodePolygon(coords json.RawMessage) (p geom.Polygon, err error) {
	ps := [][][2]float64{}
	if err = json.Unmarshal(coords, &ps); err != nil {
		return
	}
	return d.newPolygon(ps)
}

func (d Decoder) newPolygon(ps [][][2]float64) (geom.Polygon, error) {
	if len(ps) == 0 {
		return d.geoFact.NewEmptyPolygon(), nil
	}

	rings := make([][]geom.Coord, len(ps))
	for i := 0; i < len(ps); i++ {
		points := ps[i]
		ring := make([]geom.Coord, len(points))

		for j := 0; j < len(points); j++ {
			ring[j] = geom.Coord{X: points[j][0], Y: points[j][1]}
		}
		rings[i] = ring
	}

	shell := rings[0]
	holes := rings[1:]
	return d.geoFact.NewPolygon(shell, holes...)
}

func (d Decoder) decodeMultipolygon(coords json.RawMessage) (g geom.Multipolygon, err error) {
	ps := [][][][2]float64{}
	if err = json.Unmarshal(coords, &ps); err != nil {
		return
	}

	if len(ps) == 0 {
		return d.geoFact.NewEmptyMultipolygon(), nil
	}

	polys := make([]geom.Polygon, len(ps))
	var poly geom.Polygon
	for i := 0; i < len(ps); i++ {
		if poly, err = d.newPolygon(ps[i]); err != nil {
			return
		}
		polys[i] = poly
	}

	return d.geoFact.NewMultipolygon(polys...)
}
