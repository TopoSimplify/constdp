package constdp

import (
	"simplex/struct/rtree"
	"simplex/struct/deque"
)

const Epsilon = 1.0e-8
const EpsilonDist = 1.0e-5

//Convert slice of interface to ints
func as_ints(iter []interface{}) []int {
	ints := make([]int, len(iter))
	for i, o := range iter {
		ints[i] = o.(int)
	}
	return ints
}

//Rtree nodes as hull nodes
func as_hullnodes(iter []*rtree.Node) []*HullNode {
	nodes := make([]*HullNode, len(iter))
	for i, h := range iter {
		nodes[i] = h.GetItem().(*HullNode)
	}
	return nodes
}

//Hull nodes as deque
func as_deque(iter []*HullNode) *deque.Deque {
	queue := deque.NewDeque()
	for _, h := range iter {
		queue.Append(h)
	}
	return queue
}
