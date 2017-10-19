package constdp

import (
	"github.com/intdxdt/rtree"
)

//split hull based on score selected vertex
func (self *ConstDP) deform_hulls(hulldb *rtree.RTree, selections *HullNodes) {
	selections.Reverse()
	for _, hull := range selections.list {
		ha, hb := split_at_score_selection(self, hull)
		hulldb.Remove(hull)

		self.Hulls.AppendLeft(hb)
		self.Hulls.AppendLeft(ha)
	}
	//empty selections
	selections.Empty()
}
