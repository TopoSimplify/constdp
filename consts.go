package constdp

import (
    . "simplex/geom"
    . "simplex/dp"
)

 // update list of constraints with
 // intersection points with neighbours
func (self *ConstDP) updateconsts(constlist []Geometry, subgeom *LineString,
node *Node) []Geometry {

  var interlist = make([]Geometry , 0)
  var xorplns = make([]Geometry,    0)

  var avoid_self_intersection = self.opts.AvoidSelfIntersection

  var preserve_complex = self.opts.PreserveComplex
  //avoid self intersection
  if avoid_self_intersection {
    var i, j  = node.Key[0], node.Key[1]
    var len = len(self.Pln) -1
      for _, g := range self.xor_subpln(i, j, len){
          xorplns = append(xorplns, g )
      }
  }
  //preserve complex geometries
  if preserve_complex {
      for _, g := range self.self_intersections(){
            interlist = append(interlist, g)
      }
  }

  var pts []*Point
  for _, g := range constlist {
    var glist = g.AsLinear()
    for _, g := range glist {
        pts = subgeom.Intersection(g)
        for _, pt := range pts {
            interlist = append(interlist, pt)
        }
    }
  }
    for _, g := range interlist{
        constlist = append(constlist,g)
    }
    for _, g := range xorplns {
        constlist = append(constlist,g)
    }
    return constlist
}


func (self *ConstDP) context_neighbours(node *Node) []Geometry {
    var db = self.opts.Db
    var hull = node.Hull
    var neighbours = make([]Geometry, 0)
    if db != nil {
        neighbours = SearchDb(db, hull.BBox())
    }
    return neighbours;
}