package constdp

import (
	"simplex/geom"
)

//Is hull_a collapsible with respect to hull_b
//hull_a and hull_b should be contiguous
func is_contig_hull_collapsible(ha, hb *HullNode) bool {
	pln := ha.Pln.Coordinates()
	pt := func(i int) *geom.Point {
		return geom.NewPoint(pln[i][:2])
	}

	ra := ha.Range
	rb := hb.Range
	ai, aj := pt(ra.I()), pt(ra.J())
	bi, bj := pt(rb.I()), pt(rb.J())

	var c *geom.Point
	if ai.Equals2D(bi) || aj.Equals2D(bi) {
		c = bi
	} else if ai.Equals2D(bj) || aj.Equals2D(bj) {
		c = bj
	} else {
		return true
	}

	t := bj
	if c.Equals2D(t) {
		t = bi
	}
	ply, ok := ha.Geom.(*geom.Polygon)
	if !ok {
		if _, ok = ha.Geom.(*geom.LineString); ok {
			return ok
		}
	}
	return !ply.Shell.PointCompletelyInRing(t)
}
