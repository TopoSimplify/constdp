package constdp

import (
	"sort"
	"simplex/struct/rtree"
	"simplex/struct/deque"
	"simplex/struct/sset"
	"simplex/constdp/cmp"
	"simplex/constdp/ctx"
)

var Debug = false

const Epsilon = 1.0e-8
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

//Rtree nodes as hull nodes
func as_hullnodes(iter []*rtree.Node) []*HullNode {
	nodes := make([]*HullNode, len(iter))
	for i, h := range iter {
		nodes[i] = h.GetItem().(*HullNode)
	}
	return nodes
}

//Rtree nodes boxes as hull nodes
func as_hullnodes_from_boxes(iter []rtree.BoxObj) []*HullNode {
	nodes := make([]*HullNode, len(iter))
	for i, h := range iter {
		nodes[i] = h.(*HullNode)
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

//map[range]*hullnode to slice of hullnode
func map_to_slice(dict map[[2]int]*HullNode, s []*HullNode) []*HullNode {
	for _, o := range dict {
		s = append(s, o)
	}
	return s
}

func simple_hulls_as_ptset(hulls []*HullNode) *sset.SSet {
	var ptset = sset.NewSSet(cmp.IntCmp)
	for _, o := range hulls {
		ptset.Extend(o.Range.I(), o.Range.J())
	}
	return ptset
}

func pop_left_hull(que *deque.Deque) *HullNode {
	return que.PopLeft().(*HullNode)
}

func pop_hull_from_slice(ptr *[]*HullNode) {
	slice := *ptr
	n := len(slice) - 1
	slice[n] = nil
	*ptr = slice[:n]
}

func empty_hull_slice(ptr *[]*HullNode) {
	slice := *ptr
	for i := range slice{
		slice[i] = nil
	}
	*ptr = slice[:0]
}

func cast_as_context_geom(o interface{}) *ctx.CtxGeom {
	return o.(*ctx.CtxGeom)
}

func cast_as_hullnode(o interface{}) *HullNode{
	return o.(*HullNode)
}

