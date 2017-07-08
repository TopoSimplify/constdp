package constdp

import "simplex/geom"

type Seg struct {
	*geom.Segment
	I int
	J int
}

//New Segment constructor
func NewSeg(a, b *geom.Point, i, j int) *Seg {
	return &Seg{
		Segment: &geom.Segment{A: a, B: b}, I: i, J: j,
	}
}
