package constdp

import (
    "simplex/geom"
    "simplex/util/iter"
)

//generates sub polyline from generator indices
func (cdp *ConstDP) subpoly(gen *iter.Generator) []*geom.Point {
    var poly = make([]*geom.Point, 0)
    for gen.Next {
        poly = append(poly, cdp.Pln.coords[gen.Val()])
    }
    return poly
}
