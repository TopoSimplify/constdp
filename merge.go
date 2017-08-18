package constdp

import (
	"simplex/struct/sset"
	"simplex/constdp/rng"
	"simplex/struct/rtree"
)

func merge_contiguous_fragments_at_threshold(self *ConstDP, ha, hb *HullNode) *HullNode {
	m := self.merge_contiguous_fragments(ha, hb)
	_, val := self.Score(self, m.Range)
	if self.is_score_relate_valid(val) {
		return m
	}
	return nil
}

//merge contiguous hulls
func (self *ConstDP) merge_contiguous_fragments(ha, hb *HullNode) *HullNode {
	var ranges = sort_ints(append(ha.Range.AsSlice(), hb.Range.AsSlice()...))
	// i...[ha]...k...[hb]...j
	i, j := ranges[0], ranges[len(ranges)-1]
	return NewHullNode(self.Pln, rng.NewRange(i, j), rng.NewRange(i, j))
}

//merge contig hulls after split - merge line segment fragments
func (self *ConstDP) find_mergeable_contiguous_fragments(
	hulls []*HullNode, hulldb *rtree.RTree,
	vertex_set *sset.SSet, unmerged map[[2]int]*HullNode,
) ([]*HullNode, []*HullNode) {

	//@formatter:off
	var pln       = self.Polyline()
	var keep      = make([]*HullNode, 0)
	var rm        = make([]*HullNode, 0)

	var hdict    = make(map[[2]int]*HullNode, 0)
	var mrgdict  = make(map[[2]int]*HullNode, 0)

	var is_merged = func(o *rng.Range) bool {
		_, ok := mrgdict[o.AsArray()]
		return ok
	}

	for _, h := range hulls {
		hr := h.Range
		if is_merged(hr){
			continue
		}
		hdict[h.Range.AsArray()] = h

		//if hr.Size() < 4{
		if hr.Size() == 1 {
			// sort hulls for consistency
			hs := sort_hulls(
				as_hullnodes_from_boxes(find_context_hulls(hulldb, h, EpsilonDist)),
			)

			for _, s := range hs {
				sr := s.Range
				if is_merged(sr){
					continue
				}
				//test whether sr.i or sr.j is a self inter-vertex -- split point
				//not sr.i != hr.i or sr.j != hr.j without i/j being a inter-vertex
				//tests for contiguous and whether contiguous index is part of vertex set
				//if the location at which they are contiguous is not part of vertex set then
				//its mergeable : mergeable score <= threshold
				bln := (hr.J() == sr.I() && !vertex_set.Contains(sr.I())) ||
					   (hr.I() == sr.J() && !vertex_set.Contains(sr.J()))

				l := sort_ints(append(sr.AsSlice(), hr.AsSlice()...))
				r := rng.NewRange(l[0], l[len(l)-1])
				_, val      := self.Score(self, r)
				mergeable   := bln && self.is_score_relate_valid(val)

				if mergeable {
					//keep track of items merged
					mrgdict[hr.AsArray()] = h
					mrgdict[sr.AsArray()] = s

					// rm sr + hr
					delete(hdict, sr.AsArray())
					delete(hdict, hr.AsArray())

					m := NewHullNode(pln, r, r.Clone())

					// add merge
					hdict[m.Range.AsArray()] = m

					// add to remove list to remove , after merge
					rm = append(rm, s)
					rm = append(rm, h)

					//if present in umerged as fragment remove
					delete(unmerged, hr.AsArray())

					break
				} else {
					unmerged[hr.AsArray()] = h
				}
			}
		}
	}

	keep  = map_to_slice(hdict, keep)

	return keep, rm
}

