package constdp

import (
	"simplex/struct/sset"
	"simplex/constdp/ctx"
	"simplex/constdp/cmp"
	"simplex/struct/rtree"
	"simplex/constdp/opts"
	"simplex/struct/deque"
)

//constrain hulls at self intersection fragments - planar self-intersection
func (self *ConstDP) _const_at_self_intersect_fragments(hulldb *rtree.RTree,
	self_inters []*ctx.CtxGeom, at_vertex_set *sset.SSet) map[[2]int]*HullNode {
	//@formatter:off
	var fragment_size = 1

	var hsubs []*HullNode
	var hulls *HullNodes
	var idxs []int
	var unmerged = make(map[[2]int]*HullNode, 0)

	for _, inter := range self_inters {
		if !inter.IsSelfVertex() {
			continue
		}

		hulls = NewHullNodesFromBoxes(find_context_neighbs(hulldb, inter, EpsilonDist)).Sort()

		idxs = as_ints(inter.Meta.SelfVertices.Values())
		for _, hull := range hulls.list {
			hsubs = split_at_index(self, hull, idxs)

			if len(hsubs) == 0 && (hull.Range.Size() == fragment_size) {
				hsubs = append(hsubs, hull)
			}

			if len(hsubs) == 0 {
				continue
			}

			hulldb.Remove(hull)
			keep, rm := self.merge_contiguous_fragments_by_size(
				hsubs, hulldb, at_vertex_set, unmerged, fragment_size,
			)

			for _, h := range rm {
				hulldb.Remove(h)
			}

			for _, h := range keep {
				hulldb.Insert(h)
			}
		}
	}

	return unmerged
}

//Constrain for planar self-intersection
func (self *ConstDP) constrain_to_selfintersects(opts *opts.Opts, const_verts []int) (*deque.Deque, bool, *sset.SSet) {
	var at_vertex_set *sset.SSet
	if !opts.KeepSelfIntersects {
		return self.Hulls, true, at_vertex_set
	}

	var hulldb = rtree.NewRTree(8)
	var self_inters = linear_self_intersection(self.Pln)

	var data = make([]rtree.BoxObj, 0)
	for _, v := range *self.Hulls.DataView() {
		data = append(data, v.(*HullNode))
	}
	hulldb.Load(data)

	at_vertex_set = sset.NewSSet(cmp.IntCmp)
	for _, inter := range self_inters {
		if inter.IsSelfVertex() {
			at_vertex_set = at_vertex_set.Union(inter.Meta.SelfVertices)
		}
	}

	//update  const vertices if any
	//add const vertices as self inters
	for _, i := range const_verts {
		if at_vertex_set.Contains(i) { //exclude already self intersects
			continue
		}
		at_vertex_set.Add(i)
		pt := self.Pln.Coordinate(i)
		cg := ctx.NewCtxGeom(pt.Clone(), i, i).AsSelfVertex()
		cg.Meta.SelfVertices = sset.NewSSet(cmp.IntCmp, 4).Add(i)
		cg.Meta.SelfNonVertices = sset.NewSSet(cmp.IntCmp, 4)
		self_inters = append(self_inters, cg)
	}

	//constrain fragments aroud self intersects
	//try to merge fragments from first attempt
	var mcount = 2
	for mcount > 0 {
		fragments := self._const_at_self_intersect_fragments(hulldb, self_inters, at_vertex_set)
		if len(fragments) == 0 {
			break
		}
		mcount += -1
	}
	return NewHullNodesFromNodes(hulldb.All()).Sort().AsDeque(), true, at_vertex_set
}

//Constrain for self-intersection as a result of simplification
//returns boolean : is hull collapsible
func (self *ConstDP) constrain_ftclass_intersection(hull *HullNode, hulldb *rtree.RTree, selections *HullNodes) bool {
	var bln = true
	//find hull neighbours
	var hulls = self.select_ftclass_deformation_candidates(hulldb, hull)
	for _, h := range hulls {
		//if bln & selection contains current hull : bln : false
		if bln && (h == hull) {
			bln = false // cmp ref
		}
		selections.Push(h)
	}
	return bln
}

//Constrain for self-intersection as a result of simplification
//returns boolean : is hull collapsible
func (self *ConstDP) constrain_self_intersection(hull *HullNode, hulldb *rtree.RTree, selections *HullNodes) bool {
	//assume hull is valid and proof otherwise
	var bln = true
	// find hull neighbours
	hulls := self.select_deformation_candidates(hulldb, hull)
	for _, h := range hulls {
		//if bln & selection contains current hull : bln : false
		if bln && (h == hull) {
			bln = false //cmp &
		}
		selections.Push(h)
	}

	return bln
}

//Constrain for context neighbours
// finds the collapsibility of hull with respect to context hull neighbours
// if hull is deformable, its added to selections
func (self *ConstDP) constrain_context_relation(hull *HullNode, selections *HullNodes) bool {
	var bln = true

	// find context neighbours - if valid
	ctxs := find_context_neighbs(self.ContextDB, hull, self.Opts.MinDist)
	for _, contxt := range ctxs {
		if !bln {
			break
		}

		cg := cast_as_context_geom(contxt)
		if bln && self.Opts.GeomRelation {
			bln = self.is_geom_relate_valid(hull, cg)
		}

		if bln && self.Opts.DistRelation {
			bln = self.is_dist_relate_valid(hull, cg)
		}

		if bln && self.Opts.DirRelation {
			bln = self.is_dir_relate_valid(hull, cg)
		}
	}

	if !bln {
		selections.Push(hull)
	}

	return bln
}
