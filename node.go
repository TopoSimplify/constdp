package constdp

import (
	"simplex/rng"
	"simplex/seg"
	"simplex/pln"
	"github.com/intdxdt/geom"
	"github.com/intdxdt/mbr"
	"github.com/intdxdt/sset"
)

const (
	z = 2
)

//hull node
type HullNode struct {
	Pln   *pln.Polyline
	Range *rng.Range
	Geom  geom.Geometry
	PtSet *sset.SSet
	DP    *ConstDP
}

//New Hull Node
func NewHullNode(polyline *pln.Polyline, rng *rng.Range) *HullNode {
	var coords = make([]*geom.Point, 0)
	for i := range rng.Stride() {
		pt := polyline.Coordinates[i].Clone()
		pt[z] = float64(i)
		coords = append(coords, pt)
	}

	convex_hull := geom.ConvexHull(coords, false)

	ptset := sset.NewSSet(PointIndexCmp)
	for _, pt := range convex_hull {
		ptset.Add(pt)
	}

	g := hull_geom(convex_hull)
	return &HullNode{
		Pln:   polyline,
		Range: rng,
		Geom:  g,
		PtSet: ptset,
	}
}

//implements igeom interface
func (o *HullNode) Geometry() geom.Geometry {
	return o.Geom
}

//implements bbox interface
func (h *HullNode) BBox() *mbr.MBR {
	return h.Geom.BBox()
}

//stringer interface
func (h *HullNode) String() string {
	return h.Geom.WKT()
}

//stringer interface
func (h *HullNode) Coordinates() []*geom.Point {
	return h.Pln.Coordinates
}

//as segment
func (h *HullNode) Segment() *seg.Seg {
	coords := h.Coordinates()
	i, j := h.Range.I(), h.Range.J()
	return seg.NewSeg(coords[i], coords[j], i, j)
}

//as segment
func (h *HullNode) SubPolyline() *pln.Polyline {
	return h.Pln.SubPolyline(h.Range)
}

//hull geom
func hull_geom(coords []*geom.Point) geom.Geometry {
	var g geom.Geometry
	coords = coords[:]

	if len(coords) > 2 {
		g = geom.NewPolygon(coords)
	} else if len(coords) == 2 {
		g = geom.NewLineString(coords)
	} else {
		g = coords[0].Clone()
	}
	return g
}
