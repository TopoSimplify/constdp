package constdp

import "simplex/geom"

//default comparator
func PtIdxCmp(a, b interface{}) int {
	_a, _b := a.(*geom.Point), b.(*geom.Point)
	return int(_a[2] - _b[2])
}

//default comparator
func int_cmp(a, b interface{}) int {
	return a.(int) - b.(int)
}

func hpt_cmp(a, b interface{}) int {
	return a.(*HPt).Compare(b.(*HPt))
}
