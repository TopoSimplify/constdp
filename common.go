package constdp

import (
	"simplex/lnr"
	"simplex/ctx"
	"github.com/intdxdt/geom"
	"github.com/intdxdt/math"
)

const concurProcs = 8
const rtreeBucketSize = 4

func castAsContextGeom(o interface{}) *ctx.ContextGeometry {
	return o.(*ctx.ContextGeometry)
}

func castConstDP(o lnr.Linegen) *ConstDP {
	return o.(*ConstDP)
}

//hull point compare
func PointIndexCmp(a interface{}, b interface{}) int {
	var self, other = a.(*geom.Point), b.(*geom.Point)
	var d = self[2] - other[2]
	if math.FloatEqual(d, 0.0) {
		return 0
	} else if d < 0 {
		return -1
	}
	return 1
}
