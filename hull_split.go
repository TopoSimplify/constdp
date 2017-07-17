package constdp

import (
	"simplex/struct/rtree"
	"simplex/struct/sset"
	"sort"
)

//split hull at vertex with
//maximum_offset offset -- k
func (self *ConstDP) split_hull(hull *HullNode) (*HullNode, *HullNode) {
	i, j := hull.Range.i, hull.Range.j
	k, _ := self.MaximumOffset(self, hull.Range)
	// -------------------------------------------
	// i..[ha]....k...[hb].....j
	ha := NewHullNode(self.Pln, NewRange(i, k), NewRange(i, j))
	hb := NewHullNode(self.Pln, NewRange(k, j), NewRange(i, j))
	// -------------------------------------------
	return ha, hb
}

//	split hull at indexes (index, index, ...)
func (self *ConstDP) split_hull_at_index(hull *HullNode, idxs []int) []*HullNode {
	pln := self.Pln
	i, j := hull.Range.i, hull.Range.j
	subhulls := make([]*HullNode, 0)
	for _, k := range idxs {
		if i < k && k < j {
			ar, br, pr := NewRange(i, k), NewRange(k, j), NewRange(i, j)
			ha := NewHullNode(pln, ar, pr)
			hb := NewHullNode(pln, br, pr)
			subhulls = append(subhulls, ha)
			subhulls = append(subhulls, hb)
		}
	}
	return subhulls

}

//merge contig hulls after split - merge line segment fragments
func (self *ConstDP) merge_contig_fragments(
	hulls []*HullNode,
	db *rtree.RTree,
	vertex_set *sset.SSet,
) ([]*HullNode, []*HullNode) {

	pln := self.Pln
	keep, rm := make([]*HullNode, 0), make([]*HullNode, 0)

	for _, h := range hulls {
		hs := dbKNN(db, h.Geom, 1.0e-5)
		hr := h.Range
		m := h

		if hr.Size() == 1 {
			for _, _s := range hs {
				s := _s.(*HullNode)
				sr := s.Range
				bln := (hr.j == sr.i && vertex_set.Contains(sr.i)) ||
					(hr.i == sr.j && vertex_set.Contains(sr.j))

				if !bln && (hr.Contains(sr.i) || hr.Contains(sr.j)) {
					l := []int{sr.i, sr.j, hr.i, hr.j}
					sort.Ints(l)

					r := NewRange(l[0], l[len(l)-1])
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
