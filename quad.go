package constdp

import (
	"strings"
	"github.com/intdxdt/geom"
	"github.com/intdxdt/mbr"
	"github.com/intdxdt/math"
	"github.com/intdxdt/rtree"
	"simplex/pln"
	"simplex/seg"
	"simplex/ctx"
	"simplex/box"
)

//Direction Relate
func DirectionRelate(polyline *pln.Polyline, g geom.Geometry) string {
	var segdb = rtree.NewRTree(8)
	var objs = make([]rtree.BoxObj, 0)
	for _, s := range polyline.Segments() {
		objs = append(objs, ctx.NewCtxGeom(s, s.I, s.J).AsSelfSegment())
	}
	segdb.Load(objs)

	var lnbox   = polyline.BBox()
	var gbox    = g.BBox()
	var extbox  = gbox.Clone()
	extbox.ExpandIncludeMBR(lnbox)

	var delta  = math.MaxF64(extbox.Height(), extbox.Width()) / 2.0
	var uppper = [2]float64 {
		extbox.MaxX() + delta,
		extbox.MaxY() + delta,
	}
	var lower = [2]float64 {
		extbox.MinX() - delta,
		extbox.MinY() - delta,
	}

	extbox.ExpandIncludeXY(uppper[0], uppper[1])
	extbox.ExpandIncludeXY(lower [0], lower [1])

	lx, ly, ux, uy      := extbox.MinX(), extbox.MinY(), extbox.MaxX(), extbox.MaxY()
	glx, gly, gux, guy  := gbox.MinX(),   gbox.MinY(),   gbox.MaxX(),   gbox.MaxY()

	nw := box.MBRToPolygon(mbr.NewMBR(lx,  guy, glx, uy))
	nn := box.MBRToPolygon(mbr.NewMBR(glx, guy, gux, uy))
	ne := box.MBRToPolygon(mbr.NewMBR(gux, guy, ux,  uy))

	iw := box.MBRToPolygon(mbr.NewMBR(lx,  gly, glx, guy))
	ii := box.MBRToPolygon(mbr.NewMBR(glx, gly, gux, guy))
	ie := box.MBRToPolygon(mbr.NewMBR(gux, gly, ux,  guy))

	sw := box.MBRToPolygon(mbr.NewMBR(lx, ly, glx, gly))
	ss := box.MBRToPolygon(mbr.NewMBR(glx, ly, gux, gly))
	se := box.MBRToPolygon(mbr.NewMBR(gux, ly, ux, gly))

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

//find if intersects segment
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
