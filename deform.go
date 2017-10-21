package constdp

import (
	"simplex/node"
	"simplex/split"
	"github.com/intdxdt/rtree"
)

//split hull based on score selected vertex
func (self *ConstDP) deform_hulls(hulldb *rtree.RTree, selections *node.Nodes) {
	selections.Reverse()
	for _, hull := range selections.DataView() {
		ha, hb := split.AtScoreSelection(self, hull, hullGeom)
		hulldb.Remove(hull)

		self.Hulls.AppendLeft(hb)
		self.Hulls.AppendLeft(ha)
	}
	//empty selections
	selections.Empty()
}
