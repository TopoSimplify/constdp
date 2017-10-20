package constdp

import (
	"simplex/rng"
	"simplex/lnr"
	"simplex/node"
)

//split hull at vertex with
//maximum_offset offset -- k
func split_at_score_selection(self lnr.Linear, hull *node.Node) (*node.Node, *node.Node) {
	i, j := hull.Range.I(), hull.Range.J()
	k, _ := self.Score(self, hull.Range)
	// -------------------------------------------
	// i..[ha]..k..[hb]..j
	ha := node.New(self.Polyline(), rng.NewRange(i, k), hullGeom)
	hb := node.New(self.Polyline(), rng.NewRange(k, j), hullGeom)
	// -------------------------------------------
	return ha, hb
}

//split hull at indexes (index, index, ...)
func split_at_index(self lnr.Linear, hull *node.Node, idxs []int) []*node.Node {
	//formatter:off
	var pln       = self.Polyline()
	var ranges    = hull.Range.Split(idxs)
	var sub_hulls = make([]*node.Node, 0)
	for _, r := range ranges {
		sub_hulls = append(sub_hulls, node.New(pln, r, hullGeom))
	}
	return sub_hulls
}
