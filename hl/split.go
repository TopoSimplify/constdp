package hl

import (
	"sort"
	"simplex/struct/sset"
	"simplex/constdp/rng"
	"simplex/struct/rtree"
	"simplex/constdp/db"
	"simplex/constdp/ln"
)

//split hull at vertex with
//maximum_offset offset -- k
func  SplitHull(self ln.Linear, hull *HullNode) (*HullNode, *HullNode) {
	i, j := hull.Range.I(), hull.Range.J()
	k, _ := self.MaximumOffset(self, hull.Range)
	// -------------------------------------------
	// i..[ha]....k...[hb].....j
	ha := NewHullNode(self.Polyline(), rng.NewRange(i, k), rng.NewRange(i, j))
	hb := NewHullNode(self.Polyline(), rng.NewRange(k, j), rng.NewRange(i, j))
	// -------------------------------------------
	return ha, hb
}

//	split hull at indexes (index, index, ...)
func  SplitHullAtIndex(self ln.Linear, hull *HullNode, idxs []int) []*HullNode {
	pln := self.Polyline()
	i, j := hull.Range.I(), hull.Range.J()
	subhulls := make([]*HullNode, 0)
	for _, k := range idxs {
		if i < k && k < j {
			ar, br, pr := rng.NewRange(i, k), rng.NewRange(k, j), rng.NewRange(i, j)
			ha := NewHullNode(pln, ar, pr)
			hb := NewHullNode(pln, br, pr)
			subhulls = append(subhulls, ha)
			subhulls = append(subhulls, hb)
		}
	}
	return subhulls

}

//merge contig hulls after split - merge line segment fragments
func  MergeContigFragments(
	self ln.Linear,
	hulls []*HullNode,
	tree *rtree.RTree,
	vertex_set *sset.SSet,
) ([]*HullNode, []*HullNode) {

	pln := self.Polyline()
	keep, rm := make([]*HullNode, 0), make([]*HullNode, 0)

	for _, h := range hulls {
		hs := db.KNN(tree, h.Geom, 1.0e-5)
		hr := h.Range
		m := h

		if hr.Size() == 1 {
			for _, _s := range hs {
				s := _s.(*HullNode)
				sr := s.Range
				bln := (hr.J() == sr.I() && vertex_set.Contains(sr.I())) ||
					(hr.I() == sr.J() && vertex_set.Contains(sr.J()))

				if !bln && (hr.Contains(sr.I()) || hr.Contains(sr.J())) {
					l := []int{sr.I(), sr.J(), hr.I(), hr.J()}
					sort.Ints(l)

					r := rng.NewRange(l[0], l[len(l)-1])
					m = NewHullNode(pln, r, r)

					// add to remove list to remove , after merge
					rm = append(rm, s)
					break
				}
			}
		}
		keep = append(keep, m)
	}
	return keep, rm
}
