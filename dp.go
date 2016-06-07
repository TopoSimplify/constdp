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
    opts            *Options
}

//Creates a new constrained DP Simplification instance
func NewConstDP(options *Options, build bool) *ConstDP{
    var self = &ConstDP {
        NewDP(options, false),
        make([]*Point, 0),
        NewLineDeflection(),
        options,
    }

    self.opts =  options
    if build {
        self.Build()
    }
    return self
}
