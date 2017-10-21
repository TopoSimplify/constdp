package constdp

import (
	"simplex/ctx"
	"simplex/node"
	"simplex/opts"
	"github.com/intdxdt/cmp"
	"github.com/intdxdt/sset"
	"github.com/intdxdt/rtree"
	"github.com/intdxdt/deque"
	"simplex/relate"
	"simplex/split"
	"simplex/knn"
)

//constrain hulls at self intersection fragments - planar self-intersection
func (self *ConstDP) _const_at_self_intersect_fragments(hulldb *rtree.RTree,
	self_inters []*ctx.CtxGeom, at_vertex_set *sset.SSet) map[[2]int]*node.Node {
	//@formatter:off
	var fragment_size = 1
	var hsubs []*node.Node
	var hulls *node.Nodes
	var idxs []int
	var unmerged = make(map[[2]int]*node.Node, 0)

	for _, inter := range self_inters {
		if !inter.IsSelfVertex() {
			continue
		}

		hulls = nodesFromBoxes(knn.FindNeighbours(hulldb, inter, EpsilonDist)).Sort()

		idxs = as_ints(inter.Meta.SelfVertices.Values())
		for _, hull := range hulls.DataView() {
			hsubs = split.AtIndex(self, hull, idxs, hullGeom)

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
		data = append(data, v.(*node.Node))
	}
	hulldb.Load(data)

	at_vertex_set = sset.NewSSet(cmp.Int)
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
		cg.Meta.SelfVertices = sset.NewSSet(cmp.Int, 4).Add(i)
		cg.Meta.SelfNonVertices = sset.NewSSet(cmp.Int, 4)
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
	return nodesFromRtreeNodes(hulldb.All()).Sort().AsDeque(), true, at_vertex_set
}

//Constrain for self-intersection as a result of simplification
//returns boolean : is hull collapsible
func (self *ConstDP) constrain_ftclass_intersection(hull *node.Node, hulldb *rtree.RTree, selections *node.Nodes) bool {
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
func (self *ConstDP) constrain_self_intersection(hull *node.Node, hulldb *rtree.RTree, selections *node.Nodes) bool {
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
func (self *ConstDP) constrain_context_relation(hull *node.Node, selections *node.Nodes) bool {
	var bln = true

	// find context neighbours - if valid
	ctxs := knn.FindNeighbours(self.ContextDB, hull, self.Opts.MinDist)
	for _, contxt := range ctxs {
		if !bln {
			break
		}

		cg := castAsContextGeom(contxt)
		if bln && self.Opts.GeomRelation {
			bln = relate.IsGeomRelateValid(self, hull, cg)
		}

		if bln && self.Opts.DistRelation {
			bln = relate.IsDistRelateValid(self, hull, cg)
		}

		if bln && self.Opts.DirRelation {
			bln = relate.IsDirRelateValid(self, hull, cg)
		}
	}

	if !bln {
		selections.Push(hull)
	}

	return bln
}
