package constdp

import (
	"simplex/seg"
	"simplex/lnr"
	"simplex/node"
)

//hull segment
func hull_segment(self lnr.Linear, hull *node.Node) *seg.Seg {
	coords := self.Coordinates()
	a, b := coords[hull.Range.I()], coords[hull.Range.J()]
	return seg.NewSeg(a, b, hull.Range.I(), hull.Range.J())
}
