package seg

import (
	"simplex/geom"
	"simplex/constdp/rng"
)

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

func (s *Seg) Range() *rng.Range {
	return rng.NewRange(s.I, s.J)
}
