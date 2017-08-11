package constdp

import (
	"sort"
	"simplex/geom"
	"simplex/geom/mbr"
	"simplex/constdp/ln"
	"simplex/struct/sset"
	"simplex/constdp/db"
	"simplex/constdp/box"
	"simplex/constdp/rng"
	"simplex/struct/rtree"
)

func contiguous_fragments_at_threshold(self *ConstDP, ha, hb *HullNode) *HullNode {
	m := contiguous_fragments(self, ha, hb)
	_, score := self.score(self, m.Range)
	if score <= self.Opts.Threshold {
		return m
	}
	return nil
}

//merge contiguous hulls
func contiguous_fragments(self *ConstDP, ha, hb *HullNode) *HullNode {
	var l = append(ha.Range.AsSlice(), hb.Range.AsSlice()...)
	sort.Ints(l)
	// i...[ha]...k...[hb]...j
	i, j := l[0], l[len(l)-1]
	r := rng.NewRange(i, j)
	return NewHullNode(self.Pln, r, r.Clone())
}

//merge contig hulls after split - merge line segment fragments
func find_mergeable_contiguous_fragments(
	self ln.Linear, hulls []*HullNode, hulldb *rtree.RTree,
	vertex_set *sset.SSet, ) ([]*HullNode, []*HullNode) {
	pln := self.Polyline()
	keep, rm := make([]*HullNode, 0), make([]*HullNode, 0)

	hdict := make(map[[2]int]*HullNode, 0)
	for _, h := range hulls {
		hdict[h.Range.AsArray()] = h

		hr := h.Range
		//if hr.Size() < 4{
		if hr.Size() == 1 {
			//@formatter:off
			hs_knn      := find_context_hulls(hulldb, h, EpsilonDist)

			hs := make([]*HullNode, len(hs_knn))
			for i, h := range hs_knn {
				hs[i] = h.(*HullNode)
			}
			sort_hulls(hs) // sort hulls for consistency

			for _, s := range hs {
				sr := s.Range
				//test whether sr.i or sr.j is a self inter-vertex -- split point
				//not sr.i != hr.i or sr.j != hr.j without i/j being a inter-vertex
				//tests for contiguous and whether contiguous index is part of vertex set
				//if the location at which they are contiguous is not part of vertex set then
				//its mergeable
				mergeable := (hr.J() == sr.I() && !vertex_set.Contains(sr.I())) ||
							 (hr.I() == sr.J() && !vertex_set.Contains(sr.J()))

				if mergeable {
					l := append(sr.AsSlice(), hr.AsSlice()...)
					sort.Ints(l)

					// rm sr + hr
					delete(hdict, sr.AsArray())
					delete(hdict, hr.AsArray())

					r := rng.NewRange(l[0], l[len(l)-1])
					m := NewHullNode(pln, r, r.Clone())

					// add merge
					hdict[m.Range.AsArray()] = m

					// add to remove list to remove , after merge
					rm = append(rm, s)
					rm = append(rm, h)
					break
				}
			}
		}
	}

	for _, v := range hdict {
		keep = append(keep, v)
	}
	return keep, rm
}
