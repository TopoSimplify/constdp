package constdp

import (
	"sort"
	"simplex/ctx"
	"github.com/intdxdt/deque"
	"simplex/node"
	"simplex/lnr"
)

var Debug = false

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

func cast_as_context_geom(o interface{}) *ctx.CtxGeom {
	return o.(*ctx.CtxGeom)
}

func cast_as_hullnode(o interface{}) *node.Node {
	return o.(*node.Node)
}

func pop_left_hull(que *deque.Deque) *node.Node {
	return que.PopLeft().(*node.Node)
}

func cast_cdp(o lnr.SimpleAlgorithm) *ConstDP {
	return o.(*ConstDP)
}


func is_same(a, b lnr.SimpleAlgorithm) bool{
	return cast_cdp(a) == cast_cdp(b)
}

//Rtree nodes as hull nodes
//func as_hullnodes(iter []*rtree.Node) []*node.Node {
//	nodes := make([]*node.Node, len(iter))
//	for i, h := range iter {
//		nodes[i] = h.GetItem().(*node.Node)
//	}
//	return nodes
//}

//Rtree nodes boxes as hull nodes
//func as_hullnodes_from_boxes(iter []rtree.BoxObj) []*node.Node {
//	nodes := make([]*node.Node, len(iter))
//	for i, h := range iter {
//		nodes[i] = h.(*node.Node)
//	}
//	return nodes
//}

//Hull nodes as deque
//func as_deque(iter []*node.Node) *deque.Deque {
//	queue := deque.NewDeque()
//	for _, h := range iter {
//		queue.Append(h)
//	}
//	return queue
//}

//map[range]*node.Node to slice of node.Node
//func map_to_slice(dict map[[2]int]*node.Node, s []*node.Node) []*node.Node {
//	for _, o := range dict {
//		s = append(s, o)
//	}
//	return s
//}

//func simple_hulls_as_ptset(hulls []*node.Node) *sset.SSet {
//	var ptset = sset.NewSSet(cmp.Int)
//	for _, o := range hulls {
//		ptset.Extend(o.Range.I(), o.Range.J())
//	}
//	return ptset
//}

//func empty_hull_slice(ptr *[]*node.Node) {
//	slice := *ptr
//	for i := range slice {
//		slice[i] = nil
//	}
//	*ptr = slice[:0]
//}
