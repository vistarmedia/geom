package geom

import (
	"fmt"
	"sync"

	"vistarmedia.com/vistar/geom/geos-go"
	"vistarmedia.com/vistar/geom/geos-go/handle"
	"vistarmedia.com/vistar/geom/geos-go/memory"
)

type GeometryType int

const (
	POINT GeometryType = iota
	LINESTRING
	LINEARRING
	POLYGON
	MULTIPOINT
	MULTILINESTRING
	MULTIPOLYGON
	GEOMETRYCOLLECTION
)

// See OGC Simple Feature Specification for geometry operation details:
// http://portal.opengeospatial.org/files/?artifact_id=25355
type Geometry struct {
	hp handle.GeosHandleProvider
	g  *geos.Geometry
}

type toGeos interface {
	UnsafeToGeos() *geos.Geometry
}

type unaryOp func(*geos.Handle) (*geos.Geometry, error)
type unaryPredicate func(*geos.Handle) (bool, error)
type binaryOp func(*geos.Handle, *geos.Geometry) (*geos.Geometry, error)
type binaryPredicate func(*geos.Handle, *geos.Geometry) (bool, error)

func newGeometry(
	hp handle.GeosHandleProvider, g *geos.Geometry) *Geometry {

	memory.GeosGCManaged(hp, g)
	return &Geometry{
		hp: hp,
		g:  g,
	}
}

func newGeometryOrError(
	hp handle.GeosHandleProvider,
	g *geos.Geometry,
	err error) (*Geometry, error) {

	if err != nil {
		return nil, err
	} else {
		return newGeometry(hp, g), nil
	}
}

func (g *Geometry) unaryOperation(op unaryOp) (*Geometry, error) {
	h := g.hp.Get()
	geom, err := op(h)
	g.hp.Put(h)
	return newGeometryOrError(g.hp, geom, err)
}

func (g *Geometry) unaryPredicate(op unaryPredicate) (bool, error) {
	h := g.hp.Get()
	defer g.hp.Put(h)
	return op(h)
}

func (g *Geometry) binaryOperation(op binaryOp, o toGeos) (*Geometry, error) {

	h := g.hp.Get()
	geom, err := op(h, o.UnsafeToGeos())
	g.hp.Put(h)
	return newGeometryOrError(g.hp, geom, err)
}

func (g *Geometry) binaryPredicate(op binaryPredicate, o toGeos) (bool, error) {
	h := g.hp.Get()
	defer g.hp.Put(h)
	return op(h, o.UnsafeToGeos())
}

func (g *Geometry) Prepared() *PreparedGeometry {
	h := g.hp.Get()
	prep := g.g.Prepared(h)
	g.hp.Put(h)
	memory.GeosGCManaged(g.hp, prep)
	return &PreparedGeometry{
		hp:     g.hp,
		p:      prep,
		parent: g,
	}
}

// Unsafe access to the geos geometry. This geometry is still subject to GC.
// For internal use only.
func (g *Geometry) UnsafeToGeos() *geos.Geometry {
	return g.g
}

func (g *Geometry) Type() GeometryType {
	h := g.hp.Get()
	id := g.g.TypeId(h)
	g.hp.Put(h)

	switch id {
	case geos.POINT:
		return POINT
	case geos.LINESTRING:
		return LINESTRING
	case geos.LINEARRING:
		return LINEARRING
	case geos.POLYGON:
		return POLYGON
	case geos.MULTIPOINT:
		return MULTIPOINT
	case geos.MULTILINESTRING:
		return MULTILINESTRING
	case geos.MULTIPOLYGON:
		return MULTIPOLYGON
	case geos.GEOMETRYCOLLECTION:
		return GEOMETRYCOLLECTION
	default:
		panic(fmt.Sprintf("Unknown GEOS geometry id: %d", id))
	}
}

func (g *Geometry) Area() float64 {
	h := g.hp.Get()
	defer g.hp.Put(h)

	return g.g.Area(h)
}

func (g *Geometry) ClipByRect(
	xmin, ymin, xmax, ymax float64) (*Geometry, error) {

	h := g.hp.Get()
	geom, err := g.g.ClipByRect(h, xmin, ymin, xmax, ymax)
	g.hp.Put(h)
	return newGeometryOrError(g.hp, geom, err)
}

func (g *Geometry) Buffer(width float64, quadsegs int) (*Geometry, error) {
	h := g.hp.Get()
	geom, err := g.g.Buffer(h, width, quadsegs)
	g.hp.Put(h)
	return newGeometryOrError(g.hp, geom, err)
}

func (g *Geometry) Intersection(o toGeos) (*Geometry, error) {
	return g.binaryOperation(g.g.Intersection, o)
}

func (g *Geometry) Union(o toGeos) (*Geometry, error) {
	return g.binaryOperation(g.g.Union, o)
}

func (g *Geometry) Envelope() (*Geometry, error) {
	return g.unaryOperation(g.g.Envelope)
}

func (g *Geometry) Intersects(o toGeos) (bool, error) {
	return g.binaryPredicate(g.g.Intersects, o)
}

func (g *Geometry) Contains(o toGeos) (bool, error) {
	return g.binaryPredicate(g.g.Contains, o)
}

func (g *Geometry) Disjoint(o toGeos) (bool, error) {
	return g.binaryPredicate(g.g.Disjoint, o)
}

func (g *Geometry) Touches(o toGeos) (bool, error) {
	return g.binaryPredicate(g.g.Touches, o)
}

func (g *Geometry) Overlaps(o toGeos) (bool, error) {
	return g.binaryPredicate(g.g.Overlaps, o)
}

func (g *Geometry) Within(o toGeos) (bool, error) {
	return g.binaryPredicate(g.g.Within, o)
}

func (g *Geometry) IsEmpty() (bool, error) {
	return g.unaryPredicate(g.g.IsEmpty)
}

// Coerces to Point. Panics if the underlying type doesnt match.
func (g *Geometry) Point() Point {
	if id := g.Type(); id != POINT {
		panic(fmt.Sprintf("Cannot cast geom with type %d to POINT (%d)", id, POINT))
	}
	return newPoint(g)
}

// Coerces to LinearRing. Panics if the underlying type doesnt match.
func (g *Geometry) LinearRing() LinearRing {
	if id := g.Type(); id != LINEARRING {
		panic(fmt.Sprintf(
			"Cannot cast geom with type %d to LINEARRING (%d)", id, LINEARRING))
	}
	return newLinearRing(g)
}

// Coerces to Polygon. Panics if the underlying type doesnt match.
func (g *Geometry) Polygon() Polygon {
	if id := g.Type(); id != POLYGON {
		panic(fmt.Sprintf(
			"Cannot cast geom with type %d to POLYGON (%d)", id, POLYGON))
	}
	return newPolygon(g)
}

// Expensive to create, but faster predicate operations.
type PreparedGeometry struct {
	hp handle.GeosHandleProvider
	p  *geos.PreparedGeometry
	// Hold on to the parent geom so it doesnt get GCed
	parent *Geometry
	// PreparedGeoms are not thread safe. Lock operations.
	sync.Mutex
}

func (pg *PreparedGeometry) Covers(o toGeos) (bool, error) {
	h := pg.hp.Get()
	defer pg.hp.Put(h)
	pg.Lock()
	defer pg.Unlock()

	return pg.p.Covers(h, o.UnsafeToGeos())
}

// Point
type Point struct {
	*Geometry
}

func newPoint(g *Geometry) Point {
	return Point{g}
}

func (p Point) Coord() (Coord, error) {
	h := p.hp.Get()
	defer p.hp.Put(h)
	cs, err := p.g.CoordSeq(h)
	if err != nil {
		return Coord{}, err
	}
	return Coord{cs.X(h, 0), cs.Y(h, 0)}, nil
}

// LinearRing
type LinearRing struct {
	*Geometry
}

func newLinearRing(g *Geometry) LinearRing {
	return LinearRing{g}
}

func (lr LinearRing) Coords() (coords []Coord, err error) {
	h := lr.hp.Get()
	defer lr.hp.Put(h)
	cs, err := lr.g.CoordSeq(h)
	if err != nil {
		return
	}
	for i := uint(0); i < cs.Size(h); i++ {
		coords = append(coords, Coord{cs.X(h, i), cs.Y(h, i)})
	}
	return
}

// Polygon
type Polygon struct {
	*Geometry
}

func newPolygon(g *Geometry) Polygon {
	return Polygon{g}
}

func (p Polygon) Shell() (coords []Coord, err error) {
	h := p.hp.Get()
	defer p.hp.Put(h)
	geosShell, err := p.g.ExteriorRing(h)
	if err != nil {
		return
	}
	cs, err := geosShell.CoordSeq(h)
	if err != nil {
		return
	}
	// It would be nice if this could just return a LinearRing, but the polygon
	// owns the shell so we would have to clone it to be safe.
	for i := uint(0); i < cs.Size(h); i++ {
		coords = append(coords, Coord{cs.X(h, i), cs.Y(h, i)})
	}
	return
}

func (p Polygon) Holes() (coords [][]Coord, err error) {
	h := p.hp.Get()
	defer p.hp.Put(h)
	numRings, err := p.g.NumInteriorRings(h)
	if err != nil {
		return
	}
	for ringIdx := 0; ringIdx < numRings; ringIdx++ {
		var geosRing *geos.Geometry
		geosRing, err = p.g.InteriorRingN(h, ringIdx)
		if err != nil {
			return
		}
		var cs *geos.CoordSeq
		cs, err = geosRing.CoordSeq(h)
		if err != nil {
			return
		}
		var ringCoords []Coord
		for i := uint(0); i < cs.Size(h); i++ {
			ringCoords = append(ringCoords, Coord{cs.X(h, i), cs.Y(h, i)})
		}
		coords = append(coords, ringCoords)
	}
	return
}
