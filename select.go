package constdp

import (
	"simplex/geom"
	"simplex/geom/mbr"
	"simplex/constdp/db"
	"simplex/constdp/box"
	"simplex/struct/rtree"
	"simplex/constdp/opts"
)


//find context deformation list
func (self *ConstDP) select_deformation_candidates (hulldb *rtree.RTree, hull *HullNode ) []*HullNode {
	seldict    := make(map[[2]int]*HullNode, 0)
	predicate  := hull_predicate(hull,  1.e-5)
	ctxs       := db.KNN(hulldb, hull, 1.e-5, func(_, item rtree.BoxObj) float64 {
		var other geom.Geometry
		if o, ok := item.(*mbr.MBR); ok {
			other = box.MBRToPolygon(o)
		} else {
			other = item.(*HullNode).Geom
		}
		return hull.Geom.Distance(other)
	}, predicate)

	// for each item in the context list
	for _, h := range ctxs {
		// find which item to deform against current hull
		selection := make([]*HullNode, 0)

		h := h.(*HullNode)
		inters, contig, n := is_contiguous(hull, h)

		if inters {
			sels :=  []*HullNode{}
			if  contig && n > 1 {
				sels = _contiguous_candidates(hull,  h, self.Opts)
			} else if !contig {
				sels = _non_contiguous_candidates(hull,  h, self.Opts)
			}
			for _, s := range sels {
				selection = append(selection, s)
			}
		}
		// add candidate deformation hulls to selection list
		for _, s := range selection {
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
func _contiguous_candidates(a, b *HullNode, opts *opts.Opts) []*HullNode{
		var selection = make([]*HullNode, 0)
		// compute sidedness relation between contiguous hulls to avoid hull flip
		hulls := sort_hulls([]*HullNode{a, b})
		//future should not affect the past
		ha, hb := hulls[0], hulls[1]

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
func _non_contiguous_candidates(a, b *HullNode, opts *opts.Opts) []*HullNode {
	aseg := a.Pln.Segment(a.Range)
	bseg := b.Pln.Segment(b.Range)

	aln := a.Pln.SubPolyline(a.Range)
	bln := b.Pln.SubPolyline(b.Range)

	aseg_geom := aseg.Segment
	bseg_geom := bseg.Segment

	aln_geom := aln.Geom
	bln_geom := bln.Geom

	aseg_inters_bseg := aseg_geom.Intersects(bseg_geom)
	aseg_inters_bln  := aseg_geom.Intersects(bln_geom)
	bseg_inters_aln  := bseg_geom.Intersects(aln_geom)
	aln_inters_bln   := aln_geom.Intersects(bln_geom)

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
				if delta > opts.RelaxDist {
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
		if h.Range.Size() > 1 {
			*selection = append(*selection, h)
		}
	}
}