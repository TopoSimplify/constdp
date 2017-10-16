package constdp

import (
	"sort"
	"simplex/constdp/ctx"
)


//sort hulls
//func sort_hulls(hulls []*HullNode) []*HullNode {
//	sort.Sort(HullNodes(hulls))
//	return hulls
//}

//reverse sort hulls
//func sort_reverse(hulls []*HullNode) []*HullNode {
//	sort.Sort(sort.Reverse(HullNodes(hulls)))
//	return hulls
//}

//sort context geoms
func sort_context_geoms(ctxgs []*ctx.CtxGeom) []*ctx.CtxGeom {
	sort.Sort(ctx.ContextGeoms(ctxgs))
	return ctxgs
}
