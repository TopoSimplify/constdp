package constdp

import (
	"simplex/geom"
	"simplex/geom/mbr"
	"simplex/struct/sset"
	"simplex/util/math"
)

type HPt struct {
	*geom.Point
}

//implements item.Item
func (o *HPt) Compare(pt *HPt) int {
	d := o.Point[2] - pt.Point[2]
	if math.FloatEqual(d, 0.0) {
		return 0
	} else if d < 0 {
		return -1
	}
	return 1
}

//String - implements stringer
func (o *HPt) String() string {
	return o.Point.WKT()
}

//constructs a hull node
type HullNode struct {
	Pln    *Polyline
	Range  *Range
	PRange *Range
	Geom   geom.Geometry
	PtSet  *sset.SSet
}

func NewHullNode(pln *Polyline, rng, prng *Range) *HullNode {
	x, y := 0, 1
	coords := pln.coords
	hull_coords := make([]*geom.Point, 0)
	for _, i := range rng.Stride() {
		hull_coords = append(
			hull_coords,
			geom.NewPointXYZ(coords[i][x], coords[i][y], float64(i)),
		)
	}

	cvxhull := geom.ConvexHull(coords, false)
	hull_g := hull_geom(cvxhull)

	ptset := sset.NewSSet(hpt_cmp)
	for _, pt := range cvxhull {
		ptset.Add(&HPt{pt})
	}

	self := &HullNode{
		Pln:    pln,
		Range:  rng,
		PRange: prng,
		Geom:   hull_g,
	}
	return self
}

//bbox interface
func (h *HullNode) BBox() *mbr.MBR {
	return h.Geom.BBox()
}

//stringer interface
func (h *HullNode) String() string {
	return h.Geom.WKT()
}

//as segment
func (h *HullNode) Segment() *Seg {
	coords := h.Pln.coords
	i, j := h.Range.i, h.Range.j
	return NewSeg(coords[i], coords[j], i, j)
}

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
