package constdp

import (
    . "simplex/geom"
)

//self intersections
func (self *ConstDP) self_intersections() []*Point {

    if len(self.intersections) > 0 {
        return self.intersections
    }

    var selfgeom = NewLineString(self.Pln)
    if selfgeom.IsSimple() {
        self.intersections = make([]*Point, 0)
    } else {
        self.intersections = selfgeom.SelfIntersection()
    }

    return self.intersections
}
