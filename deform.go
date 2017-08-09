package constdp

import "simplex/struct/rtree"

func (self *ConstDP) deform_hull(hulldb *rtree.RTree, hulls []*HullNode) {
	// split h at maximum_offset offset
	for _, h := range hulls {
		ha, hb := split_at_offset(self, h)
		hulldb.Remove(h)

		self.Hulls.AppendLeft(hb)
		self.Hulls.AppendLeft(ha)
	}
}
