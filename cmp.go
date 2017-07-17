package constdp

import (
	"simplex/geom"
	"simplex/util/math"
)

//int cmp
func IntCmp(a, b interface{}) int {
	return a.(int) - b.(int)
}

//hullpt compare
func PointIndexCmp(a interface{} , b interface{}) int {
	self, other := a.(*geom.Point), b.(*geom.Point)
	d := self[2] - other[2]
	if math.FloatEqual(d, 0.0) {
		return 0
	} else if d < 0 {
		return -1
	}
	return 1
}

