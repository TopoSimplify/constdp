package constdp

import (
	"simplex/struct/rtree"
)

//split hull based on score selected vertex
func (self *ConstDP) deform_hull(hulldb *rtree.RTree, selections *[]*HullNode) {
	sort_reverse(*selections)
	for _, hull := range *selections {
		ha, hb := split_at_score_selection(self, hull)
		hulldb.Remove(hull)

		self.Hulls.AppendLeft(hb)
		self.Hulls.AppendLeft(ha)
	}
	//empty selections
	empty_hull_slice(selections)
}
