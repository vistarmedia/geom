// Well Known Text encoding/decoding
package wkt

import (
	"vistarmedia.com/vistar/geom"
	"vistarmedia.com/vistar/geom/geos-go"
	"vistarmedia.com/vistar/geom/geos-go/handle"
	"vistarmedia.com/vistar/geom/geos-go/memory"
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
	memory.GeosGCManaged(hp, writer)
	return &Encoder{
		hp:     hp,
		writer: writer,
	}
}

func (e *Encoder) Encode(g Encodeable) string {
	h := e.hp.Get()
	defer e.hp.Put(h)
	return e.writer.Write(h, g.UnsafeToGeos())
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
	memory.GeosGCManaged(hp, reader)
	return &Decoder{
		hp:      hp,
		reader:  reader,
		factory: fact,
	}
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
