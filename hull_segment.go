package constdp

//hull segment
func (self *ConstDP) HullSegment(hull *HullNode) *Seg {
	coords := self.Coordinates()
	a, b := coords[hull.Range.i], coords[hull.Range.j]
	return NewSeg(a, b, hull.Range.i, hull.Range.j)
}
