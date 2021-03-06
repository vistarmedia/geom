package context

import (
	"fmt"
	"testing"

	"github.com/vistarmedia/geom"
)

func TestExample(t *testing.T) {
	ctx := NewContext()
	fact := ctx.Factory()
	point, _ := fact.NewPoint(geom.Coord{4, 2})
	fmt.Println(ctx.WKTEncoder().Encode(point))
}
