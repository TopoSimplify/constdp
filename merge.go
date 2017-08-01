package constdp

import (
	"sort"
	"simplex/geom"
	"simplex/geom/mbr"
	"simplex/constdp/ln"
	"simplex/struct/sset"
	"simplex/constdp/db"
	"simplex/constdp/box"
	"simplex/constdp/rng"
	"simplex/struct/rtree"
)

//merge contig hulls after split - merge line segment fragments
func MergeContigFragments(
	self ln.Linear, hulls []*HullNode, hulldb *rtree.RTree,
	vertex_set *sset.SSet, ) ([]*HullNode, []*HullNode) {
	pln := self.Polyline()
	keep, rm := make([]*HullNode, 0), make([]*HullNode, 0)

	hdict := make(map[[2]int]*HullNode, 0)
	for _, h := range hulls {
		hdict[h.Range.AsArray()] = h

		hs_knn := db.KNN(hulldb, h, 1.0e-5, func(_, item rtree.BoxObj) float64 {
			var other geom.Geometry
			if o, ok := item.(*mbr.MBR); ok {
				other = box.MBRToPolygon(o)
			} else {
				other = item.(*HullNode).Geom
			}
			return h.Geom.Distance(other)
		})

		hs := make([]*HullNode, len(hs_knn))
		for i, h := range hs_knn {
			hs[i] = h.(*HullNode)
		}

		hr := h.Range

		//if hr.Size() < 4{
		if hr.Size() == 1 {
			hs = ExceptHull(hs, h)
			sort_hulls(hs) // sort hulls for consistency

			for _, s := range hs {
				sr := s.Range
				bln :=  (hr.J() == sr.I() && vertex_set.Contains(sr.I())) ||
						(hr.I() == sr.J() && vertex_set.Contains(sr.J()))

				if !bln && (hr.Contains(sr.I()) || hr.Contains(sr.J())) {
					l := []int{sr.I(), sr.J(), hr.I(), hr.J()}
					sort.Ints(l)

					// rm sr + hr
					delete(hdict, sr.AsArray())
					delete(hdict, hr.AsArray())

					r := rng.NewRange(l[0], l[len(l)-1])
					m := NewHullNode(pln, r, r)

					// add merge
					hdict[m.Range.AsArray()] = m

					// add to remove list to remove , after merge
					rm = append(rm, s)
					rm = append(rm, h)
					break
				}
			}
		}
	}

	for _, v := range hdict {
		keep = append(keep, v)
	}
	return keep, rm
}

