package constdp

import (
	"simplex/geom"
	"simplex/struct/sset"
	"simplex/struct/rtree"
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
			indxset: sset.NewSSet(IntCmp,8),
		}
		dict[k] = v
	}
	v.indxset.Add(index)
	v.count += 1
}

func LinearSelfIntersection(pln *Polyline) []*CtxGeom {
	var tree = *rtree.NewRTree(8)
	var dict = make(map[[2]float64]*kvCount)

	var data = make([]rtree.BoxObj, 0)
	for _, seg := range pln.Segments() {
		data = append(data, seg)
	}
	tree.Load(data)

	self_intersects := make(map[string]*selfInter, 0)

	for _, d := range data {
		seg := d.(*Seg)
		res := tree.Search(seg.BBox())
		update(dict, seg.A, seg.I)
		update(dict, seg.B, seg.J)

		for _, node := range res {
			other_seg := node.GetItem().(*Seg)

			if seg == other_seg {
				continue
			}
			seg_g, other_seg_g := seg.Segment, other_seg.Segment
			intersects := seg_g.Intersection(other_seg_g)

			if len(intersects) == 0 {
				continue
			}

			for _, pt := range intersects {
				if seg.A.Equals2D(pt) || seg.B.Equals2D(pt) {
					continue
				}
				skey := sset.NewSSet(IntCmp,8).Extend(seg.I, seg.J, other_seg.I, other_seg.J)

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

	results := make([]*CtxGeom, 0)
	for _, val := range self_intersects {
		cg := NewCtxGeom(val.point, 0, -1).AsSelfNonVertex()
		cg.Meta.SelfNonVertices = val.keyset
		results = append(results, cg)
	}

	for k, v := range dict {
		if v.count > 2 {
			cg := NewCtxGeom(geom.NewPoint(k[:]), 0, -1).AsSelfVertex()
			cg.Meta.SelfVertices = v.indxset
			results = append(results, cg)
		}
	}

	return results
}
