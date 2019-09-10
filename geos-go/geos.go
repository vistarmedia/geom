// Package geos implements bindings to the reentrant libgeos API.
// Requires libgeos 3.5.0 or greater to be linked. Not thread or memory safe.
package geos

// #cgo LDFLAGS: -lgeos_c
// #include <geos_c.h>
// #include <stdlib.h>
// extern GEOSContextHandle_t createGEOSHandle();
import "C"

import (
	"errors"
	"fmt"
	"unsafe"
)

type GeometryTypeId int

const (
	POINT              GeometryTypeId = C.GEOS_POINT
	LINESTRING         GeometryTypeId = C.GEOS_LINESTRING
	LINEARRING         GeometryTypeId = C.GEOS_LINEARRING
	POLYGON            GeometryTypeId = C.GEOS_POLYGON
	MULTIPOINT         GeometryTypeId = C.GEOS_MULTIPOINT
	MULTILINESTRING    GeometryTypeId = C.GEOS_MULTILINESTRING
	MULTIPOLYGON       GeometryTypeId = C.GEOS_MULTIPOLYGON
	GEOMETRYCOLLECTION GeometryTypeId = C.GEOS_GEOMETRYCOLLECTION
)

// Errors
var (
	ErrIndexOutOfBounds = errors.New("Index out of bounds")
	ErrGeos             = errors.New("GEOS Error")
	ErrEmptyWKB         = errors.New("Tried to read empty WKB")
)

func predicate(char C.char) (bool, error) {
	if char == 2 {
		return false, ErrGeos
	}
	return char == 1, nil
}

// Wraps a GEOS handle to provide access to the reentrant API.
// Not goroutine-safe.
type Handle struct {
	h C.GEOSContextHandle_t
}

func NewHandle() *Handle {
	return &Handle{C.createGEOSHandle()}
}

func (h *Handle) Destroy() {
	C.finishGEOS_r(h.h)
}

// A list of coordinates used to construct geometries.
// http://geos.osgeo.org/doxygen/classgeos_1_1geom_1_1CoordinateSequence.html
type CoordSeq struct {
	cs *C.GEOSCoordSequence
}

func NewCoordSeq(h *Handle, size, dims uint) *CoordSeq {
	return &CoordSeq{C.GEOSCoordSeq_create_r(h.h, C.uint(size), C.uint(dims))}
}

func (cs *CoordSeq) Destroy(h *Handle) {
	C.GEOSCoordSeq_destroy_r(h.h, cs.cs)
}

func (cs *CoordSeq) Size(h *Handle) uint {
	return cs.size(h.h)
}

func (cs *CoordSeq) size(handle C.GEOSContextHandle_t) uint {
	var size C.uint
	C.GEOSCoordSeq_getSize_r(handle, cs.cs, &size)
	return uint(size)
}

func (cs *CoordSeq) X(h *Handle, idx uint) float64 {
	var x C.double
	C.GEOSCoordSeq_getX_r(h.h, cs.cs, C.uint(idx), &x)
	return float64(x)
}

func (cs *CoordSeq) Y(h *Handle, idx uint) float64 {
	var y C.double
	C.GEOSCoordSeq_getY_r(h.h, cs.cs, C.uint(idx), &y)
	return float64(y)
}

func (cs *CoordSeq) Z(h *Handle, idx uint) float64 {
	var z C.double
	C.GEOSCoordSeq_getZ_r(h.h, cs.cs, C.uint(idx), &z)
	return float64(z)
}

func (cs *CoordSeq) SetX(h *Handle, idx uint, val float64) error {
	if err := cs.checkIdx(h.h, idx); err != nil {
		return err
	}
	C.GEOSCoordSeq_setX_r(h.h, cs.cs, C.uint(idx), C.double(val))
	return nil
}

func (cs *CoordSeq) SetY(h *Handle, idx uint, val float64) error {
	if err := cs.checkIdx(h.h, idx); err != nil {
		return err
	}
	C.GEOSCoordSeq_setY_r(h.h, cs.cs, C.uint(idx), C.double(val))
	return nil
}

func (cs *CoordSeq) SetZ(h *Handle, idx uint, val float64) error {
	if err := cs.checkIdx(h.h, idx); err != nil {
		return err
	}
	C.GEOSCoordSeq_setZ_r(h.h, cs.cs, C.uint(idx), C.double(val))
	return nil
}

func (cs *CoordSeq) Point(h *Handle) (*Geometry, error) {
	if geom := C.GEOSGeom_createPoint_r(h.h, cs.cs); geom != nil {
		return &Geometry{geom}, nil
	}
	return nil, ErrGeos
}

func (cs *CoordSeq) LinearRing(h *Handle) (*Geometry, error) {
	if geom := C.GEOSGeom_createLinearRing_r(h.h, cs.cs); geom != nil {
		return &Geometry{geom}, nil
	}
	return nil, ErrGeos
}

func (cs *CoordSeq) checkIdx(handle C.GEOSContextHandle_t, idx uint) error {
	if idx < 0 || idx >= cs.size(handle) {
		return ErrIndexOutOfBounds
	}
	return nil
}

// http://geos.osgeo.org/doxygen/classgeos_1_1geom_1_1Geometry.html
type Geometry struct {
	g *C.GEOSGeometry
}

// Clones this geometry. The caller takes ownership
func (g *Geometry) Clone(h *Handle) *Geometry {
	return &Geometry{C.GEOSGeom_clone_r(h.h, g.g)}
}

func (g *Geometry) Destroy(h *Handle) {
	C.GEOSGeom_destroy_r(h.h, g.g)
}

func NewEmptyPoint(h *Handle) *Geometry {
	return &Geometry{C.GEOSGeom_createEmptyPoint_r(h.h)}
}

func NewEmptyPolygon(h *Handle) *Geometry {
	return &Geometry{C.GEOSGeom_createEmptyPolygon_r(h.h)}
}

func NewEmptyGeometryCollection(h *Handle, geomType GeometryTypeId) *Geometry {
	i := C.int(geomType)
	return &Geometry{C.GEOSGeom_createEmptyCollection_r(h.h, i)}
}

func NewPolygon(
	h *Handle, shell *Geometry, holes []*Geometry) (*Geometry, error) {

	var geosHoles []*C.GEOSGeometry
	for _, hole := range holes {
		geosHoles = append(geosHoles, hole.g)
	}
	var holesCArray **C.GEOSGeometry
	holeCount := len(geosHoles)
	if holeCount > 0 {
		holesCArray = &geosHoles[0]
	}
	geom := C.GEOSGeom_createPolygon_r(
		h.h, shell.g, holesCArray, C.uint(holeCount))
	if geom == nil {
		return nil, ErrGeos
	}
	return &Geometry{geom}, nil
}

func NewGeometryCollection(h *Handle, geomType GeometryTypeId,
	gs []*Geometry) (*Geometry, error) {

	ngeoms := C.uint(len(gs))

	geosGeoms := make([]*C.GEOSGeometry, len(gs))
	for i := 0; i < len(gs); i++ {
		geosGeoms[i] = gs[i].g
	}

	geom := C.GEOSGeom_createCollection_r(h.h, C.int(geomType), &geosGeoms[0],
		ngeoms)

	if geom == nil {
		return nil, ErrGeos
	}
	return &Geometry{geom}, nil
}

func (g *Geometry) TypeId(h *Handle) GeometryTypeId {
	return GeometryTypeId(C.GEOSGeomTypeId_r(h.h, g.g))
}

func (g *Geometry) Area(h *Handle) float64 {
	var area C.double
	C.GEOSArea_r(h.h, g.g, &area)
	return float64(area)
}

func (g *Geometry) Prepared(h *Handle) *PreparedGeometry {
	return &PreparedGeometry{C.GEOSPrepare_r(h.h, g.g)}
}

func (g *Geometry) ClipByRect(
	h *Handle, xmin, ymin, xmax, ymax float64) (*Geometry, error) {

	geom := C.GEOSClipByRect_r(h.h, g.g, C.double(xmin), C.double(ymin),
		C.double(xmax), C.double(ymax))
	if geom == nil {
		return nil, ErrGeos
	}
	return &Geometry{geom}, nil
}

func (g *Geometry) Buffer(
	h *Handle, width float64, quadsegs int) (*Geometry, error) {

	geom := C.GEOSBuffer_r(h.h, g.g, C.double(width), C.int(quadsegs))
	if geom == nil {
		return nil, ErrGeos
	}
	return &Geometry{geom}, nil
}

// Can only be called with LineString, LinearRing, and Point. Parent Geometry
// retains ownership of the coordseq
func (g *Geometry) CoordSeq(h *Handle) (*CoordSeq, error) {
	cs := C.GEOSGeom_getCoordSeq_r(h.h, g.g)
	if cs == nil {
		return nil, ErrGeos
	}
	return &CoordSeq{cs}, nil
}

func (g *Geometry) Intersection(h *Handle, o *Geometry) (*Geometry, error) {
	geom := C.GEOSIntersection_r(h.h, g.g, o.g)
	if geom == nil {
		return nil, ErrGeos
	}
	return &Geometry{geom}, nil
}

func (g *Geometry) Union(h *Handle, o *Geometry) (*Geometry, error) {
	geom := C.GEOSUnion_r(h.h, g.g, o.g)
	if geom == nil {
		return nil, ErrGeos
	}
	return &Geometry{geom}, nil
}

func (g *Geometry) Envelope(h *Handle) (*Geometry, error) {
	geom := C.GEOSEnvelope_r(h.h, g.g)
	if geom == nil {
		return nil, ErrGeos
	}
	return &Geometry{geom}, nil
}

// Can only be called with polygons. Polygon retains ownership of ring
func (g *Geometry) ExteriorRing(h *Handle) (*Geometry, error) {
	geom := C.GEOSGetExteriorRing_r(h.h, g.g)
	if geom == nil {
		return nil, ErrGeos
	}
	return &Geometry{geom}, nil
}

// Can only be called with polygons.
func (g *Geometry) NumInteriorRings(h *Handle) (int, error) {
	num := C.GEOSGetNumInteriorRings_r(h.h, g.g)
	if num < 0 {
		return 0, ErrGeos
	}
	return int(num), nil
}

// Can only be called with polygons. Polygon retains ownership of ring
func (g *Geometry) InteriorRingN(h *Handle, n int) (*Geometry, error) {
	geom := C.GEOSGetInteriorRingN_r(h.h, g.g, C.int(n))
	if geom == nil {
		return nil, ErrGeos
	}
	return &Geometry{geom}, nil
}

func (g *Geometry) Intersects(h *Handle, o *Geometry) (bool, error) {
	return predicate(C.GEOSIntersects_r(h.h, g.g, o.g))
}

func (g *Geometry) Contains(h *Handle, o *Geometry) (bool, error) {
	return predicate(C.GEOSContains_r(h.h, g.g, o.g))
}

func (g *Geometry) Disjoint(h *Handle, o *Geometry) (bool, error) {
	return predicate(C.GEOSDisjoint_r(h.h, g.g, o.g))
}

func (g *Geometry) Touches(h *Handle, o *Geometry) (bool, error) {
	return predicate(C.GEOSTouches_r(h.h, g.g, o.g))
}

func (g *Geometry) Overlaps(h *Handle, o *Geometry) (bool, error) {
	return predicate(C.GEOSOverlaps_r(h.h, g.g, o.g))
}

func (g *Geometry) Within(h *Handle, o *Geometry) (bool, error) {
	return predicate(C.GEOSWithin_r(h.h, g.g, o.g))
}

func (g *Geometry) IsEmpty(h *Handle) (bool, error) {
	return predicate(C.GEOSisEmpty_r(h.h, g.g))
}

func (g *Geometry) NumGeometries(h *Handle) (int, error) {
	// Valid for all Geometries, return -1 on error
	i := int(C.GEOSGetNumGeometries_r(h.h, g.g))
	if i < 0 {
		return 0, ErrGeos
	}
	return i, nil
}

func (g *Geometry) GeometryN(h *Handle, n int) (*Geometry, error) {
	geom := C.GEOSGetGeometryN_r(h.h, g.g, C.int(n))
	if geom == nil {
		return nil, ErrGeos
	}
	return &Geometry{geom}, nil
}

// https://trac.osgeo.org/geos/wiki/PreparedGeometry
// Not thread safe.
type PreparedGeometry struct {
	pg *C.GEOSPreparedGeometry
}

func (pg *PreparedGeometry) Destroy(h *Handle) {
	C.GEOSPreparedGeom_destroy_r(h.h, pg.pg)
}

func (pg *PreparedGeometry) Covers(h *Handle, o *Geometry) (bool, error) {
	return predicate(C.GEOSPreparedCovers_r(h.h, pg.pg, o.g))
}

// http://geos.osgeo.org/doxygen/classgeos_1_1io_1_1WKBReader.html
// Not thread safe.
type WKBReader struct {
	r *C.GEOSWKBReader
}

func NewWKBReader(h *Handle) *WKBReader {
	return &WKBReader{C.GEOSWKBReader_create_r(h.h)}
}

func (r *WKBReader) Destroy(h *Handle) {
	C.GEOSWKBReader_destroy_r(h.h, r.r)
}

func (r *WKBReader) Read(h *Handle, wkb []byte) (*Geometry, error) {
	if len(wkb) < 1 {
		return nil, ErrEmptyWKB
	}
	d := (*C.uchar)(&wkb[0])
	length := C.size_t(len(wkb))
	geom := C.GEOSWKBReader_read_r(h.h, r.r, d, length)
	if geom == nil {
		return nil, fmt.Errorf("Malformed WKB: %s", wkb)
	}
	return &Geometry{geom}, nil
}

// http://geos.osgeo.org/doxygen/classgeos_1_1io_1_1WKBWriter.html
// Not thread safe.
type WKBWriter struct {
	w *C.GEOSWKBWriter
}

func NewWKBWriter(h *Handle) *WKBWriter {
	return &WKBWriter{C.GEOSWKBWriter_create_r(h.h)}
}

func (w *WKBWriter) Destroy(h *Handle) {
	C.GEOSWKBWriter_destroy_r(h.h, w.w)
}

func (w *WKBWriter) Write(h *Handle, g *Geometry) []byte {
	size := C.size_t(1)
	wkb := unsafe.Pointer(C.GEOSWKBWriter_write_r(h.h, w.w, g.g, &size))
	defer C.free(wkb)
	return C.GoBytes(wkb, C.int(size))
}

// http://geos.osgeo.org/doxygen/classgeos_1_1io_1_1WKTReader.html
type WKTReader struct {
	r *C.GEOSWKTReader
}

func NewWKTReader(h *Handle) *WKTReader {
	return &WKTReader{C.GEOSWKTReader_create_r(h.h)}
}

func (r *WKTReader) Destroy(h *Handle) {
	C.GEOSWKTReader_destroy_r(h.h, r.r)
}

func (r *WKTReader) Read(h *Handle, wkt string) (*Geometry, error) {
	str := C.CString(wkt)
	defer C.free(unsafe.Pointer(str))
	geom := C.GEOSWKTReader_read_r(h.h, r.r, str)
	if geom == nil {
		return nil, fmt.Errorf("Malformed WKT: %s", wkt)
	}
	return &Geometry{geom}, nil
}

// http://geos.osgeo.org/doxygen/classgeos_1_1io_1_1WKTWriter.html
type WKTWriter struct {
	w *C.GEOSWKTWriter
}

func NewWKTWriter(h *Handle) *WKTWriter {
	return &WKTWriter{C.GEOSWKTWriter_create_r(h.h)}
}

func (w *WKTWriter) Destroy(h *Handle) {
	C.GEOSWKTWriter_destroy_r(h.h, w.w)
}

func (w *WKTWriter) Write(h *Handle, g *Geometry) string {
	str := C.GEOSWKTWriter_write_r(h.h, w.w, g.g)
	defer C.free(unsafe.Pointer(str))
	return C.GoString(str)
}
