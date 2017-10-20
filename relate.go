package constdp

import (
	"github.com/intdxdt/geom"
	"simplex/pln"
	"simplex/ctx"
	"simplex/node"
)

//checks if score is valid at threshold of constrained dp
func (self *ConstDP) is_score_relate_valid(val float64) bool {
	return val <= self.Opts.Threshold
}

//geometry relate
func (self *ConstDP) is_geom_relate_valid(hull *node.Node, ctx *ctx.CtxGeom) bool {
	var seg    = hull_segment(self, hull)
	var subpln = self.Pln.SubPolyline(hull.Range)

	var ln_geom  = subpln.Geometry
	var seg_geom = seg
	var ctx_geom = ctx.Geom

	var ln_g_inter  = ln_geom.Intersects(ctx_geom)
	var seg_g_inter = seg_geom.Intersects(ctx_geom)

	var bln = true
	if (seg_g_inter && !ln_g_inter)  || (!seg_g_inter && ln_g_inter){
		bln = false
	}
	// both intersects & disjoint
	return bln
}

//distance relate
func (self *ConstDP) is_dist_relate_valid(hull *node.Node, ctx *ctx.CtxGeom) bool {
	var mindist = self.Opts.MinDist
	var seg     = hull_segment(self, hull)
	var ln_geom = hull.SubPolyline().Geometry

	var seg_geom = seg
	var ctx_geom = ctx.Geom

	var _or = ln_geom.Distance(ctx_geom) // original relate
	var dr  = seg_geom.Distance(ctx_geom) // new relate

	bln := dr >= mindist
	if (!bln) && _or < mindist {//if not bln and _or <= mindist:
		//if original violates constraint, then simple can
		// >= than original or <= original, either way should be true
		// [original & simple] <= mindist, then simple cannot be  simple >= mindist no matter
		// how many vertices introduced
		bln = true
	}
	return bln
}

//direction relate
func (self *ConstDP) is_dir_relate_valid(hull *node.Node, ctx *ctx.CtxGeom) bool {
	subpln  := self.Pln.SubPolyline(hull.Range)
	segment := pln.New([]*geom.Point{
		self.Pln.Coordinates[hull.Range.I()],
		self.Pln.Coordinates[hull.Range.J()],
	})

	lnr  := DirectionRelate(subpln,  ctx.Geom)
	segr := DirectionRelate(segment, ctx.Geom)

	return lnr == segr
}
