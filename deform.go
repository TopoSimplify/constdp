package constdp

import (
	"github.com/intdxdt/rtree"
	"simplex/split"
)

//split hull based on score selected vertex
func (self *ConstDP) deform_hulls(hulldb *rtree.RTree, selections *HullNodes) {
	selections.Reverse()
	for _, hull := range selections.list {
		ha, hb := split.AtScoreSelection(self, hull, hullGeom)
		hulldb.Remove(hull)

		self.Hulls.AppendLeft(hb)
		self.Hulls.AppendLeft(ha)
	}
	//empty selections
	selections.Empty()
}
