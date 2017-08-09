package constdp

import (
	"sort"
	"simplex/geom"
	"simplex/geom/mbr"
	"simplex/constdp/db"
	"simplex/struct/sset"
	"simplex/constdp/ctx"
	"simplex/constdp/cmp"
	"simplex/struct/rtree"
	"simplex/constdp/box"
	"simplex/constdp/opts"
	"simplex/struct/deque"
)

func (self *ConstDP) constrain_to_selfintersects(opts *opts.Opts) (*deque.Deque, bool) {
	if !opts.KeepSelfIntersects {
		return self.Hulls, true
	}

	hulldb := rtree.NewRTree(8)

	// TODO: use self intersects at constructor
	data := make([]rtree.BoxObj, 0)
	self_inters := LinearSelfIntersection(self.Pln)
	for _, v := range *self.Hulls.DataView() {
		data = append(data, v.(*HullNode))
	}
	hulldb.Load(data)
	at_vertex_set := sset.NewSSet(cmp.IntCmp)

	// TODO: use inters_vertex_set from constructor
	for _, inter := range self_inters {
		if inter.IsSelfVertex() {
			at_vertex_set.Union(inter.Meta.SelfVertices)
		}
	}

	for _, inter := range self_inters {
		if !inter.IsSelfVertex() {
			continue
		}

		hulls := db.KNN(hulldb, inter, 1.e-5, func(_, item rtree.BoxObj) float64 {
			var other geom.Geometry
			if o, ok := item.(*mbr.MBR); ok {
				other = box.MBRToPolygon(o)
			} else {
				other = item.(*HullNode).Geom
			}
			return inter.Geom.Distance(other)
		})

		for _, h := range hulls {
			hull := h.(*HullNode)
			idxs := inter.Meta.SelfVertices.Values()
			indices := make([]int, 0)
			for _, o := range idxs {
				indices = append(indices, o.(int))
			}
			hsubs := splitHullAtIndex(self, hull, indices)

			if len(hsubs) > 0 {
				hulldb.Remove(hull)
			}

			keep, rm := find_mergeable_contiguous_fragments(
				self, hsubs, hulldb, at_vertex_set,
			)

			for _, h := range rm {
				hulldb.Remove(h)
			}

			for _, h := range keep {
				hulldb.Insert(h)
			}
		}
	}

	hdata := make([]*HullNode, 0)
	for _, h := range hulldb.All() {
		hdata = append(hdata, h.GetItem().(*HullNode))
	}
	hdata = sort_hulls(hdata)

	hulls := deque.NewDeque()
	for _, h := range hdata {
		hulls.Append(h)
	}
	return hulls, true
}

func (self *ConstDP) constrain_self_intersection(hull *HullNode, hulldb *rtree.RTree, bln bool) ([]*HullNode, bool) {
	var selections = make([]*HullNode, 0)
	if bln && self.Opts.AvoidNewSelfIntersects {
		// find hull neighbours
		selections = self.select_deformation_candidates(hulldb, hull)
		for _, h := range selections {
			bln = !(h == hull)
			selections = append(selections, h)
		}
	}
	return selections, bln
}

func (self *ConstDP) constrain_context_relation(hull *HullNode, hulldb *rtree.RTree, bln bool) ([]*HullNode, bool) {
	var hlist = []*HullNode{hull}
	if !bln {
		return hlist, bln
	}

	// find context neighbours - if valid
	ctxs := db.KNN(self.CtxDB, hull, self.Opts.MinDist, func(_, item rtree.BoxObj) float64 {
		var other geom.Geometry
		if o, ok := item.(*mbr.MBR); ok {
			other = box.MBRToPolygon(o)
		} else {
			other = item.(*ctx.CtxGeom).Geom
		}
		return hull.Geom.Distance(other)
	})

	for _, contxt := range ctxs {
		if !bln {
			break
		}

		c := contxt.(*ctx.CtxGeom)
		if bln && self.Opts.GeomRelation {
			bln = self.is_geom_relate_valid(hull, c)
		}

		if bln && self.Opts.DistRelation {
			bln = self.is_dist_relate_valid(hull, c)
		}

		if bln && self.Opts.DirRelation {
			bln = self.is_dir_relate_valid(hull, c)
		}
	}

	if !bln {
		n := len(hlist) - 1
		hlist[n] = nil
		hlist = hlist[:n]
	}

	return hlist, bln
}
