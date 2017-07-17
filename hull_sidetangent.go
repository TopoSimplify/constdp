package constdp

import "simplex/struct/sset"

type HullSideTangent struct {
	aseg *Seg
	bseg *Seg
	rtan *Seg
	ltan *Seg
	hull *sset.SSet
	side *HullCollapseSidedness
}
