package constdp

import (
	"github.com/intdxdt/geom"
	"github.com/intdxdt/math"
	"github.com/TopoSimplify/node"
	"sort"
)
const CacheKeySize = 6
const ConcurProcs  = 7

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

func CacheKey(a, b *node.Node) [CacheKeySize]int {
	var o = [CacheKeySize]int{
		a.Range.I, a.Range.J,
		b.Range.I, b.Range.J,
		a.Instance.Id(), b.Instance.Id(),
	}
	sort.Ints(o[:])
	return o
}