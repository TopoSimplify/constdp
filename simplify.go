package constdp

import (
	"simplex/struct/rtree"
	"simplex/constdp/opts"
)

//homotopic simplification at a given threshold
func (self *ConstDP) Simplify(opts *opts.Opts) *ConstDP {
	self.Opts = opts
	self.Simple = make([]*HullNode, 0)
	self.Hulls = self.decompose(opts.Threshold)

	// constrain hulls to self intersects
	self.Hulls, _ = self.constrain_to_selfintersects(opts)

	// for _, h := range *self.Hulls.DataView() {
	// 	fmt.Println(h)
	// }

	var bln bool
	var hull *HullNode
	var hlist []*HullNode

	var hulldb = rtree.NewRTree(8)
	for !self.Hulls.IsEmpty() {
		// assume poped hull to be valid
		bln = true

		// pop hull in queue
		hull = self.Hulls.PopLeft().(*HullNode)

		// insert hull into hull db
		hulldb.Insert(hull)

		// self intersection constraint
		hlist, bln = self.constrain_self_intersection(hull, hulldb, bln)
		self.deform_hull(hulldb, hlist)

		if !bln {
			continue
		}

		// context_geom geometry constraint
		hlist, bln = self.constrain_context_relation(hull, bln)
		self.deform_hull(hulldb, hlist)

	}

	self.Simple = self.merge_simple_segments(hulldb)
	return self
}

func (self *ConstDP) merge_simple_segments(hulldb *rtree.RTree) []*HullNode {
	hulls := as_deque(sort_hulls(as_hullnodes(hulldb.All())))

	for !hulls.IsEmpty() {
		// from left
		hull := hulls.PopLeft().(*HullNode)

		if hull.Range.Size() == 1 {
			hulldb.Remove(hull)
			// find context neighbours
			neighbs := as_hullnodes_from_boxes(
				find_context_hulls(hulldb, hull, EpsilonDist),
			)

			// find context neighbours
			prev, nxt := extract_neighbours(hull, neighbs)

			// find mergeable neihbs contig
			var merge_prev, merge_nxt *HullNode
			if prev != nil  {
				merge_prev = contiguous_fragments_at_threshold(self, prev, hull)
			}

			if nxt != nil {
				merge_nxt = contiguous_fragments_at_threshold(self, hull, nxt)
			}

			// nxt, prev
			if (merge_nxt != nil) && self.is_merge_simplx_valid(merge_nxt, hulldb) {
				h := hulls.Get(0).(*HullNode)
				if h == nxt {
					hulls.PopLeft()
				}
				hulldb.Remove(nxt)
				hulldb.Insert(merge_nxt)
			} else if (merge_prev != nil) && self.is_merge_simplx_valid(merge_prev, hulldb) {
				// prev cannot exist since moving from left --- right
				// if hulls[-1] is prev:
				//     hulls.popleft()
				hulldb.Remove(prev)
				hulldb.Insert(merge_prev)
			} else {
				hulldb.Insert(hull)
			}
		}
	}
	return sort_hulls(as_hullnodes(hulldb.All()))
}


func (self *ConstDP) is_merge_simplx_valid(hull *HullNode, hulldb *rtree.RTree) bool {
	bln := hull != nil

	if !bln {
		return bln
	}

	// self intersection constraint
	hlist, bln := self.constrain_self_intersection(hull, hulldb, bln)
	if len(hlist) > 0 || !bln {
		return false
	}

	// context_geom geometry constraint
	hlist, bln = self.constrain_context_relation(hull, bln)
	return len(hlist) == 0 && bln
}
