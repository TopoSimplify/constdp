package constdp

import (
    "simplex/geom"
    "simplex/dp"
    . "simplex/struct/item"
)

//Build
func (self *ConstDP) Build() *ConstDP {
    fn := self.processhull
    if self.opts.Process != nil {
        fn = self.opts.Process
    }
    //use superclass simplify
    self.DP.Build(fn)
    return self
}

//process hull
func (self *ConstDP) processhull(n Item) {
    node := n.(*dp.Node)
    var pln = self.Coordinates()[node.Key[0]: node.Key[1] + 1]
    if len(pln) > 1 {
        node.Hull = geom.NewPolygon(geom.ConvexHull(pln))
    }
}