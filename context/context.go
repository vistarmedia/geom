// Provides a convenience wrapper for the geom and encoding packages.
package context

import (
	"github.com/vistarmedia/geom"
	"github.com/vistarmedia/geom/encoding/wkb"
	"github.com/vistarmedia/geom/encoding/wkt"
	"github.com/vistarmedia/geom/geos-go/handle"
)

type Context struct {
	hp handle.GeosHandleProvider
}

// Only one context should be instantiated per application.
func NewContext() Context {
	return Context{handle.NewPooledHandleProvider()}
}

func (ctx Context) Factory() geom.Factory {
	return geom.NewFactory(ctx.hp)
}

func (ctx Context) WKTEncoder() *wkt.Encoder {
	return wkt.NewEncoder(ctx.hp)
}

func (ctx Context) WKTDecoder() *wkt.Decoder {
	return wkt.NewDecoder(ctx.hp, ctx.Factory())
}

func (ctx Context) WKBEncoder() *wkb.Encoder {
	return wkb.NewEncoder(ctx.hp)
}

func (ctx Context) WKBDecoder() *wkb.Decoder {
	return wkb.NewDecoder(ctx.hp, ctx.Factory())
}
