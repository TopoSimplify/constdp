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
