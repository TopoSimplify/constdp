package offset

import (
	"simplex/geom"
	"simplex/constdp/ln"
	"simplex/constdp/rng"
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
