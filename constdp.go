package constdp

import (
	"simplex/geom"
)

//Type DP
type ConstDP struct {
	Pln []*geom.Point
}

//Creates a new constrained DP Simplification instance
func NewConstDP(options *Opts, build bool) *ConstDP {
	return &ConstDP{}
}
