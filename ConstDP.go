package constdp

import (
    . "simplex/dp"
    . "simplex/geom"
)

//Type DP
type ConstDP struct {
    *DP
    intersections   []*Point
    defln           *LineDeflection
}

//Creates a new constrained DP Simplification instance
func NewConstDP(options *Options, build bool) *ConstDP{
    var self = &ConstDP{
        NewDP(options, false),
        make([]*Point, 0),
        NewLineDeflection(),
    }

    fn := options.Process
    if build {
        self.Build(fn)
    }
    return self
}
