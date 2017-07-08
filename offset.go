package constdp

import (
	"simplex/geom"
)

//euclidean offset distance from dp - archor line [i, j] to maximum
//vertex at i < k <= j - not maximum offset is may not  be perpendicular
func MaxOffset(lnr Linear, rng *Range) (int, float64) {
	pln := lnr.Coordinates()
	seg := geom.NewSegment(pln[rng.i] , pln[rng.j])
	index, offset := rng.j, 0.0

	if rng.Size() > 1 {
		for _, k := range rng.ExclusiveStride() {
			dist := seg.DistanceToPoint(pln[k])
			if dist >= offset {
				index, offset = k, dist
			}
		}
	}
	return index, offset
}
