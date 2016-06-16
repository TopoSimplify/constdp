package constdp

import (
    . "simplex/dp"
    . "simplex/geom"
    . "simplex/homotopy"
    . "simplex/relations"
    . "simplex/interest"
)

//Type DP
type ConstDP struct {
    *DP
    intersections []*Point
    defln         *LineDeflection
    opts          *Options
    homos         *Homotopy
}

//Creates a new constrained DP Simplification instance
func NewConstDP(options *Options, build bool) *ConstDP {
    var self = &ConstDP{
        DP              : NewDP(options, false),
        intersections   : make([]*Point, 0),
        defln           : NewLineDeflection(),
        opts            : options,
        homos           : NewHomotopy(
            []*Point{}, &IntCandidates{},
            []Relations{}, []Geometry{},
        ),
    }

    self.opts = options
    if build {
        self.Build()
    }
    return self
}
