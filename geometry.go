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

func (g *Geometry) binaryOperation(
	op binaryOp, o *geos.Geometry) (*Geometry, error) {

	h := g.hp.Get()
	geom, err := op(h, o)
	g.hp.Put(h)
	return newGeometryOrError(g.hp, geom, err)
}

func (g *Geometry) binaryPredicate(
	op binaryPredicate, o *geos.Geometry) (bool, error) {

	h := g.hp.Get()
	defer g.hp.Put(h)
	return op(h, o)
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

func (g *Geometry) Intersection(o *Geometry) (*Geometry, error) {
	return g.binaryOperation(g.g.Intersection, o.g)
}

func (g *Geometry) Union(o *Geometry) (*Geometry, error) {
	return g.binaryOperation(g.g.Union, o.g)
}

func (g *Geometry) Envelope() (*Geometry, error) {
	return g.unaryOperation(g.g.Envelope)
}

func (g *Geometry) Intersects(o *Geometry) (bool, error) {
	return g.binaryPredicate(g.g.Intersects, o.g)
}

func (g *Geometry) Contains(o *Geometry) (bool, error) {
	return g.binaryPredicate(g.g.Contains, o.g)
}

func (g *Geometry) Disjoint(o *Geometry) (bool, error) {
	return g.binaryPredicate(g.g.Disjoint, o.g)
}

func (g *Geometry) Touches(o *Geometry) (bool, error) {
	return g.binaryPredicate(g.g.Touches, o.g)
}

func (g *Geometry) Overlaps(o *Geometry) (bool, error) {
	return g.binaryPredicate(g.g.Overlaps, o.g)
}

func (g *Geometry) Within(o *Geometry) (bool, error) {
	return g.binaryPredicate(g.g.Within, o.g)
}

func (g *Geometry) IsEmpty() (bool, error) {
	return g.unaryPredicate(g.g.IsEmpty)
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

func (pg *PreparedGeometry) Covers(o *Geometry) (bool, error) {
	h := pg.hp.Get()
	defer pg.hp.Put(h)
	pg.Lock()
	defer pg.Unlock()

	return pg.p.Covers(h, o.g)
}
