package constdp

func extract_neighbours(hull *HullNode, neighbs []*HullNode) (*HullNode, *HullNode) {
	var prev, nxt *HullNode
	var i, j = hull.Range.I(), hull.Range.J()
	for _, h := range neighbs {
		if h != hull {
			if i == h.Range.J() {
				prev = h
			} else if j == h.Range.I() {
				nxt = h
			}
		}
	}
	return prev, nxt
}

//returns bool (intersects), bool(is contig at vertex)
func is_contiguous(a, b *HullNode) (bool, bool, int) {
	//@formatter:off
	pln         := a.Pln
	coords      := pln.Coordinates()
	ga          := a.Geom
	gb          := b.Geom
	contig      := false
	inter_count := 0

	bln := ga.Intersects(gb)
	if bln {
		interpts := ga.Intersection(gb)

		ai_pt := coords[a.Range.I()]
		aj_pt := coords[a.Range.J()]

		bi_pt := coords[b.Range.I()]
		bj_pt := coords[b.Range.J()]

		inter_count = len(interpts)

		for _, pt := range interpts {
			bln_aseg := pt.Equals2D(ai_pt) || pt.Equals2D(aj_pt)
			bln_bseg := pt.Equals2D(bi_pt) || pt.Equals2D(bj_pt)

			if bln_aseg || bln_bseg {
				contig = aj_pt.Equals2D(bi_pt) ||
					     aj_pt.Equals2D(bj_pt) ||
					     ai_pt.Equals2D(bj_pt) ||
					     ai_pt.Equals2D(bi_pt)
			}

			if contig {
				break
			}
		}
	}

	return bln, contig, inter_count
}
