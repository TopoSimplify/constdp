package constdp

import (
    "simplex/geom"
    "simplex/util/iter"
)

//generates sub polyline from generator indices
func (self *ConstDP) subpoly(gen *iter.Generator) []*geom.Point {
    var poly = make([]*geom.Point, 0)

    for gen.Next {
        poly = append(poly, self.Pln[gen.Val()])
    }
    return poly
}
