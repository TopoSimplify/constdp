package constdp

import (
	"simplex/geom"
	"simplex/struct/rtree"
	"simplex/util/math"
	"simplex/geom/mbr"
	"strings"
)

func DirectionRelate(pln *Polyline, g geom.Geometry) string {
	segdb := rtree.NewRTree(8)
	objs := make([]rtree.BoxObj, 0)
	for _, seg := range pln.Segments() {
		ctx := NewCtxGeom(seg, seg.I, seg.J).AsSelfSegment()
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

	nw := mbr.NewMBR(lx, guy, glx, uy)
	nn := mbr.NewMBR(glx, guy, gux, uy)
	ne := mbr.NewMBR(gux, guy, ux, uy)

	iw := mbr.NewMBR(lx, gly, glx, guy)
	ii := mbr.NewMBR(glx, gly, gux, guy)
	ie := mbr.NewMBR(gux, gly, ux, guy)

	sw := mbr.NewMBR(lx, ly, glx, gly)
	ss := mbr.NewMBR(glx, ly, gux, gly)
	se := mbr.NewMBR(gux, ly, ux, gly)

	quads := make([]string, 0)
	for _, q := range []*mbr.MBR{nw, nn, ne, iw, ii, ie, sw, ss, se} {
		res := segdb.Search(q)
		if len(res) > 0 {
			quads = append(quads, "T")
		} else {
			quads = append(quads, "F")
		}
	}
	return strings.Join(quads, "")
}
