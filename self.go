package constdp

import (
	"simplex/geom"
	"simplex/constdp/ln"
	"simplex/constdp/seg"
	"simplex/struct/sset"
	"simplex/struct/rtree"
	"simplex/constdp/cmp"
	"simplex/constdp/ctx"
)

const (
	x = 0
	y = 1
)

type kvCount struct {
	count   int
	indxset *sset.SSet
}

type selfInter struct {
	keyset *sset.SSet
	point  *geom.Point
}

func update(dict map[[2]float64]*kvCount, o *geom.Point, index int) {
	k := [2]float64{o[x], o[y]}
	v, ok := dict[k]
	if !ok {
		v = &kvCount{
			count:   0,
			indxset: sset.NewSSet(cmp.IntCmp,8),
		}
		dict[k] = v
	}
	v.indxset.Add(index)
	v.count += 1
}

func LinearSelfIntersection(pln *ln.Polyline) []*ctx.CtxGeom {
	var tree = *rtree.NewRTree(8)
	var dict = make(map[[2]float64]*kvCount)

	var data = make([]rtree.BoxObj, 0)
	for _, s := range pln.Segments() {
		data = append(data, s)
	}
	tree.Load(data)

	self_intersects := make(map[string]*selfInter, 0)

	for _, d := range data {
		s := d.(*seg.Seg)
		res := tree.Search(s.BBox())
		update(dict, s.A, s.I)
		update(dict, s.B, s.J)

		for _, node := range res {
			other_seg := node.GetItem().(*seg.Seg)

			if s == other_seg {
				continue
			}
			seg_g, other_seg_g := s.Segment, other_seg.Segment
			intersects := seg_g.Intersection(other_seg_g)

			if len(intersects) == 0 {
				continue
			}

			for _, pt := range intersects {
				if s.A.Equals2D(pt) || s.B.Equals2D(pt) {
					continue
				}
				skey := sset.NewSSet(cmp.IntCmp,8).Extend(s.I, s.J, other_seg.I, other_seg.J)

				k := skey.String()
				v, ok := self_intersects[k]
				if !ok {
					v = &selfInter{
						keyset: skey,
						point:  pt,
					}
				}
				self_intersects[k] = v
			}
		}
	}

	results := make([]*ctx.CtxGeom, 0)
	for _, val := range self_intersects {
		cg := ctx.NewCtxGeom(val.point, 0, -1).AsSelfNonVertex()
		cg.Meta.SelfNonVertices = val.keyset
		results = append(results, cg)
	}

	for k, v := range dict {
		if v.count > 2 {
			cg := ctx.NewCtxGeom(geom.NewPoint(k[:]), 0, -1).AsSelfVertex()
			cg.Meta.SelfVertices = v.indxset
			results = append(results, cg)
		}
	}

	return results
}
