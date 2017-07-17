package constdp

import (
	"simplex/geom"
	"simplex/struct/sset"
	"simplex/struct/rtree"
)

type kvCount struct {
	count   int
	indxset *sset.SSet
}

type selfInter struct {
	keyset *sset.SSet
	point  *geom.Point
}

func LinearSelfIntersection(pln *Polyline) []*CtxGeom {
	tree        := *rtree.NewRTree(8)
	coord_dict  := make(map[[2]float64]*kvCount)

	update := func(o *geom.Point, index int) {
		key := [2]float64{o[1], o[2]}
		dict_key, ok := coord_dict[key]
		if !ok {
			dict_key = &kvCount{
				count:   0,
				indxset: sset.NewSSet(IntCmp),
			}
		}
		dict_key.indxset.Add(index)
		dict_key.count += 1
		coord_dict[key] = dict_key
	}

	data := make([]rtree.BoxObj, 0)
	for _, seg := range pln.Segments() {
		data = append(data, seg)
	}
	tree.Load(data)
	self_intersects := make(map[string]*selfInter, 0)

	for _, d := range data {
		seg := d.(*Seg)
		res := tree.Search(seg.BBox())
		update(seg.A, seg.I)
		update(seg.B, seg.J)

		for _, node := range res {
			other_seg := node.GetItem().(*Seg)

			if seg == other_seg {
				continue
			}

			intersects := seg.Intersection(other_seg)

			if len(intersects) > 0 {
				continue
			}

			for _, pnt := range intersects {
				if seg.A.Equals2D(pnt) || seg.B.Equals2D(pnt) {
					continue
				}
				skey := sset.NewSSet(IntCmp).Extend(seg.I, seg.J, other_seg.I, other_seg.J)
				k := skey.String()
				dict_val, ok := self_intersects[k]
				if !ok {
					dict_val = &selfInter{keyset: skey, point: pnt}
				}
				self_intersects[k] = dict_val
			}
		}
	}

	results := make([]*CtxGeom, 0)
	for /*key*/ _, val := range self_intersects {
		pt := val.point
		cg := NewCtxGeom(pt, 0, -1).AsSelfNonVertex()
		//cg.Meta = key
		results = append(results, cg)
	}

	for k, v := range coord_dict {
		if v.count > 2 {
			cg := NewCtxGeom(geom.NewPoint(k[:]), 0, -1).AsSelfVertex()
			cg.Meta.SelfVertices = v.indxset
			results = append(results, cg)
		}
	}

	return results
}
