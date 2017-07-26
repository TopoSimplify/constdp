package constdp

import (
	"sort"
	"simplex/geom"
	"simplex/constdp/hl"
	"simplex/constdp/ln"
	"simplex/constdp/db"
	"simplex/constdp/ctx"
	"simplex/constdp/cmp"
	"simplex/struct/sset"
	"simplex/struct/deque"
	"simplex/struct/rtree"
	"simplex/constdp/opts"
	"simplex/constdp/quad"
	"simplex/geom/mbr"
	"simplex/constdp/box"
)

func (self *ConstDP) split_hulls_at_selfintersects(dphulls *deque.Deque) *deque.Deque {
	hulldb := rtree.NewRTree(8)
	self_inters := LinearSelfIntersection(self.Pln)
	data := make([]rtree.BoxObj, 0)
	for _, v := range *dphulls.DataView() {
		data = append(data, v.(*hl.HullNode))
	}
	hulldb.Load(data)
	at_vertex_set := sset.NewSSet(cmp.IntCmp)

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
				other = item.(*hl.HullNode).Geom
			}
			return inter.Geom.Distance(other)
		})

		for _, h := range hulls {
			hull := h.(*hl.HullNode)
			idxs := inter.Meta.SelfVertices.Values()
			indices := make([]int, 0)
			for _, o := range idxs {
				indices = append(indices, o.(int))
			}
			hsubs := hl.SplitHullAtIndex(self, hull, indices)

			if len(hsubs) > 0 {
				hulldb.Remove(hull)
			}

			keep, rm := hl.MergeContigFragments(
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

	hdata := make([]*hl.HullNode, 0)
	for _, h := range hulldb.All() {
		hdata = append(hdata, h.GetItem().(*hl.HullNode))
	}
	sort.Sort(hl.HullNodes(hdata))
	hulls := deque.NewDeque()
	for _, hn := range hdata {
		hulls.Append(hn)
	}
	return hulls
}

//homotopic simplification at a given threshold
func (self *ConstDP) Simplify(opts *opts.Opts) *ConstDP {
	self.Simple = make([]*hl.HullNode, 0)
	self.Hulls = self.decompose(opts.Threshold)

	// split hulls by self intersects
	if opts.KeepSelfIntersects {
		self.Hulls = self.split_hulls_at_selfintersects(self.Hulls)
	}

	//for _, h := range *self.Hulls.DataView() {
	//	fmt.Println(h)
	//}

	hulldb := rtree.NewRTree(8)
	for self.Hulls.Len() > 0 {
		// assume poped hull to be valid
		bln := true

		// pop hull in queue
		hull := self.Hulls.PopLeft().(*hl.HullNode)

		// insert hull into hull db
		hulldb.Insert(hull)

		if bln && self.Opts.AvoidNewSelfIntersects {
			// find hull neighbours
			hlist := hl.FindHullDeformationList(hulldb, hull, self.Opts)
			for _, h := range hlist {
				bln = !(h == hull)
				self.deform_hull(hulldb, h)
			}
		}

		if !bln {
			continue
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

		i := 0
		for bln && i < len(ctxs) {
			c := ctxs[i].(*ctx.CtxGeom)
			if bln && self.Opts.GeomRelation {
				bln = self.is_geom_relate_valid(hull, c)
			}

			if bln && self.Opts.DistRelation {
				bln = self.is_dist_relate_valid(hull, c)
			}

			if bln && self.Opts.DirRelation {
				bln = self.is_dir_relate_valid(hull, c)
			}

			i += 1
		}

		if !bln {
			self.deform_hull(hulldb, hull)
		}
	}
	hdata := make([]*hl.HullNode, 0)
	for _, h := range hulldb.All() {
		hdata = append(hdata, h.GetItem().(*hl.HullNode))
	}
	sort.Sort(hl.HullNodes(hdata))
	self.Simple = hdata
	return self
}

func (self *ConstDP) deform_hull(hulldb *rtree.RTree, hull *hl.HullNode) {
	// split hull at maximum_offset offset
	ha, hb := hl.SplitHull(self, hull)
	hulldb.Remove(hull)

	self.Hulls.AppendLeft(hb)
	self.Hulls.AppendLeft(ha)
}

func (self *ConstDP) is_geom_relate_valid(hull *hl.HullNode, ctx *ctx.CtxGeom) bool {
	seg := hl.HullSegment(self, hull)
	subpln := self.Pln.SubPolyline(hull.Range)

	ln_geom := subpln.Geom
	seg_geom := seg
	ctx_geom := ctx.Geom

	ln_g_inter := ln_geom.Intersects(ctx_geom)
	seg_g_inter := seg_geom.Intersects(ctx_geom)

	bln := true
	if seg_g_inter && (! ln_g_inter) {
		bln = false
	} else if (! seg_g_inter) && ln_g_inter {
		bln = false
	}
	// both intersects & disjoint
	return bln
}

//is distance relate valid ?
func (self *ConstDP) is_dist_relate_valid(hull *hl.HullNode, ctx *ctx.CtxGeom) bool {
	mindist := self.Opts.MinDist
	seg := hl.HullSegment(self, hull)
	ln_geom := hull.Pln.Geom

	seg_geom := seg
	ctx_geom := ctx.Geom

	_or := ln_geom.Distance(ctx_geom) // original relate
	dr := seg_geom.Distance(ctx_geom) // new relate

	bln := dr >= mindist
	if !bln && _or < mindist {
		bln = (dr >= _or)
	}
	return bln
}

func (self *ConstDP) is_dir_relate_valid(hull *hl.HullNode, ctx *ctx.CtxGeom) bool {
	subpln := self.Pln.SubPolyline(hull.Range)
	segment := ln.NewPolyline([]*geom.Point{
		self.Pln.Coords[hull.Range.I()],
		self.Pln.Coords[hull.Range.J()],
	})

	lnr := quad.DirectionRelate(subpln, ctx.Geom)
	segr := quad.DirectionRelate(segment, ctx.Geom)

	return lnr == segr
}
