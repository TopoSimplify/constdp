package constdp

import (
	"simplex/geom"
	"simplex/constdp/ctx"
	"simplex/constdp/ln"
	"simplex/constdp/quad"
)

func (self *ConstDP) is_geom_relate_valid(hull *HullNode, ctx *ctx.CtxGeom) bool {
	seg := hullSegment(self, hull)
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
func (self *ConstDP) is_dist_relate_valid(hull *HullNode, ctx *ctx.CtxGeom) bool {
	mindist := self.Opts.MinDist
	seg := hullSegment(self, hull)
	ln_geom := hull.Pln.Geom

	seg_geom := seg
	ctx_geom := ctx.Geom

	_or := ln_geom.Distance(ctx_geom) // original relate
	dr := seg_geom.Distance(ctx_geom) // new relate

	bln := dr >= mindist
	if !bln && _or < mindist {
        //if original violates constraint, then simple can
        // >= than original or <= original, either way should be true
        // [original & simple] <= mindist, then simple cannot be  simple >= mindist no matter
        // how many vertices introduced
		bln = true
	}
	return bln
}

func (self *ConstDP) is_dir_relate_valid(hull *HullNode, ctx *ctx.CtxGeom) bool {
	subpln := self.Pln.SubPolyline(hull.Range)
	segment := ln.NewPolyline([]*geom.Point{
		self.Pln.Coords[hull.Range.I()],
		self.Pln.Coords[hull.Range.J()],
	})

	lnr := quad.DirectionRelate(subpln, ctx.Geom)
	segr := quad.DirectionRelate(segment, ctx.Geom)

	return lnr == segr
}

