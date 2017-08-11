package constdp

import (
	"simplex/constdp/rng"
	"simplex/constdp/ln"
)

//split hull at vertex with
//maximum_offset offset -- k
func split_at_offset(self ln.Linear, hull *HullNode) (*HullNode, *HullNode) {
	i, j := hull.Range.I(), hull.Range.J()
	k, _ := self.Score(self, hull.Range)
	// -------------------------------------------
	// i..[ha]..k..[hb]..j
	ha := NewHullNode(self.Polyline(), rng.NewRange(i, k), rng.NewRange(i, j))
	hb := NewHullNode(self.Polyline(), rng.NewRange(k, j), rng.NewRange(i, j))
	// -------------------------------------------
	return ha, hb
}

//split hull at indexes (index, index, ...)
func split_at_index(self ln.Linear, hull *HullNode, idxs []int) []*HullNode {
	//formatter:off
	pr          := hull.Range
	pln         := self.Polyline()
	ranges      := pr.Split(idxs)
	sub_hulls   := make([]*HullNode, 0)
	for _, r := range ranges {
		sub_hulls = append(sub_hulls, NewHullNode(pln, r, pr))
	}
	return sub_hulls
}
