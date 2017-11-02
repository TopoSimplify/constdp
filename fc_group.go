package constdp

import (
    "simplex/node"
    "github.com/intdxdt/rtree"
)

// Group hulls in hulldb by instance of ConstDP
func groupHullsByFC(hulldb *rtree.RTree) {
    var ok bool
    var hull *node.Node
    var selfs = make([]*ConstDP, 0)
    var smap = make(map[string]*node.Nodes)
    for _, h := range nodesFromRtreeNodes(hulldb.All()).DataView() {
        var lst *node.Nodes
        var self = castConstDP(h.Instance)
        var id = self.Id()
        if lst, ok = smap[id]; !ok {
            lst = node.NewNodes()
        }
        lst.Push(h)
        smap[id] = lst
    }

    for _, lst := range smap {
        var self = castConstDP(lst.Get(0).Instance)
        self.Hulls.Clear()
        for _, h := range lst.Sort().DataView() {
            self.Hulls.Append(h)
        }
        selfs = append(selfs, self)
    }

    for _, self := range selfs {
        self.SimpleSet.Empty() //update new simple
        for _, h := range *self.Hulls.DataView() {
            hull = castAsNode(h)
            self.SimpleSet.Extend(hull.Range.I(), hull.Range.J())
        }
    }
}
