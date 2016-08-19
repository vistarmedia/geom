# Libgeos bindings for golang

* Requires libgeos 3.5.0 or greater to be linked.

Provides a small wrapper over the Cgo generated bindings. Converts C types and
errors in to native go types. These are not thread or memory safe. The user must
take care to not share handle, WKB reader/writers, and prepared geometries
across goroutines. The user is also responsible for knowing libgeos ownership
semantics and destroying objects when they are no longer in use. For example:

    h := NewHandle()
    defer h.Destroy()

    cs := NewCoordSeq(h, 1, 2)
    cs.SetX(h, 0, 10)
    cs.SetY(h, 0, 20)
    point := cs.Point(h)
    defer point.Destroy(h)

In this example, the point takes ownership of the CoordSeq so destroying both
the CoordSeq and Point would cause a double free.

## Resources

Libgeos:
https://trac.osgeo.org/geos/

Libgeos reentrant API:
https://raw.githubusercontent.com/libgeos/libgeos/6741c5e725d90cc979b08d6337dd1cdce1ff59d1/capi/geos\_ts\_c.cpp
