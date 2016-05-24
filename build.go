package constdp

import (
    "simplex/geom"
    . "simplex/struct/item"
    "simplex/dp"
)
/*
 description  build self
 */
func (self *ConstDP) Build(process ...func(item Item)) *ConstDP {
    fn := self.processhull
    if len(process) > 0 && process[0] != nil {
        fn = process[0]
    }
    //use superclass simplify
    self.DP.Build(fn)
    return self
}
/*
 description process hull
 private
 */
func (self *ConstDP) processhull(n Item) {
    node := n.(*dp.Node)
    var pln = self.Coordinates()[node.Key[0]: node.Key[1] + 1]
    if len(pln) > 1 {
        node.Hull = geom.NewPolygon(geom.ConvexHull(pln))
    }
}