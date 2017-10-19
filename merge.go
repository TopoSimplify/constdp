package constdp

import (
	"github.com/intdxdt/sset"
	"simplex/rng"
	"github.com/intdxdt/rtree"
)

//Merge two ranges
func merge_range(ra, rb *rng.Range) *rng.Range {
	var ranges = sort_ints(append(ra.AsSlice(), rb.AsSlice()...))
	// i...[ra]...k...[rb]...j
	return rng.NewRange(ranges[0], ranges[len(ranges)-1])
}

//Merge contiguous fragments based combined score
func merge_contiguous_fragments_at_threshold(self *ConstDP, ha, hb *HullNode) *HullNode {
	_, val := self.Score(self, merge_range(ha.Range, hb.Range))
	if self.is_score_relate_valid(val) {
		return self.merge_contiguous_fragments(ha, hb)
	}
	return nil
}

//Merge contiguous hulls
func (self *ConstDP) merge_contiguous_fragments(ha, hb *HullNode) *HullNode {
	var r = merge_range(ha.Range, hb.Range)
	// i...[ha]...k...[hb]...j
	return NewHullNode(self.Pln, r)
}

//Merge contiguous hulls by fragment size
func (self *ConstDP) merge_contiguous_fragments_by_size(
	hulls []*HullNode, hulldb *rtree.RTree, vertex_set *sset.SSet,
	unmerged map[[2]int]*HullNode, fragment_size int,
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

		if hr.Size() != fragment_size {
			continue
		}

		// sort hulls for consistency
		var hs = NewHullNodesFromBoxes(find_context_hulls(hulldb, h, EpsilonDist)).Sort()

		for _, s := range hs.list {
			sr := s.Range
			if is_merged(sr){
				continue
			}

			//merged range
			r := merge_range(sr, hr)

			//test whether sr.i or sr.j is a self inter-vertex -- split point
			//not sr.i != hr.i or sr.j != hr.j without i/j being a inter-vertex
			//tests for contiguous and whether contiguous index is part of vertex set
			//if the location at which they are contiguous is not part of vertex set then
			//its mergeable : mergeable score <= threshold
			mergeable := (hr.J() == sr.I() && !vertex_set.Contains(sr.I())) ||
				         (hr.I() == sr.J() && !vertex_set.Contains(sr.J()))

			if mergeable {
				_, val      := self.Score(self, r)
				mergeable   = self.is_score_relate_valid(val)
			}

			if !mergeable {
				unmerged[hr.AsArray()] = h
				continue
			}

			//keep track of items merged
			mrgdict[hr.AsArray()] = h
			mrgdict[sr.AsArray()] = s

			// rm sr + hr
			delete(hdict, sr.AsArray())
			delete(hdict, hr.AsArray())

			// add merge
			hdict[r.AsArray()] = NewHullNode(pln, r)

			// add to remove list to remove , after merge
			rm = append(rm, s)
			rm = append(rm, h)

			//if present in umerged as fragment remove
			delete(unmerged, hr.AsArray())
			break
		}
	}

	for _, o := range hdict {
		keep = append(keep, o)
	}
	return keep, rm
}

