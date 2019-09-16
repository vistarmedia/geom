// Well Known Text encoding/decoding
package wkt

import (
	"runtime"

	"github.com/vistarmedia/geom"
	"github.com/vistarmedia/geom/geos-go"
	"github.com/vistarmedia/geom/geos-go/handle"
)

type Encodeable interface {
	UnsafeToGeos() *geos.Geometry
}

type Encoder struct {
	writer *geos.WKTWriter
	hp     handle.GeosHandleProvider
}

func NewEncoder(hp handle.GeosHandleProvider) *Encoder {
	h := hp.Get()
	writer := geos.NewWKTWriter(h)
	hp.Put(h)
	encoder := &Encoder{
		hp:     hp,
		writer: writer,
	}
	runtime.SetFinalizer(encoder, func(encoder1 *Encoder) {
		h := encoder1.hp.Get()
		encoder1.writer.Destroy(h)
		encoder1.hp.Put(h)
	})
	return encoder
}

func (e *Encoder) Encode(g Encodeable) string {
	h := e.hp.Get()
	defer e.hp.Put(h)
	wkt := e.writer.Write(h, g.UnsafeToGeos())
	runtime.KeepAlive(g)
	return wkt
}

type Decoder struct {
	reader  *geos.WKTReader
	hp      handle.GeosHandleProvider
	factory geom.Factory
}

func NewDecoder(hp handle.GeosHandleProvider, fact geom.Factory) *Decoder {
	h := hp.Get()
	reader := geos.NewWKTReader(h)
	hp.Put(h)
	decoder := &Decoder{
		hp:      hp,
		reader:  reader,
		factory: fact,
	}
	runtime.SetFinalizer(decoder, func(decoder1 *Decoder) {
		h := decoder1.hp.Get()
		decoder1.reader.Destroy(h)
		decoder1.hp.Put(h)
	})
	return decoder
}

func (d *Decoder) Decode(wkt string) (*geom.Geometry, error) {
	h := d.hp.Get()
	geom, err := d.reader.Read(h, wkt)
	d.hp.Put(h)
	if err != nil {
		return nil, err
	}
	return d.factory.FromGeos(geom), nil
}
