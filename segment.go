package constdp

import (
	"simplex/seg"
	"simplex/lnr"
)

//hull segment
func hull_segment(self lnr.Linear, hull *HullNode) *seg.Seg {
	coords := self.Coordinates()
	a, b := coords[hull.Range.I()], coords[hull.Range.J()]
	return seg.NewSeg(a, b, hull.Range.I(), hull.Range.J())
}
