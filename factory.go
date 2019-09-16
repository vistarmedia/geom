package geom

import (
	"errors"
	"runtime"

	"github.com/vistarmedia/geom/geos-go"
	"github.com/vistarmedia/geom/geos-go/handle"
)

var (
	ErrEmptyCoords = errors.New("Empty coordinates")
)

type Coord struct {
	X, Y float64
}

func newGeosLinearRing(h *geos.Handle, coords []Coord) (*geos.Geometry, error) {
	coordsLen := len(coords)
	if coordsLen == 0 {
		return nil, ErrEmptyCoords
	}
	cs := geos.NewCoordSeq(h, uint(coordsLen), 2)
	for i, c := range coords {
		if err := cs.SetX(h, uint(i), c.X); err != nil {
			cs.Destroy(h)
			return nil, err
		}
		if err := cs.SetY(h, uint(i), c.Y); err != nil {
			cs.Destroy(h)
			return nil, err
		}
	}
	// LinearRing destructor will destroy the coord seq even on error
	return cs.LinearRing(h)
}

// For creating geometries
type Factory struct {
	hp handle.GeosHandleProvider
}

func NewFactory(hp handle.GeosHandleProvider) Factory {
	return Factory{hp}
}

func (f Factory) FromGeos(g *geos.Geometry) *Geometry {
	return newGeometry(f.hp, g)
}

func (f Factory) NewEmptyPoint() Point {
	h := f.hp.Get()
	defer f.hp.Put(h)
	return newPoint(newGeometry(f.hp, geos.NewEmptyPoint(h)))
}

func (f Factory) NewEmptyPolygon() Polygon {
	h := f.hp.Get()
	defer f.hp.Put(h)
	return newPolygon(newGeometry(f.hp, geos.NewEmptyPolygon(h)))
}

func (f Factory) NewEmptyMultipolygon() Multipolygon {
	h := f.hp.Get()
	defer f.hp.Put(h)
	mp := geos.NewEmptyGeometryCollection(h, geos.MULTIPOLYGON)
	return newMultiPolygon(newGeometry(f.hp, mp))
}

func (f Factory) NewPoint(c Coord) (p Point, err error) {
	h := f.hp.Get()
	defer f.hp.Put(h)
	cs := geos.NewCoordSeq(h, 1, 2)
	if err = cs.SetX(h, 0, c.X); err != nil {
		cs.Destroy(h)
		return
	}
	if err = cs.SetY(h, 0, c.Y); err != nil {
		cs.Destroy(h)
		return
	}
	if point, err := cs.Point(h); err == nil {
		p = newPoint(newGeometry(f.hp, point))
	}
	return
}

func (f Factory) NewLinearRing(coords []Coord) (lr LinearRing, err error) {
	h := f.hp.Get()
	g, err := newGeosLinearRing(h, coords)
	f.hp.Put(h)
	if err != nil {
		return
	}
	lr = newLinearRing(newGeometry(f.hp, g))
	return
}

func (f Factory) NewPolygon(
	shell []Coord, holes ...[]Coord) (p Polygon, err error) {

	h := f.hp.Get()
	defer f.hp.Put(h)
	shellRing, err := newGeosLinearRing(h, shell)
	if err != nil {
		return
	}
	var holeRings []*geos.Geometry
	for _, hole := range holes {
		var holeRing *geos.Geometry
		holeRing, err = newGeosLinearRing(h, hole)
		if err != nil {
			return
		}
		holeRings = append(holeRings, holeRing)
	}
	g, err := geos.NewPolygon(h, shellRing, holeRings)
	if err != nil {
		return
	}
	p = newPolygon(newGeometry(f.hp, g))
	return
}

// Create a MULTIPOLYGON from some POLYGONs. This will clone the passed POLYGONs
// and assume ownership of those clones. Arguments continue to be managed
// independently by GC
func (f Factory) NewMultipolygon(ps ...Polygon) (mp Multipolygon, err error) {
	h := f.hp.Get()
	defer f.hp.Put(h)

	geoms := make([]*geos.Geometry, len(ps))
	for i, poly := range ps {
		geoms[i] = poly.g.Clone(h)
		runtime.KeepAlive(poly)
	}

	g, err := geos.NewGeometryCollection(h, geos.MULTIPOLYGON, geoms)
	if err != nil {
		return
	}
	mp = newMultiPolygon(newGeometry(f.hp, g))
	return
}
