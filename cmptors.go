package constdp

import (
    . "simplex/geom"
    . "simplex/constrelate"
)


/*
 description check if sub geom is valid
 param subgeom
 param comparators
 return {boolean}
 private
 */
func (self *ConstDP) _isvalid(g Geometry , comparators []Comparator) bool {
  //make true , proof otherwise
  var bln = true
  for  i := 0; bln && i < len(comparators); i++ {
      bln = bln && comparators[i](g)
  }
  return bln
}


/*
 description gen cmp functors
 param polygeom
 param constlist
 param options
 returns {Array}
 private
 */
func (self *ConstDP)_cmptors(g Geometry, constlist []Geometry) []Comparator {
    var consts = self.opts.Constraints
    var comparators = make([]Comparator, len(consts))
    for i := 0; i < len(consts); i++ {
        comparators[i] = consts[i](g, constlist) //return cmptor
    }
    return comparators
}

