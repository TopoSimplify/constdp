package constdp

import (
	"simplex/geom/mbr"
	"simplex/constdp/box"
	"simplex/constdp/ctx"
	"simplex/constdp/db"
	"simplex/struct/rtree"
	"simplex/constdp/igeom"
)

func scorer(query igeom.IGeom) func(_, item rtree.BoxObj) float64 {
	return func(_, item rtree.BoxObj) float64 {
		var ok bool
		var mb *mbr.MBR
		var other igeom.IGeom
		//item is box from rtree
		if mb, ok = item.(*mbr.MBR); ok {
			other = box.MBRToPolygon(mb)
		} else { //item is either ctxgeom or hullnode
			if other, ok = item.(*ctx.CtxGeom); !ok {
				other = item.(*HullNode)
			}
		}
		return query.Geometry().Distance(other.Geometry())
	}
}

func find_context_neighbs(database *rtree.RTree, query igeom.IGeom, dist float64) []rtree.BoxObj {
	return db.KNN(database, query.Geometry(), dist, scorer(query))
}

func find_context_hulls(hulldb *rtree.RTree, hull *HullNode, dist float64) []rtree.BoxObj {
	predicate   := hull_predicate(hull, dist)
	return db.KNN(hulldb, hull.Geometry(), dist, scorer(hull), predicate)
}