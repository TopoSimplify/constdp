package constdp

import (
	"simplex/knn"
	"simplex/opts"
	"simplex/node"
	"simplex/merge"
	"github.com/intdxdt/sset"
	"github.com/intdxdt/rtree"
	"simplex/dp"
)

//checks if score is valid at threshold of constrained dp
func (self *ConstDP) is_score_relate_valid(val float64) bool {
	return val <= self.Opts.Threshold
}


//Homotopic simplification at a given threshold
func (self *ConstDP) Simplify(opts *opts.Opts, const_vertices ...[]int) *ConstDP {
	var const_vertex_set *sset.SSet
	var const_verts = []int{}

	if len(const_vertices) > 0 {
		const_verts = const_vertices[0]
	}

	self.SimpleSet.Empty()
	self.Opts  = opts
	self.Hulls = self.Decompose()

	//debug_print_ptset(self.Hulls)

	//debug_print_hulls(self.Hulls)
	// constrain hulls to self intersects
	self.Hulls, _, const_vertex_set = self.constrain_to_selfintersects(opts, const_verts)
	//debug_print_ptset(self.Hulls)

	var bln bool
	var hull *node.Node
	var selections  = node.NewNodes()

	var hulldb = rtree.NewRTree(8)
	for !self.Hulls.IsEmpty() {
		// assume popped hull to be valid
		bln = true

		// pop hull in queue
		hull = popLeftHull(self.Hulls)

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
	self.SimpleSet.Empty()
	for _, h := range nodesFromRtreeNodes(hulldb.All()).Sort().DataView() {
		self.Hulls.Append(h)
		self.SimpleSet.Extend(h.Range.I(), h.Range.J())
	}
	return self
}

//Merge segment fragments where possible
func (self *ConstDP) merge_simple_segments(hulldb *rtree.RTree, const_vertex_set *sset.SSet) {
	var fragment_size = 1
	var hull *node.Node
	var neighbs *node.Nodes
	var cache = make(map[[4]int]bool)
	var hulls = nodesFromRtreeNodes(hulldb.All()).Sort().AsDeque()

	//fmt.Println("After constraints:")
	//DebugPrintPtSet(hulls)

	for !hulls.IsEmpty() {
		// from left
		hull = popLeftHull(hulls)

		if hull.Range.Size() != fragment_size {
			continue
		}

		//make sure hull index is not part of vertex with degree > 2
		if const_vertex_set.Contains(hull.Range.I()) || const_vertex_set.Contains(hull.Range.J()) {
			continue
		}

		hulldb.Remove(hull)

		// find context neighbours
		neighbs = nodesFromBoxes(knn.FindNodeNeighbours(hulldb, hull, EpsilonDist))

		// find context neighbours
		prev, nxt := node.Neighbours(hull, neighbs)

		// find mergeable neihbs contig
		var key [4]int
		var merge_prev, merge_nxt *node.Node

		if prev != nil {
			key = cache_key(prev, hull)
			if !cache[key] {
				add_to_merge_cache(cache, &key)
				merge_prev = merge.ContiguousFragmentsAtThreshold(self, prev, hull,self.is_score_relate_valid, dp.NodeGeometry)
			}
		}

		if nxt != nil {
			key = cache_key(hull, nxt)
			if !cache[key] {
				add_to_merge_cache(cache, &key)
				merge_nxt = merge.ContiguousFragmentsAtThreshold(self, hull, nxt, self.is_score_relate_valid,  dp.NodeGeometry)
			}
		}

		var merged bool
		//nxt, prev
		if !merged && merge_nxt != nil {
			hulldb.Remove(nxt)
			if self.is_merge_simplx_valid(merge_nxt, hulldb) {
				var h = castAsNode(hulls.First())
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

func (self *ConstDP) is_merge_simplx_valid(hull *node.Node, hulldb *rtree.RTree) bool {
	var bln = true
	var side_effects = node.NewNodes()

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

func cache_key(a, b *node.Node) [4]int {
	ij := [4]int{a.Range.I(), a.Range.J(), b.Range.I(), b.Range.J()}
	sort_ints(ij[:])
	return ij
}

func add_to_merge_cache(cache map[[4]int]bool, key *[4]int) {
	cache[*key] = true
}
