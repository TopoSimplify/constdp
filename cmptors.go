package constdp

import (
    . "simplex/geom"
    . "simplex/relations"
)


 //check if sub geom is valid
func (self *ConstDP) _isvalid(g Geometry, comparators []Comparator) bool {
    //make true , proof otherwise
    var bln = true
    for i := 0; bln && i < len(comparators); i++ {
        bln = bln && comparators[i](g)
    }
    return bln
}


// gen cmp functors
func (self *ConstDP) _cmptors(g Geometry, constlist []Geometry) []Comparator {
    var relates = self.opts.Relations
    var comparators = make([]Comparator, len(relates))

    for i := 0; i < len(relates); i++ {
        var fn Relations = relates[i]
        comparators[i] = fn.Relate(g, constlist) //return cmptor
    }
    return comparators
}

