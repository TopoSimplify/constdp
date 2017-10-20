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

	pln := ha.Coordinates()
	pt_at := func(i int) *geom.Point {
		return geom.NewPoint(pln[i][:2])
	}

	ra := ha.Range
	rb := hb.Range
	ai, aj := pt_at(ra.I()), pt_at(ra.J())
	bi, bj := pt_at(rb.I()), pt_at(rb.J())

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
	ply := ha.Geom.(*geom.Polygon)
	return !ply.Shell.PointCompletelyInRing(t)
}
