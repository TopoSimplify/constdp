package constdp

import (
	"sort"
	"simplex/lnr"
	"simplex/ctx"
	"simplex/node"
	"github.com/intdxdt/rtree"
	"github.com/intdxdt/deque"
)

const EpsilonDist = 1.0e-5

func sort_ints(iter []int) []int {
	sort.Ints(iter)
	return iter
}

//Convert slice of interface to ints
func as_ints(iter []interface{}) []int {
	ints := make([]int, len(iter))
	for i, o := range iter {
		ints[i] = o.(int)
	}
	return ints
}

func castAsContextGeom(o interface{}) *ctx.ContextGeometry {
	return o.(*ctx.ContextGeometry)
}

func castAsNode(o interface{}) *node.Node {
	return o.(*node.Node)
}

func popLeftHull(que *deque.Deque) *node.Node {
	return que.PopLeft().(*node.Node)
}

func castConstDP(o lnr.SimpleAlgorithm) *ConstDP {
	return o.(*ConstDP)
}

func isSame(a, b lnr.SimpleAlgorithm) bool {
	return castConstDP(a) == castConstDP(b)
}

//node.Nodes from Rtree boxes
func nodesFromBoxes(iter []rtree.BoxObj) *node.Nodes {
	var self = node.NewNodes(len(iter))
	for _, h := range iter {
		self.Push(h.(*node.Node))
	}
	return self
}

//node.Nodes from Rtree nodes
func nodesFromRtreeNodes(iter []*rtree.Node) *node.Nodes {
	var self = node.NewNodes(len(iter))
	for _, h := range iter {
		self.Push(h.GetItem().(*node.Node))
	}
	return self
}
