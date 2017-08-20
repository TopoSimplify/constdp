package offset

import (
	"simplex/geom"
	"simplex/constdp/ln"
	"simplex/constdp/rng"
	"simplex/vect"
)

//euclidean offset distance from dp - archor line [i, j] to maximum
//vertex at i < k <= j - not maximum offset is may not  be perpendicular
func MaxOffset(lnr ln.Linear, rng *rng.Range) (int, float64) {
	pln := lnr.Coordinates()
	seg := geom.NewSegment(pln[rng.I()], pln[rng.J()])
	index, offset := rng.J(), 0.0

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
	t               := 2
	pln             := lnr.Coordinates()
	index, offset   := rng.J(), 0.0
	a, b            := pln[rng.I()], pln[rng.J()]
	segvect         := vect.NewVect(&vect.Options{
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
