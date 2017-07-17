package constdp

import (
	"sort"
	"simplex/struct/rtree"
)

//select deform hull
func sel_deform_hull(a, b *HullNode, opts *Opts) []*HullNode {
	aseg := a.Pln.Segment(a.Range)
	bseg := b.Pln.Segment(b.Range)

	aln := a.Pln.SubPolyline(a.Range).geom
	bln := b.Pln.SubPolyline(b.Range).geom

	aseg_inters_bseg := aseg.Intersects(bseg)
	aseg_inters_bln := aseg.Intersects(bln)
	bseg_inters_aln := bseg.Intersects(aln)
	aln_inters_bln := aln.Intersects(bln)

	if aseg_inters_bseg && aseg_inters_bln && (!aln_inters_bln) {
		return []*HullNode{a}
	} else if aseg_inters_bseg && bseg_inters_aln && (!aln_inters_bln) {
		return []*HullNode{b}
	} else if aln_inters_bln {

		// find out whether is a shared vertex or overlap
		// is aseg inter bset  --- dist --- aln inter bln > relax dist
		pt_lns := aln.Intersection(bln)
		at_seg := aseg.Intersection(bseg)

		// if segs are disjoint but lines intersect, deform a&b
		if len(at_seg) == 0 && len(pt_lns) > 0 {
			return []*HullNode{a, b}
		}

		for _, ptln := range pt_lns {
			for _, ptseg := range at_seg {
				delta := ptln.Distance(ptseg)
				if delta > opts.RelaxDist {
					return []*HullNode{a, b}
				}
			}
		}
	}
	return []*HullNode{}
}

//returns bool (intersects), bool(is contig at vertex)
func is_hull_contiguous_at_vertex(a, b *HullNode) (bool, bool, int) {
	pln := a.Pln
	ga := a.Geom
	gb := b.Geom
	bln_at_vertex := false
	inter_count := 0

	bln := ga.Intersects(gb)
	if bln {
		interpts := ga.Intersection(gb)

		ai_pt := pln.coords[a.Range.i]
		aj_pt := pln.coords[a.Range.j]

		bi_pt := pln.coords[b.Range.i]
		bj_pt := pln.coords[b.Range.j]

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

func select_hulls_to_deform(a, b *HullNode, opts *Opts) []*HullNode {
	deformlist := make([]*HullNode, 0)
	intersects, at_contig_vertex, n := is_hull_contiguous_at_vertex(a, b)

	if intersects && (! at_contig_vertex) {
		sels := sel_deform_hull(a, b, opts)
		for _, s := range sels {
			deformlist = append(deformlist, s)
		}
	} else if intersects && at_contig_vertex && n > 1 {
		// compute sidedness relation between contiguous hulls to avoid hull flip
		hulls := []*HullNode{a, b}
		sort.Sort(HullNodes(hulls))
		//future should not affect the past
		ha, hb := hulls[0], hulls[1]
		hs := NewHullCollapseSidedness(hb.PtSet)
		bln := hs.IsValid(ha.PtSet)
		if !bln {
			deformlist = append(deformlist, hb)
		}
	}

	return deformlist
}

//find context deformation list
func find_hull_deformation_list(hulldb *rtree.RTree, hull *HullNode, opts *Opts) []*HullNode {
	ctxs := dbKNN(hulldb, hull, 1.e-5)
	selections := make(map[[2]int]*HullNode, 0)

	// for each item in the context list
	for _, ctx := range ctxs {
		h := ctx.(*HullNode)
		// find which item to deform against current hull
		selns := select_hulls_to_deform(hull, h, opts)
		// add candidate deformation hulls to selection list
		for _, s := range selns {
			selections[s.Range.AsArray()] = s
		}
	}

	items := make([]*HullNode, 0)
	for _, v := range selections {
		items = append(items, v)
	}
	return items
}
