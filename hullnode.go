package constdp

import (
	"simplex/geom"
	"simplex/geom/mbr"
	"simplex/struct/sset"
)

//constructs a hull node
type HullNode struct {
	Pln    *Polyline
	Range  *Range
	PRange *Range
	Geom   geom.Geometry
	PtSet  *sset.SSet
}

//New Hull Node
func NewHullNode(pln *Polyline, rng, prng *Range) *HullNode {
	coords := make([]*geom.Point, 0)
	for _, i := range rng.Stride() {
		x, y, idx := pln.Coords[i][0], pln.Coords[i][1], float64(i)
		coords = append(coords, geom.NewPointXYZ(x, y, idx))
	}

	convex_hull := geom.ConvexHull(coords, false)

	ptset := sset.NewSSet(PointIndexCmp)
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
func (h *HullNode) Segment() *Seg {
	coords := h.Coordinates()
	i, j := h.Range.i, h.Range.j
	return NewSeg(coords[i], coords[j], i, j)
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
