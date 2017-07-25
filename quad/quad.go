package quad

import (
	"strings"
	"simplex/geom"
	"simplex/geom/mbr"
	"simplex/util/math"
	"simplex/struct/rtree"
	"simplex/constdp/ln"
	"simplex/constdp/seg"
	"simplex/constdp/ctx"
)

func DirectionRelate(pln *ln.Polyline, g geom.Geometry) string {
	segdb := rtree.NewRTree(8)
	objs := make([]rtree.BoxObj, 0)
	for _, s := range pln.Segments() {
		ctx := ctx.NewCtxGeom(s, s.I, s.J).AsSelfSegment()
		objs = append(objs, ctx)
	}
	segdb.Load(objs)

	lnbox := pln.BBox()
	gbox := g.BBox()

	extbox := gbox.Clone()
	extbox.ExpandIncludeMBR(lnbox)

	delta := math.MaxF64(extbox.Height(), extbox.Width()) / 2.0
	uppper := [2]float64{
		extbox.MaxX() + delta,
		extbox.MaxY() + delta,
	}
	lower := [2]float64{
		extbox.MinX() - delta,
		extbox.MinY() - delta,
	}

	extbox.ExpandIncludeXY(uppper[0], uppper[1])
	extbox.ExpandIncludeXY(lower [0], lower [1])

	lx, ly, ux, uy := extbox.MinX(), extbox.MinY(), extbox.MaxX(), extbox.MaxY()
	glx, gly, gux, guy := gbox.MinX(), gbox.MinY(), gbox.MaxX(), gbox.MaxY()

	nw := asPolygon(mbr.NewMBR(lx, guy, glx, uy))
	nn := asPolygon(mbr.NewMBR(glx, guy, gux, uy))
	ne := asPolygon(mbr.NewMBR(gux, guy, ux, uy))

	iw := asPolygon(mbr.NewMBR(lx, gly, glx, guy))
	ii := asPolygon(mbr.NewMBR(glx, gly, gux, guy))
	ie := asPolygon(mbr.NewMBR(gux, gly, ux, guy))

	sw := asPolygon(mbr.NewMBR(lx, ly, glx, gly))
	ss := asPolygon(mbr.NewMBR(glx, ly, gux, gly))
	se := asPolygon(mbr.NewMBR(gux, ly, ux, gly))

	quads := make([]string, 0)
	for _, q := range []*geom.Polygon{nw, nn, ne, iw, ii, ie, sw, ss, se} {
		res := segdb.Search(q.BBox())
		if len(res) > 0 {
			if intersects_quad(q, res) {
				quads = append(quads, "T")
			} else {
				quads = append(quads, "F")
			}
		} else {
			quads = append(quads, "F")
		}
	}
	return strings.Join(quads, "")
}

func asPolygon(box *mbr.MBR) *geom.Polygon {
	array := box.AsPolyArray()
	coords := make([]*geom.Point, 0)
	for _, arr := range array {
		coords = append(coords, geom.NewPoint(arr[:]))
	}
	return geom.NewPolygon(coords)
}

func intersects_quad(q geom.Geometry, res []*rtree.Node) bool {
	bln := false
	for _, node := range res {
		c := node.GetItem().(*ctx.CtxGeom)
		s := c.Geom.(*seg.Seg)
		if q.Intersects(s.Segment) {
			bln = true
			break
		}
	}
	return bln
}
