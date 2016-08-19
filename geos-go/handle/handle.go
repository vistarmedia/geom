// Provdies a thread safe interface for obtaining a libgeos handle
package handle

import (
	"runtime"
	"sync"

	"vistarmedia.com/vistar/geom/geos-go"
)

// Leases GEOS handles in a thread safe manner. GEOS handles can not be shared
// between goroutines. Any implementation of GeosHandleProvider must enforce
// that invariant.
type GeosHandleProvider interface {
	Get() *geos.Handle
	Put(*geos.Handle)
}

// sync.Pool backed GeosHandleProvider
type PooledHandleProvider struct {
	pool *sync.Pool
}

func NewPooledHandleProvider() GeosHandleProvider {
	return PooledHandleProvider{
		&sync.Pool{New: func() interface{} {
			h := geos.NewHandle()
			runtime.SetFinalizer(h, func(h1 *geos.Handle) {
				h1.Destroy()
			})
			return h
		}},
	}
}

func (hp PooledHandleProvider) Get() *geos.Handle {
	return hp.pool.Get().(*geos.Handle)
}

func (hp PooledHandleProvider) Put(h *geos.Handle) {
	hp.pool.Put(h)
}
