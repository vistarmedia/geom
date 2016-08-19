# Libgeos memory management

Utility package for the consumers of the libgeos bindings. Provides a function
to manage libgeos objects with the go's GC.

    hp := NewPooledHandleProvider()
    h := hp.Get()
    geom := geos.NewEmptyPolygon(h)
    GeosGCManaged(h)
    hp.Put(h)
