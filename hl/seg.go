package hl

import (
	"simplex/constdp/ln"
	"simplex/constdp/seg"
)

//hull segment
func HullSegment(self ln.Linear, hull *HullNode) *seg.Seg {
	coords := self.Coordinates()
	a, b := coords[hull.Range.I()], coords[hull.Range.J()]
	return seg.NewSeg(a, b, hull.Range.I(), hull.Range.J())
}
