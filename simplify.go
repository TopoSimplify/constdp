package constdp

import (
	"simplex/struct/rtree"
	"simplex/constdp/opts"
	"simplex/struct/sset"
)

//Homotopic simplification at a given threshold
func (self *ConstDP) Simplify(opts *opts.Opts, const_vertices ...[]int) *ConstDP {
	var const_vertex_set *sset.SSet
	var const_verts = make([]int, 0)
	if len(const_vertices) > 0 {
		const_verts = const_vertices[0]
	}

	self.Simple.Empty()
	self.Opts = opts
	self.Hulls = self.decompose()

	//debug_print_ptset(self.Hulls)

	//debug_print_hulls(self.Hulls)
	// constrain hulls to self intersects
	self.Hulls, _, const_vertex_set = self.constrain_to_selfintersects(opts, const_verts)
	//debug_print_ptset(self.Hulls)

	var bln bool
	var hull *HullNode
	var selections  = NewHullNodes()

	var hulldb = rtree.NewRTree(8)
	for !self.Hulls.IsEmpty() {
		// assume popped hull to be valid
		bln = true

		// pop hull in queue
		hull = pop_left_hull(self.Hulls)

		// insert hull into hull db
		hulldb.Insert(hull)

		// self intersection constraint
		if bln && self.Opts.AvoidNewSelfIntersects {
			bln = self.constrain_self_intersection(hull, hulldb, selections)
		}

		if !selections.IsEmpty() {
			self.deform_hulls(hulldb, selections)
		}

		if !bln {
			continue
		}

		// context_geom geometry constraint
		bln = self.constrain_context_relation(hull, selections)
		if !selections.IsEmpty() {
			self.deform_hulls(hulldb, selections)
		}
	}

	self.merge_simple_segments(hulldb, const_vertex_set)

	self.Hulls.Clear()
	self.Simple.Empty()
	for _, h := range NewHullNodesFromNodes(hulldb.All()).Sort().list {
		self.Hulls.Append(h)
		self.Simple.Extend(h.Range.I(), h.Range.J())
	}
	return self
}

//Merge segment fragments where possible
func (self *ConstDP) merge_simple_segments(hulldb *rtree.RTree, const_vertex_set *sset.SSet) {
	var fragment_size = 1
	var hull *HullNode
	var neighbs *HullNodes
	var cache = make(map[[4]int]bool)
	var hulls = NewHullNodesFromNodes(hulldb.All()).Sort().AsDeque()

	//fmt.Println("After constraints:")
	//DebugPrintPtSet(hulls)

	for !hulls.IsEmpty() {
		// from left
		hull = pop_left_hull(hulls)

		if hull.Range.Size() != fragment_size {
			continue
		}

		//make sure hull index is not part of vertex with degree > 2
		if const_vertex_set.Contains(hull.Range.I()) || const_vertex_set.Contains(hull.Range.J()) {
			continue
		}

		hulldb.Remove(hull)

		// find context neighbours
		neighbs = NewHullNodesFromBoxes(find_context_hulls(hulldb, hull, EpsilonDist))

		// find context neighbours
		prev, nxt := extract_neighbours(hull, neighbs)

		// find mergeable neihbs contig
		var key [4]int
		var merge_prev, merge_nxt *HullNode

		if prev != nil {
			key = cache_key(prev, hull)
			if !cache[key] {
				add_to_merge_cache(cache, &key)
				merge_prev = merge_contiguous_fragments_at_threshold(self, prev, hull)
			}
		}

		if nxt != nil {
			key = cache_key(hull, nxt)
			if !cache[key] {
				add_to_merge_cache(cache, &key)
				merge_nxt = merge_contiguous_fragments_at_threshold(self, hull, nxt)
			}
		}

		var merged bool
		//nxt, prev
		if !merged && merge_nxt != nil {
			hulldb.Remove(nxt)
			if self.is_merge_simplx_valid(merge_nxt, hulldb) {
				var h = cast_as_hullnode(hulls.First())
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
			//prev cannot exist since moving from left --- right
			if self.is_merge_simplx_valid(merge_prev, hulldb) {
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
}

func (self *ConstDP) is_merge_simplx_valid(hull *HullNode, hulldb *rtree.RTree) bool {
	var bln = true
	var side_effects = NewHullNodes()

	if bln && self.Opts.AvoidNewSelfIntersects {
		// self intersection constraint
		bln = self.constrain_self_intersection(hull, hulldb, side_effects)
	}

	if !side_effects.IsEmpty() || !bln {
		return false
	}

	// context geometry constraint
	bln = self.constrain_context_relation(hull, side_effects)
	return side_effects.IsEmpty() && bln
}

func cache_key(a, b *HullNode) [4]int {
	ij := [4]int{a.Range.I(), a.Range.J(), b.Range.I(), b.Range.J()}
	sort_ints(ij[:])
	return ij
}

func add_to_merge_cache(cache map[[4]int]bool, key *[4]int) {
	cache[*key] = true
}
