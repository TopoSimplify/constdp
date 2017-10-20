package constdp

import (
	"simplex/node"
	"github.com/intdxdt/geom"
)

//Is hull_a collapsible with respect to hull_b
//hull_a and hull_b should be contiguous
func is_contig_hull_collapsible(ha, hb *node.Node) bool {
	//segments are already collapsed
	if ha.Range.Size() == 1 {
		return true
	}

	var pln = ha.Coordinates()
	var pt_at = func(i int) *geom.Point {
		return geom.NewPoint(pln[i][:2])
	}

	var ra = ha.Range
	var rb = hb.Range
	var ai, aj = pt_at(ra.I()), pt_at(ra.J())
	var bi, bj = pt_at(rb.I()), pt_at(rb.J())

	var c *geom.Point
	if ai.Equals2D(bi) || aj.Equals2D(bi) {
		c = bi
	} else if ai.Equals2D(bj) || aj.Equals2D(bj) {
		c = bj
	} else {
		return true
	}

	var t = bj
	if c.Equals2D(t) {
		t = bi
	}
	var ply = ha.Geom.(*geom.Polygon)
	return !ply.Shell.PointCompletelyInRing(t)
}
