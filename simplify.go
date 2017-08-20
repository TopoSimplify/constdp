package constdp

import (
	"simplex/struct/rtree"
	"simplex/constdp/opts"
	"simplex/struct/sset"
)

//homotopic simplification at a given threshold
func (self *ConstDP) Simplify(opts *opts.Opts) *sset.SSet {
	self.Opts = opts
	self.Simple = make([]*HullNode, 0)

	self.Hulls = self.decompose()
	//debug_print_ptset(self.Hulls)

	//debug_print_hulls(self.Hulls)
	// constrain hulls to self intersects
	self.Hulls, _ = self.constrain_to_selfintersects(opts)
	//debug_print_ptset(self.Hulls)

	var bln bool
	var hull *HullNode
	var selections =  make([]*HullNode, 0)

	var hulldb = rtree.NewRTree(8)
	for !self.Hulls.IsEmpty() {
		// assume poped hull to be valid
		bln = true

		// pop hull in queue
		hull = pop_left_hull(self.Hulls)

		// insert hull into hull db
		hulldb.Insert(hull)

		// self intersection constraint
		if bln && self.Opts.AvoidNewSelfIntersects {
			bln = self.constrain_self_intersection(hull, hulldb, &selections)
		}
		if len(selections) > 0 {
			self.deform_hull(hulldb, &selections)
		}

		if !bln {
			continue
		}

		// context_geom geometry constraint
		bln = self.constrain_context_relation(hull,  &selections)
		if len(selections) > 0 {
			self.deform_hull(hulldb, &selections)
		}
	}

	return simple_hulls_as_ptset(self.merge_simple_segments(hulldb), )
}

func (self *ConstDP) merge_simple_segments(hulldb *rtree.RTree) []*HullNode {
	var hull *HullNode
	var hulls = as_deque(sort_hulls(as_hullnodes(hulldb.All())))
	var cache = make(map[[4]int]bool)
	var fragment_size = 1

	//fmt.Println("After constraints:")
	//DebugPrintPtSet(hulls)

	for !hulls.IsEmpty() {
		// from left
		hull = pop_left_hull(hulls)

		if hull.Range.Size() != fragment_size {
			continue
		}

		hulldb.Remove(hull)

		// find context neighbours
		neighbs := as_hullnodes_from_boxes(find_context_hulls(hulldb, hull, EpsilonDist))

		// find context neighbours
		prev, nxt := extract_neighbours(hull, neighbs)

		// find mergeable neihbs contig
		var merge_prev, merge_nxt *HullNode
		if prev != nil {
			key := cache_key(prev, hull)
			if !cache[key] {
				add_to_merge_cache(cache, &key)
				merge_prev = merge_contiguous_fragments_at_threshold(self, prev, hull)
			}
		}

		if nxt != nil {
			key := cache_key(hull, nxt)
			if !cache[key] {
				add_to_merge_cache(cache, &key)
				merge_nxt = merge_contiguous_fragments_at_threshold(self, hull, nxt)
			}
		}

		var merged = false

		//nxt, prev
		if !merged && merge_nxt != nil {
			hulldb.Remove(nxt)
			if self.is_merge_simplx_valid(merge_nxt, hulldb) {
				h := cast_as_hullnode(hulls.First())
				if h == nxt {
					hulls.PopLeft()
				}
				hulldb.Insert(merge_nxt)
				merged = true
			} else {
				merged = false
				hulldb.Insert(nxt)
			}
		}

		if !merged && merge_prev != nil {
			hulldb.Remove(prev)
			if self.is_merge_simplx_valid(merge_prev, hulldb) {
				//prev cannot exist since moving from left --- right
				hulldb.Insert(merge_prev)
				merged = true
			} else {
				merged = false
				hulldb.Insert(prev)
			}
		}

		if !merged {
			hulldb.Insert(hull)
		}

	}

	return sort_hulls(as_hullnodes(hulldb.All()))
}

func (self *ConstDP) is_merge_simplx_valid(hull *HullNode, hulldb *rtree.RTree) bool {
	var bln  = true
	var side_effects = make([]*HullNode, 0)

	if bln && self.Opts.AvoidNewSelfIntersects {
		// self intersection constraint
		bln = self.constrain_self_intersection(hull, hulldb,  &side_effects)
	}

	if len(side_effects) > 0 || !bln {
		return false
	}

	// context geometry constraint
	bln = self.constrain_context_relation(hull,  &side_effects)
	return len(side_effects) == 0 && bln
}

func cache_key(a, b *HullNode) [4]int {
	ij := [4]int{a.Range.I(), a.Range.J(), b.Range.I(), b.Range.J()}
	sort_ints(ij[:])
	return ij
}

func add_to_merge_cache(cache map[[4]int]bool, key *[4]int) {
	cache[*key] = true
}
