package constdp

import (
    "simplex/geom"
    . "simplex/struct/item"
    "simplex/dp"
)
/*
 description  build self
 */
func (self *ConstDP) Build() *ConstDP {
    fn := self.processhull
    if self.opts.Process != nil {
        fn = self.opts.Process
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