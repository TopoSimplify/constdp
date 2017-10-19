package offset

import (
	"simplex/constdp/ln"
	"simplex/constdp/rng"
	"github.com/intdxdt/geom"
	"github.com/intdxdt/vect"
)

//@formatter:off

//euclidean offset distance from dp - archor line [i, j] to maximum
//vertex at i < k <= j - not maximum offset is may not  be perpendicular
func MaxOffset(lnr ln.Linear, rng *rng.Range) (int, float64) {
	var pln           = lnr.Coordinates()
	var seg           = geom.NewSegment(pln[rng.I()], pln[rng.J()])
	var index, offset = rng.J(), 0.0

	if rng.Size() > 1 {
		for _, k := range rng.ExclusiveStride(1) {
			dist := seg.DistanceToPoint(pln[k])
			if dist >= offset {
				index, offset = k, dist
			}
		}
	}
	return index, offset
}


//computes Synchronized Euclidean Distance
func MaxSEDOffset(lnr ln.Linear, rng *rng.Range) (int, float64) {
	var t               = 2
	var pln             = lnr.Coordinates()
	var index, offset   = rng.J(), 0.0
	var a, b            = pln[rng.I()], pln[rng.J()]
	var segvect         = vect.NewVect(&vect.Options{
		A: a, B: b, At: &a[t], Bt: &b[t],
	})

	if rng.Size() > 1 {
		for _, k := range rng.ExclusiveStride(1) {
			sedvect := segvect.SEDVector(pln[k], pln[k][t])
			dist := sedvect.Magnitude()
			if dist >= offset {
				index, offset = k, dist
			}
		}
	}
	return index, offset
}
