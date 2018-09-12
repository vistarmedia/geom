// Provides memory management utilities for libgeos bindings
package memory

import (
	"runtime"

	"github.com/vistarmedia/geom/geos-go"
	"github.com/vistarmedia/geom/geos-go/handle"
)

type GeosDestroyable interface {
	Destroy(h *geos.Handle)
}

// Attaches a finalizer to a geos object
func GeosGCManaged(hp handle.GeosHandleProvider, d GeosDestroyable) {
	runtime.SetFinalizer(d, func(d1 GeosDestroyable) {
		h := hp.Get()
		d1.Destroy(h)
		hp.Put(h)
	})
}
