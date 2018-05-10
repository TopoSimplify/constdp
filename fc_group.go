package constdp

import (
	"github.com/intdxdt/rtree"
	"github.com/TopoSimplify/node"
	"github.com/TopoSimplify/common"
)

// Group hulls in hulldb by instance of ConstDP
func groupHullsByFC(hulldb *rtree.RTree) {
	var ok bool
	var selfs = []*ConstDP{}
	var smap = make(map[string][]*node.Node)

	for _, h := range common.NodesFromRtreeNodes(hulldb.All()) {
		var lst []*node.Node
		var self = h.Instance.(*ConstDP)
		var id = self.Id()
		if lst, ok = smap[id]; !ok {
			lst = make([]*node.Node, 0)
		}
		lst = append(lst, h)
		smap[id] = lst
	}

	for _, lst := range smap {
		var self = (lst[0].Instance).(*ConstDP)
		node.Clear(&self.Hulls)
		for _, h := range lst {
			self.Hulls = append(self.Hulls, h)
		}
		selfs = append(selfs, self)
	}

	for _, self := range selfs {
		self.SimpleSet.Empty() //update new simple
		for _, hull := range self.Hulls {
			self.SimpleSet.Extend(hull.Range.I, hull.Range.J)
		}
	}
}
