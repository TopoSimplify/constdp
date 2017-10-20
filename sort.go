package constdp

import (
	"sort"
	"simplex/ctx"
)


//sort context geoms
func sort_context_geoms(ctxgs []*ctx.CtxGeom) []*ctx.CtxGeom {
	sort.Sort(ctx.ContextGeoms(ctxgs))
	return ctxgs
}
