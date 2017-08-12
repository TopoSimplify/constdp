package constdp

import "simplex/struct/rtree"

//split hull based on score selected vertex
func (self *ConstDP) deform_hull(hulldb *rtree.RTree, hulls []*HullNode) {
	for _, hull := range hulls {
		ha, hb := split_at_score_selection(self, hull)
		hulldb.Remove(hull)

		self.Hulls.AppendLeft(hb)
		self.Hulls.AppendLeft(ha)
	}
}
