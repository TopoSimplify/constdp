package constdp

import (
	"simplex/struct/rtree"
)

//find context_geom deformable hulls
func (self *ConstDP) select_ftclass_deformation_candidates(hulldb *rtree.RTree, hull *HullNode) []*HullNode {
	var n int
	var inters, contig bool

	seldict := make(map[[2]int]*HullNode, 0)
	ctx_hulls := find_context_hulls(hulldb, hull, EpsilonDist)
	// for each item in the context_geom list
	for _, cn := range ctx_hulls {
		n = 0
		h := cast_as_hullnode(cn)

		same_feature := hull.DP == h.DP
		// find which item to deform against current hull
		if same_feature { // check for contiguity
			inters, contig, n = is_contiguous(hull, h)
		} else {
			// contiguity is by default false for different features
			contig = false
			ga, gb := hull.Geom, h.Geom

			inters = ga.Intersects(gb)
			if inters {
				interpts := ga.Intersection(gb)
				inters = len(interpts) > 0
				n = len(interpts)
			}
		}

		if !inters { // disjoint : nothing to do, continue
			continue
		}

		var sels = make([]*HullNode, 0)
		if contig && n > 1 { // contiguity with overlap greater than a vertex
			sels = self._contiguous_candidates(hull, h)
		} else if !contig {
			sels = self._non_contiguous_candidates(hull, h)
		}

		// add candidate deformation hulls to selection list
		for _, s := range sels {
			seldict[s.Range.AsArray()] = s
		}
	}

	items := make([]*HullNode, 0)
	for _, v := range seldict {
		items = append(items, v)
	}
	return items
}

//find context deformation list
func (self *ConstDP) select_deformation_candidates(hulldb *rtree.RTree, hull *HullNode) []*HullNode {
	seldict := make(map[[2]int]*HullNode, 0)
	ctx_hulls := find_context_hulls(hulldb, hull, EpsilonDist)

	// for each item in the context list
	for _, cn := range ctx_hulls {
		// find which item to deform against current hull
		h := cast_as_hullnode(cn)
		inters, contig, n := is_contiguous(hull, h)

		if !inters {
			continue
		}

		sels := make([]*HullNode, 0)
		if contig && n > 1 { //contiguity with overlap greater than a vertex
			sels = self._contiguous_candidates(hull, h)
		} else if !contig {
			sels = self._non_contiguous_candidates(hull, h)
		}

		// add candidate deformation hulls to selection list
		for _, s := range sels {
			seldict[s.Range.AsArray()] = s
		}
	}

	items := make([]*HullNode, 0)
	for _, v := range seldict {
		items = append(items, v)
	}
	return items
}

//select contiguous candidates
func (self *ConstDP) _contiguous_candidates(a, b *HullNode) []*HullNode {
	var selection = make([]*HullNode, 0)
	// compute sidedness relation between contiguous hulls to avoid hull flip
	hulls := NewHullNodes().Extend(a, b).Sort()
	//future should not affect the past
	ha, hb := hulls.list[0], hulls.list[1]

	//all hulls that are simple should be collapsible
	// if not collapsible -- add to selection for deformation
	// to reach collapsibility

	//& the present should not affect the future
	bln := is_contig_hull_collapsible(ha, hb)
	if !bln {
		selection = append(selection, ha)
	}

	//future should not affect the present
	bln = is_contig_hull_collapsible(hb, ha)
	if !bln {
		selection = append(selection, hb)
	}
	return selection
}

//select non-contiguous candidates
func (self *ConstDP) _non_contiguous_candidates(a, b *HullNode) []*HullNode {
	aseg := a.Segment()
	bseg := b.Segment()

	aln := a.SubPolyline()
	bln := b.SubPolyline()

	aseg_geom := aseg.Segment
	bseg_geom := bseg.Segment

	aln_geom := aln.Geom
	bln_geom := bln.Geom

	aseg_inters_bseg := aseg_geom.Intersects(bseg_geom)
	aseg_inters_bln := aseg_geom.Intersects(bln_geom)
	bseg_inters_aln := bseg_geom.Intersects(aln_geom)
	aln_inters_bln := aln_geom.Intersects(bln_geom)

	selection := []*HullNode{}
	if aseg_inters_bseg && aseg_inters_bln && (!aln_inters_bln) {
		_add_to_selection(&selection, a)
	} else if aseg_inters_bseg && bseg_inters_aln && (!aln_inters_bln) {
		_add_to_selection(&selection, b)
	} else if aln_inters_bln {
		// find out whether is a shared vertex or overlap
		// is aseg inter bset  --- dist --- aln inter bln > relax dist
		pt_lns := aln_geom.Intersection(bln_geom)
		at_seg := aseg.Intersection(bseg_geom)

		// if segs are disjoint but lines intersect, deform a&b
		if len(at_seg) == 0 && len(pt_lns) > 0 {
			_add_to_selection(&selection, a, b)
			return selection
		}

		for _, ptln := range pt_lns {
			for _, ptseg := range at_seg {
				delta := ptln.Distance(ptseg)
				if delta > self.Opts.RelaxDist {
					_add_to_selection(&selection, a, b)
					return selection
				}
			}
		}
	}
	return selection
}

//add to hull selection based on range size
// add if range size is greater than 1 : not a segment
func _add_to_selection(selection *[]*HullNode, hulls ...*HullNode) {
	for _, h := range hulls {
		//add to selection for deformation - if polygon
		if h.Range.Size() > 1 {
			*selection = append(*selection, h)
		}
	}
}
