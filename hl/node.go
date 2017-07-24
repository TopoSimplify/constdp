package hl
import (
	"simplex/geom"
	"simplex/geom/mbr"
	"simplex/struct/sset"
	"simplex/constdp/rng"
	"simplex/constdp/ln"
	"simplex/constdp/cmp"
	"simplex/constdp/seg"
)

//hull node
type HullNode struct {
	Pln    *ln.Polyline
	Range  *rng.Range
	PRange *rng.Range
	Geom   geom.Geometry
	PtSet  *sset.SSet
}

//New Hull Node
func NewHullNode(pln *ln.Polyline, rng, prng *rng.Range) *HullNode {
	coords := make([]*geom.Point, 0)
	for _, i := range rng.Stride() {
		x, y, idx := pln.Coords[i][0], pln.Coords[i][1], float64(i)
		coords = append(coords, geom.NewPointXYZ(x, y, idx))
	}

	convex_hull := geom.ConvexHull(coords, false)

	ptset := sset.NewSSet(cmp.PointIndexCmp)
	for _, pt := range convex_hull {
		ptset.Add(pt)
	}

	g := hull_geom(convex_hull)
	return &HullNode{
		Pln:    pln,
		Range:  rng,
		PRange: prng,
		Geom:   g,
		PtSet:  ptset,
	}
}

//bbox interface
func (h *HullNode) BBox() *mbr.MBR {
	return h.Geom.BBox()
}

//stringer interface
func (h *HullNode) String() string {
	return h.Geom.WKT()
}

//stringer interface
func (h *HullNode) Coordinates() []*geom.Point {
	return h.Pln.Coords
}

//as segment
func (h *HullNode) Segment() *seg.Seg {
	coords := h.Coordinates()
	i, j := h.Range.I(), h.Range.J()
	return seg.NewSeg(coords[i], coords[j], i, j)
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
