package constdp

import (
	"simplex/struct/rtree"
	"simplex/constdp/opts"
	"strings"
	"fmt"
)

//homotopic simplification at a given threshold
func (self *ConstDP) Simplify(opts *opts.Opts) *ConstDP {
	self.Opts = opts
	self.Simple = make([]*HullNode, 0)
	self.Hulls = self.decompose()

	// constrain hulls to self intersects
	self.Hulls, _ = self.constrain_to_selfintersects(opts)

	//for _, h := range *self.Hulls.DataView() {
	//	fmt.Println(h)
	//}
	//fmt.Println(strings.Repeat("-", 80))

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
		if len(hlist) > 0 {
			self.deform_hull(hulldb, hlist)
		}

		if !bln {
			continue
		}

		// context_geom geometry constraint
		hlist, bln = self.constrain_context_relation(hull, bln)
		if len(hlist) > 0 {
			self.deform_hull(hulldb, hlist)
		}

	}

	self.Simple = self.merge_simple_segments(hulldb)
	return self
}

func (self *ConstDP) merge_simple_segments(hulldb *rtree.RTree) []*HullNode {
	hulls := as_deque(sort_hulls(as_hullnodes(hulldb.All())))
	cache := make(map[[4]int]bool)

	for _, h := range *hulls.DataView() {
		fmt.Println(h)
	}
	fmt.Println(strings.Repeat("-", 80))

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

			merged := false
			// nxt, prev
			if merge_nxt != nil {
				hulldb.Remove(nxt)
				if self.is_merge_simplx_valid(merge_nxt, hulldb) {
					h := hulls.First().(*HullNode)
					if h == nxt {
						hulls.PopLeft()
					}
					hulldb.Insert(merge_nxt)
					merged = true
				} else {
					merged = false
					hulldb.Insert(nxt)
				}
			} else if merge_prev != nil {
				hulldb.Remove(prev)
				if self.is_merge_simplx_valid(merge_prev, hulldb) {
					// prev cannot exist since moving from left --- right
					// if hulls[-1] is prev:
					//     hulls.popleft()
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

	return sort_hulls(as_hullnodes(hulldb.All()))
}

func (self *ConstDP) is_merge_simplx_valid(hull *HullNode, hulldb *rtree.RTree) bool {
	bln := hull != nil

	if !bln {
		return bln
	}

	// self intersection constraint
	side_effects, bln := self.constrain_self_intersection(hull, hulldb, bln)
	if len(side_effects) > 0 || !bln {
		return false
	}

	// context geometry constraint
	side_effects, bln = self.constrain_context_relation(hull, bln)
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
