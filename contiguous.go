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
	pln           := a.Pln
	ga            := a.Geom
	gb            := b.Geom
	bln_at_vertex := false
	inter_count   := 0

	bln := ga.Intersects(gb)
	if bln {
		interpts := ga.Intersection(gb)

		ai_pt := pln.Coords[a.Range.I()]
		aj_pt := pln.Coords[a.Range.J()]

		bi_pt := pln.Coords[b.Range.I()]
		bj_pt := pln.Coords[b.Range.J()]

		inter_count = len(interpts)

		for _, pt := range interpts {
			bln_aseg := pt.Equals2D(ai_pt) || pt.Equals2D(aj_pt)
			bln_bseg := pt.Equals2D(bi_pt) || pt.Equals2D(bj_pt)

			if bln_aseg || bln_bseg {
				bln_at_vertex = aj_pt.Equals2D(bi_pt) ||
								aj_pt.Equals2D(bj_pt) ||
								ai_pt.Equals2D(bj_pt) ||
								ai_pt.Equals2D(bi_pt)
			}

			if bln_at_vertex {
				break
			}
		}
	}

	return bln, bln_at_vertex, inter_count
}
