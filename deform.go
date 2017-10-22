package constdp

import (
	"simplex/node"
	"simplex/split"
	"github.com/intdxdt/rtree"
	"simplex/dp"
)

//split hull based on score selected vertex
func (self *ConstDP) deform_hulls(nodedb *rtree.RTree, selections *node.Nodes) {
	selections.Reverse()
	for _, hull := range selections.DataView() {
		var ha, hb = split.AtScoreSelection(self, hull, dp.NodeGeometry)
		nodedb.Remove(hull)

		self.Hulls.AppendLeft(hb)
		self.Hulls.AppendLeft(ha)
	}
	//empty selections
	selections.Empty()
}
