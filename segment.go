package constdp

import (
	"simplex/seg"
	"simplex/lnr"
	"simplex/node"
)

//hull segment
func hull_segment(self lnr.Linear, hull *node.Node) *seg.Seg {
	var i, j = hull.Range.I(), hull.Range.J()
	var coords = self.Coordinates()
	return seg.NewSeg(coords[i], coords[j], i, j)
}
