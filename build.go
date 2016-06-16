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
    process_hull(self.Coordinates(), n.(*dp.Node))
}

//process hull
func process_hull(pln []*geom.Point, node *dp.Node){
    pln = pln[node.Key[0]: node.Key[1] + 1]
    if len(pln) > 1 {
        node.Hull = geom.NewPolygon(geom.ConvexHull(pln))
    }
}