package constdp

import (
	"simplex/rng"
	"simplex/lnr"
)

//split hull at vertex with
//maximum_offset offset -- k
func split_at_score_selection(self lnr.Linear, hull *HullNode) (*HullNode, *HullNode) {
	i, j := hull.Range.I(), hull.Range.J()
	k, _ := self.Score(self, hull.Range)
	// -------------------------------------------
	// i..[ha]..k..[hb]..j
	ha := NewHullNode(self.Polyline(), rng.NewRange(i, k))
	hb := NewHullNode(self.Polyline(), rng.NewRange(k, j))
	// -------------------------------------------
	return ha, hb
}

//split hull at indexes (index, index, ...)
func split_at_index(self lnr.Linear, hull *HullNode, idxs []int) []*HullNode {
	//formatter:off
	pln         := self.Polyline()
	ranges      := hull.Range.Split(idxs)
	sub_hulls   := make([]*HullNode, 0)
	for _, r := range ranges {
		sub_hulls = append(sub_hulls, NewHullNode(pln, r))
	}
	return sub_hulls
}
