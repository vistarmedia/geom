// Well Known Binary encoding/decoding. Assumes machine endianness.
package wkb

import (
	"vistarmedia.com/vistar/geom"
	"vistarmedia.com/vistar/geom/geos-go"
	"vistarmedia.com/vistar/geom/geos-go/handle"
)

type Encodeable interface {
	UnsafeToGeos() *geos.Geometry
}

type Encoder struct {
	hp handle.GeosHandleProvider
}

func NewEncoder(hp handle.GeosHandleProvider) *Encoder {
	return &Encoder{hp}
}

func (e *Encoder) Encode(g Encodeable) []byte {
	h := e.hp.Get()
	defer e.hp.Put(h)
	// Unlike WKT, the geos WKB reader and writer is not thread safe. The object
	// should be cheap to construct so we just create and destory a writer every
	// time.
	writer := geos.NewWKBWriter(h)
	defer writer.Destroy(h)
	return writer.Write(h, g.UnsafeToGeos())
}

type Decoder struct {
	hp      handle.GeosHandleProvider
	factory geom.Factory
}

func NewDecoder(hp handle.GeosHandleProvider, fact geom.Factory) *Decoder {
	return &Decoder{
		hp:      hp,
		factory: fact,
	}
}

func (d *Decoder) Decode(wkb []byte) (*geom.Geometry, error) {
	h := d.hp.Get()
	reader := geos.NewWKBReader(h)
	defer reader.Destroy(h)
	geom, err := reader.Read(h, wkb)
	d.hp.Put(h)
	if err != nil {
		return nil, err
	}
	return d.factory.FromGeos(geom), nil
}
