package constdp



/*
 description update list of constraints with
 * intersection points with neighbours
 param constlist
 param subgeom
 param node
 param options{Object}
 private
 */
//func (self *ConstDP) updateconsts(constlist []geom.Geometry, subgeom geom.Geometry, node *Node, options) {
//
//  var interlist = [], xorplns = []
//
//  var avoid_self_intersection = options &&
//                                (options.avoid_self || options.self)
//
//  var preserve_complex = options &&
//                         (options.preserve_complex || options.complex)
//  //avoid self intersection
//  if avoid_self_intersection {
//    var i = node[node._key][0]
//    var j = node[node._key][1]
//    var len = self.root[node._key][1]
//    xorplns.append.apply(xorplns, self._xor_subpln(i, j, len))
//  }
//  //preserve complex geometries
//  if preserve_complex {
//    interlist.append.apply(interlist, self._self_intersection())
//  }
//
//  var pts
//  _.each(constlist, func (g) {
//    var glist = subgeom._lineother(g)
//    _.each(glist, func (g) {
//      pts = subgeom.intersection(g)
//      interlist.append.apply(interlist, pts)
//    })
//  })
//  constlist.append.apply(constlist, interlist)
//  constlist.append.apply(constlist, xorplns)
//}
