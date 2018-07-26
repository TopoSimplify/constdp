package constdp

import (
	"github.com/TopoSimplify/node"
	"github.com/TopoSimplify/hdb"
)

// Group hulls in hulldb by instance of ConstDP
func groupHullsByFC(hulldb *hdb.Hdb) {
	var ok bool
	var selfs = []*ConstDP{}
	var smap = make(map[string][]*node.Node)
	var nodes = hulldb.All()
	for _, h := range nodes {
		var lst []*node.Node
		var self = h.Instance.(*ConstDP)
		var id = self.Id()
		if lst, ok = smap[id]; !ok {
			lst = []*node.Node{}
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
